package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`

	TokenHash string `json:"-" gorm:"type:char(64);uniqueIndex;not null"`

	TokenType TokenType `json:"token_type" gorm:"type:varchar(50);not null;index"`

	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	UsedAt    *time.Time `json:"used_at" gorm:""`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	User User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (UserToken) TableName() string {
	return "user_tokens"
}

func (t *UserToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}