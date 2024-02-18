package main

import (
	"github.com/diegofsousa/rinha-de-backend-q1/internal/infra/configuration"
	"github.com/labstack/gommon/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := configuration.Init()
	go configuration.NewApp(c).Start()
	shutdownServer()
}

func shutdownServer() {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	switch <-signalChannel {
	case syscall.SIGINT:
		log.Info("Received SIGINT, stopping")
	case syscall.SIGTERM:
		log.Info("Received SIGINT, stopping")
	}
}
