package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	//DB
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Server
	ServerPort string

	// Migrations
	MigrationsDir string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return &Config{}, err
	}

	cfg := &Config{
		// DB
		DBUser:     os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
		DBHost:     os.Getenv("POSTGRES_HOST"),
		DBPort:     os.Getenv("POSTGRES_PORT"),
		DBName:     os.Getenv("POSTGRES_DB"),

		// SERVER
		ServerPort: os.Getenv("SERVER_PORT"),

		MigrationsDir: os.Getenv("MIGRATIONS_DIR"),
	}

	if cfg.DBUser == "" || cfg.DBPassword == "" {
		err := errors.New("DB_USER or DB_PASSWORD is empty")
		return nil, err
	}

	if cfg.ServerPort == "" {
		err := errors.New("SERVER_PORT is empty")
		return nil, err
	}

	return cfg, nil
}
