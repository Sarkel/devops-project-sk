package main

import (
	"context"
	"devops/seeder/internal/export"
	"devops/seeder/internal/seeder"
	"embed"
	"fmt"
	"log"

	"devops/common/config"
	"devops/common/db"
	"devops/common/logger"
	_ "devops/seeder/seeds"
)

//go:embed seeds/.keep
var baseFS embed.FS

func main() {
	if err := RunSeeder(); err != nil {
		log.Fatal(err)
	}
}

func RunSeeder() error {
	ctx := context.Background()

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

	s := seeder.New(seeder.Dependencies{
		Db:  conManager,
		L:   l,
		Dir: baseFS,
	})

	if err := s.Run(); err != nil {
		return fmt.Errorf("failed to seed database: %w", err)
	}

	exporter := export.New(export.Dependencies{
		Db: conManager,
		L:  l,
	})

	if err := exporter.Run(ctx); err != nil {
		return fmt.Errorf("failed to export data: %w", err)
	}
	return nil
}
