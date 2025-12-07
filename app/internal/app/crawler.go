package app

import (
	"context"
	"devops-project-sk/internal/config"
	"devops-project-sk/internal/core/crawler"
	"devops-project-sk/internal/core/meteo"
	"devops-project-sk/internal/db"
	"devops-project-sk/internal/logger"
	"devops-project-sk/internal/mqtt"
	"fmt"
)

func RunCrawler() error {
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

	meteoClient := meteo.NewOpenMeteoClient(&meteo.OpenMeteoDependencies{})

	broker, err := mqtt.NewMosquittoClient(mqtt.Dependencies{
		Logger: log,
		Config: &cfg.MQTTBroker,
	})

	if err != nil {
		return fmt.Errorf("failed to create mqtt client: %w", err)
	}

	defer broker.Close()

	crawlerService := crawler.NewService(&crawler.ServiceDependencies{
		DB:          conManager,
		Logger:      log,
		MeteoClient: meteoClient,
		Broker:      broker,
	})

	rootCtx := context.Background()
	defer rootCtx.Done()

	if err := crawlerService.Crawl(rootCtx); err != nil {
		return fmt.Errorf("failed to crawl: %w", err)
	}
	return nil
}
