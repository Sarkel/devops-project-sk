package app

import (
	"devops-project-sk/internal/config"
	"devops-project-sk/internal/core/location"
	"devops-project-sk/internal/core/sensor"
	"devops-project-sk/internal/db"
	"devops-project-sk/internal/http"
	v1 "devops-project-sk/internal/http/handlers/v1"
	"devops-project-sk/internal/http/interfaces"
	"devops-project-sk/internal/logger"
	"fmt"
)

func StartApi() error {
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

	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}

	defer db.Close(conManager, log)

	sensorsSvr := sensor.NewService(sensor.Dependencies{
		Db: conManager,
	})

	sensorsCtrl := v1.NewSensorsCtrl(v1.SensorsCtrlDependencies{
		Service: sensorsSvr,
	})

	locationSvr := location.NewService(location.Dependencies{
		Db: conManager,
	})

	locationCtrl := v1.NewLocationCtrl(v1.LocationCtrlDependencies{
		Service: locationSvr,
	})

	ctrls := []interfaces.Controller{
		sensorsCtrl, locationCtrl,
	}

	r := http.NewRouter(&http.RouterDependencies{
		Controllers:  ctrls,
		AuthConfig:   &cfg.Auth,
		ServerConfig: &cfg.Server,
	})

	svr := http.NewServer(http.ServerDependencies{
		Router: r,
		Logger: log,
		Config: &cfg.Server,
	})

	if err := svr.Start(); err != nil {
		return err
	}
	return nil
}
