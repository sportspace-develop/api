package config

import (
	"errors"
	"fmt"
	"os"

	"sport-space/internal/adapter/sender"
	"sport-space/internal/adapter/storage"
	"sport-space/internal/adapter/storage/database"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Store     storage.Config
	Sender    sender.Config
	Address   string `env:"HTTP_ADDRESS"`
	SecretKey string `env:"SECRET_KEY"`
}

func Init() (*Config, error) {
	cfg := &Config{
		Store: storage.Config{
			Database: &database.Config{},
		},
		Sender: sender.Config{},
	}

	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return cfg, fmt.Errorf("failed load enviorements from file: %w", err)
	}

	if err := env.Parse(cfg); err != nil {
		return cfg, fmt.Errorf("failed parse env: %w", err)
	}

	return cfg, nil
}
