package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PeriodReport struct {
    ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID          uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	PeriodType 		PeriodType `gorm:"type:period_type_enum;not null"`
    PeriodStart       time.Time `json:"period_start" gorm:"type:date;not null"`
    PeriodEnd         time.Time `json:"period_end" gorm:"type:date;not null"`
    ReportData      map[string]interface{} `json:"report_data" gorm:"type:jsonb"`
    TotalIncome     float64   `json:"total_income" gorm:"type:decimal(15,2);default:0"`
    TotalExpense    float64   `json:"total_expense" gorm:"type:decimal(15,2);default:0"`
    GeneratedAt     time.Time `json:"generated_at" gorm:"default:CURRENT_TIMESTAMP"`
    CreatedAt       time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
    UpdatedAt       time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

    // Relations
    User User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (PeriodReport) TableName() string {
    return "period_reports"
}

func (pr *PeriodReport) BeforeCreate(tx *gorm.DB) error {
    if pr.ID == uuid.Nil {
        pr.ID = uuid.New()
    }
    return nil
}