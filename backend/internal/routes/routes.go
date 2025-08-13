package routes

import (
	"digital-wellbeing-backend/internal/config"
	"digital-wellbeing-backend/internal/handlers"
	"digital-wellbeing-backend/internal/middleware"
	"digital-wellbeing-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(services *services.Services, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(services.Auth)
	userHandler := handlers.NewUserHandler(services.User)
	goalHandler := handlers.NewGoalHandler(services.Goal)
	checkInHandler := handlers.NewCheckInHandler(services.CheckIn)
	chatHandler := handlers.NewChatHandler(services.Chat)
	usageHandler := handlers.NewUsageHandler(services.Usage)

	// CORS middleware
	router.Use(middleware.CORSMiddleware(cfg))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
			}

			// Goal routes
			goals := protected.Group("/goals")
			{
				goals.GET("/", goalHandler.GetGoals)
				goals.POST("/", goalHandler.CreateGoal)
				goals.PUT("/:id", goalHandler.UpdateGoal)
				goals.DELETE("/:id", goalHandler.DeleteGoal)
				goals.GET("/progress", goalHandler.GetProgress)
			}

			// Check-in routes
			checkins := protected.Group("/checkins")
			{
				checkins.POST("/", checkInHandler.CreateCheckIn)
				checkins.GET("/", checkInHandler.GetCheckIns)
				checkins.GET("/today", checkInHandler.GetTodayCheckIn)
				checkins.GET("/stats", checkInHandler.GetMoodStats)
			}

			// Chat routes
			chat := protected.Group("/chat")
			{
				chat.POST("/message", chatHandler.SendMessage)
				chat.GET("/history", chatHandler.GetHistory)
			}

			// Usage routes
			usage := protected.Group("/usage")
			{
				usage.POST("/log", usageHandler.LogUsage)
				usage.GET("/daily", usageHandler.GetDailyUsage)
				usage.GET("/weekly", usageHandler.GetWeeklyUsage)
				usage.GET("/stats", usageHandler.GetStats)
			}

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/users", userHandler.GetAllUsers)
			}
		}
	}

	return router
}
