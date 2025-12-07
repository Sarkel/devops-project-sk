package http

import (
	"crypto/subtle"
	"devops-project-sk/internal/config"
	"devops-project-sk/internal/http/interfaces"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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
	r.registerMiddlewares()
	r.registerHealthCheck()
	r.registerControllers()
}

func (r *Router) registerMiddlewares() {
	r.e.Use(middleware.Recover())
	r.e.Use(middleware.RequestID())
	r.e.Use(middleware.Logger())
	r.e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {

		isUsernameValid := 1 == subtle.ConstantTimeCompare([]byte(username), []byte(r.authCfg.Username))
		isPasswordValid := 1 == subtle.ConstantTimeCompare([]byte(password), []byte(r.authCfg.Password))

		return isUsernameValid && isPasswordValid, nil
	}))
}

func (r *Router) registerHealthCheck() {
	r.e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
}

func (r *Router) registerControllers() {
	api := r.e.Group("/api")

	v1 := api.Group("/v1")

	for _, ctrl := range r.ctrls {
		ctrl.RegisterRoutes(v1)
	}
}
