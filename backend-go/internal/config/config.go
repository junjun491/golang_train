package config

import (
	"fmt"
	"os"
)

type Config struct {
	JWTSecret  string
	DatabaseURL string
}

func Load() (*Config, error) {
	cfg := &Config{
		JWTSecret:  os.Getenv("JWT_SECRET"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	return cfg, nil
}