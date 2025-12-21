package export

import gen "devops/seeder/internal/db/gen"

type Data struct {
	Locations       []gen.TempCheckerLocation       `json:"locations"`
	LocationSensors []gen.TempCheckerLocationSensor `json:"location_sensors"`
	SensorData      []gen.TempCheckerSensorData     `json:"sensor_data"`
}
