package server

import (
	"lobby/server/context"
	"lobby/server/handlers"
	"lobby/server/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Server struct {
	echo       *echo.Echo
	metadata   *context.Metadata
	components *context.Components
	listenAt   string
}

func NewServer(
	echo *echo.Echo,
	metadata *context.Metadata,
	components *context.Components,
	listenAt string,
) *Server {
	return &Server{
		echo:       echo,
		metadata:   metadata,
		components: components,
		listenAt:   listenAt,
	}
}

func (s *Server) ConvertContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.NewContext(
			c,
			s.metadata,
			s.components,
		)
		return next(ctx)
	}
}

func (s *Server) Run() {
	s.echo.Use(s.ConvertContext)
	s.echo.Validator = validator.NewValidator()
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.Logger())

	s.echo.GET("/", handlers.Root)
	s.echo.GET("/lobby/list", handlers.LobbyList)
	s.echo.POST("/lobby/detail", handlers.LobbyDetail)
	s.echo.POST("/lobby/create", handlers.LobbyCreate)
	s.echo.POST("/lobby/join", handlers.LobbyJoin)
	s.echo.GET("/lobby/listen/:lobby", handlers.LobbyListen)

	s.echo.Logger.SetLevel(log.INFO)
	s.echo.Logger.Fatal(s.echo.Start(s.listenAt))
}
