package player

import (
	libuuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	id   string
	name string

	connection *websocket.Conn
}

func NewPlayer(name string, connection *websocket.Conn) *Player {
	return &Player{
		id:         libuuid.NewString(),
		name:       name,
		connection: connection,
	}
}

func (p *Player) Id() string {
	return p.id
}

func (p *Player) Name() string {
	return p.name
}
