package v1

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"devops/app/internal/core/sensor"
	genDb "devops/app/internal/db/gen"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSensorService struct {
	mock.Mock
}

func (m *MockSensorService) GetSummary(ctx context.Context, params sensor.SummaryQs) (sensor.Summary, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(sensor.Summary), args.Error(1)
}

func (m *MockSensorService) GetData(ctx context.Context, params sensor.DataQs) ([]sensor.DataPoint, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sensor.DataPoint), args.Error(1)
}

func TestNewSensorsCtrl(t *testing.T) {
	deps := SensorsCtrlDependencies{
		Service: &sensor.Service{},
	}

	ctrl := NewSensorsCtrl(deps)

	assert.NotNil(t, ctrl)
	assert.NotNil(t, ctrl.s)
}

func TestSensorsCtrl_RegisterRoutes(t *testing.T) {
	e := echo.New()
	group := e.Group("/v1")

	ctrl := &SensorsCtrl{}
	ctrl.RegisterRoutes(group)

	// Verify that routes are registered
	routes := e.Routes()

	// Look for the registered routes
	foundSummary := false
	foundData := false

	for _, route := range routes {
		if route.Path == "/v1/sensors/summary" && route.Method == "GET" {
			foundSummary = true
		}
		if route.Path == "/v1/sensors/data" && route.Method == "GET" {
			foundData = true
		}
	}

	assert.True(t, foundSummary, "Summary route should be registered")
	assert.True(t, foundData, "Data route should be registered")
}

func TestSensorDataResponse_JSON(t *testing.T) {
	now := time.Now()
	data := []sensor.DataPoint{
		{
			Type:        genDb.TempCheckerSensorTypeLocal,
			Timestamp:   now,
			Temperature: 22.5,
		},
	}

	jsonData, err := json.Marshal(data)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded []sensor.DataPoint
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)
	assert.Len(t, decoded, 1)
	assert.Equal(t, 22.5, decoded[0].Temperature)
}

func TestSensorSummaryResponse_JSON(t *testing.T) {
	now := time.Now()
	summary := sensor.Summary{
		Local: &sensor.SummaryItem{
			Timestamp:   now,
			Temperature: 22.5,
			Trend:       0,
		},
		Api: &sensor.SummaryItem{
			Timestamp:   now,
			Temperature: 23.0,
			Trend:       0,
		},
	}

	jsonData, err := json.Marshal(summary)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded sensor.Summary
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)
	assert.NotNil(t, decoded.Local)
	assert.NotNil(t, decoded.Api)
	assert.Equal(t, 22.5, decoded.Local.Temperature)
	assert.Equal(t, 23.0, decoded.Api.Temperature)
}

func TestSensorDataPoint_DTO(t *testing.T) {
	now := time.Now()
	dp := sensor.DataPoint{
		Type:        genDb.TempCheckerSensorTypeLocal,
		Timestamp:   now,
		Temperature: 22.5,
	}

	assert.Equal(t, genDb.TempCheckerSensorTypeLocal, dp.Type)
	assert.Equal(t, 22.5, dp.Temperature)
	assert.Equal(t, now, dp.Timestamp)
}

func TestSensorsCtrlDependencies_Structure(t *testing.T) {
	mockService := new(MockSensorService)
	deps := SensorsCtrlDependencies{Service: mockService}
	assert.NotNil(t, deps.Service)
}

func TestSensorService_Interface(t *testing.T) {
	// Verify that sensor.Service implements SensorService interface
	var _ SensorService = (*sensor.Service)(nil)
}
