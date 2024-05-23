package mikro

type ButtonMessage struct {
	pressed uint64 // bitflag of all currently pressed buttons

	encoderTouched  bool
	encoderPosition uint8

	stripPos1 uint8
	stripPos2 uint8
}

func (b *ButtonMessage) Type() MessageType {
	return ButtonType
}

func (b *ButtonMessage) IsButtonPressed(btn Button) bool {
	return b.pressed&(1<<btn) != 0
}

func (b *ButtonMessage) PressedButtons() []Button {
	var pressed []Button

	for idx := 0; idx < 64; idx++ {
		if b.pressed&(1<<idx) != 0 {
			pressed = append(pressed, Button(idx))
		}
	}

	return pressed
}

func (b *ButtonMessage) IsEncoderTouched() bool {
	return b.encoderTouched
}

func (b *ButtonMessage) EncoderPosition() uint8 {
	return b.encoderPosition
}

func (b *ButtonMessage) StripPosition() uint8 {
	return b.stripPos1
}

func (b *ButtonMessage) StripSecondPosition() uint8 {
	return b.stripPos2
}
