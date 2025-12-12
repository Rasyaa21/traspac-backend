package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserBudget struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	IncomeWeekly float64   `json:"income_weekly" gorm:"type:decimal(15,2);not null;default:0"`

	NeedsBudget   float64   `json:"needs_budget" gorm:"type:decimal(15,2);not null;default:0"`
	WantsBudget   float64   `json:"wants_budget" gorm:"type:decimal(15,2);not null;default:0"`
	SavingsBudget float64   `json:"savings_budget" gorm:"type:decimal(15,2);not null;default:0"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	User     User     `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (UserBudget) TableName() string {
	return "user_budgets"
}

func (ub *UserBudget) BeforeCreate(tx *gorm.DB) error {
	if ub.ID == uuid.Nil {
		ub.ID = uuid.New()
	}
	return nil
}

// Validation methods
func (ub *UserBudget) ValidateBudgetDistribution() error {
	total := ub.NeedsBudget + ub.WantsBudget + ub.SavingsBudget
	if total > ub.IncomeWeekly {
		return fmt.Errorf("total budget (%.2f) exceeds weekly income (%.2f)", total, ub.IncomeWeekly)
	}
	return nil
}

func (ub *UserBudget) CalculateRemainingIncome() float64 {
	total := ub.NeedsBudget + ub.WantsBudget + ub.SavingsBudget
	return ub.IncomeWeekly - total
}