package lobby

import (
	"lobby/generics"
	"lobby/lobby/message"
	"lobby/lobby/player"
	"lobby/logger"

	libuuid "github.com/google/uuid"
)

type MemLobby struct {
	id   string
	name string

	playerMap *generics.TypedMap[player.Player]
	logger    logger.Logger
}

func NewMemLobby(name string, logger logger.Logger) *MemLobby {
	return &MemLobby{
		id:        libuuid.NewString(),
		name:      name,
		playerMap: generics.NewTypedMap[player.Player](),
		logger:    logger,
	}
}

func (l *MemLobby) Id() string {
	return l.id
}

func (l *MemLobby) Name() string {
	return l.name
}

func (l *MemLobby) AddPlayer(p *player.Player) {
	l.playerMap.AddPtr(p.Id(), p)
}

func (l *MemLobby) FindPlayer(id string) (*player.Player, error) {
	return l.playerMap.ItemPtr(id)
}

func (l *MemLobby) BroadcastNotification(n *message.Notification) error {
	return l.playerMap.RangePtr(func(p *player.Player) error {
		if !p.HasConnection() {
			return nil
		}

		conn := p.Connection()
		if err := conn.WriteJSON(n); err != nil {
			l.logger.Warn(err)
			defer conn.Close()
			p.SetConnection(nil)
		}
		return nil
	})
}
