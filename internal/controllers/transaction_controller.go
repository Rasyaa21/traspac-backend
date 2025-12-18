package controllers

import (
	"gin-backend-app/internal/dto/common"
	"gin-backend-app/internal/dto/request"
	"gin-backend-app/internal/models"
	"gin-backend-app/internal/services"
	"gin-backend-app/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionController struct {
	TransactionService *services.TransactionService
}

func NewTransactionController(transactionService *services.TransactionService) *TransactionController {
	return &TransactionController{
		TransactionService: transactionService,
	}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction for the authenticated user with optional photo upload
// @Tags Transaction Management
// @Accept multipart/form-data
// @Produce json
// @Param type formData string true "Transaction type" Enums(income, expense)
// @Param amount formData integer true "Transaction amount (must be greater than 0)"
// @Param description formData string false "Transaction description"
// @Param date formData string true "Transaction date (YYYY-MM-DDTHH:MM:SSZ format)"
// @Param budget_category formData string false "Budget category" Enums(needs, wants, savings)
// @Param photo formData file false "Transaction photo/receipt (max 5MB, jpg/png/gif/webp)"
// @Success 201 {object} common.Response "Transaction created successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request data"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "Category not found"
// @Security BearerAuth
// @Router /transactions [post]
func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	var req request.TransactionRequest

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Use ShouldBind for multipart form data
	if err := c.ShouldBind(&req); err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// Validate photo if uploaded
	if req.Photo != nil {
		if !utils.IsValidImageFile(req.Photo) {
			common.SendError(c, http.StatusBadRequest, "Invalid photo format. Only JPEG, PNG, GIF, and WebP are allowed")
			return
		}

		// Check file size (max 5MB)
		if req.Photo.Size > 5*1024*1024 {
			common.SendError(c, http.StatusBadRequest, "Photo size must be less than 5MB")
			return
		}
	}

	// Create photo path: user/transaction/:userId
	photoPath := "user/transaction/" + userID.String()

	transaction, err := tc.TransactionService.CreateTransaction(userID, &req, photoPath)
	if err != nil {
		if err.Error() == "category not found" || err.Error() == "category does not belong to user" {
			common.SendError(c, http.StatusNotFound, err.Error())
			return
		}
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusCreated, gin.H{
		"transaction": transaction,
	}, "Transaction created successfully")
}

// GetTransactionByID godoc
// @Summary Get transaction by ID
// @Description Get a specific transaction by ID for the authenticated user
// @Tags Transaction Management
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID (UUID format)"
// @Success 200 {object} common.Response "Transaction retrieved successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid transaction ID"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "Transaction not found"
// @Security BearerAuth
// @Router /transactions/{id} [get]
func (tc *TransactionController) GetTransactionByID(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid transaction ID format")
		return
	}

	transaction, err := tc.TransactionService.GetTransactionByID(userID, transactionID)
	if err != nil {
		if err.Error() == "transaction not found" {
			common.SendError(c, http.StatusNotFound, err.Error())
			return
		}
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"transaction": transaction,
	}, "Transaction retrieved successfully")
}

// UpdateTransaction godoc
// @Summary Update an existing transaction
// @Description Update transaction details for the authenticated user with optional photo upload
// @Tags Transaction Management
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Transaction ID (UUID format)"
// @Param type formData string false "Transaction type" Enums(income, expense)
// @Param amount formData integer false "Transaction amount (must be greater than 0)"
// @Param description formData string false "Transaction description"
// @Param date formData string false "Transaction date (YYYY-MM-DDTHH:MM:SSZ format)"
// @Param budget_category formData string false "Budget category" Enums(needs, wants, savings)
// @Param photo formData file false "Transaction photo/receipt (max 5MB, jpg/png/gif/webp)"
// @Success 200 {object} common.Response "Transaction updated successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request data"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "Transaction not found"
// @Security BearerAuth
// @Router /transactions/{id} [put]
func (tc *TransactionController) UpdateTransaction(c *gin.Context) {
	var req request.TransactionUpdateRequest

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid transaction ID format")
		return
	}

	// Use ShouldBind for multipart form data
	if err := c.ShouldBind(&req); err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// Validate photo if uploaded
	if req.Photo != nil {
		if !utils.IsValidImageFile(req.Photo) {
			common.SendError(c, http.StatusBadRequest, "Invalid photo format. Only JPEG, PNG, GIF, and WebP are allowed")
			return
		}

		// Check file size (max 5MB)
		if req.Photo.Size > 5*1024*1024 {
			common.SendError(c, http.StatusBadRequest, "Photo size must be less than 5MB")
			return
		}
	}

	// Create photo path: user/transaction/:userId
	photoPath := "user/transaction/" + userID.String()

	transaction, err := tc.TransactionService.UpdateTransaction(userID, transactionID, &req, photoPath)
	if err != nil {
		if err.Error() == "category not found" || err.Error() == "category does not belong to user" || err.Error() == "transaction not found" {
			common.SendError(c, http.StatusNotFound, err.Error())
			return
		}
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"transaction": transaction,
	}, "Transaction updated successfully")
}

// DeleteTransaction godoc
// @Summary Delete a transaction
// @Description Delete a specific transaction for the authenticated user including its photo
// @Tags Transaction Management
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID (UUID format)"
// @Success 200 {object} common.Response "Transaction deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid transaction ID"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "Transaction not found"
// @Security BearerAuth
// @Router /transactions/{id} [delete]
func (tc *TransactionController) DeleteTransaction(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid transaction ID format")
		return
	}

	err = tc.TransactionService.DeleteTransaction(userID, transactionID)
	if err != nil {
		if err.Error() == "transaction not found" {
			common.SendError(c, http.StatusNotFound, err.Error())
			return
		}
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"transaction_id": transactionID,
		"deleted_at":     time.Now(),
		"status":         "deleted",
	}, "Transaction deleted successfully")
}

// GetAllTransactions godoc
// @Summary Get all user transactions
// @Description Retrieve all transactions for the authenticated user including photos
// @Tags Transaction Management
// @Accept json
// @Produce json
// @Success 200 {object} common.Response "Transactions retrieved successfully"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Security BearerAuth
// @Router /transactions [get]
func (tc *TransactionController) GetAllTransactions(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	transactions, err := tc.TransactionService.GetAllUserTransactions(userID)
	if err != nil {
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"transactions": transactions,
		"count":        len(transactions),
	}, "Transactions retrieved successfully")
}

// GetTransactionsByPeriod godoc
// @Summary Get transactions by period
// @Description Retrieve transactions grouped by specified period (daily, weekly, monthly)
// @Tags Transaction Management
// @Accept json
// @Produce json
// @Param period_type query string true "Period type" Enums(daily, weekly, monthly)
// @Success 200 {object} common.Response "Transactions by period retrieved successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Security BearerAuth
// @Router /transactions/period [get]
func (tc *TransactionController) GetTransactionsByPeriod(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	periodTypeStr := c.Query("period_type")

	if periodTypeStr == "" {
		common.SendError(c, http.StatusBadRequest, "period_type are required")
		return
	}

	periodType := models.PeriodType(periodTypeStr)
	if !isValidPeriodType(periodType) {
		common.SendError(c, http.StatusBadRequest, "Invalid period type. Must be one of: daily, weekly, monthly")
		return
	}

	result, err := tc.TransactionService.GetTransactionsByPeriod(userID, periodType)
	if err != nil {
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, result, "Transactions by period retrieved successfully")
}

func isValidPeriodType(periodType models.PeriodType) bool {
	validTypes := []models.PeriodType{
		models.PeriodDaily,
		models.PeriodWeekly,
		models.PeriodMonthly,
	}

	for _, validType := range validTypes {
		if periodType == validType {
			return true
		}
	}
	return false
}