package message

func NewPlayerJoinNotification(name string, id string) *Envelope {
	return &Envelope{
		Direction: Notification,
		Flag:      Join,
		Messages: []Message{
			{Key: PlayerName, Value: name},
			{Key: PlayerId, Value: id},
		},
	}
}

func NewChatMessageRequest(msg string) *Envelope {
	return &Envelope{
		Direction: Request,
		Flag:      Chat,
		Messages:  []Message{{Key: ChatMessage, Value: msg}},
	}
}

func NewChatMessageNotification(name string, id string, msg string) *Envelope {
	return &Envelope{
		Direction: Notification,
		Flag:      Chat,
		Messages: []Message{
			{Key: PlayerName, Value: name},
			{Key: PlayerId, Value: id},
			{Key: ChatMessage, Value: msg},
		},
	}
}
