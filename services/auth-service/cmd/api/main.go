package main

import (
	"log"

	"github.com/gin-gonic/gin"
	//"github.com/thuanvu301103/auth-service/internal/auth"
	"github.com/thuanvu301103/auth-service/internal/config"
	//"github.com/thuanvu301103/auth-service/internal/database"
)

func main() {
	// 1. LOAD CONFIGURATION (Viper)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 2. CONNECT TO DATABASE (GORM)
	/*db, err := database.InitMySQL(cfg.DbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}*/

	// 3. INITIALIZE LAYERS (Dependency Injection)
	// Repository -> Service -> Controller
	/*authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg.JwtSecret)
	authController := auth.NewController(authService)*/

	// 4. SETUP GIN ROUTER
	r := gin.Default()

	// 5. DEFINE ROUTES
	/*api := r.Group("/api/v1")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authController.Register)
			authRoutes.POST("/login", authController.Login)
		}
	}*/

	// 6. START SERVER
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
