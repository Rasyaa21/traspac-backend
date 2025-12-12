package controllers

import (
	"gin-backend-app/internal/dto/common"
	"gin-backend-app/internal/dto/request"
	"gin-backend-app/internal/services"
	"gin-backend-app/pkg/utils"
	"net/http"
	"strconv"

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
	// If no custom percentages provided, allocation will be nil and service will use default

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
	savingsPercent := (budget.SavingsBudget / budget.IncomeWeekly) * 100
	wantsPercent := (budget.WantsBudget / budget.IncomeWeekly) * 100
	needsPercent := (budget.NeedsBudget / budget.IncomeWeekly) * 100

	responseData := gin.H{
		"budget": gin.H{
			"id":             budget.ID,
			"user_id":        budget.UserID,
			"weekly_income":  budget.IncomeWeekly,
			"needs_budget":   budget.NeedsBudget,
			"wants_budget":   budget.WantsBudget,
			"savings_budget": budget.SavingsBudget,
			"remaining":      budget.CalculateRemainingIncome(),
		},
		"allocation": gin.H{
			"savings_percent": savingsPercent,
			"wants_percent":   wantsPercent,
			"needs_percent":   needsPercent,
		},
		"breakdown": gin.H{
			"savings_amount": budget.SavingsBudget,
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

	// Determine allocation (default or custom)
	var allocation *services.BudgetAllocation
	isCustom := false

	if req.SavingsPercent != nil && req.WantsPercent != nil && req.NeedsPercent != nil {
		// Custom allocation provided
		customAlloc := services.NewCustomAllocation(*req.SavingsPercent, *req.WantsPercent, *req.NeedsPercent)
		allocation = &customAlloc
		isCustom = true
	}
	// If no custom percentages provided, allocation will be nil and service will use default

	budget, err := bc.UserBudgetService.UpdateUserBudget(userID, req.WeeklyIncome, allocation)
	if err != nil {
		if err.Error() == "user budget not found" {
			common.SendError(c, http.StatusNotFound, "Budget not found. Use POST to create a new budget")
			return
		}
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Calculate percentages for response
	savingsPercent := (budget.SavingsBudget / budget.IncomeWeekly) * 100
	wantsPercent := (budget.WantsBudget / budget.IncomeWeekly) * 100
	needsPercent := (budget.NeedsBudget / budget.IncomeWeekly) * 100

	responseData := gin.H{
		"budget": gin.H{
			"id":             budget.ID,
			"user_id":        budget.UserID,
			"weekly_income":  budget.IncomeWeekly,
			"needs_budget":   budget.NeedsBudget,
			"wants_budget":   budget.WantsBudget,
			"savings_budget": budget.SavingsBudget,
			"remaining":      budget.CalculateRemainingIncome(),
		},
		"allocation": gin.H{
			"savings_percent": savingsPercent,
			"wants_percent":   wantsPercent,
			"needs_percent":   needsPercent,
		},
		"breakdown": gin.H{
			"savings_amount": budget.SavingsBudget,
			"wants_amount":   budget.WantsBudget,
			"needs_amount":   budget.NeedsBudget,
		},
		"is_custom": isCustom,
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

	// Calculate percentages for display
	savingsPercent := (budget.SavingsBudget / budget.IncomeWeekly) * 100
	wantsPercent := (budget.WantsBudget / budget.IncomeWeekly) * 100
	needsPercent := (budget.NeedsBudget / budget.IncomeWeekly) * 100

	// Determine if it's using default allocation
	isDefault := (savingsPercent == 50.0 && wantsPercent == 30.0 && needsPercent == 20.0)

	common.SendResponse(c, http.StatusOK, gin.H{
		"budget": gin.H{
			"id":             budget.ID,
			"user_id":        budget.UserID,
			"weekly_income":  budget.IncomeWeekly,
			"needs_budget":   budget.NeedsBudget,
			"wants_budget":   budget.WantsBudget,
			"savings_budget": budget.SavingsBudget,
			"remaining":      budget.CalculateRemainingIncome(),
			"created_at":     budget.CreatedAt,
			"updated_at":     budget.UpdatedAt,
		},
		"allocation": gin.H{
			"savings_percent": savingsPercent,
			"wants_percent":   wantsPercent,
			"needs_percent":   needsPercent,
		},
		"allocation_type": gin.H{
			"is_default": isDefault,
			"is_custom":  !isDefault,
		},
	}, "Budget details retrieved successfully")
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
		"user_id":        userID,
		"deleted_at":     "now",
		"status":         "deleted",
	}, "Budget deleted successfully")
}


// CalculateBudget godoc
// @Summary Calculate budget preview
// @Description Calculate budget allocation preview without saving to database
// @Tags Budget Management
// @Accept json
// @Produce json
// @Param income query float64 true "Weekly income amount"
// @Param savings_percent query float64 false "Custom savings percentage (default: 50)"
// @Param wants_percent query float64 false "Custom wants percentage (default: 30)"
// @Param needs_percent query float64 false "Custom needs percentage (default: 20)"
// @Success 200 {object} common.Response "Budget calculation completed successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid parameters"
// @Router /budget/calculate [get]
func (bc *UserBudgetController) CalculateBudget(c *gin.Context) {
    incomeStr := c.Query("income")
    if incomeStr == "" {
        common.SendError(c, http.StatusBadRequest, "Income parameter is required")
        return
    }

    income, err := strconv.ParseFloat(incomeStr, 64)
    if err != nil || income <= 0 {
        common.SendError(c, http.StatusBadRequest, "Invalid income amount")
        return
    }

    // Check for custom percentages
    var allocation *services.BudgetAllocation

    savingsStr := c.Query("savings_percent")
    wantsStr := c.Query("wants_percent")
    needsStr := c.Query("needs_percent")

    if savingsStr != "" && wantsStr != "" && needsStr != "" {
        savings, err1 := strconv.ParseFloat(savingsStr, 64)
        wants, err2 := strconv.ParseFloat(wantsStr, 64)
        needs, err3 := strconv.ParseFloat(needsStr, 64)

        if err1 != nil || err2 != nil || err3 != nil {
            common.SendError(c, http.StatusBadRequest, "Invalid percentage values")
            return
        }

        customAlloc := services.NewCustomAllocation(savings, wants, needs)
        allocation = &customAlloc
    }
    // If no custom percentages provided, allocation will be nil and service will use default

    preview, err := bc.UserBudgetService.CalculateBudgetPreview(income, allocation)
    if err != nil {
        common.SendError(c, http.StatusBadRequest, err.Error())
        return
    }

    common.SendResponse(c, http.StatusOK, preview, "Budget calculation completed successfully")
}