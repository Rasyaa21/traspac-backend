package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserBudget struct {
    ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
    CategoryID   uuid.UUID `json:"category_id" gorm:"type:uuid;not null;index"`
    Amount       float64   `json:"amount" gorm:"type:decimal(15,2);not null"`
    PeriodType   string    `json:"period_type" gorm:"type:period_type_enum;not null"`
    PeriodValue  int       `json:"period_value" gorm:"not null"`
    StartDate    time.Time `json:"start_date" gorm:"type:date;not null"`
    EndDate      *time.Time `json:"end_date" gorm:"type:date"`
    IsActive     bool      `json:"is_active" gorm:"default:true"`
    CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
    UpdatedAt    time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

    // Relations
    User     User     `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
    Category Category `json:"category" gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:CASCADE"`
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