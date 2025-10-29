package main

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/tmozzze/SubChecker/internal/config"
	"github.com/tmozzze/SubChecker/internal/db"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config error:", err)
	}

	ctx := context.Background()

	db, err := db.NewDB(ctx, cfg)
	if err != nil {
		log.Fatal("db failed:", err)
	}
	defer db.Pool.Close()

	logger.Infof("Connected to %s", cfg.DBName)

}
