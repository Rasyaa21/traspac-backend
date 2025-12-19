package response

import (
	"gin-backend-app/internal/models"
	"time"

	"github.com/google/uuid"
)

type EStatement struct {
	GeneratedAt  time.Time       `json:"generated_at"` 
	Period       PeriodWindow    `json:"period"`
	Balances     BalanceSection  `json:"balances"`
	Transactions []TransactionLine `json:"transactions"`
}

type PeriodWindow struct {
	Start time.Time `json:"start"` 
	End   time.Time `json:"end"`
}

type BalanceSection struct {
	StartingBalance int64 `json:"starting_balance"`
	EndingBalance   int64 `json:"ending_balance"`
	Categories      map[models.CategoryType]BucketBalance `json:"categories"`
	TotalBuckets    int64 `json:"total_buckets"`
}

type BucketBalance struct {
	Starting int64 `json:"starting"`
	In       int64 `json:"in"`  
	Out      int64 `json:"out"` 
	Ending   int64 `json:"ending"`
}

type TransactionLine struct {
	ID          uuid.UUID                   `json:"id"`
	Date        time.Time                   `json:"date"`       
	Description string                      `json:"description"` 
	Type        models.TransactionGroupType `json:"type"`        
	Bucket      models.CategoryType         `json:"bucket"`    
	Amount      int64                       `json:"amount"` 
}

// PeriodReportResponse represents the period report response without user data
type PeriodReportResponse struct {
	ID          uuid.UUID   `json:"id"`
	UserID      uuid.UUID   `json:"user_id"`
	PeriodStart time.Time   `json:"period_start"`
	PeriodEnd   time.Time   `json:"period_end"`
	PDFReport   *string     `json:"pdf_report"`
	ReportData  EStatement  `json:"report_data"`
	GeneratedAt time.Time   `json:"generated_at"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// PeriodReportsListResponse represents list of period reports
type PeriodReportsListResponse struct {
	Reports []PeriodReportResponse `json:"reports"`
}

