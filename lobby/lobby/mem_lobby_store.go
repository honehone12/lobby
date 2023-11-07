package lobby

import (
	"lobby/generics"
	"lobby/logger"
	"time"
)

type MemLobbyStore struct {
	lobbyMap *generics.TypedMap[Lobby]
	ticker   *time.Ticker
	closeCh  chan bool
	logger   logger.Logger
}

func NewMemLobyStore(logger logger.Logger) *MemLobbyStore {
	s := &MemLobbyStore{
		lobbyMap: generics.NewTypedMap[Lobby](),
		ticker:   time.NewTicker(LobbyCleanUpInterval),
		closeCh:  make(chan bool),
		logger:   logger,
	}
	go s.cleanUp()
	return s
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

func (s *MemLobbyStore) DeleteLobby(id string) {
	s.lobbyMap.Delete(id)
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
	if err != nil {
		return nil, err.E
	}
	return buff, nil
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

func (s *MemLobbyStore) recoverCleanUp() {
	if r := recover(); r != nil {
		s.logger.Warn("recover clean up goroutine")
		go s.cleanUp()
	}
}

func (s *MemLobbyStore) cleanUp() {
	defer s.recoverCleanUp()

LOOP:
	for {
		select {
		case <-s.ticker.C:
			cleanUpBuff := make(map[string]Lobby)
			err := s.lobbyMap.Range(func(l Lobby) error {
				if l.ActiveCount() == 0 {
					cleanUpBuff[l.Id()] = l
				}
				return nil
			})
			if err != nil {
				s.logger.Panic(err)
			}

			for k, l := range cleanUpBuff {
				l.Delete()
				s.lobbyMap.Delete(k)
				s.logger.Warnf("deleted the lobby because no player are active")
			}

			s.logger.Debugf("[cleaning up] %d lobby stored", s.lobbyMap.Count())
		case <-s.closeCh:
			break LOOP
		}
	}

	s.logger.Info("cleaning up goroutine of the memlobbystore has been stopped")
}
