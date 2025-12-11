package repositories

import (
	"errors"
	"gin-backend-app/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
    DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
    err := r.DB.Model(&models.User{}).Create(user).Error
    if err != nil {
        return err
    }
    return nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
    var user models.User
    err := r.DB.Model(&models.User{}).First(&user, "id = ?", id).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }

    return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.DB.Model(&models.User{}).Where("email = ?", email).First(&user).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }

    return &user, nil
}

func (r *UserRepository) FindByEmailOrUsername(email, username string) (*models.User, error) {
    var user models.User
    err := r.DB.Model(&models.User{}).Where("email = ? OR name = ?", email, username).First(&user).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &user, err
}

func (r *UserRepository) Update(user *models.User) error {
    err := r.DB.Save(user).Error
    if err != nil {
        return err
    }
    return nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
    tx := r.DB.Delete(&models.User{}, "id = ?", id)
    err := tx.Error
    if err != nil {
        return err
    }

    if tx.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }

    return nil
}

func (r *UserRepository) List(limit, offset int) ([]*models.User, error) {
    var users []*models.User
    err := r.DB.Limit(limit).Offset(offset).Find(&users).Error
    if err != nil {
        return nil, err
    }

    return users, nil
}

func (r *UserRepository) Count() (int64, error) {
    var count int64
    err := r.DB.Model(&models.User{}).Count(&count).Error
    if err != nil {
        return 0, err
    }

    return count, nil
}

func (r *UserRepository) MarkEmailVerified(userId uuid.UUID) error {
    now := time.Now()
    err := r.DB.Model(&models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
        "is_email_verified": true,
        "email_verified_at": &now,
    }).Error
    if err != nil {
        return err
    }
    return nil
}

func (r *UserRepository) ChangePassword (hashedPassword string, userId uuid.UUID) error {
    err := r.DB.Model(&models.User{}).Where("id = ?", userId).Update("password", hashedPassword).Error
    if err != nil {
        return err
    } 
    return nil
}