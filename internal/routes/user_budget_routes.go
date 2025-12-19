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
		budgetGroup.Use(middleware.AuthMiddleware(), middleware.RequireEmailVerified())
		{
			budgetGroup.POST("", userBudgetController.CreateUserBudget)    
			budgetGroup.PUT("", userBudgetController.UpdateUserBudget)     
			budgetGroup.GET("", userBudgetController.GetUserBudget)        
			budgetGroup.DELETE("", userBudgetController.DeleteUserBudget)  
			budgetGroup.GET("/summary", userBudgetController.GetBudgetSummary)    
		}
	}
}