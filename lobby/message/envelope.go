package message

import (
	"slices"
)

// this will make things simple for C#, C++ clients
type Message struct {
	Key   string
	Value string
}

type Envelope struct {
	Direction byte
	Flag      byte
	Messages  []Message
}

func NewEnvelope(direction byte) *Envelope {
	return &Envelope{
		Direction: direction,
		Flag:      0,
		Messages:  make([]Message, 0),
	}
}

func (e *Envelope) SetFlag(whitchBit byte) {
	e.Flag |= whitchBit
}

func (e *Envelope) GetFlag(whitchBit byte) bool {
	return e.Flag&whitchBit == whitchBit
}

func (e *Envelope) SetMessage(key string, value string) {
	e.Messages = append(e.Messages, Message{Key: key, Value: value})
}

func (e *Envelope) GetMessage(key string) (string, bool) {
	idx := slices.IndexFunc(e.Messages, func(m Message) bool {
		return m.Key == key
	})
	if idx < 0 {
		return "", false
	}

	return e.Messages[idx].Value, true
}
