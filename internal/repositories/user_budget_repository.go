package repositories

import (
	"gin-backend-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserBudgetRepository struct {
	DB *gorm.DB
}

func NewUserBudgetRepository(db *gorm.DB) *UserBudgetRepository {
	return &UserBudgetRepository{DB: db}
}

func (r *UserBudgetRepository) CreateUserBudget(userBudget *models.UserBudget) error {
	err := r.DB.Model(&models.UserBudget{}).Create(userBudget).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserBudgetRepository) EdituserBudget (needs, savings, wants, income float64, userId uuid.UUID) error {
	err := r.DB.Model(&models.UserBudget{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"needs_budget": needs,
        "wants_budget": wants,
        "savings_budget": savings,
		"income_weekly": income,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserBudgetRepository) DeleteUserBudget(userId uuid.UUID) error {
	err := r.DB.Model(&models.UserBudget{}).Where("user_id = ?", userId).Delete(&models.UserBudget{}).Error
	if err != nil {
		return err
	} 
	return nil
}

func (r *UserBudgetRepository) GetUserBudget(userId uuid.UUID) (*models.UserBudget, error) {
	var userBudget models.UserBudget
	err := r.DB.
		Where("user_id = ?", userId).
		First(&userBudget).Error
	if err != nil {
		return nil, err
	}
	return &userBudget, err
}

