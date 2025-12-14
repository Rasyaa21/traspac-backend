package request

import "gin-backend-app/internal/models"

type CreateCategoryRequest struct {
    Name         string                `json:"name" binding:"required,max=100" example:"Food & Dining"`
    CategoryType models.CategoryType   `json:"category_type" binding:"required" example:"needs"`
    Description  *string              `json:"description,omitempty" binding:"omitempty,max=500" example:"Daily meals and dining expenses"`
}

type UpdateCategoryRequest struct {
    Name         *string               `json:"name,omitempty" binding:"omitempty,max=100" example:"Food & Dining"`
    CategoryType *models.CategoryType  `json:"category_type,omitempty" example:"needs"`
    Description  *string              `json:"description,omitempty" binding:"omitempty,max=500" example:"Daily meals and dining expenses"`
}