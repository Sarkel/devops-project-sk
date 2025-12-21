package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"devops/app/internal/http/interfaces"
	"devops/common/config"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockController struct {
	mock.Mock
}

func (m *MockController) RegisterRoutes(e *echo.Group) {
	m.Called(e)
}

func TestNewRouter(t *testing.T) {
	authConfig := &config.AuthConfig{
		KeyName: "X-API-Key",
		KeyVal:  "test-key",
	}

	serverConfig := &config.ServerConfig{
		Port: "8080",
	}

	mockCtrl := &MockController{}
	mockCtrl.On("RegisterRoutes", mock.Anything).Return()

	deps := &RouterDependencies{
		Controllers:  []interfaces.Controller{mockCtrl},
		AuthConfig:   authConfig,
		ServerConfig: serverConfig,
	}

	router := NewRouter(deps)

	assert.NotNil(t, router)
	assert.NotNil(t, router.e)
	assert.Equal(t, authConfig, router.authCfg)
	assert.Equal(t, serverConfig, router.serverCfg)
	assert.Len(t, router.ctrls, 1)

	mockCtrl.AssertExpectations(t)
}

func TestRouter_GetRouterInstance(t *testing.T) {
	authConfig := &config.AuthConfig{
		KeyName: "X-API-Key",
		KeyVal:  "test-key",
	}

	serverConfig := &config.ServerConfig{
		Port: "8080",
	}

	deps := &RouterDependencies{
		Controllers:  []interfaces.Controller{},
		AuthConfig:   authConfig,
		ServerConfig: serverConfig,
	}

	router := NewRouter(deps)
	instance := router.GetRouterInstance()

	assert.NotNil(t, instance)
	assert.IsType(t, &echo.Echo{}, instance)
}

func TestRouter_HealthCheck(t *testing.T) {
	authConfig := &config.AuthConfig{
		KeyName: "X-API-Key",
		KeyVal:  "test-key",
	}

	serverConfig := &config.ServerConfig{
		Port: "8080",
	}

	deps := &RouterDependencies{
		Controllers:  []interfaces.Controller{},
		AuthConfig:   authConfig,
		ServerConfig: serverConfig,
	}

	router := NewRouter(deps)
	e := router.GetRouterInstance()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("X-API-Key", "test-key") // Auth is required globally
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())
}

func TestRouter_AuthMiddleware_ValidKey(t *testing.T) {
	authConfig := &config.AuthConfig{
		KeyName: "X-API-Key",
		KeyVal:  "test-secret-key",
	}

	serverConfig := &config.ServerConfig{
		Port: "8080",
	}

	mockCtrl := &MockController{}
	mockCtrl.On("RegisterRoutes", mock.Anything).Run(func(args mock.Arguments) {
		group := args.Get(0).(*echo.Group)
		group.GET("/test", func(c echo.Context) error {
			return c.String(http.StatusOK, "authenticated")
		})
	})

	deps := &RouterDependencies{
		Controllers:  []interfaces.Controller{mockCtrl},
		AuthConfig:   authConfig,
		ServerConfig: serverConfig,
	}

	router := NewRouter(deps)
	e := router.GetRouterInstance()

	req := httptest.NewRequest(http.MethodGet, "/v1/test", nil)
	req.Header.Set("X-API-Key", "test-secret-key")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "authenticated", rec.Body.String())
}

func TestRouter_AuthMiddleware_InvalidKey(t *testing.T) {
	authConfig := &config.AuthConfig{
		KeyName: "X-API-Key",
		KeyVal:  "test-secret-key",
	}

	serverConfig := &config.ServerConfig{
		Port: "8080",
	}

	mockCtrl := &MockController{}
	mockCtrl.On("RegisterRoutes", mock.Anything).Run(func(args mock.Arguments) {
		group := args.Get(0).(*echo.Group)
		group.GET("/test", func(c echo.Context) error {
			return c.String(http.StatusOK, "authenticated")
		})
	})

	deps := &RouterDependencies{
		Controllers:  []interfaces.Controller{mockCtrl},
		AuthConfig:   authConfig,
		ServerConfig: serverConfig,
	}

	router := NewRouter(deps)
	e := router.GetRouterInstance()

	req := httptest.NewRequest(http.MethodGet, "/v1/test", nil)
	req.Header.Set("X-API-Key", "wrong-key")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRouter_AuthMiddleware_MissingKey(t *testing.T) {
	authConfig := &config.AuthConfig{
		KeyName: "X-API-Key",
		KeyVal:  "test-secret-key",
	}

	serverConfig := &config.ServerConfig{
		Port: "8080",
	}

	mockCtrl := &MockController{}
	mockCtrl.On("RegisterRoutes", mock.Anything).Run(func(args mock.Arguments) {
		group := args.Get(0).(*echo.Group)
		group.GET("/test", func(c echo.Context) error {
			return c.String(http.StatusOK, "authenticated")
		})
	})

	deps := &RouterDependencies{
		Controllers:  []interfaces.Controller{mockCtrl},
		AuthConfig:   authConfig,
		ServerConfig: serverConfig,
	}

	router := NewRouter(deps)
	e := router.GetRouterInstance()

	req := httptest.NewRequest(http.MethodGet, "/v1/test", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRouter_RegisterControllers(t *testing.T) {
	authConfig := &config.AuthConfig{
		KeyName: "X-API-Key",
		KeyVal:  "test-key",
	}

	serverConfig := &config.ServerConfig{
		Port: "8080",
	}

	mockCtrl1 := &MockController{}
	mockCtrl2 := &MockController{}

	mockCtrl1.On("RegisterRoutes", mock.Anything).Return()
	mockCtrl2.On("RegisterRoutes", mock.Anything).Return()

	deps := &RouterDependencies{
		Controllers:  []interfaces.Controller{mockCtrl1, mockCtrl2},
		AuthConfig:   authConfig,
		ServerConfig: serverConfig,
	}

	router := NewRouter(deps)

	assert.NotNil(t, router)

	mockCtrl1.AssertExpectations(t)
	mockCtrl2.AssertExpectations(t)
}

func TestRouter_ValidatorSetup(t *testing.T) {
	authConfig := &config.AuthConfig{
		KeyName: "X-API-Key",
		KeyVal:  "test-key",
	}

	serverConfig := &config.ServerConfig{
		Port: "8080",
	}

	deps := &RouterDependencies{
		Controllers:  []interfaces.Controller{},
		AuthConfig:   authConfig,
		ServerConfig: serverConfig,
	}

	router := NewRouter(deps)
	e := router.GetRouterInstance()

	assert.NotNil(t, e.Validator)
	assert.IsType(t, &CustomValidator{}, e.Validator)
}
