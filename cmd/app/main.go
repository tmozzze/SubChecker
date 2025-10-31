package main

import (
	"context"
	"fmt"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tmozzze/SubChecker/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tmozzze/SubChecker/internal/config"
	database "github.com/tmozzze/SubChecker/internal/db"
	httpHandler "github.com/tmozzze/SubChecker/internal/http"
	"github.com/tmozzze/SubChecker/internal/repository"
	"github.com/tmozzze/SubChecker/internal/service"
)

// @title SubChecker API
// @version 1.0
// @description REST service for subscription aggregation
// @host localhost:8080
// @BasePath /
func main() {
	// Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// Config
	cfg, err := config.Load()
	if err != nil {
		logger.WithError(err).Fatal("failed to load config")
	}
	logger.Info("Config is OK")

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// DB
	db, err := database.NewDB(ctx, cfg)
	if err != nil {
		logger.WithError(err).Fatal("db failed:", err)
	}
	defer db.Pool.Close()

	logger.Infof("Connected to %s on port %s", cfg.DBName, cfg.DBPort)

	// Migrations
	err = database.RunMigration(ctx, db.Pool, cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("failed to migrate db")
	}

	// Repository
	repo := repository.NewSubRepository(db.Pool, logger)

	// Service
	svc := service.NewSubService(repo, logger)

	// Hanlders
	handler := httpHandler.NewSubHandler(svc, logger)

	// Router
	router := gin.Default()

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// CRUD + SUM
	subs := router.Group("/subs")
	{
		subs.POST("", handler.CreateSub)
		subs.GET("", handler.ListSubs)
		subs.GET("/:sub_id", handler.GetSubById)
		subs.PUT("/:sub_id", handler.UpdateSub)
		subs.DELETE("/:sub_id", handler.DeleteSub)
		subs.GET("/sum", handler.SumCost)
	}

	// Start
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	logger.Infof("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		logger.Fatal("failed to start server:", err)
	}
}
