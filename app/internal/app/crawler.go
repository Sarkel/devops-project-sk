package app

import (
	"context"
	"devops/app/internal/core/crawler"
	"devops/app/internal/core/meteo"
	"devops/common/config"
	"devops/common/db"
	"devops/common/logger"
	"devops/common/mqtt"
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
