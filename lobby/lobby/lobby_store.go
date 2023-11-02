package lobby

type LobbyStore interface {
	LobbyCount() uint
	AddLobby(Lobby)
	FindLobby(string) (Lobby, error)
	DeleteLobby(string)

	GetSummaries() ([]LobbySummary, error)
	GetDetail(string) (*LobbyDetail, error)
}
