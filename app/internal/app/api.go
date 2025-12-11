package app

import (
	"devops/app/internal/core/location"
	"devops/app/internal/core/sensor"
	"devops/app/internal/http"
	v1 "devops/app/internal/http/handlers/v1"
	"devops/app/internal/http/interfaces"
	"devops/common/config"
	"devops/common/db"
	"devops/common/logger"
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
