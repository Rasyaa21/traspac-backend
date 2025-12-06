package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AILog struct {
    ID            uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID        uuid.UUID              `json:"user_id" gorm:"type:uuid;not null;index"`
    TransactionID *uuid.UUID             `json:"transaction_id" gorm:"type:uuid;index"`
    CategoryID    *uuid.UUID             `json:"category_id" gorm:"type:uuid;index"`
	AnalysisType  AiAnalysisType `gorm:"type:ai_analysis_type_enum;not null"`
    InputData     map[string]interface{} `json:"input_data" gorm:"type:jsonb"`
    OutputData    map[string]interface{} `json:"output_data" gorm:"type:jsonb"`
    CreatedAt     time.Time              `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`

    // Relations
    User        User         `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
    Transaction *Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID;references:ID;constraint:OnDelete:SET NULL"`
    Category    *Category    `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:SET NULL"`
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