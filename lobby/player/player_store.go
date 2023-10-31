package player

type PlayerStore interface {
	AddPlayer(string, *Player)
	FindPlayer(string) (*Player, error)
}
