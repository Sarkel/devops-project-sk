package seeder

import (
	cDB "devops/common/db"
	"embed"
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
)

const gooseSeedTable = "goose_seed_version"

type Dependencies struct {
	Db  *cDB.ConManager
	L   *slog.Logger
	Dir embed.FS
}

type Seeder struct {
	db  *cDB.ConManager
	l   *slog.Logger
	dir embed.FS
}

func New(deps Dependencies) *Seeder {
	return &Seeder{
		db:  deps.Db,
		l:   deps.L,
		dir: deps.Dir,
	}
}

func (s *Seeder) Run() error {
	goose.SetBaseFS(s.dir)
	goose.SetTableName(gooseSeedTable)
	goose.WithSlog(s.l)

	if err := goose.Up(s.db.GetDB(), "seeds"); err != nil {
		return fmt.Errorf("failed to run seeds: %w", err)
	}

	s.l.Info("Seeding completed successfully!")
	return nil
}
