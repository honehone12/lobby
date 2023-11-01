package lobby

import (
	"lobby/generics"
	"lobby/lobby/player"
)

type MemLobbyStore struct {
	lobbyMap *generics.TypedMap[player.PlayerStore]
}

func NewMemLobyStore() *MemLobbyStore {
	return &MemLobbyStore{
		lobbyMap: generics.NewTypedMap[player.PlayerStore](),
	}
}

func (m *MemLobbyStore) AddLobby(l player.PlayerStore) {
	m.lobbyMap.Add(l.Id(), l)
}

func (m *MemLobbyStore) FindLobby(id string) (player.PlayerStore, error) {
	return m.lobbyMap.ItemOrDefault(id, nil)
}
