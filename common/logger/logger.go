package logger

import (
	"context"
	"devops/common/config"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/tracelog"
)

type Dependencies struct {
	Config config.LoggerConfig
}

func New(deps Dependencies) *slog.Logger {
	l := getLevel(deps.Config.Level)

	ops := &slog.HandlerOptions{
		AddSource:   false,
		Level:       l,
		ReplaceAttr: nil,
	}

	handler := getHandler(deps.Config.Format, ops)

	return slog.New(handler)
}

func getHandler(format string, ops *slog.HandlerOptions) (h slog.Handler) {
	switch format {
	case "json":
		h = slog.NewJSONHandler(os.Stdout, ops)
	case "text":
		h = slog.NewTextHandler(os.Stdout, ops)
	default:
		h = slog.NewTextHandler(os.Stdout, ops)
	}
	return
}

func getLevel(level string) (l slog.Level) {
	switch level {
	case "info":
		l = slog.LevelInfo
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}
	return
}

func MapDBLogLevels(dbLogLevel tracelog.LogLevel) (slogLevel slog.Level) {
	switch dbLogLevel {
	case tracelog.LogLevelTrace:
		slogLevel = slog.LevelDebug
	case tracelog.LogLevelDebug:
		slogLevel = slog.LevelDebug
	case tracelog.LogLevelInfo:
		slogLevel = slog.LevelInfo
	case tracelog.LogLevelWarn:
		slogLevel = slog.LevelWarn
	case tracelog.LogLevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}
	return
}

func TraceDBLogs(log *slog.Logger) func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	return func(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
		slogLevel := MapDBLogLevels(level)
		log.Log(ctx, slogLevel, msg, slog.Any("pgx", data))
	}
}
