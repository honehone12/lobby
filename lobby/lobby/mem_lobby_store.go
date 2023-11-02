package lobby

import (
	"lobby/generics"
)

type MemLobbyStore struct {
	lobbyMap *generics.TypedMap[Lobby]
}

func NewMemLobyStore() *MemLobbyStore {
	return &MemLobbyStore{
		lobbyMap: generics.NewTypedMap[Lobby](),
	}
}

func (s *MemLobbyStore) LobbyCount() uint {
	return uint(s.lobbyMap.Count())
}

func (s *MemLobbyStore) AddLobby(l Lobby) {
	s.lobbyMap.Add(l.Id(), l)
}

func (s *MemLobbyStore) FindLobby(id string) (Lobby, error) {
	return s.lobbyMap.ItemOrDefault(id, nil)
}

func (s *MemLobbyStore) GetSummaries() ([]LobbySummary, error) {
	buff := make([]LobbySummary, s.lobbyMap.Count())
	iter := 0
	err := s.lobbyMap.Range(func(l Lobby) error {
		buff[iter].Id = l.Id()
		buff[iter].Name = l.Name()
		buff[iter].PlayerCount = l.PlayerCount()
		buff[iter].ActiveCount = l.ActiveCount()
		iter++
		return nil
	})
	return buff, err
}

func (s *MemLobbyStore) GetDetail(id string) (*LobbyDetail, error) {
	l, err := s.lobbyMap.ItemOrDefault(id, nil)
	if err != nil {
		return nil, err
	}

	list, err := l.GetPlayers()
	if err != nil {
		return nil, err
	}

	return &LobbyDetail{
		Id:          l.Id(),
		Name:        l.Name(),
		PlayerCount: l.PlayerCount(),
		ActiveCount: l.ActiveCount(),
		PlayerList:  list,
	}, nil
}
