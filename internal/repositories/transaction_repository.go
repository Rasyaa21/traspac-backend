package repositories

import (
	"gin-backend-app/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	return r.DB.Create(transaction).Error
}

// func (r *TransactionRepository) 