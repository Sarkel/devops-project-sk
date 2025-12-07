package location

import (
	"context"
	"devops-project-sk/internal/db"
)

type Dependencies struct {
	Db *db.ConManager
}

type Service struct {
	db *db.ConManager
}

func NewService(deps Dependencies) *Service {
	return &Service{
		db: deps.Db,
	}
}

func (s *Service) GetLocations(ctx context.Context) ([]Location, error) {
	q := s.db.WithQ()

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
