package db

import (
	sqlc "devops/app/internal/db/gen"
	cDB "devops/common/db"
)

func WithQ(c *cDB.ConManager) *sqlc.Queries {
	return sqlc.New(c.GetDB())
}
