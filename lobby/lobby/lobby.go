package lobby

import (
	"lobby/lobby/message"
	"lobby/lobby/player"
	"time"
)

const (
	LobbyPingInterval = time.Second
)

type Lobby interface {
	Id() string
	Name() string
	PlayerCount() uint
	ActiveCount() uint

	GetPlayers() ([]player.PlayerInfo, error)
	AddPlayer(*player.Player)
	FindPlayer(string) (*player.Player, error)

	BroadcastNotification(*message.Notification) error
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
