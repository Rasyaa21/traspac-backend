package routes

import (
	"gin-backend-app/internal/controllers"
	"gin-backend-app/internal/middleware"
	"gin-backend-app/internal/repositories"
	"gin-backend-app/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CategoryRoutes(api *gin.RouterGroup, db *gorm.DB) {
	categoryGroup := api.Group("/categories")

	// Initialize dependencies
	categoryRepo := repositories.NewCategoryRepository(db)
	userRepo := repositories.NewUserRepository(db)
	categoryService := services.NewCategoryService(categoryRepo, userRepo)
	categoryController := controllers.NewCategoryController(categoryService)

	// All routes require authentication
	categoryGroup.Use(middleware.AuthMiddleware())
	{
		categoryGroup.POST("", categoryController.CreateCategory)      // Create category
		categoryGroup.GET("", categoryController.GetCategories)        // Get all categories
		categoryGroup.GET("/:id", categoryController.GetCategoryByID)  // Get category by ID
		categoryGroup.PUT("/:id", categoryController.UpdateCategory)   // Update category
		categoryGroup.DELETE("/:id", categoryController.DeleteCategory) // Delete category
	}
}