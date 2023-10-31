package context

import (
	"lobby/lobby/lobby"

	"github.com/gorilla/websocket"
)

type Components struct {
	lobbyStore        lobby.LobbyStore
	webSocketSwticher *websocket.Upgrader
}

func NewComponents(s lobby.LobbyStore) *Components {
	return &Components{
		lobbyStore:        s,
		webSocketSwticher: &websocket.Upgrader{},
	}
}

func (c *Components) LobbyStore() lobby.LobbyStore {
	return c.lobbyStore
}

func (c *Components) WebSocketSwitcher() *websocket.Upgrader {
	return c.webSocketSwticher
}
