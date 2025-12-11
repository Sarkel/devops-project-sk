package sensor

import (
	"context"
	"devops/app/internal/db"
	genDb "devops/app/internal/db/gen"
	cDB "devops/common/db"
	"errors"
	"fmt"
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

func (s *Service) GetSummary(ctx context.Context, params SummaryQs) (Summary, error) {
	q := db.WithQ(s.db)

	locExist, err := q.LocationExistBySid(ctx, params.LocationSid)

	if err != nil {
		return Summary{}, err
	}

	if locExist == 0 {
		return Summary{}, errors.New("location not found")
	}

	sum, err := q.GetTodaySensorsSummary(ctx, params.LocationSid)

	if err != nil {
		return Summary{}, err
	}

	if len(sum) > 2 {
		return Summary{}, errors.New("unexpected sensors summary")
	}

	res := Summary{}

	for _, item := range sum {
		switch item.Type {
		case genDb.TempCheckerSensorTypeLocal:
			res.Local = &SummaryItem{
				Timestamp:   item.Date,
				Temperature: item.AvgTemperature,
				Trend:       0,
			}
			break
		case genDb.TempCheckerSensorTypeApi:
			res.Api = &SummaryItem{
				Timestamp:   item.Date,
				Temperature: item.AvgTemperature,
				Trend:       0,
			}
			break
		}
	}

	return res, nil
}

func (s *Service) GetData(ctx context.Context, params DataQs) ([]DataPoint, error) {
	q := db.WithQ(s.db)

	res, err := q.GetSensorDataPoints(ctx, genDb.GetSensorDataPointsParams{
		Aggregation:   params.Aggregation,
		LocationSid:   params.LocationSid,
		Types:         params.Types,
		StartDatetime: params.StartDatetime,
		EndDatetime:   params.EndDatetime,
	})

	if err != nil {
		return nil, fmt.Errorf("get sensor data points: %w", err)
	}

	points := make([]DataPoint, len(res))

	for i, r := range res {
		points[i] = DataPoint{
			Type:        r.Type,
			Timestamp:   r.TimeDim,
			Temperature: r.AvgTemperature,
		}
	}

	return points, nil
}
