package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/thuanvu301103/auth-service/internal/auth"
	"github.com/thuanvu301103/auth-service/internal/config"
	"github.com/thuanvu301103/auth-service/internal/database"

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
	// 1. LOAD CONFIGURATION (Viper)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 2. CONNECT TO DATABASE (GORM)
	db, err := database.InitPostgres(cfg.DbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// 3. AUTO MIGRATE
	if cfg.DbMigrate {
		log.Println("Option DB_AUTO_MIGRATE is ON. Starting migration...")
		err := db.AutoMigrate(
			&auth.User{},
		)

		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Database migration completed!")
	} else {
		log.Println("Option DB_AUTO_MIGRATE is OFF. Skipping migration.")
	}

	// 4. SETUP GIN ROUTER
	r := gin.Default()

	// 5. MAP ROUTER
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth.MapRoutes(r, db)

	// 6. START SERVER
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
