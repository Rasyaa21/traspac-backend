package main

import (
	_ "gin-backend-app/cmd/server/docs"
	"gin-backend-app/internal/config"
	"gin-backend-app/internal/routes"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Traspac Backend API
// @version 1.0
// @description Backend API for Traspac Competition
// @host localhost:8080
// @BasePath /api/v1

func main() {
	// Set Gin mode based on environment
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Println("🚀 Starting Traspac Backend...")

	// Initialize database FIRST (before starting server)
	log.Println("📦 Initializing database connection...")
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Get SQL DB instance for connection management
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("❌ Failed to get SQL DB instance:", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("❌ Error closing database connection: %v", err)
		} else {
			log.Println("✅ Database connection closed successfully")
		}
	}()

	// Test database connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("❌ Failed to ping database:", err)
	}
	log.Println("✅ Database connected and ping successful")
	log.Printf("📊 Database Stats - Max Open Connections: %d", sqlDB.Stats().MaxOpenConnections)

	// Setup Gin router
	router := gin.Default()
	log.Println("🌐 Gin router initialized")

	// Basic health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "up",
			"service": "traspac-backend",
		})
	})

	// Database health check
	router.GET("/health/db", func(c *gin.Context) {
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "down",
				"database": "disconnected",
				"error":    err.Error(),
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"status":   "up",
			"database": "connected",
			"service":  "traspac-backend",
		})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Traspac Backend API",
			"version": "v1.0.0",
			"status":  "running",
			"docs":    "/swagger/index.html",
		})
	})

	// Setup all routes
	routes.SetupRoutes(router, db)
	
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	log.Println("🛣️ Routes configured successfully")

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🌟 Server starting on port :%s", port)
	log.Printf("🔗 Health check: http://localhost:%s/health", port)
	log.Printf("🔗 DB Health check: http://localhost:%s/health/db", port)
	log.Printf("📚 Swagger docs: http://localhost:%s/swagger/index.html", port)

	// Start the server (this blocks)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to run server: %v", err)
	}
}