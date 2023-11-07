package player

import (
	"time"

	libuuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type PlayerInfo struct {
	Id     string
	Name   string
	Active bool
}

type Player struct {
	id   string
	name string

	joinedAt   int64
	connection *websocket.Conn
}

func NewPlayer(name string) *Player {
	return &Player{
		id:         libuuid.NewString(),
		name:       name,
		joinedAt:   0,
		connection: nil,
	}
}

func (p *Player) Id() string {
	return p.id
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) JoinedAt() int64 {
	return p.joinedAt
}

func (p *Player) HasConnection() bool {
	return p.connection != nil
}

func (p *Player) Connection() *websocket.Conn {
	return p.connection
}

func (p *Player) SetJoinedAtNow() {
	p.joinedAt = time.Now().Unix()
}

func (p *Player) SetConnection(conn *websocket.Conn) {
	if conn == nil {
		panic("use close() for setting connection nil")
	}

	p.connection = conn
}

func (p *Player) Close() {
	if p.connection != nil {
		defer p.connection.Close()
		p.connection = nil
	}
}
