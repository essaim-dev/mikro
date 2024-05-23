package mikro

import (
	"errors"

	bp "github.com/antoi-ne/mikro/api/mk3"
	"github.com/karalabe/hid"
)

const (
	Mk3VID = 0x17cc
	Mk3PID = 0x1700
)

type Mk3 struct {
	device hid.Device
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
	}, nil
}

func (m *Mk3) Close() error {
	return m.device.Close()
}

func (m *Mk3) ReadMessage() (Message, error) {
	buf, err := m.read()
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		return nil, errors.New("empty message")
	}

	switch MessageType(buf[0]) {
	case ButtonType:
		return m.decodeButtonMessage(buf), nil
	case PadType:
		return m.decodePadMessage(buf), nil
	default:
		return nil, errors.New("unknown message type")
	}
}

func (m *Mk3) read() ([]byte, error) {
	buf := make([]byte, 1024)

	n, err := m.device.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

func (m *Mk3) decodeButtonMessage(buf []byte) *ButtonMessage {
	report := bp.ButtonReport{}
	report.Decode(buf)

	return &ButtonMessage{
		pressed:         report.PressedButtons,
		encoderPosition: report.EncoderValue,
		encoderTouched:  report.EncoderTouched,
		stripPos1:       report.StripValue1,
		stripPos2:       report.StripValue2,
	}
}

func (m *Mk3) decodePadMessage(buf []byte) *PadMessage {
	report := bp.PadReport{}
	report.Decode(buf)

	return &PadMessage{
		pad:      Pad(report.Pad),
		velocity: report.Velocity,
		action:   PadAction(report.Action),
	}
}
