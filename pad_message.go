package mikro

type PadAction uint8

const (
	PadActionPressed PadAction = iota + 1
	PadActionTouched
	PadActionReleased
)

type PadMessage struct {
	pad      Pad
	velocity uint16
	action   PadAction
}

func (p *PadMessage) Type() MessageType {
	return PadType
}

func (p *PadMessage) Pad() Pad {
	return p.pad
}

func (p *PadMessage) Velocity() uint16 {
	return p.velocity
}

func (p *PadMessage) Action() PadAction {
	return p.action
}
