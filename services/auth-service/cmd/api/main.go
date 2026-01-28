package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/thuanvu301103/auth-service/internal/auth"
	"github.com/thuanvu301103/auth-service/internal/config"
	"github.com/thuanvu301103/auth-service/internal/database"
	"github.com/thuanvu301103/auth-service/internal/infrastructure/kafka"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/thuanvu301103/auth-service/docs"
)

// @title           Gopher Swagger Example API
// @version         1.0
// @description     This is a sample server for Auth Service.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  vungocthuan1234@gmail.com

// @host      localhost:3000
// @BasePath  /api/v1/

func main() {
	// 1. INITIALIZE STRUCTURED LOGGER
	// This replaces log.Printf with structured JSON for better monitoring
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 2. LOAD CONFIGURATION
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// 3. CONNECT TO DATABASE
	db, err := database.InitPostgres(cfg.DbURL)
	if err != nil {
		slog.Error("Database connection failed", "url", cfg.DbURL, "error", err)
		os.Exit(1)
	}

	// 4. AUTO MIGRATE
	if cfg.DbMigrate {
		slog.Info("Database auto-migration is enabled")
		if err := db.AutoMigrate(&auth.User{}); err != nil {
			slog.Error("Migration failed", "error", err)
			os.Exit(1)
		}
		slog.Info("Database migration successful")
	}

	// 5. SETUP INFRASTRUCTURE
	// We use the NewProducer which now returns (*Producer, error)
	kafkaProducer := kafka.NewProducer(cfg)
	/*if err != nil {
		slog.Error("Kafka infrastructure setup failed", "error", err)
		os.Exit(1)
	}*/
	defer kafkaProducer.Close()

	// 6. SETUP GIN ROUTER
	gin.SetMode(gin.ReleaseMode) // Optional: cleaner logs in production
	r := gin.Default()

	// 7. MAP ROUTER
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth.MapRoutes(r, db, kafkaProducer)

	// 8. START SERVER
	slog.Info("Application server is starting", "port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		slog.Error("Server shutdown unexpectedly", "error", err)
		os.Exit(1)
	}
}
