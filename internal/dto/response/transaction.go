package response

import (
	"gin-backend-app/internal/models"
	"time"

	"github.com/google/uuid"
)

type TransactionSummary struct {
	Income  int64 `json:"income"`
	Expense int64 `json:"expense"`
	Total   int64 `json:"total"` 
}

type PeriodTransactionGroup struct {
	Period       time.Time            `json:"period"`
	Transactions []models.Transaction `json:"transactions"`
	Summary      TransactionSummary   `json:"summary"`
}

type GetAllTransactionByPeriodResponse struct {
	Groups []PeriodTransactionGroup `json:"groups"`
	Total  TransactionSummary       `json:"total"`
}

type CategoryPeriodSummary struct {
	Period     time.Time `json:"period"`
	CategoryID uuid.UUID `json:"category_id"`
	Category   string    `json:"category"`
	Type       string    `json:"type"` 
	Total      int64     `json:"total"`
}

type GetByCategoryResponse struct {
	Items []CategoryPeriodSummary `json:"items"`
	Total TransactionSummary      `json:"total"`
}