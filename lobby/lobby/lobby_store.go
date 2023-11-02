package lobby

type LobbyStore interface {
	AddLobby(Lobby)
	FindLobby(string) (Lobby, error)
	LobbyCount() uint

	GetSummaries() ([]LobbySummary, error)
	GetDetail(string) (*LobbyDetail, error)
}
