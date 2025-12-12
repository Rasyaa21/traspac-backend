package routes

import (
	"gin-backend-app/internal/controllers"
	"gin-backend-app/internal/middleware"
	"gin-backend-app/internal/repositories"
	"gin-backend-app/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserBudgetRoutes(api *gin.RouterGroup, db *gorm.DB) {
	budgetGroup := api.Group("/budget")

	userBudgetRepo := repositories.NewUserBudgetRepository(db)
	userRepo := repositories.NewUserRepository(db)
	userBudgetService := services.NewUserBudgetService(userRepo, userBudgetRepo)
	userBudgetController := controllers.NewUserBudgetController(userBudgetService)
	{
		// Public route for budget calculation preview
		budgetGroup.GET("/calculate", userBudgetController.CalculateBudget)

		// Protected routes - require authentication
		budgetGroup.Use(middleware.AuthMiddleware())
		{
			// Separated budget operations
			budgetGroup.POST("", userBudgetController.CreateUserBudget)    // Create new budget
			budgetGroup.PUT("", userBudgetController.UpdateUserBudget)     // Update existing budget
			budgetGroup.GET("", userBudgetController.GetUserBudget)        // Get budget details
			budgetGroup.DELETE("", userBudgetController.DeleteUserBudget)  // Delete budget
		}
	}
}