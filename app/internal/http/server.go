package http

import (
	"context"
	"devops-project-sk/internal/config"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

type ServerDependencies struct {
	Router *Router
	Logger *slog.Logger
	Config *config.ServerConfig
}
type Server struct {
	e   *echo.Echo
	log *slog.Logger
	cfg *config.ServerConfig
}

func NewServer(deps ServerDependencies) *Server {
	s := &Server{
		log: deps.Logger,
		cfg: deps.Config,
		e:   deps.Router.GetRouterInstance(),
	}

	return s
}

func (s *Server) Start() error {
	svrErr := make(chan error, 1)

	go func() {
		s.log.Info("Starting HTTP server", slog.String("port", s.cfg.Port))
		if err := s.e.Start(":" + s.cfg.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("Server startup failed", slog.String("error", err.Error()))
			svrErr <- err
		}
	}()

	return s.waitForShutdown(svrErr)
}

func (s *Server) waitForShutdown(svrErr chan error) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		s.log.Info("Shutdown signal received, starting graceful shutdown...")
		return s.shutdown()
	case err := <-svrErr:
		return err
	}
}

func (s *Server) shutdown() error {
	// Create context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	s.log.Info("Shutting down HTTP server...")
	if err := s.e.Shutdown(shutdownCtx); err != nil {
		s.log.Error("HTTP server forced to shutdown", slog.String("error", err.Error()))
		return err
	}

	s.log.Info("HTTP server stopped gracefully")

	s.log.Info("Application shutdown complete")
	return nil
}

// Shutdown Force shutdown (for testing)
func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
