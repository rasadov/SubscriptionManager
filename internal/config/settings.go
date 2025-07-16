package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type LogConfig struct {
	Level string
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnvInt("SERVER_PORT", 8080),
			Host: getEnvString("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnvString("POSTGRES_HOST", "localhost"),
			Port:     getEnvInt("POSTGRES_PORT", 5432),
			User:     getEnvString("POSTGRES_USER", "postgres"),
			Password: getEnvString("POSTGRES_PASSWORD", "password"),
			DBName:   getEnvString("POSTGRES_DB", "subscriptions"),
			SSLMode:  getEnvString("POSTGRES_SSLMODE", "disable"),
		},
		Log: LogConfig{
			Level: getEnvString("LOG_LEVEL", "info"),
		},
	}

	return config, nil
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
