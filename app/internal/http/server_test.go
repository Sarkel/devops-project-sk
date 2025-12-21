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
