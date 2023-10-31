package memstore

import (
	"lobby/lobby/lobby"
	"sync"
)

type Memstore struct {
	lobbyCount int
	lobbyMap   *sync.Map
}

func NewMemstore() *Memstore {
	return &Memstore{
		lobbyCount: 0,
		lobbyMap:   &sync.Map{},
	}
}

func (m *Memstore) AddLobby(id string, l *lobby.Lobby) {
	if _, exists := m.lobbyMap.LoadOrStore(id, l); !exists {
		m.lobbyCount++
	}
}
