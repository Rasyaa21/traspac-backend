package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PeriodReport struct {
    ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID          uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
    PeriodStart       time.Time `json:"period_start" gorm:"type:date;not null"`
    PeriodEnd        time.Time `json:"period_end" gorm:"type:date;not null"`
    PdfReport       *string `json:"pdf_report" gorm:"type:string"`
    ReportData datatypes.JSON `json:"report_data" gorm:"type:jsonb;not null"`
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