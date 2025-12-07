package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func LoggerMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			level := slog.LevelInfo
			if c.Response().Status >= 400 && c.Response().Status < 500 {
				level = slog.LevelWarn
			} else if c.Response().Status >= 500 {
				level = slog.LevelError
			}

			attrs := []slog.Attr{
				slog.String("method", c.Request().Method),
				slog.String("uri", c.Request().RequestURI),
				slog.Int("status", c.Response().Status),
				slog.Duration("latency", time.Since(start)),
				slog.String("remote_ip", c.RealIP()),
				slog.Int64("bytes_in", c.Request().ContentLength),
				slog.Int64("bytes_out", c.Response().Size),
			}

			if err != nil {
				attrs = append(attrs, slog.String("error", err.Error()))
				level = slog.LevelError
			}

			logger.LogAttrs(c.Request().Context(), level, "HTTP Request", attrs...)

			return err
		}
	}
}
