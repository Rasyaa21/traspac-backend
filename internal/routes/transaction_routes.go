package routes

import (
	"gin-backend-app/internal/controllers"
	"gin-backend-app/internal/middleware"
	"gin-backend-app/internal/repositories"
	"gin-backend-app/internal/services"
	"gin-backend-app/pkg/utils"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TransactionRoutes(api *gin.RouterGroup, db *gorm.DB) {
	transactionGroup := api.Group("/transactions")

	transactionRepo := repositories.NewTransactionRepository(db)
	userRepo := repositories.NewUserRepository(db)
	userBudgetRepo := repositories.NewUserBudgetRepository(db)

	r2Client, err := utils.NewR2Client()
	if err != nil {
		log.Printf("Warning: R2 client initialization failed: %v. Photo upload will be disabled.", err)
		r2Client = nil 
	}

	transactionService := services.NewTransactionService(transactionRepo, userRepo, userBudgetRepo, r2Client)

	transactionController := controllers.NewTransactionController(transactionService)

	{
		transactionGroup.Use(middleware.AuthMiddleware(), middleware.RequireEmailVerified())
		{
			transactionGroup.POST("", transactionController.CreateTransaction)                // Create transaction with photo
			transactionGroup.GET("", transactionController.GetAllTransactions)               // Get all user transactions
			transactionGroup.GET("/:id", transactionController.GetTransactionByID)           // Get transaction by ID
			transactionGroup.PUT("/:id", transactionController.UpdateTransaction)            // Update transaction with photo
			transactionGroup.DELETE("/:id", transactionController.DeleteTransaction)         // Delete transaction and photo

			transactionGroup.GET("/period", transactionController.GetTransactionsByPeriod)   // Get transactions by period
		}
	}
}

