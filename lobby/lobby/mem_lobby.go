package lobby

import (
	"lobby/generics"
	"lobby/lobby/message"
	"lobby/lobby/player"
	"lobby/logger"
	"time"

	libuuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MemLobby struct {
	id   string
	name string

	activeCount uint
	playerMap   *generics.TypedMap[player.Player]
	ticker      time.Ticker
	logger      logger.Logger
	errCh       chan error
}

func NewMemLobby(name string, logger logger.Logger) *MemLobby {
	l := &MemLobby{
		id:          libuuid.NewString(),
		name:        name,
		activeCount: 0,
		playerMap:   generics.NewTypedMap[player.Player](),
		ticker:      *time.NewTicker(LobbyPingInterval),
		logger:      logger,
		errCh:       make(chan error),
	}
	go l.ping()
	return l
}

func (l *MemLobby) Id() string {
	return l.id
}

func (l *MemLobby) Name() string {
	return l.name
}

func (l *MemLobby) PlayerCount() uint {
	return uint(l.playerMap.Count())
}

func (l *MemLobby) ActiveCount() uint {
	return l.activeCount
}

func (l *MemLobby) GetPlayers() ([]player.PlayerInfo, error) {
	buff := make([]player.PlayerInfo, l.playerMap.Count())
	iter := 0
	err := l.playerMap.RangePtr(func(p *player.Player) error {
		buff[iter].Id = p.Id()
		buff[iter].Name = p.Name()
		buff[iter].Active = p.HasConnection()
		iter++
		return nil
	})
	return buff, err
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

		if err := p.Connection().WriteJSON(n); err != nil {
			l.logger.Warn(err)
		}
		return nil
	})
}

func (l *MemLobby) ping() {
	for range l.ticker.C {
		activeCount := uint(0)
		if err := l.playerMap.RangePtr(func(p *player.Player) error {
			if !p.HasConnection() {
				return nil
			}

			conn := p.Connection()
			if err := conn.WriteMessage(
				websocket.PingMessage,
				message.PingBytes,
			); err != nil {
				l.logger.Warn(err)
				defer conn.Close()
				p.SetDisconnected()
			} else {
				activeCount++
			}

			return nil
		}); err != nil {
			l.errCh <- err
		} else {
			l.activeCount = activeCount
		}
	}
}
