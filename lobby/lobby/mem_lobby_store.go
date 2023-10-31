package lobby

import (
	"lobby/generics"
	"lobby/lobby/player"
)

type MemLobbyStore struct {
	lobbyCount int
	lobbyMap   *generics.TypedMap[player.PlayerStore]
}

func NewMemLobyStore() *MemLobbyStore {
	return &MemLobbyStore{
		lobbyCount: 0,
		lobbyMap:   generics.NewTypedMap[player.PlayerStore](),
	}
}

func (m *MemLobbyStore) AddLobby(id string, l player.PlayerStore) {
	m.lobbyMap.Add(id, l)
}

func (m *MemLobbyStore) FindLobby(id string) (player.PlayerStore, error) {
	return m.lobbyMap.ItemOrDefault(id, nil)
}
