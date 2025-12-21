package seeds

import (
	"context"
	"database/sql"
	"fmt"

	gen "devops/seeder/internal/db/gen"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upSeedLocations, downSeedLocations)
}

func upSeedLocations(ctx context.Context, tx *sql.Tx) error {
	//fake := faker.New()
	queries := gen.New(tx)

	locations := []struct {
		sid       string
		name      string
		latitude  float64
		longitude float64
	}{
		//{fmt.Sprintf("LOC%07d", fake.RandomDigitNotNull()), fake.Address().City(), fake.Address().Latitude(), fake.Address().Longitude()},
		//{fmt.Sprintf("LOC%07d", fake.RandomDigitNotNull()), fake.Address().City(), fake.Address().Latitude(), fake.Address().Longitude()},
	}

	for _, loc := range locations {
		_, err := queries.CreateLocation(ctx, gen.CreateLocationParams{
			LocationName: loc.name,
			Latitude:     loc.latitude,
			Longitude:    loc.longitude,
			LocationSid:  loc.sid,
		})
		if err != nil {
			return fmt.Errorf("failed to insert location %s: %w", loc.name, err)
		}
	}

	return nil
}

func downSeedLocations(ctx context.Context, tx *sql.Tx) error {
	queries := gen.New(tx)
	return queries.DeleteAllLocations(ctx)
}
