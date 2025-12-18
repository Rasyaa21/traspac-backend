package request

import (
	"gin-backend-app/internal/models"
	"mime/multipart"
	"time"
)

type TransactionRequest struct {
	Type           models.TransactionGroupType     `form:"type" binding:"required,oneof=income expense"`
	Amount         int64                           `form:"amount" binding:"required,gt=0"`
	Description    *string                         `form:"description"`
	Date           time.Time                       `form:"date" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	BudgetCategory *models.CategoryType            `form:"budget_category" binding:"omitempty,oneof=needs wants savings"`
	Photo          *multipart.FileHeader           `form:"photo"` // File upload from form-data
}

type TransactionUpdateRequest struct {
	Type           *models.TransactionGroupType    `form:"type" binding:"omitempty,oneof=income expense"`
	Amount         *int64                          `form:"amount" binding:"omitempty,gt=0"`
	Description    *string                         `form:"description"`
	Date           *time.Time                      `form:"date" time_format:"2006-01-02T15:04:05Z07:00"`
	BudgetCategory *models.CategoryType            `form:"budget_category" binding:"omitempty,oneof=needs wants savings"`
	Photo          *multipart.FileHeader           `form:"photo"`    // File upload from form-data
	PhotoURL       *string                         `json:"photo_url"` // Internal use for storing URL
}

