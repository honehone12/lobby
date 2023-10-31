package context

import (
	"lobby/storage"

	"github.com/gorilla/websocket"
)

type Components struct {
	lobbyStore        storage.LobbyStore
	webSocketSwticher *websocket.Upgrader
}

func NewComponents(s storage.LobbyStore) *Components {
	return &Components{
		lobbyStore:        s,
		webSocketSwticher: &websocket.Upgrader{},
	}
}

func (c *Components) LobbyStore() storage.LobbyStore {
	return c.lobbyStore
}

func (c *Components) WebSocketSwitcher() *websocket.Upgrader {
	return c.webSocketSwticher
}
