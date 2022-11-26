package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/eazygood/getground-app/cmd/app/server"
	"github.com/eazygood/getground-app/internal/config"
	"github.com/eazygood/getground-app/internal/infrastructure/log"
	logger "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		logger.WithError(err).Fatal("failed to load config")
	}

	log.Init(cfg.Log)
	server.Start(contextWithTermSignal(), *cfg)
	logger.Info("server started")
}

// contextWithTermSignal returns a context that will be cancelled whenever a SIGTERM is received.
func contextWithTermSignal() context.Context {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		<-signals
		logger.Info("Received TERM signal, stopping service...")
		cancelFunc()
	}()
	return ctx
}
