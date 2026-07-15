package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort  string
	StorageType string
	DatabaseDSN string
	BaseUrl     string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		StorageType: getEnv("STORAGE_TYPE", "memory"),
		DatabaseDSN: getEnv("DATABASE_DSN", "postgres://user:pass@localhost:5432/shortener?sslmode=disable"),
		BaseUrl:     getEnv("BASE_URL", "http://localhost:8080"),
	}

	if cfg.StorageType != "memory" && cfg.StorageType != "postgres" {
		return nil, fmt.Errorf("Invalid storage type: %s", cfg.StorageType)
	}
	return cfg, nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
