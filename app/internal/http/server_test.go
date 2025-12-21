package http

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"devops/common/config"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestServer_Shutdown(t *testing.T) {
	e := echo.New()
	r := &Router{e: e}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := &config.ServerConfig{
		Port: "8080",
	}

	deps := ServerDependencies{
		Router: r,
		Logger: logger,
		Config: cfg,
	}

	server := NewServer(deps)

	// Start server in background
	go func() {
		// Ignore error from Start since we're testing Shutdown
		_ = e.Start(":0")
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	// Shutdown might return ErrServerClosed which is expected
	if err != nil && err != http.ErrServerClosed {
		t.Errorf("Unexpected shutdown error: %v", err)
	}
}

func TestServer_ShutdownWithoutStart(t *testing.T) {
	e := echo.New()
	r := &Router{e: e}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg := &config.ServerConfig{
		Port: "8080",
	}

	deps := ServerDependencies{
		Router: r,
		Logger: logger,
		Config: cfg,
	}

	server := NewServer(deps)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	// Should not error if server was never started
	assert.NoError(t, err)
}

func TestServer_ConfigurationPassing(t *testing.T) {
	e := echo.New()
	r := &Router{e: e}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	testCases := []struct {
		port string
	}{
		{"8080"},
		{"3000"},
		{"9999"},
	}

	for _, tc := range testCases {
		cfg := &config.ServerConfig{
			Port: tc.port,
		}

		deps := ServerDependencies{
			Router: r,
			Logger: logger,
			Config: cfg,
		}

		server := NewServer(deps)

		assert.Equal(t, tc.port, server.cfg.Port)
	}
}

func TestNewServer(t *testing.T) {
	e := echo.New()
	r := &Router{e: e}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.ServerConfig{Port: "8080"}

	deps := ServerDependencies{
		Router: r,
		Logger: logger,
		Config: cfg,
	}

	server := NewServer(deps)

	assert.NotNil(t, server)
	assert.NotNil(t, server.e)
	assert.NotNil(t, server.log)
	assert.NotNil(t, server.cfg)
	assert.Equal(t, e, server.e)
	assert.Equal(t, logger, server.log)
	assert.Equal(t, cfg, server.cfg)
}

func TestServer_WaitForShutdownWithServerError(t *testing.T) {
	e := echo.New()
	r := &Router{e: e}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.ServerConfig{Port: "8080"}

	deps := ServerDependencies{
		Router: r,
		Logger: logger,
		Config: cfg,
	}

	server := NewServer(deps)

	svrErr := make(chan error, 1)
	expectedErr := assert.AnError
	svrErr <- expectedErr

	err := server.waitForShutdown(svrErr)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestServer_ShutdownMethod(t *testing.T) {
	e := echo.New()
	r := &Router{e: e}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.ServerConfig{Port: "8080"}

	deps := ServerDependencies{
		Router: r,
		Logger: logger,
		Config: cfg,
	}

	server := NewServer(deps)

	// Start server in background so we have something to shutdown
	go func() {
		_ = e.Start(":0")
	}()

	time.Sleep(50 * time.Millisecond)

	err := server.shutdown()

	// Should not error or should return ErrServerClosed
	if err != nil && err != http.ErrServerClosed {
		t.Errorf("Unexpected shutdown error: %v", err)
	}
}
