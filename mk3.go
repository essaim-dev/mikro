package mikro

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"

	bp "essaim.dev/mikro/api/mk3"
	"github.com/karalabe/hid"
)

const (
	Mk3VID = 0x17cc
	Mk3PID = 0x1700

	buttonReport byte = 1
	padReport    byte = 2
)

type Mk3 struct {
	device hid.Device

	onPadFunc    func(msg PadMessage)
	onButtonFunc func(msg ButtonMessage)

	lights   Lights
	lightsMu sync.RWMutex
}

func OpenMk3() (*Mk3, error) {
	hids, _ := hid.Enumerate(Mk3VID, Mk3PID)
	if len(hids) == 0 {
		return nil, errors.New("mikro: no device found")
	}

	device, err := hids[0].Open()
	if err != nil {
		return nil, err
	}

	return &Mk3{
		device: device,
		lights: NewLights(),
	}, nil
}

func (m *Mk3) Close() error {
	return m.device.Close()
}

func (m *Mk3) SetOnPadFunc(fn func(msg PadMessage)) {
	m.onPadFunc = fn
}

func (m *Mk3) SetOnButtonFunc(fn func(msg ButtonMessage)) {
	m.onButtonFunc = fn
}

func (m *Mk3) Lights() Lights {
	m.lightsMu.RLock()
	defer m.lightsMu.RUnlock()

	return m.lights
}

func (m *Mk3) SetLights(lights Lights) error {
	m.lightsMu.Lock()

	m.lights = lights
	state := m.lights.toBP()

	m.lightsMu.Unlock()

	if _, err := m.device.Write(state.Encode()); err != nil {
		return fmt.Errorf("could not write updated light state: %w", err)
	}

	return nil
}

func (m *Mk3) SetScreen(img image.Image) error {
	i := image.NewPaletted(image.Rect(0, 0, 128, 32), color.Palette{color.Black, color.White})

	draw.FloydSteinberg.Draw(i, i.Bounds(), img, image.Pt(0, 0))

	bitPixels := imageToBit(i)

	stateHigh := bp.ScreenState{
		Magic1:        [3]byte{0xe0, 0x00, 0x00},
		ScreenPortion: 0x00,
		Magic2:        [5]byte{0x00, 0x80, 0x00, 0x02, 0x0},
		Pixels:        [256]byte(bitPixels[:256]),
	}
	stateLow := bp.ScreenState{
		Magic1:        [3]byte{0xe0, 0x00, 0x00},
		ScreenPortion: 0x02,
		Magic2:        [5]byte{0x00, 0x80, 0x00, 0x02, 0x0},
		Pixels:        [256]byte(bitPixels[256:]),
	}

	if _, err := m.device.Write(stateHigh.Encode()); err != nil {
		return fmt.Errorf("could not write updated higher screen state: %w", err)
	}
	if _, err := m.device.Write(stateLow.Encode()); err != nil {
		return fmt.Errorf("could not write updated lower screen state: %w", err)
	}

	return nil
}

func (m *Mk3) Run(ctx context.Context) error {
	b := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, err := m.device.Read(b)
			if err != nil {
				return fmt.Errorf("could not read message from device: %w", err)
			}
			switch b[0] {
			case buttonReport:
				if m.onButtonFunc != nil {
					m.onButtonFunc(m.decodeButtonMessage(b))
				}
			case padReport:
				if m.onPadFunc != nil {
					m.onPadFunc(m.decodePadMessage(b))
				}
			default:
				return fmt.Errorf("received unknown message type: %b", b[0])
			}
		}
	}
}

func (m *Mk3) decodeButtonMessage(buf []byte) ButtonMessage {
	report := bp.ButtonReport{}
	report.Decode(buf)

	return ButtonMessage{
		pressed:         report.PressedButtons,
		encoderPosition: report.EncoderValue,
		encoderTouched:  report.EncoderTouched,
		stripPos1:       report.StripValue1,
		stripPos2:       report.StripValue2,
	}
}

func (m *Mk3) decodePadMessage(buf []byte) PadMessage {
	report := bp.PadReport{}
	report.Decode(buf)

	return PadMessage{
		pad:      Pad(report.Pad),
		velocity: report.Velocity,
		actions:  report.Action,
	}
}

// Converts an image.Image to a byte slice where each pixel is represented by 1 bit
func imageToBit(img *image.Paletted) []byte {
	output := make([]byte, 512)

	for i := range 128 {
		for line := range 4 {
			byteVal := byte(0)
			for bit := range 8 {
				byteVal = byteVal << 1
				if img.At(i, (8*line)+(7-bit)) == color.Black {
					byteVal += 1
				}
			}
			output[128*line+i] = byteVal
		}
	}
	return output
}
