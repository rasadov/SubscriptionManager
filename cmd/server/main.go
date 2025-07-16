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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/rasadov/subscription-manager/internal/config"
	"github.com/rasadov/subscription-manager/internal/handlers"
	"github.com/rasadov/subscription-manager/internal/models"
	"github.com/rasadov/subscription-manager/internal/repository"
	"github.com/rasadov/subscription-manager/internal/service"
	"github.com/rasadov/subscription-manager/pkg/database"
	"github.com/rasadov/subscription-manager/pkg/logger"

	_ "github.com/rasadov/subscription-manager/docs"
)

// @title Subscription Manager API
// @version 1.0
// @description REST API for managing user subscriptions
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
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

	// Run migrations
	err = db.AutoMigrate(&models.Subscription{})
	if err != nil {
		log.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	log.Info("Database migrations completed")

	// Initialize repository, service and handlers
	subscriptionRepo := repository.NewSubscriptionRepositiry(db)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, log)

	// Setup Gin router
	if cfg.Server.Host == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Setup API routes
	api := router.Group("/api/v1")
	{
		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.POST("", subscriptionHandler.CreateSubscription)
			subscriptions.GET("", subscriptionHandler.ListSubscriptions)
			subscriptions.GET("/:id", subscriptionHandler.GetSubscription)
			subscriptions.PUT("/:id", subscriptionHandler.UpdateSubscription)
			subscriptions.DELETE("/:id", subscriptionHandler.DeleteSubscription)
			subscriptions.GET("/total-cost", subscriptionHandler.CalculateTotalCost)
		}
	}

	// Setup Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
		})
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		log.Info("Server starting", "port", cfg.Server.Port, "swagger_url", fmt.Sprintf("http://localhost:%d/swagger/index.html", cfg.Server.Port))
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
