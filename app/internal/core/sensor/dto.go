package sensor

import (
	genDb "devops/app/internal/db/gen"
	"time"
)

type SummaryQs struct {
	LocationSid string `query:"location_sid" validate:"required"`
}

type Summary struct {
	Api   *SummaryItem `json:"api"`
	Local *SummaryItem `json:"local"`
}

type SummaryItem struct {
	Timestamp   time.Time `json:"timestamp"`
	Temperature float64   `json:"temperature"`
	Trend       int       `json:"trend"`
}

type DataQs struct {
	LocationSid   string                        `query:"location_sid" validate:"required"`
	StartDatetime time.Time                     `query:"start_datetime" validate:"required"`
	EndDatetime   time.Time                     `query:"end_datetime" validate:"required"`
	Aggregation   string                        `query:"aggregation" validate:"omitempty,oneof=day"`
	Types         []genDb.TempCheckerSensorType `query:"types" validate:"required,dive,required,oneof=api local"`
}

type DataPoint struct {
	Type        genDb.TempCheckerSensorType `json:"type"`
	Timestamp   time.Time                   `json:"timestamp"`
	Temperature float64                     `json:"temperature"`
}
