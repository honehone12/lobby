package handlers

import (
	"lobby/lobby/lobby"
	"lobby/lobby/player"
	"lobby/server/context"
	"lobby/server/errres"
	"lobby/server/form"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LobbyCreateForm struct {
	PlyerName string `form:"name" validate:"required,alphanum,min=2,max=32"`
	LobbyName string `form:"name" validate:"required,alphanum,min=2,max=32"`
}

type LobbyCreateResponse struct {
	PlayerId string
	LobbyId  string
}

func LobbyCreate(c echo.Context) error {
	formData := &LobbyCreateForm{}
	if err := form.ProcessFormData(c, formData); err != nil {
		return errres.BadRequest(err, c.Logger())
	}

	ctx, err := context.FromEchoCtx(c)
	if err != nil {
		return err
	}

	conn, err := ctx.WebSocketSwitcher().Upgrade(
		c.Response(),
		c.Request(),
		nil,
	)
	if err != nil {
		return err
	}

	l := lobby.NewLobby(formData.LobbyName)
	lid := l.Id()
	p := player.NewPlayer(formData.PlyerName, conn)
	pid := p.Id()
	l.AddPlayer(pid, p)
	ctx.LobbyStore().AddLobby(lid, l)

	return c.JSON(http.StatusOK, &LobbyCreateResponse{
		PlayerId: pid,
		LobbyId:  lid,
	})
}

func LobbyJoin(c echo.Context) error {
	return nil
}
