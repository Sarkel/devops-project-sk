package seeds

import (
	"context"
	"database/sql"
	"fmt"

	gen "devops/seeder/internal/db/gen"

	"github.com/jaswdr/faker/v2"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upSeedSensors, downSeedSensors)
}

func upSeedSensors(ctx context.Context, tx *sql.Tx) error {
	queries := gen.New(tx)
	fake := faker.New()

	// Get all location IDs
	locationIDs, err := queries.GetAllLocationIDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get locations: %w", err)
	}

	// Create 2 sensors per location (1 local, 1 api)
	sensorTypes := []gen.TempCheckerSensorType{
		gen.TempCheckerSensorTypeLocal,
		gen.TempCheckerSensorTypeApi,
	}

	for _, locationID := range locationIDs {
		for _, sensorType := range sensorTypes {
			sensorSID := fmt.Sprintf("SEN-%05d", fake.RandomDigitNotNull())

			_, err := queries.CreateLocationSensor(ctx, gen.CreateLocationSensorParams{
				LocationID: locationID,
				SensorSid:  sensorSID,
				Type:       sensorType,
			})
			if err != nil {
				return fmt.Errorf("failed to insert sensor for location %d: %w", locationID, err)
			}
		}
	}

	return nil
}

func downSeedSensors(ctx context.Context, tx *sql.Tx) error {
	queries := gen.New(tx)
	return queries.DeleteAllLocationSensors(ctx)
}
