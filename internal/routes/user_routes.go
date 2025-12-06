package routes

import (
	_ "gin-backend-app/cmd/server/docs"
	"gin-backend-app/internal/controllers"
	"gin-backend-app/internal/repositories"
	"gin-backend-app/internal/services"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(api *gin.RouterGroup, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
    userController := controllers.NewUserController(userService)

    // Auth endpoints
    auth := api.Group("/auth")
    {
        auth.POST("/register", userController.RegisterUser)
        auth.POST("/login", userController.LoginUser)
    }

}