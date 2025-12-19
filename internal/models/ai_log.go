package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AILog struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	PeriodID  *uuid.UUID     `json:"period_id,omitempty" gorm:"type:uuid;index"`

	OutputData map[string]interface{} `json:"output_data" gorm:"type:jsonb;not null"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	User         User          `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	PeriodReport *PeriodReport `json:"period_report,omitempty" gorm:"foreignKey:PeriodID;references:ID;constraint:OnDelete:SET NULL"`
}

func (AILog) TableName() string {
	return "ai_logs"
}

func (al *AILog) BeforeCreate(tx *gorm.DB) error {
	if al.ID == uuid.Nil {
		al.ID = uuid.New()
	}
	return nil
}