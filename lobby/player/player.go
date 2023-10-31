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

func NewPlayer(name string) *Player {
	return &Player{
		id:         libuuid.NewString(),
		name:       name,
		connection: nil,
	}
}

func (p *Player) Id() string {
	return p.id
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) HasConnection() bool {
	return p.connection != nil
}

func (p *Player) Connection() *websocket.Conn {
	return p.connection
}

func (p *Player) SetConnection(conn *websocket.Conn) {
	p.connection = conn
}
