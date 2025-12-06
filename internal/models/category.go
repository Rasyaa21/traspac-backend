package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
    Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	GroupType 		TransactionGroupType `gorm:"column:group_type;type:transaction_group_enum;not null"`
    Color       string    `json:"color" gorm:"type:varchar(7);default:'#000000'"`
    Icon        *string   `json:"icon" gorm:"type:varchar(50)"`
    Description *string   `json:"description" gorm:"type:text"`
    IsDefault   bool      `json:"is_default" gorm:"default:false"`
    CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

    // Relations
    User User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Category) TableName() string {
    return "categories"
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
    if c.ID == uuid.Nil {
        c.ID = uuid.New()
    }
    return nil
}