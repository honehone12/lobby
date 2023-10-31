package lobby

import "lobby/lobby/player"

type LobbyStore interface {
	AddLobby(string, player.PlayerStore)
	FindLobby(string) (player.PlayerStore, error)
}
