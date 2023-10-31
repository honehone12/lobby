package lobby

import (
	"lobby/generics"
)

type MemLobbyStore struct {
	lobbyCount int
	lobbyMap   *generics.TypedMap[MemLobby]
}

func NewMemLobyStore() *MemLobbyStore {
	return &MemLobbyStore{
		lobbyCount: 0,
		lobbyMap:   generics.NewTypedMap[MemLobby](),
	}
}

func (m *MemLobbyStore) AddLobby(id string, l *MemLobby) {
	m.lobbyMap.Add(id, l)
}

func (m *MemLobbyStore) FindLobby(id string) (*MemLobby, error) {
	return m.lobbyMap.Item(id)
}
