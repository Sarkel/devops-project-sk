package http

import (
	"crypto/subtle"
	"devops/app/internal/http/interfaces"
	"devops/common/config"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const healthPath = "/health"

type RouterDependencies struct {
	Controllers  []interfaces.Controller
	AuthConfig   *config.AuthConfig
	ServerConfig *config.ServerConfig
}

type Router struct {
	e         *echo.Echo
	ctrls     []interfaces.Controller
	authCfg   *config.AuthConfig
	serverCfg *config.ServerConfig
}

func (r *Router) GetRouterInstance() *echo.Echo {
	return r.e
}

func NewRouter(deps *RouterDependencies) *Router {
	r := &Router{
		e:         echo.New(),
		ctrls:     deps.Controllers,
		authCfg:   deps.AuthConfig,
		serverCfg: deps.ServerConfig,
	}

	r.setup()

	return r
}

func (r *Router) setup() {
	r.e.Validator = &CustomValidator{validator: validator.New()}
	r.registerHealthCheck()
	r.registerMiddlewares()
	r.registerControllers()
}

func (r *Router) registerMiddlewares() {
	r.e.Use(middleware.Recover())
	r.e.Use(middleware.RequestID())
	r.e.Use(middleware.Logger())
	r.e.Use(
		middleware.KeyAuthWithConfig(
			middleware.KeyAuthConfig{
				KeyLookup: fmt.Sprintf("header:%s", r.authCfg.KeyName),
				Validator: func(key string, c echo.Context) (bool, error) {
					return 1 == subtle.ConstantTimeCompare(
						[]byte(key),
						[]byte(r.authCfg.KeyVal),
					), nil
				},
			},
		),
	)
}

func (r *Router) registerHealthCheck() {
	r.e.GET(healthPath, func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
}

func (r *Router) registerControllers() {
	v1 := r.e.Group("/v1")

	for _, ctrl := range r.ctrls {
		ctrl.RegisterRoutes(v1)
	}
}
