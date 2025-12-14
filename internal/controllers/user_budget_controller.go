package controllers

import (
	"gin-backend-app/internal/dto/common"
	"gin-backend-app/internal/dto/request"
	"gin-backend-app/internal/dto/response"
	"gin-backend-app/internal/services"
	"gin-backend-app/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserBudgetController struct {
	UserBudgetService *services.UserBudgetService
}

func NewUserBudgetController(userBudgetService *services.UserBudgetService) *UserBudgetController {
	return &UserBudgetController{
		UserBudgetService: userBudgetService,
	}
}

// CreateUserBudget godoc
// @Summary Create user budget
// @Description Create budget for authenticated user with default (50-30-20) or custom allocation
// @Tags Budget Management
// @Accept json
// @Produce json
// @Param request body request.BudgetRequest true "Budget data (custom percentages optional)"
// @Success 201 {object} common.Response "Budget created successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request data"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 409 {object} common.ErrorResponse "Budget already exists"
// @Security BearerAuth
// @Router /budget [post]
func (bc *UserBudgetController) CreateUserBudget(c *gin.Context) {
	var req request.BudgetRequest

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.WeeklyIncome <= 0 {
		common.SendError(c, http.StatusBadRequest, "Weekly income must be greater than 0")
		return
	}

	// Determine allocation (default or custom)
	var allocation *services.BudgetAllocation
	isCustom := false

	if req.SavingsPercent != nil && req.WantsPercent != nil && req.NeedsPercent != nil {
		// Custom allocation provided
		customAlloc := services.NewCustomAllocation(*req.SavingsPercent, *req.WantsPercent, *req.NeedsPercent)
		allocation = &customAlloc
		isCustom = true
	}

	budget, err := bc.UserBudgetService.CreateUserBudget(userID, req.WeeklyIncome, allocation)
	if err != nil {
		if err.Error() == "user budget already exists, use update instead" {
			common.SendError(c, http.StatusConflict, "Budget already exists. Use PUT to update existing budget")
			return
		}
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Calculate percentages for response
	savingsPercent := float64(budget.SavedMoney) / float64(budget.IncomeWeekly) * 100
	wantsPercent := float64(budget.WantsBudget) / float64(budget.IncomeWeekly) * 100
	needsPercent := float64(budget.NeedsBudget) / float64(budget.IncomeWeekly) * 100

	responseData := gin.H{
		"budget": gin.H{
			"id":             budget.ID,
			"user_id":        budget.UserID,
			"weekly_income":  budget.IncomeWeekly,
			"needs_budget":   budget.NeedsBudget,
			"wants_budget":   budget.WantsBudget,
			"savings_budget": budget.SavedMoney,
		},
		"allocation": gin.H{
			"savings_percent": savingsPercent,
			"wants_percent":   wantsPercent,
			"needs_percent":   needsPercent,
		},
		"breakdown": gin.H{
			"savings_amount": budget.SavedMoney,
			"wants_amount":   budget.WantsBudget,
			"needs_amount":   budget.NeedsBudget,
		},
		"is_custom": isCustom,
	}

	message := "Budget created successfully"
	if !isCustom {
		message += " with 50-30-20 rule"
	} else {
		message += " with custom allocation"
	}

	common.SendResponse(c, http.StatusCreated, responseData, message)
}

// UpdateUserBudget godoc
// @Summary Update user budget
// @Description Update existing budget for authenticated user with default (50-30-20) or custom allocation
// @Tags Budget Management
// @Accept json
// @Produce json
// @Param request body request.BudgetRequest true "Budget data (custom percentages optional)"
// @Success 200 {object} common.Response "Budget updated successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request data"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "Budget not found"
// @Security BearerAuth
// @Router /budget [put]
func (bc *UserBudgetController) UpdateUserBudget(c *gin.Context) {
	var req request.BudgetRequest

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.WeeklyIncome <= 0 {
		common.SendError(c, http.StatusBadRequest, "Weekly income must be greater than 0")
		return
	}

	var allocation *services.BudgetAllocation
	isCustom := false

	if req.SavingsPercent != nil && req.WantsPercent != nil && req.NeedsPercent != nil {
		customAlloc := services.NewCustomAllocation(*req.SavingsPercent, *req.WantsPercent, *req.NeedsPercent)
		allocation = &customAlloc
		isCustom = true
	}

	budget, err := bc.UserBudgetService.UpdateUserBudget(userID, req.WeeklyIncome, allocation)
	if err != nil {
		if err.Error() == "user budget not found" {
			common.SendError(c, http.StatusNotFound, "Budget not found. Use POST to create a new budget")
			return
		}
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	savingsPercent := float64(budget.SavedMoney) / float64(budget.IncomeWeekly) * 100
	wantsPercent := float64(budget.WantsBudget) / float64(budget.IncomeWeekly) * 100
	needsPercent := float64(budget.NeedsBudget) / float64(budget.IncomeWeekly) * 100

	responseData := response.UserBudgetResponse{
	Budget: response.BudgetResponse{
		ID:            budget.ID,
		UserID:        budget.UserID,
		WeeklyIncome:  budget.IncomeWeekly,
		NeedsBudget:   budget.NeedsBudget,
		WantsBudget:   budget.WantsBudget,
		SavingsBudget: budget.SavedMoney,
	},
	Allocation: response.AllocationResponse{
		NeedsPercent:   int(needsPercent),
		WantsPercent:   int(wantsPercent),
		SavingsPercent: int(savingsPercent),
	},
	Breakdown: response.BreakdownResponse{
		NeedsAmount:   budget.NeedsBudget,
		WantsAmount:   budget.WantsBudget,
		SavingsAmount: budget.SavedMoney,
	},
	IsCustom: isCustom,
}



	message := "Budget updated successfully"
	if !isCustom {
		message += " with 50-30-20 rule"
	} else {
		message += " with custom allocation"
	}

	common.SendResponse(c, http.StatusOK, responseData, message)
}

// GetUserBudget godoc
// @Summary Get user budget details
// @Description Retrieve detailed budget information for authenticated user
// @Tags Budget Management
// @Accept json
// @Produce json
// @Success 200 {object} common.Response "Budget details retrieved successfully"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "Budget not found"
// @Security BearerAuth
// @Router /budget [get]
func (bc *UserBudgetController) GetUserBudget(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	budget, err := bc.UserBudgetService.GetUserBudget(userID)
	if err != nil {
		if err.Error() == "user budget not found" {
			common.SendError(c, http.StatusNotFound, err.Error())
			return
		}
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	savingsPercent := float64(budget.SavedMoney) / float64(budget.IncomeWeekly) * 100
	wantsPercent := float64(budget.WantsBudget) / float64(budget.IncomeWeekly) * 100
	needsPercent := float64(budget.NeedsBudget) / float64(budget.IncomeWeekly) * 100

	responseData := response.UserBudgetResponse{
	Budget: response.BudgetResponse{
		ID:            budget.ID,
		UserID:        budget.UserID,
		WeeklyIncome:  budget.IncomeWeekly,
		NeedsBudget:   budget.NeedsBudget,
		WantsBudget:   budget.WantsBudget,
		SavingsBudget: budget.SavedMoney,
	},
	Allocation: response.AllocationResponse{
		NeedsPercent:   int(needsPercent),
		WantsPercent:   int(wantsPercent),
		SavingsPercent: int(savingsPercent),
	},
	Breakdown: response.BreakdownResponse{
		NeedsAmount:   budget.NeedsBudget,
		WantsAmount:   budget.WantsBudget,
		SavingsAmount: budget.SavedMoney,
	},
}


	common.SendResponse(c, http.StatusOK, responseData, "Budget details retrieved successfully")
}

// DeleteUserBudget godoc
// @Summary Delete user budget
// @Description Delete the budget for authenticated user
// @Tags Budget Management
// @Accept json
// @Produce json
// @Success 200 {object} common.Response "Budget deleted successfully"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "Budget not found"
// @Security BearerAuth
// @Router /budget [delete]
func (bc *UserBudgetController) DeleteUserBudget(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	err = bc.UserBudgetService.DeleteUserBudget(userID)
	if err != nil {
		if err.Error() == "user budget not found" {
			common.SendError(c, http.StatusNotFound, err.Error())
			return
		}
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"user_id":    userID,
		"deleted_at": "now",
		"status":     "deleted",
	}, "Budget deleted successfully")
}

// GetBudgetSummary godoc
// @Summary Get detailed budget summary with usage
// @Description Retrieve detailed budget summary including allocation, usage, and remaining amounts
// @Tags Budget Management
// @Accept json
// @Produce json
// @Success 200 {object} common.Response "Budget summary retrieved successfully"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "Budget not found"
// @Security BearerAuth
// @Router /budget/summary [get]
func (bc *UserBudgetController) GetBudgetSummary(c *gin.Context) {
    userID, err := utils.GetUserIDFromContext(c)
    if err != nil {
        common.SendError(c, http.StatusUnauthorized, "Authentication required")
        return
    }

    summary, err := bc.UserBudgetService.GetBudgetSummary(userID)
    if err != nil {
        if err.Error() == "user budget not found" {
            common.SendError(c, http.StatusNotFound, err.Error())
            return
        }
        common.SendError(c, http.StatusBadRequest, err.Error())
        return
    }

    common.SendResponse(c, http.StatusOK, summary, "Budget summary retrieved successfully")
}