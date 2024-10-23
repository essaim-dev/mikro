package mikro

import bp "github.com/antoi-ne/mikro/api/mk3"

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=Intensity -trimprefix=Intensity
type Intensity uint8

const (
	IntensityLow Intensity = iota
	IntensityOff
	IntensityMedium
	IntensityHigh
)

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=Color -trimprefix=Color
type Color uint8

const (
	ColorOff Color = iota
	ColorRed
	ColorOrange
	ColorLightOrange
	ColorWarmYellow
	ColorYellow
	ColorLime
	ColorGreen
	ColorMint
	ColorCyan
	ColorTurquoise
	ColorBlue
	ColorPlum
	ColorViolet
	ColorPurple
	ColorMagenta
	ColorFuchsia
	ColorWhite
)

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=ColorLevel -trimprefix=ColorLevel
type ColorLevel uint8

const (
	ColorLevelLow ColorLevel = iota
	ColorLevelMedium
	ColorLevelHigh
	ColorLevelFaded
)

type ColoredLight struct {
	Level ColorLevel
	Color Color
}

type Lights struct {
	Buttons [39]Intensity
	Pads    [16]ColoredLight
	Strip   [35]ColoredLight
}

func NewLights() Lights {
	l := Lights{}

	for idx := range l.Buttons {
		l.Buttons[idx] = IntensityOff
	}

	return l
}

func (l Lights) toBP() bp.LightState {
	state := bp.LightState{
		Magic: 0x80,
	}

	for idx, intensity := range l.Buttons {
		state.Buttons[idx] = bp.ColoredLight{
			Intensity: uint8(intensity),
			Color:     uint8(1),
		}
	}

	for idx, light := range l.Pads {
		state.Pads[idx] = bp.ColoredLight{
			Intensity: uint8(light.Level),
			Color:     uint8(light.Color),
		}
	}

	for idx, light := range l.Strip {
		state.Strip[idx] = bp.ColoredLight{
			Intensity: uint8(light.Level),
			Color:     uint8(light.Color),
		}
	}

	return state
}
