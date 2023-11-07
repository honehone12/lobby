package lobby

import (
	"lobby/lobby/message"
	"lobby/lobby/player"
)

func ProcessEnvelope(l Lobby, p *player.Player, e *message.Envelope) error {
	if e.GetFlag(message.Chat) {
		msg, ok := e.GetMessage(message.ChatMessage)
		if !ok || len(msg) == 0 {
			return nil
		}

		n := message.NewChatMessageNotification(p.Name(), p.Id(), msg)
		if err := l.BroadcastMessage(n); err != nil {
			return err
		}
	}

	return nil
}
