package command

type Command interface {
	Write() []byte
	OnReply(buffer []byte) (bool, error)
}
