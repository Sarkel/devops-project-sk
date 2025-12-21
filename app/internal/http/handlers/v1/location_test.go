package v1

import (
	"context"
	"encoding/json"
	"testing"

	"devops/app/internal/core/location"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLocationService struct {
	mock.Mock
}

func (m *MockLocationService) GetLocations(ctx context.Context) ([]location.Location, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]location.Location), args.Error(1)
}

func TestNewLocationCtrl(t *testing.T) {
	deps := LocationCtrlDependencies{
		Service: &location.Service{},
	}

	ctrl := NewLocationCtrl(deps)

	assert.NotNil(t, ctrl)
	assert.NotNil(t, ctrl.s)
}

func TestLocationCtrl_RegisterRoutes(t *testing.T) {
	ctrl := &LocationCtrl{}

	e := echo.New()
	group := e.Group("/v1")

	ctrl.RegisterRoutes(group)

	// Verify that routes are registered
	routes := e.Routes()

	found := false
	for _, route := range routes {
		if route.Path == "/v1/locations" && route.Method == "GET" {
			found = true
			break
		}
	}

	assert.True(t, found, "Locations route should be registered")
}

func TestLocationResponse_JSON(t *testing.T) {
	locations := []location.Location{
		{Name: "Warsaw", Sid: "warsaw-sid"},
		{Name: "Krakow", Sid: "krakow-sid"},
	}

	jsonData, err := json.Marshal(locations)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var decoded []location.Location
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)
	assert.Len(t, decoded, 2)
	assert.Equal(t, "Warsaw", decoded[0].Name)
	assert.Equal(t, "warsaw-sid", decoded[0].Sid)
}

func TestLocationCtrlDependencies_Structure(t *testing.T) {
	mockService := new(MockLocationService)
	deps := LocationCtrlDependencies{Service: mockService}
	assert.NotNil(t, deps.Service)
}

func TestLocationService_Interface(t *testing.T) {
	// Verify that location.Service implements LocationService interface
	var _ LocationService = (*location.Service)(nil)
}
