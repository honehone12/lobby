package handlers

import (
	"lobby/server/context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RootResponse struct {
	Name    string
	Version string
}

func Root(c echo.Context) error {
	ctx, err := context.FromEchoCtx(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &RootResponse{
		Name:    ctx.Name(),
		Version: ctx.Version(),
	})
}
