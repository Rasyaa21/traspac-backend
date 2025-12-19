package routes

import (
	"gin-backend-app/internal/controllers"
	"gin-backend-app/internal/middleware"
	"gin-backend-app/internal/repositories"
	"gin-backend-app/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PeriodReportRoutes(api *gin.RouterGroup, db *gorm.DB) {
	reportGroup := api.Group("/reports")

	periodReportRepo := repositories.NewPeriodReportRepository(db)
	periodReportService := services.NewPeriodReportService(periodReportRepo)
	periodReportController := controllers.NewPeriodReportController(periodReportService)

	{
		reportGroup.Use(middleware.AuthMiddleware(), middleware.RequireEmailVerified())
		{
			reportGroup.POST("/period", periodReportController.CreatePeriodReport) 
			reportGroup.POST("/:id", periodReportController.GetUserReportById) 
			reportGroup.GET("/user", periodReportController.GetAllUserReports)    
		}
	}
}