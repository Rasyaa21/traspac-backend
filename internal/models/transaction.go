package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
    CategoryID  uuid.UUID `json:"category_id" gorm:"type:uuid;not null;index"`
	Type 		TransactionGroupType `json:"type" gorm:"type:transaction_group_enum;not null"`
    Amount      int64   `json:"amount" gorm:"type:decimal(15,2);not null"`
    Description *string   `json:"description" gorm:"type:text"`
    Photo       *string `json:"photo" gorm:"type:text;default:null"`
    Date        time.Time `json:"date" gorm:"type:date;not null;index"`
    
    BudgetCategory *CategoryType `json:"budget_category" gorm:"type:category_type_enum"`
    
    CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

    // Relations
    User     User     `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
    Category Category `json:"category" gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Transaction) TableName() string {
    return "transactions"
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
    if t.ID == uuid.Nil {
        t.ID = uuid.New()
    }
    return nil
}