package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rasadov/subscription-manager/internal/config"
	"github.com/rasadov/subscription-manager/pkg/database"
	"github.com/rasadov/subscription-manager/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	log := logger.NewLogger(cfg.Log.Level)
	log.Info("Starting subscription service...")

	db, err := database.NewPostgresDB(
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	log.Info("Database connected successfully")

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		log.Info("Server starting", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Error("Failed to get database connection", "error", err)
	} else {
		if err := sqlDb.Close(); err != nil {
			log.Error("Failed to close database connection", "error", err)
		} else {
			log.Info("Database connection closed successfully")
		}
	}

	log.Info("Server exited gracefully")
}
