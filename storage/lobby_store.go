package storage

import "lobby/lobby/lobby"

type LobbyStore interface {
	AddLobby(id string, l *lobby.Lobby)
}
