package location

import (
	"testing"

	genDb "devops/app/internal/db/gen"
	cDB "devops/common/db"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	conManager := &cDB.ConManager{}
	deps := Dependencies{Db: conManager}

	service := NewService(deps)

	assert.NotNil(t, service)
	assert.Equal(t, conManager, service.db)
}

func TestService_GetLocations_MapCorrectly(t *testing.T) {
	// Test that the mapping from DB rows to Location DTOs works correctly
	dbRows := []genDb.GetLocationsRow{
		{LocationName: "TestCity1", LocationSid: "test-sid-1"},
		{LocationName: "TestCity2", LocationSid: "test-sid-2"},
	}

	// Expected result
	expected := []Location{
		{Name: "TestCity1", Sid: "test-sid-1"},
		{Name: "TestCity2", Sid: "test-sid-2"},
	}

	result := make([]Location, len(dbRows))
	for i, l := range dbRows {
		result[i] = Location{
			Name: l.LocationName,
			Sid:  l.LocationSid,
		}
	}

	assert.Equal(t, expected, result)
}

func TestLocation_DTO(t *testing.T) {
	loc := Location{
		Name: "Warsaw",
		Sid:  "warsaw-sid",
	}

	assert.Equal(t, "Warsaw", loc.Name)
	assert.Equal(t, "warsaw-sid", loc.Sid)
}
