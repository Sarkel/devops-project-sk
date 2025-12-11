package db

import (
	"database/sql"
	"devops/common/config"
	"devops/common/logger"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
)

type Dependencies struct {
	Logger *slog.Logger
	Config *config.DatabaseConfig
}

type ConManager struct {
	db  *sql.DB
	log *slog.Logger
}

func NewConManager(deps Dependencies) (*ConManager, error) {
	cfg := deps.Config
	log := deps.Logger

	pgxCfg, err := pgx.ParseConfig(cfg.URL)

	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	if cfg.Debug {
		pgxCfg.Tracer = &tracelog.TraceLog{
			Logger:   tracelog.LoggerFunc(logger.TraceDBLogs(log)),
			LogLevel: tracelog.LogLevelDebug,
		}
	}

	connector := stdlib.GetConnector(*pgxCfg)

	db := sql.OpenDB(connector)

	db.SetMaxOpenConns(cfg.ConPool)

	return &ConManager{
		db:  db,
		log: deps.Logger,
	}, nil
}

func (c *ConManager) Close() error {
	return c.db.Close()
}

func (c *ConManager) GetDB() *sql.DB {
	return c.db
}

func Close(conManager *ConManager, log *slog.Logger) {
	if err := conManager.Close(); err != nil {
		log.Error("failed to close database connection", err)
	}
}
