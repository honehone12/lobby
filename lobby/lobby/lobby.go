package lobby

import (
	"lobby/lobby/message"
	"lobby/lobby/player"
	"time"
)

const (
	LobbyPingInterval    = time.Second
	LobbyCleanUpInterval = time.Second * 10
)

type Lobby interface {
	Id() string
	Name() string
	PlayerCount() uint
	ActiveCount() uint

	GetPlayers() ([]player.PlayerInfo, error)
	AddPlayer(*player.Player)
	FindPlayer(string) (*player.Player, error)
	DeletePlayer(string)

	BroadcastNotification(*message.Notification) error

	Delete()
}

type LobbySummary struct {
	Id          string
	Name        string
	PlayerCount uint
	ActiveCount uint
}

type LobbyDetail struct {
	Id          string
	Name        string
	PlayerCount uint
	ActiveCount uint
	PlayerList  []player.PlayerInfo
}
