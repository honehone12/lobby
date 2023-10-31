package service

import (
	"lobby/lobby/lobby"
	"lobby/server"
	"lobby/server/context"

	"github.com/labstack/echo/v4"
)

func Run() {
	e := echo.New()

	server.NewServer(
		e,
		context.NewMetadata("LobyService", "0.0.1"),
		context.NewComponents(
			lobby.NewMemLobyStore(),
		),
		"127.0.0.1:9990",
	).Run()
}
