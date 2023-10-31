package lobby

import (
	"lobby/generics"
	"lobby/lobby/player"

	libuuid "github.com/google/uuid"
)

type MemLobby struct {
	id   string
	name string

	playerMap *generics.TypedMap[player.Player]
}

func NewLobby(name string) *MemLobby {
	return &MemLobby{
		id:        libuuid.NewString(),
		name:      name,
		playerMap: generics.NewTypedMap[player.Player](),
	}
}

func (l *MemLobby) Id() string {
	return l.id
}

func (l *MemLobby) Name() string {
	return l.name
}

func (l *MemLobby) AddPlayer(id string, p *player.Player) {
	l.playerMap.Add(id, p)
}

func (l *MemLobby) FindPlayer(id string) (*player.Player, error) {
	return l.playerMap.Item(id)
}
