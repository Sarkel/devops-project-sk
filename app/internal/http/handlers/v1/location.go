package v1

import (
	"context"
	"devops/app/internal/core/location"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LocationService interface {
	GetLocations(ctx context.Context) ([]location.Location, error)
}

type LocationCtrlDependencies struct {
	Service LocationService
}

type LocationCtrl struct {
	s LocationService
}

func NewLocationCtrl(deps LocationCtrlDependencies) *LocationCtrl {
	return &LocationCtrl{
		s: deps.Service,
	}
}

func (c *LocationCtrl) getLocations(ctx echo.Context) error {
	locations, err := c.s.GetLocations(ctx.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(200, locations)
}

func (c *LocationCtrl) RegisterRoutes(e *echo.Group) {
	s := e.Group("/locations")

	s.GET("", c.getLocations)
}
