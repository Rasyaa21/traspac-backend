package controllers

import (
	"gin-backend-app/internal/dto/common"
	"gin-backend-app/internal/dto/request"
	"gin-backend-app/internal/services"
	"gin-backend-app/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PeriodReportController struct {
	PeriodReportService *services.PeriodReportService
}

func NewPeriodReportController(periodReportService *services.PeriodReportService) *PeriodReportController {
	return &PeriodReportController{
		PeriodReportService: periodReportService,
	}
}

// CreatePeriodReport godoc
// @Summary Create period report
// @Description Generate period report for user within date range
// @Tags Period Report
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /reports/period [post]
func (c *PeriodReportController) CreatePeriodReport(ctx *gin.Context) {
	userId, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		common.SendError(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req request.PeriodRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid request parameters: "+err.Error())
		return
	}

	if req.StartDate.After(req.EndDate) {
		common.SendError(ctx, http.StatusBadRequest, "start_date cannot be after end_date")
		return
	}

	report, err := c.PeriodReportService.CreatePeriodReport(userId, req.StartDate, req.EndDate)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	common.SendResponse(ctx, http.StatusOK, report, "Period report generated successfully")
}

// GetAllUserReports godoc
// @Summary Get all user reports
// @Description Get all period reports for the authenticated user
// @Tags Period Report
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /reports/user [get]
func (c *PeriodReportController) GetAllUserReports(ctx *gin.Context) {
	userId, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		common.SendError(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	reports, err := c.PeriodReportService.GetAllUserReports(userId)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	common.SendResponse(ctx, http.StatusOK, reports, "User reports retrieved successfully")
}


// GetUserReportById godoc
// @Summary Get user report by ID
// @Description Get specific period report by report ID for the authenticated user
// @Tags Period Report
// @Accept json
// @Produce json
// @Param id path string true "Report ID (UUID)"
// @Security BearerAuth
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /reports/{id} [get]
func (c *PeriodReportController) GetUserReportById(ctx *gin.Context) {

	userId, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		common.SendError(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	reportId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		common.SendError(ctx, http.StatusBadRequest, "Invalid transaction ID format")
		return
	}

	reports, err := c.PeriodReportService.GetUserReportById(userId, reportId)
	if err != nil {
		common.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	common.SendResponse(ctx, http.StatusOK, reports, "User reports retrieved successfully")
}