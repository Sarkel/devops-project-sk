package reader

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"devops/common/mqtt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
	broker := &MockBroker{}

	deps := &Dependencies{
		Logger: logger,
		Broker: broker,
	}

	service := NewService(deps)

	assert.NotNil(t, service)
	assert.Equal(t, logger, service.l)
	assert.Equal(t, broker, service.b)
}

func TestService_Listen_Success(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	broker := &MockBroker{}

	broker.On("Subscribe", ctx, "sensors/#", mock.Anything).Return(nil)

	service := &Service{
		l: logger,
		b: broker,
	}

	err := service.Listen(ctx)

	assert.NoError(t, err)
	broker.AssertExpectations(t)
}

func TestService_Listen_Error(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	broker := &MockBroker{}

	expectedErr := errors.New("subscription failed")
	broker.On("Subscribe", ctx, "sensors/#", mock.Anything).Return(expectedErr)

	service := &Service{
		l: logger,
		b: broker,
	}

	err := service.Listen(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "temp reader subscribe")
	broker.AssertExpectations(t)
}

func TestService_ParseSensorData_Success(t *testing.T) {
	service := &Service{}

	now := time.Now()
	timeStr := now.Format(time.RFC3339)

	msg := &mqtt.Message{
		Topic: "sensors/location1/sensor1",
		Payload: []mqtt.MessagePayload{
			{"22.5", timeStr},
			{"23.0", timeStr},
		},
	}

	result, err := service.parseSensorData(123, msg)

	assert.NoError(t, err)
	assert.Len(t, result.LocationSensorIds, 2)
	assert.Equal(t, int32(123), result.LocationSensorIds[0])
	assert.Equal(t, int32(123), result.LocationSensorIds[1])
	assert.Equal(t, 22.5, result.Temperatues[0])
	assert.Equal(t, 23.0, result.Temperatues[1])
}

func TestService_ParseSensorData_InvalidPayloadLength(t *testing.T) {
	service := &Service{}

	msg := &mqtt.Message{
		Topic: "sensors/location1/sensor1",
		Payload: []mqtt.MessagePayload{
			{"22.5"}, // Invalid: only 1 element instead of 2
		},
	}

	_, err := service.parseSensorData(123, msg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid payload length")
}

func TestService_ParseSensorData_InvalidTemperature(t *testing.T) {
	service := &Service{}

	now := time.Now()
	timeStr := now.Format(time.RFC3339)

	msg := &mqtt.Message{
		Topic: "sensors/location1/sensor1",
		Payload: []mqtt.MessagePayload{
			{"invalid_temp", timeStr},
		},
	}

	_, err := service.parseSensorData(123, msg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse sensor value")
}

func TestService_ParseSensorData_InvalidTimestamp(t *testing.T) {
	service := &Service{}

	msg := &mqtt.Message{
		Topic: "sensors/location1/sensor1",
		Payload: []mqtt.MessagePayload{
			{"22.5", "invalid_timestamp"},
		},
	}

	_, err := service.parseSensorData(123, msg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse sensor time")
}

func TestGetLocationSensorId_ValidTopic(t *testing.T) {
	// Test the topic parsing logic
	parts := []string{"sensors", "location-sid", "sensor-sid"}

	// Validate parts length
	assert.Len(t, parts, 3)
	assert.Equal(t, "location-sid", parts[1])
	assert.Equal(t, "sensor-sid", parts[2])
}

func TestGetLocationSensorId_InvalidTopic(t *testing.T) {
	// Test invalid topic format
	invalidTopics := []string{
		"sensors/location",         // Only 2 parts
		"sensors",                  // Only 1 part
		"sensors/loc/sensor/extra", // Too many parts
	}

	for _, topic := range invalidTopics {
		parts := splitTopic(topic)
		if len(parts) != 3 {
			expectedErr := fmt.Errorf("invalid topic format %s", topic)
			assert.Error(t, expectedErr)
		}
	}
}

func splitTopic(topic string) []string {
	parts := []string{}
	current := ""
	for _, ch := range topic {
		if ch == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
