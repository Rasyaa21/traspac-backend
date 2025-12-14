package repositories

import (
	"gin-backend-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) CreateCategory(category *models.Category) error {
	return r.DB.Create(category).Error
}

func (r *CategoryRepository) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.DB.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetCategoriesByUserID(userID uuid.UUID) ([]*models.Category, error) {
	var categories []*models.Category
	err := r.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) GetCategoriesByUserIDAndType(userID uuid.UUID, categoryType models.CategoryType) ([]*models.Category, error) {
	var categories []*models.Category
	err := r.DB.Where("user_id = ? AND category_type = ?", userID, categoryType).Order("created_at DESC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) UpdateCategory(category *models.Category) error {
	return r.DB.Save(category).Error
}

func (r *CategoryRepository) DeleteCategory(id uuid.UUID, userID uuid.UUID) error {
	return r.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Category{}).Error
}

func (r *CategoryRepository) GetCategoryByName(userID uuid.UUID, name string) (*models.Category, error) {
	var category models.Category
	err := r.DB.Where("user_id = ? AND name = ?", userID, name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) CheckCategoryOwnership(categoryID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&models.Category{}).Where("id = ? AND user_id = ?", categoryID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
