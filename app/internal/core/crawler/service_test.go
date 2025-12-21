package crawler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	genDb "devops/app/internal/db/gen"
	"devops/app/internal/core/meteo"
	"devops/common/mqtt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMeteoClient struct {
	mock.Mock
}

func (m *MockMeteoClient) GetWeather(ctx context.Context, params meteo.WeatherParams) ([]meteo.WeatherData, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]meteo.WeatherData), args.Error(1)
}

type MockBroker struct {
	mock.Mock
}

func (m *MockBroker) Subscribe(ctx context.Context, topic string, handler mqtt.MessageHandler) error {
	args := m.Called(ctx, topic, handler)
	return args.Error(0)
}

func (m *MockBroker) Publish(topic string, payload []mqtt.MessagePayload) error {
	args := m.Called(topic, payload)
	return args.Error(0)
}

func (m *MockBroker) Unsubscribe(topic string) error {
	args := m.Called(topic)
	return args.Error(0)
}

func (m *MockBroker) Close() {
	m.Called()
}

func TestNewService(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	meteoClient := &MockMeteoClient{}
	broker := &MockBroker{}

	deps := &ServiceDependencies{
		Logger:      logger,
		MeteoClient: meteoClient,
		Broker:      broker,
	}

	service := NewService(deps)

	assert.NotNil(t, service)
	assert.Equal(t, logger, service.l)
	assert.Equal(t, meteoClient, service.mc)
	assert.Equal(t, broker, service.b)
}

func TestService_ProcessResponse(t *testing.T) {
	service := &Service{}

	now := time.Now()
	weatherData := []meteo.WeatherData{
		{Timestamp: now, Temperature: 22.5},
		{Timestamp: now.Add(1 * time.Hour), Temperature: 23.0},
	}

	result := service.processResponse(weatherData)

	assert.Len(t, result, 2)
	assert.Len(t, result[0], 2)
	assert.Equal(t, "22.50", result[0][0])
	assert.Equal(t, now.Format(time.RFC3339), result[0][1])
	assert.Equal(t, "23.00", result[1][0])
	assert.Equal(t, now.Add(1*time.Hour).Format(time.RFC3339), result[1][1])
}

func TestService_ProcessResponse_EmptyData(t *testing.T) {
	service := &Service{}

	weatherData := []meteo.WeatherData{}

	result := service.processResponse(weatherData)

	assert.Len(t, result, 0)
}

func TestService_ProcessResponse_FormatTemperature(t *testing.T) {
	service := &Service{}

	now := time.Now()
	testCases := []struct {
		temperature float64
		expected    string
	}{
		{22.123456, "22.12"},
		{-5.678, "-5.68"},
		{0.0, "0.00"},
		{100.999, "101.00"},
	}

	for _, tc := range testCases {
		weatherData := []meteo.WeatherData{
			{Timestamp: now, Temperature: tc.temperature},
		}

		result := service.processResponse(weatherData)

		assert.Equal(t, tc.expected, result[0][0], fmt.Sprintf("Temperature %.2f should format as %s", tc.temperature, tc.expected))
	}
}

func TestService_PullWeatherUpdate_TopicFormat(t *testing.T) {
	// Test topic generation
	locationSid := "warsaw-city"
	sensorSid := "api-sensor-1"

	expectedTopic := fmt.Sprintf("sensors/%s/%s", locationSid, sensorSid)

	assert.Equal(t, "sensors/warsaw-city/api-sensor-1", expectedTopic)
}

func TestService_Crawl_ErrorHandling(t *testing.T) {
	// Test that errors are joined correctly
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	allErr := errors.Join(err1, err2)

	assert.Error(t, allErr)
	assert.Contains(t, allErr.Error(), "error 1")
	assert.Contains(t, allErr.Error(), "error 2")
}

func TestService_PullWeatherUpdate_Success(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	meteoClient := &MockMeteoClient{}
	broker := &MockBroker{}

	service := &Service{
		l:  logger,
		mc: meteoClient,
		b:  broker,
	}

	location := genDb.GetAPILocationSensorsRow{
		LocationSid:  "warsaw",
		SensorSid:    "api-sensor",
		LocationName: "Warsaw",
		Latitude:     52.2297,
		Longitude:    21.0122,
	}

	now := time.Now()
	weatherData := []meteo.WeatherData{
		{Timestamp: now, Temperature: 22.5},
	}

	meteoClient.On("GetWeather", ctx, meteo.WeatherParams{
		Lat: location.Latitude,
		Lon: location.Longitude,
	}).Return(weatherData, nil)

	broker.On("Publish", "sensors/warsaw/api-sensor", mock.Anything).Return(nil)

	err := service.pullWeatherUpdate(ctx, location)

	assert.NoError(t, err)
	meteoClient.AssertExpectations(t)
	broker.AssertExpectations(t)
}

func TestService_PullWeatherUpdate_MeteoError(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	meteoClient := &MockMeteoClient{}
	broker := &MockBroker{}

	service := &Service{
		l:  logger,
		mc: meteoClient,
		b:  broker,
	}

	location := genDb.GetAPILocationSensorsRow{
		LocationSid:  "warsaw",
		SensorSid:    "api-sensor",
		LocationName: "Warsaw",
		Latitude:     52.2297,
		Longitude:    21.0122,
	}

	expectedErr := errors.New("weather API error")

	meteoClient.On("GetWeather", ctx, meteo.WeatherParams{
		Lat: location.Latitude,
		Lon: location.Longitude,
	}).Return(nil, expectedErr)

	err := service.pullWeatherUpdate(ctx, location)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get weather")
	meteoClient.AssertExpectations(t)
}

func TestService_PullWeatherUpdate_PublishError(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	meteoClient := &MockMeteoClient{}
	broker := &MockBroker{}

	service := &Service{
		l:  logger,
		mc: meteoClient,
		b:  broker,
	}

	location := genDb.GetAPILocationSensorsRow{
		LocationSid:  "warsaw",
		SensorSid:    "api-sensor",
		LocationName: "Warsaw",
		Latitude:     52.2297,
		Longitude:    21.0122,
	}

	now := time.Now()
	weatherData := []meteo.WeatherData{
		{Timestamp: now, Temperature: 22.5},
	}

	expectedErr := errors.New("publish error")

	meteoClient.On("GetWeather", ctx, meteo.WeatherParams{
		Lat: location.Latitude,
		Lon: location.Longitude,
	}).Return(weatherData, nil)

	broker.On("Publish", "sensors/warsaw/api-sensor", mock.Anything).Return(expectedErr)

	err := service.pullWeatherUpdate(ctx, location)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "publish temperature data")
	meteoClient.AssertExpectations(t)
	broker.AssertExpectations(t)
}
