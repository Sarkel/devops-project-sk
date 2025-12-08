package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Logger      LoggerConfig
	MQTTBroker  MQTTBrokerConfig
	Auth        AuthConfig
}

type ServerConfig struct {
	Port string
}

type AuthConfig struct {
	Username string
	Password string
}

type DatabaseConfig struct {
	URL     string
	Debug   bool
	ConPool int
}

type LoggerConfig struct {
	Level  string
	Format string
}

type MQTTBrokerConfig struct {
	URL              string
	Username         string
	Password         string
	ClientID         string
	PayloadSeparator string
}

func Load() (*Config, error) {
	env := os.Getenv("GO_ENV")

	if "" == env {
		env = os.Getenv("ENV")
	}

	if "" == env {
		env = "development"
	}

	_ = godotenv.Load(".env")
	_ = godotenv.Load(".env." + env)
	_ = godotenv.Load(".env.defaults")

	dbUrl, err := getDBConnectionString()

	if err != nil {
		return nil, err
	}

	dbConPool, err := getIntEnv("DB_CON_POOL", 10)

	if err != nil {
		return nil, err
	}

	mqttUrl, err := getMqttBrokerURL()

	if err != nil {
		return nil, err
	}

	// todo: consider adding validation of loaded envs
	config := &Config{
		Environment: env,
		Server: ServerConfig{
			Port: os.Getenv("API_PORT"),
		},
		Database: DatabaseConfig{
			URL:     dbUrl,
			Debug:   getBoolEnv("DB_DEBUG"),
			ConPool: dbConPool,
		},
		Logger: LoggerConfig{
			Level:  os.Getenv("LOG_LEVEL"),
			Format: os.Getenv("LOG_FORMAT"),
		},
		MQTTBroker: MQTTBrokerConfig{
			URL:              mqttUrl,
			Username:         os.Getenv("MQTT_BROKER_USERNAME"),
			Password:         os.Getenv("MQTT_BROKER_PASSWORD"),
			ClientID:         os.Getenv("MQTT_BROKER_CLIENT_ID"),
			PayloadSeparator: os.Getenv("MQTT_BROKER_PAYLOAD_SEPARATOR"),
		},
		Auth: AuthConfig{
			Username: os.Getenv("AUTH_USERNAME"),
			Password: os.Getenv("AUTH_PASSWORD"),
		},
	}

	return config, nil
}

func getIntEnv(key string, def int) (int, error) {
	val := os.Getenv(key)
	if "" == val {
		return def, nil
	}

	atoi, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("cannot parse %s env: %w", key, err)
	}
	return int(atoi), nil
}

func getBoolEnv(key string) bool {
	return os.Getenv(key) == "true"
}

func getDBConnectionString() (string, error) {
	u := os.Getenv("DB_USER")
	pwd := os.Getenv("DB_PASSWORD")
	db := os.Getenv("DB_NAME")
	h := os.Getenv("DB_HOST")

	p, err := getIntEnv("DB_PORT", -1)

	if -1 == p || err != nil {
		return "", fmt.Errorf("cannot parse db port :%w", err)
	}

	ssl := os.Getenv("DB_SSL_MODE")

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", u, pwd, h, p, db, ssl), nil
}

func getMqttBrokerURL() (string, error) {
	h := os.Getenv("MQTT_BROKER_HOST")
	p, err := getIntEnv("MQTT_BROKER_PORT", -1)

	if -1 == p || err != nil {
		return "", fmt.Errorf("cannot parse mqtt port :%w", err)
	}

	return fmt.Sprintf("tcp://%s:%d", h, p), nil
}
