package service

import (
	"lobby/server"
	"lobby/server/context"
	"lobby/storage/memstore"

	"github.com/labstack/echo/v4"
)

func Run() {
	e := echo.New()

	server.NewServer(
		e,
		context.NewMetadata("LobyService", "0.0.1"),
		context.NewComponents(
			memstore.NewMemstore(),
		),
		"127.0.0.1:9990",
	).Run()
}
