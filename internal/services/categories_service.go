package services

import (
	"errors"
	"fmt"
	"gin-backend-app/internal/models"
	"gin-backend-app/internal/repositories"
	"strings"

	"github.com/google/uuid"
)

type CategoryService struct {
	CategoryRepo *repositories.CategoryRepository
	UserRepo     *repositories.UserRepository
}

func NewCategoryService(categoryRepo *repositories.CategoryRepository, userRepo *repositories.UserRepository) *CategoryService {
	return &CategoryService{
		CategoryRepo: categoryRepo,
		UserRepo:     userRepo,
	}
}

// CreateCategory creates a new category for user
func (s *CategoryService) CreateCategory(userID uuid.UUID, name string, categoryType models.CategoryType, description *string) (*models.Category, error) {
	// Validate user exists
	_, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate category name
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("category name is required")
	}

	// Check if category name already exists for this user
	existingCategory, _ := s.CategoryRepo.GetCategoryByName(userID, strings.TrimSpace(name))
	if existingCategory != nil {
		return nil, errors.New("category with this name already exists")
	}

	// Validate category type
	if !s.isValidCategoryType(categoryType) {
		return nil, errors.New("invalid category type")
	}

	category := &models.Category{
		UserID:       userID,
		Name:         strings.TrimSpace(name),
		CategoryType: categoryType,
		Description:  description,
		IsDefault:    false,
	}

	if err := s.CategoryRepo.CreateCategory(category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// GetCategoryByID retrieves category by ID with ownership check
func (s *CategoryService) GetCategoryByID(categoryID, userID uuid.UUID) (*models.Category, error) {
	// Validate user exists
	_, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	category, err := s.CategoryRepo.GetCategoryByID(categoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Check ownership
	if category.UserID != userID {
		return nil, errors.New("access denied: category not owned by user")
	}

	return category, nil
}

// GetUserCategories retrieves all categories for a user
func (s *CategoryService) GetUserCategories(userID uuid.UUID) ([]*models.Category, error) {
	// Validate user exists
	_, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	categories, err := s.CategoryRepo.GetCategoriesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve categories: %w", err)
	}

	return categories, nil
}

// GetUserCategoriesByType retrieves categories by type for a user
func (s *CategoryService) GetUserCategoriesByType(userID uuid.UUID, categoryType models.CategoryType) ([]*models.Category, error) {
	// Validate user exists
	_, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate category type
	if !s.isValidCategoryType(categoryType) {
		return nil, errors.New("invalid category type")
	}

	categories, err := s.CategoryRepo.GetCategoriesByUserIDAndType(userID, categoryType)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve categories: %w", err)
	}

	return categories, nil
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(categoryID, userID uuid.UUID, name *string, categoryType *models.CategoryType, description *string) (*models.Category, error) {
	// Validate user exists
	_, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get existing category
	category, err := s.CategoryRepo.GetCategoryByID(categoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Check ownership
	if category.UserID != userID {
		return nil, errors.New("access denied: category not owned by user")
	}

	// Prevent updating default categories
	if category.IsDefault {
		return nil, errors.New("default categories cannot be updated")
	}

	// Update name if provided
	if name != nil {
		trimmedName := strings.TrimSpace(*name)
		if trimmedName == "" {
			return nil, errors.New("category name cannot be empty")
		}

		// Check if new name already exists (exclude current category)
		existingCategory, _ := s.CategoryRepo.GetCategoryByName(userID, trimmedName)
		if existingCategory != nil && existingCategory.ID != categoryID {
			return nil, errors.New("category with this name already exists")
		}

		category.Name = trimmedName
	}

	// Update category type if provided
	if categoryType != nil {
		if !s.isValidCategoryType(*categoryType) {
			return nil, errors.New("invalid category type")
		}
		category.CategoryType = *categoryType
	}

	// Update description if provided
	if description != nil {
		category.Description = description
	}

	if err := s.CategoryRepo.UpdateCategory(category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(categoryID, userID uuid.UUID) error {
	// Validate user exists
	_, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Get category to check if it exists and ownership
	category, err := s.CategoryRepo.GetCategoryByID(categoryID)
	if err != nil {
		return errors.New("category not found")
	}

	// Check ownership
	if category.UserID != userID {
		return errors.New("access denied: category not owned by user")
	}

	// Prevent deleting default categories
	if category.IsDefault {
		return errors.New("default categories cannot be deleted")
	}

	if err := s.CategoryRepo.DeleteCategory(categoryID, userID); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// Helper method to validate category type
func (s *CategoryService) isValidCategoryType(categoryType models.CategoryType) bool {
	validTypes := []models.CategoryType{
		models.CategoryTypeNeeds,
		models.CategoryTypeWants,
		models.CategoryTypeSavings,
	}

	for _, validType := range validTypes {
		if categoryType == validType {
			return true
		}
	}
	return false
}