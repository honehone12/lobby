package handlers

import (
	"errors"
	"lobby/lobby/lobby"
	"lobby/lobby/message"
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

type LobbyListenParasm struct {
	LobbyId  string `validate:"required,uuid4,min=36,max=36"`
	PlayerId string `validate:"required,uuid4,min=36,max=36"`
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

	l := lobby.NewMemLobby(formData.LobbyName, c.Logger())
	ctx.LobbyStore().AddLobby(l)

	return c.JSON(http.StatusOK, &LobbyCreateResponse{
		LobbyId: l.Id(),
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

	n := message.NewNotification()
	n.SetFlag(message.JoinBit)
	n.SetMessage("player-name", formData.PlayerName)
	l.BroadcastNotification(n)

	p := player.NewPlayer(formData.PlayerName)
	l.AddPlayer(p)

	return c.JSON(http.StatusOK, &LobbyJoinResponse{
		PlayerId: p.Id(),
	})
}

func LobbyListen(c echo.Context) error {
	params := &LobbyListenParasm{
		LobbyId:  c.Param("lobby"),
		PlayerId: c.QueryParam("player"),
	}

	ctx, err := context.FromEchoCtx(c)
	if err != nil {
		return errres.ServiceError(err, c.Logger())
	}

	if err = ctx.Validate(params); err != nil {
		return errres.BadRequest(ErrorInvalidUuid, c.Logger())
	}

	l, err := ctx.Components.LobbyStore().FindLobby(params.LobbyId)
	if err != nil {
		return errres.BadRequest(err, c.Logger())
	}

	p, err := l.FindPlayer(params.PlayerId)
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

	return nil
}
