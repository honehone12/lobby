package message

func NewPlayerJoinNotification(name string) *Envelope {
	return &Envelope{
		Direction: Notification,
		Flag:      Join,
		Messages:  []Message{{Key: "player-name", Value: name}},
	}
}

func NewChatMessageRequest(message string) *Envelope {
	return &Envelope{
		Direction: Notification,
		Flag:      Chat,
		Messages:  []Message{{Key: "chat-message", Value: message}},
	}
}
