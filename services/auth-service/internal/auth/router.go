package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MapRoutes initializes all layers and defines the API routing
func MapRoutes(r *gin.Engine, db *gorm.DB) {
	// 1. Initialize Layers (Dependency Injection)
	authRepo := NewRepository(db)
	authService := NewService(authRepo)
	authController := NewController(authService)

	// 2. Define Groups
	api := r.Group("/api/v1")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authController.Register)
			authRoutes.POST("/login", authController.Login)
		}
	}
}
