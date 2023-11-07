package lobby

import (
	"lobby/lobby/message"
	"lobby/lobby/player"
	"time"
)

const (
	LobbyConnectionTimeOut = time.Second
	LobbyPingInterval      = time.Second
	LobbyCleanUpInterval   = time.Second * 10
)

type Lobby interface {
	Id() string
	Name() string
	CreatedAt() int64
	PlayerCount() uint
	ActiveCount() uint

	GetPlayers() ([]player.PlayerInfo, error)
	AddPlayer(*player.Player)
	FindPlayer(string) (*player.Player, error)
	DeletePlayer(string)

	BroadcastMessage(*message.Envelope) error

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
