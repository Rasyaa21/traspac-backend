package response

import "github.com/google/uuid"

// BudgetResponse represents user's weekly budget
type BudgetResponse struct {
	ID            uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID        uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	WeeklyIncome  int64     `json:"weekly_income" example:"1000000"`
	NeedsBudget   int64     `json:"needs_budget" example:"500000"`
	WantsBudget   int64     `json:"wants_budget" example:"300000"`
	SavingsBudget int64     `json:"savings_budget" example:"200000"`
}

// AllocationResponse represents budget percentage allocation
type AllocationResponse struct {
	NeedsPercent   int `json:"needs_percent" example:"50"`
	WantsPercent   int `json:"wants_percent" example:"30"`
	SavingsPercent int `json:"savings_percent" example:"20"`
}

// BreakdownResponse represents actual budget amount per category
type BreakdownResponse struct {
	NeedsAmount   int64 `json:"needs_amount" example:"500000"`
	WantsAmount   int64 `json:"wants_amount" example:"300000"`
	SavingsAmount int64 `json:"savings_amount" example:"200000"`
}

// UserBudgetResponse represents full budget response
type UserBudgetResponse struct {
	Budget     BudgetResponse     `json:"budget"`
	Allocation AllocationResponse `json:"allocation"`
	Breakdown  BreakdownResponse  `json:"breakdown"`
	IsCustom   bool               `json:"is_custom" example:"false"`
}