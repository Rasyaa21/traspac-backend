package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gin-backend-app/cmd/server/docs"
	"gin-backend-app/internal/config"
	"gin-backend-app/internal/cron"
	"gin-backend-app/internal/repositories"
	"gin-backend-app/internal/routes"
	"gin-backend-app/internal/services"
	"gin-backend-app/pkg/utils"
	customValidator "gin-backend-app/pkg/validator"
)

// @title Traspac Backend API
// @version 1.0
// @description Backend API for Traspac Competition
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Set Gin mode based on environment
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Println("🚀 Starting Traspac Backend...")

	// ======================================================================
	// 0. Setup custom validators
	// ======================================================================
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password_strength", customValidator.ValidatePasswordStrength)
	}

	// ======================================================================
	// 1. Initialize database
	// ======================================================================
	log.Println("📦 Initializing database connection...")
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

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

	// ======================================================================
	// 2. Initialize repositories & services
	// ======================================================================

	baseURL := os.Getenv("BASE_URL")
	mailer := utils.NewSMTPMailerFromEnv()

	userRepo := repositories.NewUserRepository(db)
	userTokenRepo := repositories.NewUserTokenRepository(db)
	userBudgetRepo := repositories.NewUserBudgetRepository(db)

	userBudgetService := services.NewUserBudgetService(userRepo, userBudgetRepo)
	emailTokenService := services.NewEmailVerificationService(userRepo, userTokenRepo, mailer, baseURL)

	// ======================================================================
	// 3. Initialize cron scheduler
	// ======================================================================
	scheduler := cron.NewScheduler(emailTokenService, userBudgetService)
	scheduler.Start()
	defer scheduler.Stop()

	// ======================================================================
	// 4. Setup Gin router & routes
	// ======================================================================
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

	// API routes
	routes.SetupRoutes(router, db)


	log.Println("🛣️ Routes configured successfully")

	// Swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// ======================================================================
	// 5. Start HTTP server
	// ======================================================================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🌟 Server starting on port :%s", port)
	log.Printf("🔗 Health check: http://localhost:%s/health", port)
	log.Printf("🔗 DB Health check: http://localhost:%s/health/db", port)
	log.Printf("📚 Swagger docs: http://localhost:%s/swagger/index.html", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to run server: %v", err)
	}
}