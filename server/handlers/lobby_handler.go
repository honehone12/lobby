package handlers

import (
	"errors"
	"lobby/lobby/lobby"
	"lobby/lobby/player"
	"lobby/server/context"
	"lobby/server/errres"
	"lobby/server/form"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LobbyCreateForm struct {
	LobbyName string `form:"lobby-name" validate:"required,alphanum,min=2,max=32"`
}

type LobbyJoinForm struct {
	LobbyId    string `form:"lobby-id" validate:"required,uuid4,min=36,max=36"`
	PlayerName string `form:"player-name" validate:"required,alphanum,min=2,max=32"`
}

type LobbyCreateResponse struct {
	LobbyId string
}

type LobbyJoinResponse struct {
	PlayerId string
}

const (
	IdLen = 36
)

var (
	ErrorInvalidUuid          = errors.New("invalid uuid format")
	ErrorDuplicatedConnection = errors.New("attempt of duplicated connection")
)

func LobbyCreate(c echo.Context) error {
	formData := &LobbyCreateForm{}
	if err := form.ProcessFormData(c, formData); err != nil {
		return errres.BadRequest(err, c.Logger())
	}

	ctx, err := context.FromEchoCtx(c)
	if err != nil {
		return errres.ServiceError(err, c.Logger())
	}

	l := lobby.NewLobby(formData.LobbyName)
	lid := l.Id()
	ctx.LobbyStore().AddLobby(lid, l)

	return c.JSON(http.StatusOK, &LobbyCreateResponse{
		LobbyId: lid,
	})
}

func LobbyJoin(c echo.Context) error {
	formData := &LobbyJoinForm{}
	if err := form.ProcessFormData(c, formData); err != nil {
		return errres.BadRequest(err, c.Logger())
	}

	ctx, err := context.FromEchoCtx(c)
	if err != nil {
		return errres.ServiceError(err, c.Logger())
	}

	l, err := ctx.LobbyStore().FindLobby(formData.LobbyId)
	if err != nil {
		return errres.ServiceError(err, c.Logger())
	}

	p := player.NewPlayer(formData.PlayerName)
	pid := p.Id()
	l.AddPlayer(pid, p)

	return c.JSON(http.StatusOK, &LobbyJoinResponse{
		PlayerId: pid,
	})
}

func LobbyListen(c echo.Context) error {
	lobbyId := c.Param("lobby")
	if len(lobbyId) != IdLen {
		return errres.BadRequest(ErrorInvalidUuid, c.Logger())
	}

	playerId := c.QueryParam("player")
	if len(playerId) != IdLen {
		return errres.BadRequest(ErrorInvalidUuid, c.Logger())
	}

	ctx, err := context.FromEchoCtx(c)
	if err != nil {
		return errres.ServiceError(err, c.Logger())
	}

	l, err := ctx.Components.LobbyStore().FindLobby(lobbyId)
	if err != nil {
		return errres.BadRequest(err, c.Logger())
	}

	p, err := l.FindPlayer(playerId)
	if err != nil {
		return errres.BadRequest(err, c.Logger())
	}

	if p.HasConnection() {
		return errres.BadRequest(ErrorDuplicatedConnection, c.Logger())
	}

	conn, err := ctx.WebSocketSwitcher().Upgrade(
		c.Response(),
		c.Request(),
		nil,
	)
	if err != nil {
		return errres.ServiceError(err, c.Logger())
	}

	p.SetConnection(conn)

	return c.NoContent(http.StatusOK)
}
