package config

import (
	"fmt"
	"os"

	"go-monolith/pkg/logger"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Environment string
	Logger      *logger.Logger
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

// NewConfig creates a new configuration instance
func NewConfig() (*Config, error) {
	// Load .env.local file only in development environment
	env := getEnvOrDefault("APP_ENV", "development")
	if env == "development" {
		if err := godotenv.Load(".env.local"); err != nil {
			return nil, fmt.Errorf("error loading .env.local file: %w", err)
		}
	}

	// Initialize logger first
	logConfig := logger.Config{
		Environment:    env,
		EnableHTTPLogs: getEnvOrDefault("ENABLE_HTTP_LOGS", "true") == "true",
		LogLevel:       getEnvOrDefault("LOG_LEVEL", "debug"),
		Format:         getEnvOrDefault("LOG_FORMAT", "console"),
	}

	log, err := logger.Initialize(logConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Create config
	config := &Config{
		Environment: logConfig.Environment,
		Logger:      log,
		Server: ServerConfig{
			Port:           getEnvOrDefault("SERVER_PORT", "8080"),
			EnableHTTPLogs: logConfig.EnableHTTPLogs,
		},
		DB: DBConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
			DBName:   getEnvOrDefault("DB_NAME", "monolith"),
		},
	}

	return config, nil
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
