package message

type Notification struct {
	Flag    byte
	Message map[string]string
}

const (
	JoinBit = 0b0000_0001
)

func NewNotification() *Notification {
	return &Notification{
		Flag:    0,
		Message: make(map[string]string),
	}
}

func (n *Notification) SetFlag(whitchBit byte) {
	n.Flag |= whitchBit
}

func (n *Notification) SetMessage(key string, value string) {
	n.Message[key] = value
}
