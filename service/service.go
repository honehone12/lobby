package service

import (
	"lobby/lobby/lobby"
	"lobby/server"
	"lobby/server/context"

	"github.com/labstack/echo/v4"
)

func Run() {
	e := echo.New()
	s := lobby.NewMemLobyStore()
	server.NewServer(
		e,
		context.NewMetadata("LobyService", "0.0.1"),
		context.NewComponents(s),
		"127.0.0.1:9990",
	).Run()
}
