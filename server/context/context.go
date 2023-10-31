package context

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	*Metadata
	*Components
}

var (
	ErrorCastFail = errors.New("failed to cast the context")
)

func NewContext(e echo.Context, m *Metadata, c *Components) *Context {
	return &Context{
		Context:    e,
		Metadata:   m,
		Components: c,
	}
}

func FromEchoCtx(e echo.Context) (*Context, error) {
	ctx, ok := e.(*Context)
	if !ok {
		return nil, ErrorCastFail
	}

	return ctx, nil
}
