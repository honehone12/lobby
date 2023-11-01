package player

import (
	"lobby/lobby/message"
)

type PlayerStore interface {
	Id() string
	Name() string

	AddPlayer(*Player)
	FindPlayer(string) (*Player, error)

	BroadcastNotification(*message.Notification) error
}
