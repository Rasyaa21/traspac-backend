package services

import (
	"errors"
	"gin-backend-app/internal/dto/request"
	"gin-backend-app/internal/dto/response"
	"gin-backend-app/internal/models"
	"gin-backend-app/internal/repositories"
	"gin-backend-app/pkg/utils"
	"os"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService struct {
	TransactionRepo *repositories.TransactionRepository
	UserRepo        *repositories.UserRepository
	UserBudgetRepo  *repositories.UserBudgetRepository
	R2Client        *utils.R2Client
}

// NewTransactionService - Creates new instance of TransactionService
// Tujuan: Membuat instance baru dari TransactionService dengan dependency injection
// Parameter: transactionRepo (*repositories.TransactionRepository), userRepo (*repositories.UserRepository), userBudgetRepo (*repositories.UserBudgetRepository), r2Client (*utils.R2Client)
// Return: *TransactionService - instance baru dari service
// Penjelasan: Constructor pattern untuk inisialisasi service dengan semua dependencies termasuk R2 client
func NewTransactionService(
	transactionRepo *repositories.TransactionRepository,
	userRepo *repositories.UserRepository,
	userBudgetRepo *repositories.UserBudgetRepository,
	r2Client *utils.R2Client,
) *TransactionService {
	return &TransactionService{
		TransactionRepo: transactionRepo,
		UserRepo:        userRepo,
		UserBudgetRepo:  userBudgetRepo,
		R2Client:        r2Client,
	}
}

// CreateTransaction - Creates a new transaction for user
// Tujuan: Membuat transaksi baru untuk user dengan validasi dan budget update, termasuk photo upload
// Parameter: userID (uuid.UUID), req (*request.TransactionRequest), photoPath (string) - path untuk storage photo
// Return: (*models.Transaction, error) - transaksi yang dibuat atau error jika gagal
// Penjelasan: Validasi user sebelum membuat transaksi baru, upload photo jika ada, kemudian update budget jika expense
func (s *TransactionService) CreateTransaction(userID uuid.UUID, req *request.TransactionRequest, photoPath string) (*models.Transaction, error) {
	// Validate user exists
	_, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate amount
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	transaction := &models.Transaction{
		UserID:         userID,
		Type:           req.Type,
		Amount:         req.Amount,
		Description:    req.Description,
		Date:           req.Date,
		BudgetCategory: req.BudgetCategory,
	}

	if req.Photo != nil && s.R2Client != nil {
		photoURL, err := s.R2Client.UploadFile(req.Photo, photoPath)
		if err != nil {
			// Don't fail transaction creation if photo upload fails
			// You might want to log this error
		} else {
			transaction.Photo = &photoURL
		}
	}

	createdTransaction, err := s.TransactionRepo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	if req.Type == models.TransactionGroupExpense && req.BudgetCategory != nil {
		if err := s.UserBudgetRepo.Spend(userID, req.Amount, *req.BudgetCategory); err != nil {
			return nil, err
		}
	}

	if req.Type == models.TransactionGroupIncome && req.BudgetCategory != nil {
		if err := s.UserBudgetRepo.AddIncomeToBucket(userID, req.Amount, *req.BudgetCategory); err != nil {
			return nil, err
		}
	}

	return createdTransaction, nil
}

// GetTransactionByID - Get transaction by ID with ownership validation
// Tujuan: Mengambil transaksi berdasarkan ID dengan validasi kepemilikan
// Parameter: userID (uuid.UUID), transactionID (uuid.UUID)
// Return: (*models.Transaction, error) - data transaksi atau error jika tidak ditemukan
// Penjelasan: Validasi user ada kemudian ambil transaksi yang dimiliki user
func (s *TransactionService) GetTransactionByID(userID uuid.UUID, transactionID uuid.UUID) (*models.Transaction, error) {
	// Validate user exists
	_, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	transaction, err := s.TransactionRepo.GetTransactionByID(userID, transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	return transaction, nil
}

// UpdateTransaction - Updates existing transaction
// Tujuan: Mengupdate transaksi yang sudah ada dengan validasi kepemilikan, termasuk photo update
// Parameter: userID (uuid.UUID), transactionID (uuid.UUID), req (*request.TransactionUpdateRequest), photoPath (string)
// Return: (*models.Transaction, error) - transaksi yang diupdate atau error jika gagal
// Penjelasan: Validasi user dan transaksi ada, upload photo baru jika ada, kemudian update dengan data baru yang diberikan
func (s *TransactionService) UpdateTransaction(userID uuid.UUID, transactionID uuid.UUID, req *request.TransactionUpdateRequest, photoPath string) (*models.Transaction, error) {
	// Validate user exists
	_, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get existing transaction
	existingTransaction, err := s.TransactionRepo.GetTransactionByID(userID, transactionID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	// Validate amount if provided
	if req.Amount != nil && *req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	if req.Photo != nil && s.R2Client != nil {
		if existingTransaction.Photo != nil {
			oldPhotoKey := s.extractKeyFromURL(*existingTransaction.Photo)
			if oldPhotoKey != "" {
				s.R2Client.DeleteFile(oldPhotoKey) 
			}
		}

		// Upload new photo
		photoURL, err := s.R2Client.UploadFile(req.Photo, photoPath)
		if err == nil {
			req.PhotoURL = &photoURL
		}
	}

	return s.TransactionRepo.EditTransaction(userID, transactionID, req)
}

// DeleteTransaction - Deletes a transaction and its photo
// Tujuan: Menghapus transaksi berdasarkan ID dengan validasi kepemilikan, termasuk menghapus photo
// Parameter: userID (uuid.UUID), transactionID (uuid.UUID)
// Return: error - error jika gagal menghapus atau nil jika berhasil
// Penjelasan: Validasi user ada, ambil data transaksi untuk mendapat photo URL, hapus photo dari storage, kemudian hapus transaksi
func (s *TransactionService) DeleteTransaction(userID uuid.UUID, transactionID uuid.UUID) error {
	// Validate user exists
	_, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Get transaction to check if it has photo
	transaction, err := s.TransactionRepo.GetTransactionByID(userID, transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Delete photo if exists
	if transaction.Photo != nil && s.R2Client != nil {
		photoKey := s.extractKeyFromURL(*transaction.Photo)
		if photoKey != "" {
			s.R2Client.DeleteFile(photoKey) // Ignore error for now
		}
	}

	return s.TransactionRepo.DeleteTransaction(userID, transactionID)
}

// GetAllUserTransactions - Retrieves all transactions for a user
// Tujuan: Mengambil semua transaksi milik user dengan urutan terbaru
// Parameter: userID (uuid.UUID)
// Return: ([]models.Transaction, error) - daftar transaksi atau error jika gagal
// Penjelasan: Validasi user ada kemudian ambil semua transaksi dari repository
func (s *TransactionService) GetAllUserTransactions(userID uuid.UUID) ([]models.Transaction, error) {
	// Validate user exists
	_, err := s.UserRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.TransactionRepo.GetAllUserTransaction(userID)
}

// GetTransactionsByPeriod - Retrieves transactions grouped by period
// Tujuan: Mengambil transaksi yang dikelompokkan berdasarkan periode waktu tertentu
// Parameter: userID (uuid.UUID), periodType (models.PeriodType), startDate (time.Time), endDate (time.Time)
// Return: (*response.GetAllTransactionByPeriodResponse, error) - data transaksi berdasarkan periode atau error
// Penjelasan: Validasi user dan periode kemudian ambil transaksi dengan grouping berdasarkan periode
func (s *TransactionService) GetTransactionsByPeriod(
	userID uuid.UUID,
	periodType models.PeriodType,
) (*response.GetAllTransactionByPeriodResponse, error) {
	// Validate user exists
	_, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.TransactionRepo.GetAllByPeriod(userID, periodType)
}

// extractKeyFromURL - Extract storage key from full URL
// Tujuan: Mengekstrak key storage dari URL penuh untuk keperluan delete file
// Parameter: url (string) - URL penuh dari photo
// Return: string - key storage atau empty string jika gagal
// Penjelasan: Memisahkan base URL dengan key untuk mendapatkan path file di storage
func (s *TransactionService) extractKeyFromURL(url string) string {
	// Extract key from URL by removing base URL
	// Example: https://storage.com/user/transaction/userId/filename.jpg -> user/transaction/userId/filename.jpg
	baseURL := os.Getenv("OBJECT_STORAGE_URL")
	if strings.HasPrefix(url, baseURL+"/") {
		return strings.TrimPrefix(url, baseURL+"/")
	}
	return ""
}