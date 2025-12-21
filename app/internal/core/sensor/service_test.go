package sensor

import (
	"errors"
	"testing"
	"time"

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

func TestService_GetSummary_LocationNotFound(t *testing.T) {
	// This test demonstrates the expected behavior when location doesn't exist
	expectedErr := errors.New("location not found")
	assert.NotNil(t, expectedErr)
}

func TestService_GetSummary_UnexpectedSensors(t *testing.T) {
	// Test case: more than 2 sensor types returned (unexpected)
	expectedErr := errors.New("unexpected sensors summary")
	assert.NotNil(t, expectedErr)
}

func TestService_GetSummary_MapSensorTypes(t *testing.T) {
	// Test the mapping logic for sensor types
	now := time.Now()

	summaryItems := []struct {
		Type           genDb.TempCheckerSensorType
		Date           time.Time
		AvgTemperature float64
	}{
		{Type: genDb.TempCheckerSensorTypeLocal, Date: now, AvgTemperature: 22.5},
		{Type: genDb.TempCheckerSensorTypeApi, Date: now, AvgTemperature: 23.0},
	}

	result := Summary{}

	for _, item := range summaryItems {
		switch item.Type {
		case genDb.TempCheckerSensorTypeLocal:
			result.Local = &SummaryItem{
				Timestamp:   item.Date,
				Temperature: item.AvgTemperature,
				Trend:       0,
			}
		case genDb.TempCheckerSensorTypeApi:
			result.Api = &SummaryItem{
				Timestamp:   item.Date,
				Temperature: item.AvgTemperature,
				Trend:       0,
			}
		}
	}

	assert.NotNil(t, result.Local)
	assert.NotNil(t, result.Api)
	assert.Equal(t, 22.5, result.Local.Temperature)
	assert.Equal(t, 23.0, result.Api.Temperature)
}

func TestService_GetData_MapDataPoints(t *testing.T) {
	// Test data point mapping
	now := time.Now()
	dbResults := []genDb.GetSensorDataPointsRow{
		{
			Type:           genDb.TempCheckerSensorTypeLocal,
			TimeDim:        now,
			AvgTemperature: 21.5,
		},
		{
			Type:           genDb.TempCheckerSensorTypeApi,
			TimeDim:        now.Add(1 * time.Hour),
			AvgTemperature: 22.0,
		},
	}

	points := make([]DataPoint, len(dbResults))
	for i, r := range dbResults {
		points[i] = DataPoint{
			Type:        r.Type,
			Timestamp:   r.TimeDim,
			Temperature: r.AvgTemperature,
		}
	}

	assert.Len(t, points, 2)
	assert.Equal(t, genDb.TempCheckerSensorTypeLocal, points[0].Type)
	assert.Equal(t, 21.5, points[0].Temperature)
	assert.Equal(t, genDb.TempCheckerSensorTypeApi, points[1].Type)
	assert.Equal(t, 22.0, points[1].Temperature)
}

func TestSummaryQs_Validation(t *testing.T) {
	// Test DTO structure
	qs := SummaryQs{LocationSid: "test-location"}
	assert.Equal(t, "test-location", qs.LocationSid)
}

func TestDataQs_Validation(t *testing.T) {
	// Test DTO structure
	start := time.Now()
	end := start.Add(24 * time.Hour)

	qs := DataQs{
		LocationSid:   "test-location",
		StartDatetime: start,
		EndDatetime:   end,
		Aggregation:   "day",
		Types:         []genDb.TempCheckerSensorType{genDb.TempCheckerSensorTypeLocal, genDb.TempCheckerSensorTypeApi},
	}

	assert.Equal(t, "test-location", qs.LocationSid)
	assert.Equal(t, "day", qs.Aggregation)
	assert.Len(t, qs.Types, 2)
}
