package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	Email     string    `json:"email" gorm:"type:varchar(100);not null;unique"`
	Password  string    `json:"password" gorm:"type:varchar(255);not null"`
	IsEmailVerified bool       `json:"is_email_verified" gorm:"default:false;not null"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" gorm:""`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func (User) TableName() string {
	return "users"
}

// BeforeCreate hook untuk generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}