package request

// BudgetRequest represents request for creating/updating budget (unified)
type BudgetRequest struct {
    WeeklyIncome   float64  `json:"weekly_income" binding:"required,gt=0" example:"1000000"`
    SavingsPercent *float64 `json:"savings_percent,omitempty" binding:"omitempty,gte=0,lte=100" example:"50"`
    WantsPercent   *float64 `json:"wants_percent,omitempty" binding:"omitempty,gte=0,lte=100" example:"30"`
    NeedsPercent   *float64 `json:"needs_percent,omitempty" binding:"omitempty,gte=0,lte=100" example:"20"`
}