package configuration

import (
	"context"
	"github.com/diegofsousa/rinha-de-backend-q1/internal/application/usecase"
	api "github.com/diegofsousa/rinha-de-backend-q1/internal/infra/api/consumer"
	"github.com/diegofsousa/rinha-de-backend-q1/internal/infra/repository"
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
	consumerUsecase := usecase.NewConsumer(repository.NewConsumer(a.config.GetString("database.url")))
	clients := api.NewConsumer(consumerUsecase)
	clients.Register(a.server)
	log.Info("Running at ", a.config.GetString("server.host"))

	log.Info(context.Background(), a.server.Start(a.config.GetString("server.host")), "Server fatal error")
}
