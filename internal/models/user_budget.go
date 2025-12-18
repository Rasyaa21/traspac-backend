package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserBudget struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`

	IncomeWeekly int64 `json:"income_weekly" gorm:"type:bigint;not null;default:0"`

	NeedsBudget int64 `json:"needs_budget" gorm:"type:bigint;not null;default:0"`
	WantsBudget int64 `json:"wants_budget" gorm:"type:bigint;not null;default:0"`

	NeedsUsed   int64 `json:"needs_used" gorm:"type:bigint;not null;default:0"`
	WantsUsed   int64 `json:"wants_used" gorm:"type:bigint;not null;default:0"`
	SavingsUsed int64 `json:"savings_used" gorm:"type:bigint;not null;default:0"`

	SavedMoney int64 `json:"saved_money" gorm:"type:bigint;not null;default:0"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (UserBudget) TableName() string { return "user_budgets" }

func (ub *UserBudget) BeforeCreate(tx *gorm.DB) error {
	if ub.ID == uuid.Nil {
		ub.ID = uuid.New()
	}
	return nil
}

// ---------- Validation ----------

func (ub *UserBudget) ValidateNonNegative() error {
	if ub.IncomeWeekly < 0 || ub.NeedsBudget < 0 || ub.WantsBudget < 0 || ub.SavedMoney < 0 ||
		ub.NeedsUsed < 0 || ub.WantsUsed < 0 || ub.SavingsUsed < 0 {
		return errors.New("budget fields must not be negative")
	}
	return nil
}

func (ub *UserBudget) ValidateBudgetDistribution() error {
	// budget allocation for the week must not exceed weekly income
	total := ub.NeedsBudget + ub.WantsBudget + ub.SavedMoney
	if total > ub.IncomeWeekly {
		return fmt.Errorf("total budget (%d) exceeds weekly income (%d)", total, ub.IncomeWeekly)
	}
	return nil
}

// ---------- Remaining ----------

func (ub *UserBudget) RemainingNeeds() int64 {
	return ub.NeedsBudget - ub.NeedsUsed
}

func (ub *UserBudget) RemainingWants() int64 {
	return ub.WantsBudget - ub.WantsUsed
}

func (ub *UserBudget) RemainingSavings() int64 {
	return ub.SavedMoney - ub.SavingsUsed
}

// ---------- Spend / Refund ----------

func (ub *UserBudget) CanSpendFromCategory(categoryType CategoryType, amount int64) bool {
	if amount <= 0 {
		return false
	}

	switch categoryType {
	case CategoryTypeNeeds:
		return ub.RemainingNeeds() >= amount
	case CategoryTypeWants:
		return ub.RemainingWants() >= amount
	case CategoryTypeSavings:
		return ub.RemainingSavings() >= amount
	default:
		return false
	}
}

func (ub *UserBudget) SpendFromCategory(categoryType CategoryType, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if !ub.CanSpendFromCategory(categoryType, amount) {
		return fmt.Errorf("insufficient budget for %s category", categoryType)
	}

	switch categoryType {
	case CategoryTypeNeeds:
		ub.NeedsUsed += amount
	case CategoryTypeWants:
		ub.WantsUsed += amount
	case CategoryTypeSavings:
		ub.SavingsUsed += amount
	default:
		return errors.New("invalid category")
	}

	return nil
}

func (ub *UserBudget) RefundToCategory(categoryType CategoryType, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	switch categoryType {
	case CategoryTypeNeeds:
		if amount > ub.NeedsUsed {
			amount = ub.NeedsUsed
		}
		ub.NeedsUsed -= amount

	case CategoryTypeWants:
		if amount > ub.WantsUsed {
			amount = ub.WantsUsed
		}
		ub.WantsUsed -= amount

	case CategoryTypeSavings:
		if amount > ub.SavingsUsed {
			amount = ub.SavingsUsed
		}
		ub.SavingsUsed -= amount

	default:
		return errors.New("invalid category")
	}

	return nil
}

// AddIncomeToCategory: ini sebenarnya "re-allocate / top up budget", bukan income murni.
// Tetap gue pertahankan signature-nya, tapi dibuat aman.
func (ub *UserBudget) AddIncomeToCategory(categoryType CategoryType, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	switch categoryType {
	case CategoryTypeNeeds:
		ub.NeedsBudget += amount
	case CategoryTypeWants:
		ub.WantsBudget += amount
	case CategoryTypeSavings:
		ub.SavedMoney += amount
	default:
		return errors.New("invalid category")
	}

	ub.IncomeWeekly += amount
	return nil
}

func (ub *UserBudget) ResetWeeklyUsage() {
	ub.NeedsUsed = 0
	ub.WantsUsed = 0
	ub.SavingsUsed = 0
}