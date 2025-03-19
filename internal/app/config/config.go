package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"go-monolith/pkg/logger"
	"go-monolith/pkg/metrics"
)

// Config holds all configuration for the application
type Config struct {
	Environment string
	Logger      logger.Config
	Metrics     metrics.Config
	Server      ServerConfig
	DB          DBConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port           string
	EnableHTTPLogs bool
}

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Host     string  `env:"METRICS_HOST" envDefault:"localhost"`
	Port     int     `env:"METRICS_PORT" envDefault:"8125"`
	Prefix   string  `env:"METRICS_PREFIX" envDefault:"go-monolith"`
	Sampling float64 `env:"METRICS_SAMPLING" envDefault:"1.0"`
	Enabled  bool    `env:"METRICS_ENABLED" envDefault:"true"`
}

// validate checks if all required fields are set
func (c *DBConfig) validate() error {
	required := map[string]string{
		"Host":     c.Host,
		"Port":     c.Port,
		"User":     c.User,
		"Password": c.Password,
		"DBName":   c.DBName,
	}

	for field, value := range required {
		if value == "" {
			return fmt.Errorf("database %s is required", field)
		}
	}
	return nil
}

// NewConfig creates a new configuration instance
func NewConfig() (*Config, error) {
	// Load environment variables
	env := os.Getenv("APP_ENV")
	if env == "" {
		return nil, errors.New("APP_ENV is required")
	}

	// Create logger config
	logConfig := logger.Config{
		Environment:    env,
		EnableHTTPLogs: getEnvOrDefault("ENABLE_HTTP_LOGS", "true") == "true",
		LogLevel:       getEnvOrDefault("LOG_LEVEL", "debug"),
		Format:         getEnvOrDefault("LOG_FORMAT", "console"),
	}

	// Create metrics config
	metricsPort, err := strconv.Atoi(getEnvOrDefault("STATSD_PORT", "8125"))
	if err != nil {
		return nil, fmt.Errorf("invalid STATSD_PORT: %w", err)
	}

	sampling, err := strconv.ParseFloat(getEnvOrDefault("STATSD_SAMPLING", "1.0"), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid STATSD_SAMPLING: %w", err)
	}

	metricsConfig := metrics.Config{
		Host:     getEnvOrDefault("STATSD_HOST", "localhost"),
		Port:     metricsPort,
		Prefix:   getEnvOrDefault("STATSD_PREFIX", "go-monolith"),
		Sampling: sampling,
		Enabled:  getEnvOrDefault("METRICS_ENABLED", "true") == "true",
	}

	dbConfig := DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
	}

	if err := dbConfig.validate(); err != nil {
		return nil, err
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		return nil, errors.New("SERVER_PORT is required")
	}

	serverConfig := ServerConfig{
		Port:           ":" + serverPort,
		EnableHTTPLogs: logConfig.EnableHTTPLogs,
	}

	return &Config{
		Environment: env,
		Logger:      logConfig,
		Metrics:     metricsConfig,
		Server:      serverConfig,
		DB:          dbConfig,
	}, nil
}

// DSN returns the database connection string
func (c *DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.User, c.Password, c.Host, c.Port, c.DBName)
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
