package lobby

import "lobby/lobby/player"

type LobbyStore interface {
	AddLobby(player.PlayerStore)
	FindLobby(string) (player.PlayerStore, error)
}
