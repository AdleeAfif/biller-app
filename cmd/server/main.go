package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nkamil/biller-app/internal/admin"
	"github.com/nkamil/biller-app/internal/auth"
	"github.com/nkamil/biller-app/internal/commitment"
	"github.com/nkamil/biller-app/internal/config"
	"github.com/nkamil/biller-app/internal/middleware"
	"github.com/nkamil/biller-app/internal/summary"
	"github.com/nkamil/biller-app/internal/user"
	"github.com/nkamil/biller-app/pkg/db"
	"github.com/nkamil/biller-app/pkg/jwt"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	mongodb, err := db.NewMongoDB(cfg.MongoDBURI, cfg.DatabaseName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Close()

	// Initialize JWT manager
	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)

	// Initialize handlers
	authHandler := auth.NewHandler(mongodb.Database, jwtManager)
	userHandler := user.NewHandler(mongodb.Database)
	commitmentHandler := commitment.NewHandler(mongodb.Database)
	summaryHandler := summary.NewHandler(mongodb.Database)
	adminHandler := admin.NewHandler(mongodb.Database)

	// Setup Gin router
	router := gin.Default()

	// Public routes
	api := router.Group("/api")
	{
		// Authentication
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := api.Group("/users/me")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			// Salary management
			protected.PUT("/salary/default", userHandler.SetDefaultSalary)
			protected.PUT("/salary/:year/:month", userHandler.SetMonthlySalary)

			// Commitments
			protected.POST("/commitments/default", commitmentHandler.SetDefaultCommitments)
			protected.POST("/commitments/:year/:month", commitmentHandler.SetMonthlyCommitments)
			protected.PATCH("/commitments/:year/:month/:commitment_id", commitmentHandler.UpdateCommitmentPaidStatus)

			// Summaries
			protected.GET("/summary/monthly/:year/:month", summaryHandler.GetMonthlySummary)
			protected.GET("/summary/yearly/:year", summaryHandler.GetYearlySummary)
		}

		// Admin routes
		adminGroup := api.Group("/admin")
		adminGroup.Use(middleware.AuthMiddleware(jwtManager), middleware.AdminMiddleware())
		{
			adminGroup.GET("/users", adminHandler.ListUsers)
			adminGroup.PUT("/users/:user_id", adminHandler.UpdateUser)
			adminGroup.DELETE("/users/:user_id", adminHandler.DeleteUser)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
