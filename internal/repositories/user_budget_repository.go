package repositories

import (
	"fmt"
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

func (r *UserBudgetRepository) EdituserBudget(needsBudget, savingsBudget, wantsBudget, weeklyIncome int64, userId uuid.UUID) error {
	err := r.DB.Model(&models.UserBudget{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"needs_budget": needsBudget,
		"wants_budget": wantsBudget,
		"savings_budget": savingsBudget,
		"income_weekly": weeklyIncome,

		"needs_used" : gorm.Expr("LEAST(needs_used, ?)", needsBudget),
		"wants_used" : gorm.Expr("LEAST(wants_used, ?)", wantsBudget),
		"savings_used" : gorm.Expr("LEAST(savings_used, ?)", savingsBudget),
		"updated_at": gorm.Expr("NOW()"),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserBudgetRepository) Spend(userId uuid.UUID, amount int64, categoryType models.CategoryType) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	var usedCol, budgetCol string
	switch categoryType {
	case models.CategoryTypeNeeds:
		usedCol, budgetCol = "needs_used", "needs_budget"
	case models.CategoryTypeWants:
		usedCol, budgetCol = "wants_used", "wants_budget"
	case models.CategoryTypeSavings:
		usedCol, budgetCol = "savings_used", "saved_money"
	default:
		return fmt.Errorf("invalid categoryType: %s", categoryType)
	}

	tx := r.DB.Model(&models.UserBudget{}).
		Where("user_id = ?", userId).
		Where(usedCol+" + ? <= "+budgetCol, amount).
		Updates(map[string]interface{}{
			usedCol:      gorm.Expr(usedCol+" + ?", amount),
			"updated_at": gorm.Expr("NOW()"),
		})

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("insufficient budget for %s", categoryType)
	}
	return nil
}

func (r *UserBudgetRepository) AddIncomeToBucket(userId uuid.UUID, amount int64, categoryType models.CategoryType) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	var budgetCol string
	switch categoryType {
	case models.CategoryTypeNeeds:
		budgetCol = "needs_budget"
	case models.CategoryTypeWants:
		budgetCol = "wants_budget"
	case models.CategoryTypeSavings:
		budgetCol = "saved_money"
	default:
		return fmt.Errorf("invalid categoryType: %s", categoryType)
	}

	err := r.DB.Model(&models.UserBudget{}).
		Where("user_id = ?", userId).
		Updates(map[string]interface{}{
			"income_weekly": gorm.Expr("income_weekly + ?", amount),
			budgetCol:       gorm.Expr(budgetCol+" + ?", amount),
			"updated_at":    gorm.Expr("NOW()"),
		}).Error

	if err != nil {
		return err
	}
	return nil
}

func (r *UserBudgetRepository) ResetWeeklyBudget() error {
	err := r.DB.Model(&models.UserBudget{}).Updates(map[string]interface{}{
		"saved_money" : gorm.Expr("saved_money + savings_budget"), 
		"needs_used":   int64(0),
		"wants_used":   int64(0),
		"savings_used": int64(0),

		"updated_at": gorm.Expr("NOW()"),
	}).Error
	return err
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

