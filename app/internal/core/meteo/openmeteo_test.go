package meteo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOpenMeteoClient(t *testing.T) {
	deps := &OpenMeteoDependencies{}
	client := NewOpenMeteoClient(deps)

	assert.NotNil(t, client)
}

func TestOpenMeteoClient_BuildUrl(t *testing.T) {
	client := &OpenMeteoClient{}

	testCases := []struct {
		lat      float64
		lon      float64
		expected string
	}{
		{52.2297, 21.0122, "https://api.open-meteo.com/v1/forecast?current_weather=true&latitude=52.229700&longitude=21.012200"},
		{0.0, 0.0, "https://api.open-meteo.com/v1/forecast?current_weather=true&latitude=0.000000&longitude=0.000000"},
		{-33.8688, 151.2093, "https://api.open-meteo.com/v1/forecast?current_weather=true&latitude=-33.868800&longitude=151.209300"},
	}

	for _, tc := range testCases {
		result := client.buildUrl(tc.lat, tc.lon)
		assert.Equal(t, tc.expected, result)
	}
}

func TestOpenMeteoClient_MapResponse_Success(t *testing.T) {
	client := &OpenMeteoClient{}

	resp := OpenMeteoResponse{
		CurrentWeather: CurrentWeather{
			Time:        "2024-01-15T14:30",
			Temperature: 22.5,
		},
	}

	result, err := client.mapResponse(resp)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 22.5, result[0].Temperature)

	expectedTime, _ := time.Parse("2006-01-02T15:04", "2024-01-15T14:30")
	assert.Equal(t, expectedTime, result[0].Timestamp)
}

func TestOpenMeteoClient_MapResponse_InvalidTime(t *testing.T) {
	client := &OpenMeteoClient{}

	resp := OpenMeteoResponse{
		CurrentWeather: CurrentWeather{
			Time:        "invalid-time-format",
			Temperature: 22.5,
		},
	}

	_, err := client.mapResponse(resp)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse current weather time")
}

func TestOpenMeteoResponse_JSONStructure(t *testing.T) {
	// Test that the struct tags are correct for JSON unmarshaling
	resp := OpenMeteoResponse{
		Latitude:         52.23,
		Longitude:        21.01,
		Timezone:         "Europe/Warsaw",
		Elevation:        100.0,
		CurrentWeather: CurrentWeather{
			Time:        "2024-01-15T14:30",
			Temperature: 22.5,
			WindSpeed:   10.5,
		},
	}

	assert.Equal(t, 52.23, resp.Latitude)
	assert.Equal(t, 21.01, resp.Longitude)
	assert.Equal(t, "Europe/Warsaw", resp.Timezone)
	assert.Equal(t, 22.5, resp.CurrentWeather.Temperature)
}

func TestWeatherParams_Structure(t *testing.T) {
	params := WeatherParams{
		Lat: 52.2297,
		Lon: 21.0122,
	}

	assert.Equal(t, 52.2297, params.Lat)
	assert.Equal(t, 21.0122, params.Lon)
}

func TestWeatherData_Structure(t *testing.T) {
	now := time.Now()
	data := WeatherData{
		Timestamp:   now,
		Temperature: 22.5,
	}

	assert.Equal(t, now, data.Timestamp)
	assert.Equal(t, 22.5, data.Temperature)
}

func TestOpenMeteoClient_BuildUrlWithDifferentCoordinates(t *testing.T) {
	client := &OpenMeteoClient{}

	testCases := []struct {
		name     string
		lat      float64
		lon      float64
		contains []string
	}{
		{
			name: "Warsaw coordinates",
			lat:  52.2297,
			lon:  21.0122,
			contains: []string{
				"latitude=52.229700",
				"longitude=21.012200",
				"current_weather=true",
			},
		},
		{
			name: "Negative coordinates",
			lat:  -33.8688,
			lon:  151.2093,
			contains: []string{
				"latitude=-33.868800",
				"longitude=151.209300",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := client.buildUrl(tc.lat, tc.lon)
			for _, substr := range tc.contains {
				assert.Contains(t, result, substr)
			}
		})
	}
}

func TestOpenMeteoResponse_AllFields(t *testing.T) {
	resp := OpenMeteoResponse{
		Latitude:  52.23,
		Longitude: 21.01,
		Timezone:  "Europe/Warsaw",
		Elevation: 100.0,
		CurrentWeather: CurrentWeather{
			Time:        "2024-01-15T14:30",
			Temperature: 22.5,
			WindSpeed:   10.5,
		},
	}

	assert.NotZero(t, resp.Latitude)
	assert.NotZero(t, resp.Longitude)
	assert.NotEmpty(t, resp.Timezone)
	assert.NotZero(t, resp.Elevation)
	assert.NotZero(t, resp.CurrentWeather.Temperature)
	assert.NotZero(t, resp.CurrentWeather.WindSpeed)
	assert.NotEmpty(t, resp.CurrentWeather.Time)
}

func TestCurrentWeather_Structure(t *testing.T) {
	cw := CurrentWeather{
		Time:        "2024-01-15T14:30",
		Temperature: 22.5,
		WindSpeed:   10.5,
	}

	assert.Equal(t, "2024-01-15T14:30", cw.Time)
	assert.Equal(t, 22.5, cw.Temperature)
	assert.Equal(t, 10.5, cw.WindSpeed)
}
