package routes

import (
	_ "gin-backend-app/cmd/server/docs"
	"gin-backend-app/internal/controllers"
	"gin-backend-app/internal/middleware"
	"gin-backend-app/internal/repositories"
	"gin-backend-app/internal/services"
	"gin-backend-app/pkg/utils"
	"os"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(api *gin.RouterGroup, db *gorm.DB) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// Initialize mailer
	mailer := utils.NewSMTPMailerFromEnv()

	userRepo := repositories.NewUserRepository(db)
	userEmailVerificationRepo := repositories.NewUserTokenRepository(db)

	// Fix the service initialization - pass the mailer and baseURL
	userEmailVerificationService := services.NewEmailVerificationService(userRepo, userEmailVerificationRepo, mailer, baseURL)

	userService := services.NewUserService(userRepo, userEmailVerificationRepo, userEmailVerificationService)
	userController := controllers.NewUserController(userService)
	userEmailVerificationController := controllers.NewEmailVerificationController(userEmailVerificationService)

	auth := api.Group("/auth")
	{
		auth.POST("/register", userController.RegisterUser)
		auth.POST("/login", userController.LoginUser)
		
		auth.POST("/request-change-password", userController.SendEmailPasswordReset)
		auth.POST("/verify-otp-password-change", userController.GenerateAndSetVerificationToken)
		auth.POST("/change-password", userController.ValidateAndChangePassword)

		protected := auth.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/verify-email", userEmailVerificationController.VerifyEmail)
			protected.POST("/resend-verification", userEmailVerificationController.ResendEmailVerification)
			protected.GET("/verification-status", userEmailVerificationController.CheckVerificationStatus)


		}
	}
}