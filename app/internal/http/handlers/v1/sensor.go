package v1

import (
	"devops-project-sk/internal/core/sensor"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SensorsCtrlDependencies struct {
	Service *sensor.Service
}

type SensorsCtrl struct {
	s *sensor.Service
}

func NewSensorsCtrl(deps SensorsCtrlDependencies) *SensorsCtrl {
	return &SensorsCtrl{
		s: deps.Service,
	}
}

func (c *SensorsCtrl) getSensorsSummary(ctx echo.Context) error {
	var params sensor.SummaryQs

	if err := ctx.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(&params); err != nil {
		return err
	}

	// todo: add distinct between 4xx and 5xx errors
	res, err := c.s.GetSummary(ctx.Request().Context(), params)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(200, res)
}

func (c *SensorsCtrl) getSensorsData(ctx echo.Context) error {
	var params sensor.DataQs

	if err := ctx.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(&params); err != nil {
		return err
	}

	res, err := c.s.GetData(ctx.Request().Context(), params)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(200, res)
}

func (c *SensorsCtrl) RegisterRoutes(e *echo.Group) {
	s := e.Group("/sensors")

	s.GET("/summary", c.getSensorsSummary)

	s.GET("/data", c.getSensorsData)
}
