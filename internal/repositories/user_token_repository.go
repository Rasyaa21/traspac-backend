package repositories

import (
	"gin-backend-app/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserTokenRepository struct {
	DB *gorm.DB
}

func NewUserTokenRepository (db *gorm.DB) *UserTokenRepository {
	return &UserTokenRepository{
		DB: db,
	}
}

func (r *UserTokenRepository) Create (token *models.UserToken) error {
	result := r.DB.Create(token)
	if result.Error != nil {
		return result.Error
	} 
	return nil
}

func (r *UserTokenRepository) FindTokenByHash (tokenHash string) (*models.UserToken, error) {
	var token models.UserToken
	err := r.DB.Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *UserTokenRepository) MarkAsUsedToken (tokenID uuid.UUID) error {
	now := time.Now()
	return r.DB.Model(&models.UserToken{}).
		Where("id = ?", tokenID).
		Update("used_at", &now).Error
}

