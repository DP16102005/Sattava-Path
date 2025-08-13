package main

import (
	"log"
	"os"

	"digital-wellbeing-backend/internal/config"
	"digital-wellbeing-backend/internal/database"
	"digital-wellbeing-backend/internal/routes"
	"digital-wellbeing-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Redis
	redis := database.InitializeRedis(cfg)

	// Initialize services
	services := services.NewServices(db, redis, cfg)

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := routes.SetupRoutes(services, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
