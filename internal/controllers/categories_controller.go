package controllers

import (
	"gin-backend-app/internal/dto/common"
	"gin-backend-app/internal/dto/request"
	"gin-backend-app/internal/models"
	"gin-backend-app/internal/services"
	"gin-backend-app/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryController struct {
	CategoryService *services.CategoryService
}

func NewCategoryController(categoryService *services.CategoryService) *CategoryController {
	return &CategoryController{
		CategoryService: categoryService,
	}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category for authenticated user
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body request.CreateCategoryRequest true "Category creation data"
// @Success 201 {object} common.Response "Category created successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request data"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Security BearerAuth
// @Router /categories [post]
func (cc *CategoryController) CreateCategory(c *gin.Context) {
	var req request.CreateCategoryRequest

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	category, err := cc.CategoryService.CreateCategory(
		userID,
		req.Name,
		req.CategoryType,
		req.Description,
	)
	if err != nil {
		if err.Error() == "category with this name already exists" {
			common.SendError(c, http.StatusConflict, err.Error())
			return
		}
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusCreated, gin.H{
		"category": gin.H{
			"id":            category.ID,
			"user_id":       category.UserID,
			"name":          category.Name,
			"description":   category.Description,
			"is_default":    category.IsDefault,
			"created_at":    category.CreatedAt,
		},
	}, "Category created successfully")
}

// GetCategories godoc
// @Summary Get all categories
// @Description Get all categories for authenticated user
// @Tags Categories
// @Accept json
// @Produce json
// @Param type query string false "Filter by category type (needs, wants, savings)"
// @Success 200 {object} common.Response "Categories retrieved successfully"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Security BearerAuth
// @Router /categories [get]
func (cc *CategoryController) GetCategories(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	typeQuery := c.Query("type")

	var categories []*models.Category

	if typeQuery != "" {
		categoryType := models.CategoryType(typeQuery)
		categories, err = cc.CategoryService.GetUserCategoriesByType(userID, categoryType)
	} else {
		categories, err = cc.CategoryService.GetUserCategories(userID)
	}

	if err != nil {
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Format response
	categoryData := make([]gin.H, len(categories))
	for i, category := range categories {
		categoryData[i] = gin.H{
			"id":            category.ID,
			"user_id":       category.UserID,
			"name":          category.Name,
			"description":   category.Description,
			"is_default":    category.IsDefault,
			"created_at":    category.CreatedAt,
			"updated_at":    category.UpdatedAt,
		}
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"categories": categoryData,
		"total":      len(categories),
		"filter": gin.H{
			"type": typeQuery,
		},
	}, "Categories retrieved successfully")
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get category details by ID for authenticated user
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} common.Response "Category retrieved successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid category ID"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "Category not found"
// @Security BearerAuth
// @Router /categories/{id} [get]
func (cc *CategoryController) GetCategoryByID(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := cc.CategoryService.GetCategoryByID(categoryID, userID)
	if err != nil {
		if err.Error() == "category not found" || err.Error() == "access denied: category not owned by user" {
			common.SendError(c, http.StatusNotFound, "Category not found")
			return
		}
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"category": gin.H{
			"id":            category.ID,
			"user_id":       category.UserID,
			"name":          category.Name,
			"description":   category.Description,
			"is_default":    category.IsDefault,
			"created_at":    category.CreatedAt,
			"updated_at":    category.UpdatedAt,
		},
	}, "Category retrieved successfully")
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update category for authenticated user
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body request.UpdateCategoryRequest true "Category update data"
// @Success 200 {object} common.Response "Category updated successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request data"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "Category not found"
// @Security BearerAuth
// @Router /categories/{id} [put]
func (cc *CategoryController) UpdateCategory(c *gin.Context) {
	var req request.UpdateCategoryRequest

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	category, err := cc.CategoryService.UpdateCategory(
		categoryID,
		userID,
		req.Name,
		req.CategoryType,
		req.Description,
	)
	if err != nil {
		if err.Error() == "category not found" || err.Error() == "access denied: category not owned by user" {
			common.SendError(c, http.StatusNotFound, "Category not found")
			return
		}
		if err.Error() == "category with this name already exists" {
			common.SendError(c, http.StatusConflict, err.Error())
			return
		}
		if err.Error() == "default categories cannot be updated" {
			common.SendError(c, http.StatusForbidden, err.Error())
			return
		}
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"category": gin.H{
			"id":            category.ID,
			"user_id":       category.UserID,
			"name":          category.Name,
			"description":   category.Description,
			"is_default":    category.IsDefault,
			"created_at":    category.CreatedAt,
			"updated_at":    category.UpdatedAt,
		},
	}, "Category updated successfully")
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete category for authenticated user
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} common.Response "Category deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid category ID"
// @Failure 401 {object} common.ErrorResponse "Authentication required"
// @Failure 404 {object} common.ErrorResponse "Category not found"
// @Security BearerAuth
// @Router /categories/{id} [delete]
func (cc *CategoryController) DeleteCategory(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		common.SendError(c, http.StatusUnauthorized, "Authentication required")
		return
	}

	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		common.SendError(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	err = cc.CategoryService.DeleteCategory(categoryID, userID)
	if err != nil {
		if err.Error() == "category not found" || err.Error() == "access denied: category not owned by user" {
			common.SendError(c, http.StatusNotFound, "Category not found")
			return
		}
		if err.Error() == "default categories cannot be deleted" {
			common.SendError(c, http.StatusForbidden, err.Error())
			return
		}
		common.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	common.SendResponse(c, http.StatusOK, gin.H{
		"category_id": categoryID,
		"deleted_at":  "now",
		"status":      "deleted",
	}, "Category deleted successfully")
}