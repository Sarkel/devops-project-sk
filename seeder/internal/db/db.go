package db

import (
	cDB "devops/common/db"
	sqlc "devops/seeder/internal/db/gen"
)

func WithQ(c *cDB.ConManager) *sqlc.Queries {
	return sqlc.New(c.GetDB())
}
