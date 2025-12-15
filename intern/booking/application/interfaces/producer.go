package interfaces

type Producer interface {
	Send(message []byte)
}
