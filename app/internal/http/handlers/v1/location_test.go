package v1

import (
	"encoding/json"
	"testing"

	"devops/app/internal/core/location"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

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
