package repositories

import (
	"errors"
	"gin-backend-app/internal/models"
	"gin-backend-app/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserTokenRepository struct {
    DB *gorm.DB
}

func NewUserTokenRepository(db *gorm.DB) *UserTokenRepository {
    return &UserTokenRepository{
        DB: db,
    }
}

func (r *UserTokenRepository) Create(token *models.UserToken) error {
    err := r.DB.Model(&models.UserToken{}).Create(token).Error
    if err != nil {
        return err
    }
    return nil
}

func (r *UserTokenRepository) FindTokenByOTP(otp string, tokenType models.TokenType, userId uuid.UUID) (*models.UserToken, error) {
    var token models.UserToken
    err := r.DB.Model(&models.UserToken{}).Where("token_otp = ? AND token_type = ? AND user_id = ?", otp, tokenType, userId).First(&token).Error
    if err != nil {
        return nil, err
    }
    return &token, nil
}

func (r *UserTokenRepository) FindTokenByUserId(userId uuid.UUID) (*models.UserToken, error) {
    var token models.UserToken
    err := r.DB.Model(&models.UserToken{}).Where("user_id = ?", userId).First(&token).Error
    if err != nil {
        return nil, err
    }
    return &token, nil
}

func (r *UserTokenRepository) MarkAsUsedToken(tokenID uuid.UUID) error {
    now := time.Now()
    err := r.DB.Model(&models.UserToken{}).
        Where("id = ?", tokenID).
        Update("used_at", &now).Error
    return err
}

func (r *UserTokenRepository) CheckUserTokenCreatedAt(userId uuid.UUID) bool {
    oneMinuteAgo := time.Now().Add(-1 * time.Minute)
    var token models.UserToken
    err := r.DB.Model(&models.UserToken{}).
        Where("user_id = ? AND token_type = ? AND created_at >= ?", userId, models.TokenTypeEmailVerification, oneMinuteAgo).
        Order("created_at DESC").
        First(&token).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return true
        }
        return false
    }
    return false
}

func (r *UserTokenRepository) DeleteExpired(now time.Time) (int64, error) {
    // Need to get RowsAffected before getting Error
    tx := r.DB.Model(&models.UserToken{}).
        Where("expires_at < ? AND used_at IS NULL", now).
        Delete(&models.UserToken{})
    
    err := tx.Error
    if err != nil {
        return 0, err
    }
    return tx.RowsAffected, nil
}

func (r *UserTokenRepository) GenerateAndSetVerificationTokenByOTP(otp, email string, tokenType models.TokenType, userId uuid.UUID) (string, error) {
    verifyToken, err := utils.GenerateRandomToken(32)
    if err != nil {
        return "", err
    }

	tx := r.DB.Model(&models.UserToken{}).
        Where("user_id = ?", userId).
        Where("token_otp = ?", otp).
        Where("token_type = ?", tokenType).
        Where("expires_at > ?", time.Now()).
        Where("verify_token IS NULL").
        Update("verify_token", verifyToken)

    err = tx.Error
    if err != nil {
        return "", err
    }

    if tx.RowsAffected == 0 {
        return "", errors.New("invalid or expired otp, or verify token already generated")
    }

    return verifyToken, nil
}

func (r *UserTokenRepository) ValidateTokenAndGetUser (verificationToken string) (*models.User, error) {
    var token models.UserToken

    err := r.DB.
        Where("verify_token = ?", verificationToken).
        Where("expires_at > ?", time.Now()).
        Where("used_at IS NULL").
        First(&token).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil ,errors.New("invalid or expired verification token")
        }
        return &token.User, nil
    }
    return &token.User, nil
}