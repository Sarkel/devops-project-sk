package app

import (
	"context"
	"devops-project-sk/internal/config"
	"devops-project-sk/internal/core/reader"
	"devops-project-sk/internal/db"
	"devops-project-sk/internal/logger"
	"devops-project-sk/internal/mqtt"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func StartReader() error {
	ctx := context.Background()

	cfg, err := config.Load()

	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	log := logger.New(logger.Dependencies{
		Config: cfg.Logger,
	})

	conManager, err := db.NewConManager(db.Dependencies{
		Logger: log,
		Config: &cfg.Database,
	})

	defer db.Close(conManager, log)

	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}

	broker, err := mqtt.NewMosquittoClient(mqtt.Dependencies{
		Logger: log,
		Config: &cfg.MQTTBroker,
	})

	if err != nil {
		return fmt.Errorf("failed to create mqtt client: %w", err)
	}

	defer broker.Close()

	readerService := reader.NewService(&reader.Dependencies{
		DB:     conManager,
		Logger: log,
		Broker: broker,
	})

	if err := readerService.Listen(ctx); err != nil {
		return fmt.Errorf("failed to listen to mqtt broker: %w", err)
	}

	log.Info("reader service running...")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh // wait for Ctrl+C / SIGTERM

	log.Info("reader service stopped")
	return nil
}
