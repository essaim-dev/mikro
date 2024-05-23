package mikro

type MessageType byte

const (
	ButtonType MessageType = iota + 1
	PadType
)

type Message interface {
	Type() MessageType
}
