package lobby

import (
	"fmt"
	"lobby/generics"
	"lobby/lobby/message"
	"lobby/lobby/message/request"
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

	pingTicker    *time.Ticker
	closeChPing   chan bool
	closeChListen chan bool

	logger logger.Logger
}

const (
	SleepOnNoConnection = time.Millisecond * 100
)

func connectionErr(err error) error {
	return fmt.Errorf("disconnected the peer because of the previous error => %s", err)
}

func NewMemLobby(name string, logger logger.Logger) *MemLobby {
	l := &MemLobby{
		id:            libuuid.NewString(),
		name:          name,
		activeCount:   0,
		playerMap:     generics.NewTypedMap[player.Player](),
		pingTicker:    time.NewTicker(LobbyPingInterval),
		closeChPing:   make(chan bool),
		closeChListen: make(chan bool),
		logger:        logger,
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
	if err != nil {
		return nil, err.E
	}
	return buff, nil
}

func (l *MemLobby) AddPlayer(p *player.Player) {
	l.playerMap.AddPtr(p.Id(), p)
	go l.listen(p)
}

func (l *MemLobby) FindPlayer(id string) (*player.Player, error) {
	return l.playerMap.ItemPtr(id)
}

func (l *MemLobby) DeletePlayer(id string) {
	l.playerMap.Delete(id)
}

func (l *MemLobby) BroadcastMessage(e *message.Envelope) error {
	err := l.playerMap.RangePtr(func(p *player.Player) error {
		if !p.HasConnection() {
			return nil
		}

		conn := p.Connection()
		t := time.Now().Add(LobbyConnectionTimeOut)
		if err := conn.SetWriteDeadline(t); err != nil {
			p.Close()
			l.logger.Error(connectionErr(err))
			return nil
		}

		if err := conn.WriteJSON(e); err != nil {
			p.Close()
			l.logger.Warn(connectionErr(err))
		}

		return nil
	})
	if err != nil {
		return err.E
	}
	return nil
}

func (l *MemLobby) listen(p *player.Player) {
LOOP:
	for {
		select {
		case <-l.closeChListen:
			break LOOP
		default:
			if !p.HasConnection() {
				time.Sleep(SleepOnNoConnection)
				continue
			}

			envelope := message.Envelope{}
			if err := p.Connection().ReadJSON(&envelope); err != nil {
				p.Close()
				l.logger.Warn(connectionErr(err))
				continue
			}

			if envelope.Direction != message.Request {
				p.Close()
				l.logger.Warn("disconnected the peer because of the malformated message")
				continue
			}

			if envelope.GetFlag(request.Chat) {

			}
		}
	}

	l.logger.Info("listening goroutine of the memlobby has been stopped")
}

func (l *MemLobby) ping() {
LOOP:
	for {
		select {
		case <-l.pingTicker.C:
			l.logger.Infof("[ping] %d players connected", l.playerMap.Count())

			activeCount := uint(0)
			err := l.playerMap.RangePtr(func(p *player.Player) error {
				if !p.HasConnection() {
					return nil
				}

				conn := p.Connection()
				t := time.Now().Add(LobbyConnectionTimeOut)
				if err := conn.SetWriteDeadline(t); err != nil {
					p.Close()
					l.logger.Error(connectionErr(err))
					return nil
				}

				err := conn.WriteMessage(websocket.PingMessage, message.PingBytes)
				if err != nil {
					p.Close()
					l.logger.Warn(connectionErr(err))
				} else {
					activeCount++
				}

				return nil
			})
			if err != nil {
				l.playerMap.DeleteRaw(err.K)
				l.logger.Warnf("deleted the player because of the previous error => %s", err.E)
				continue
			}

			l.activeCount = activeCount
		case <-l.closeChPing:
			break LOOP
		}
	}

	l.logger.Info("ping goroutine of the memlobby has been stopped")
}

func (l *MemLobby) Delete() {
	l.pingTicker.Stop()
	l.closeChListen <- true
	l.closeChPing <- true
}
