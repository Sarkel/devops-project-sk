package seeds

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	gen "devops/seeder/internal/db/gen"

	"github.com/jaswdr/faker/v2"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upSeedTemperatureData, downSeedTemperatureData)
}

func upSeedTemperatureData(ctx context.Context, tx *sql.Tx) error {
	fake := faker.New()
	queries := gen.New(tx)

	// Get all sensor IDs
	sensorIDs, err := queries.GetAllLocationSensorIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sensors: %w", err)
	}

	// Generate temperature data for the last 30 days
	now := time.Now()
	for _, sensorID := range sensorIDs {
		for i := 0; i < 30; i++ {
			timestamp := now.Add(-time.Duration(i) * 24 * time.Hour)
			// Generate realistic temperature between -10 and 20 degrees Celsius
			temperature := fake.Float64(2, -10.0, 20.0)

			err := queries.CreateSensorData(ctx, gen.CreateSensorDataParams{
				LocationSensorID: sensorID,
				Temperature:      temperature,
				Timestamp:        timestamp,
			})
			if err != nil {
				return fmt.Errorf("failed to insert temperature data for sensor %d: %w", sensorID, err)
			}
		}
	}

	return nil
}

func downSeedTemperatureData(ctx context.Context, tx *sql.Tx) error {
	queries := gen.New(tx)
	return queries.DeleteAllSensorData(ctx)
}
