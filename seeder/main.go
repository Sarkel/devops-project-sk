package main

import (
	"devops/common/config"
	"devops/common/db"
	"devops/common/logger"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
)

func main() {
	if err := RunSeeder(); err != nil {
		log.Fatal(err)
	}
}

func RunSeeder() error {
	cfg, err := config.Load()

	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	l := logger.New(logger.Dependencies{
		Config: cfg.Logger,
	})

	conManager, err := db.NewConManager(db.Dependencies{
		Logger: l,
		Config: &cfg.Database,
	})

	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}

	defer db.Close(conManager, l)

	goose.SetTableName("goose_seed_version")

	if err := goose.Up(conManager.GetDB(), "."); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return generateReport(conManager)
}

func generateReport(conManager *db.ConManager) error {
	return nil
}
