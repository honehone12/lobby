package lobby

import (
	"lobby/lobby/player"
	"sync"

	libuuid "github.com/google/uuid"
)

type Lobby struct {
	id   string
	name string

	playerCount int
	playerMap   *sync.Map
}

func NewLobby(name string) *Lobby {
	return &Lobby{
		id:          libuuid.NewString(),
		name:        name,
		playerCount: 0,
		playerMap:   &sync.Map{},
	}
}

func (l *Lobby) Id() string {
	return l.id
}

func (l *Lobby) Name() string {
	return l.name
}

func (l *Lobby) AddPlayer(id string, p *player.Player) {
	if _, exists := l.playerMap.LoadOrStore(id, p); !exists {
		l.playerCount++
	}
}
