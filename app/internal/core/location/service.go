package location

import (
	"context"
	"devops/app/internal/db"
	cDB "devops/common/db"
)

type Dependencies struct {
	Db *cDB.ConManager
}

type Service struct {
	db *cDB.ConManager
}

func NewService(deps Dependencies) *Service {
	return &Service{
		db: deps.Db,
	}
}

func (s *Service) GetLocations(ctx context.Context) ([]Location, error) {
	q := db.WithQ(s.db)

	locations, err := q.GetLocations(ctx)

	if err != nil {
		return nil, err
	}

	res := make([]Location, len(locations))
	for i, l := range locations {
		res[i] = Location{
			Name: l.LocationName,
			Sid:  l.LocationSid,
		}
	}

	return res, nil
}
