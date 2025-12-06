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
    result := r.DB.Create(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
    var user models.User
    result := r.DB.First(&user, "id = ?", id)
    
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, result.Error
    }
    
    return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    result := r.DB.Where("email = ?", email).First(&user)
    
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, nil 
        }
        return nil, result.Error
    }
    
    return &user, nil
}

func (r *UserRepository) FindByEmailOrUsername(email, username string) (*models.User, error) {
    var user models.User
    err := r.DB.Where("email = ? OR name = ?", email, username).First(&user).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &user, err
}

func (r *UserRepository) Update(user *models.User) error {
    result := r.DB.Save(user)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
    result := r.DB.Delete(&models.User{}, "id = ?", id)
    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    
    return nil
}

func (r *UserRepository) List(limit, offset int) ([]*models.User, error) {
    var users []*models.User
    result := r.DB.Limit(limit).Offset(offset).Find(&users)
    
    if result.Error != nil {
        return nil, result.Error
    }
    
    return users, nil
}

func (r *UserRepository) Count() (int64, error) {
    var count int64
    result := r.DB.Model(&models.User{}).Count(&count)
    
    if result.Error != nil {
        return 0, result.Error
    }
    
    return count, nil
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
    var count int64
    result := r.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
    
    if result.Error != nil {
        return false, result.Error
    }
    
    return count > 0, nil
}

func (r *UserRepository) MarkEmailVerified(userId uuid.UUID) error {
    now := time.Now()
    result := r.DB.Model(&models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
        "is_email_verified":true,
        "email_verified_at": &now,
    })

    if result.Error != nil {
        return result.Error
    }

    return nil
}