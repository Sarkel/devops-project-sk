package export

import (
	"context"
	cDB "devops/common/db"
	"devops/seeder/internal/db"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type Dependencies struct {
	Db *cDB.ConManager
	L  *slog.Logger
}

type Service struct {
	db  *cDB.ConManager
	l   *slog.Logger
	dir string
}

func New(deps Dependencies) *Service {
	return &Service{
		db:  deps.Db,
		l:   deps.L,
		dir: "./data/export",
	}
}

func (s *Service) Run(ctx context.Context) error {
	s.l.Info("Generating report - fetching data from database...")
	data, err := s.fetchData(ctx)

	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}

	if err := s.saveData(data); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	return nil
}

func (s *Service) fetchData(ctx context.Context) (*Data, error) {
	q := db.WithQ(s.db)

	s.l.Info("Fetching locations...")
	locations, err := q.GetAllLocations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}

	s.l.Info("Fetching location sensors...")
	locationSensors, err := q.GetAllLocationSensors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get location sensors: %w", err)
	}

	s.l.Info("Fetching sensor data...")
	sensorData, err := q.GetAllSensorData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sensor data: %w", err)
	}

	return &Data{
		Locations:       locations,
		LocationSensors: locationSensors,
		SensorData:      sensorData,
	}, nil
}

func (s *Service) saveData(d *Data) error {
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	allDataPath := filepath.Join(s.dir, "all_data.json")
	if err := writeJSON(allDataPath, d); err != nil {
		return fmt.Errorf("failed to write all data: %w", err)
	}
	s.l.Info(fmt.Sprintf("Exported all data to %s", allDataPath))

	return nil
}

func writeJSON(filepath string, data interface{}) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
