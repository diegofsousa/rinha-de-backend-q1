package configuration

import (
	"context"
	api "github.com/diegofsousa/rinha-de-backend-q1/internal/infra/api/clients"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type App struct {
	server *echo.Echo
	config *viper.Viper
}

func NewApp(config *viper.Viper) *App {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true
	server.Use(middleware.Recover())

	return &App{
		server: server,
		config: config,
	}
}

func (a *App) Start() {
	clients := api.NewClients()
	clients.Register(a.server)

	log.Info(context.Background(), a.server.Start(a.config.GetString("server.host")), "Server fatal error")
}
