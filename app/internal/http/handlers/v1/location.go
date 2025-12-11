package v1

import (
	"devops/app/internal/core/location"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LocationCtrlDependencies struct {
	Service *location.Service
}

type LocationCtrl struct {
	s *location.Service
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
