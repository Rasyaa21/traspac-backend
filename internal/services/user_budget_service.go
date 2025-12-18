package services

import (
	"errors"
	"fmt"
	"gin-backend-app/internal/models"
	"gin-backend-app/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserBudgetService struct {
    UserRepo       *repositories.UserRepository
    UserBudgetRepo *repositories.UserBudgetRepository
}

// BudgetAllocation represents budget allocation percentages
type BudgetAllocation struct {
    SavingsPercent float64
    WantsPercent   float64
    NeedsPercent   float64
}

// BudgetAmounts represents calculated budget amounts in currency
type BudgetAmounts struct {
    Savings float64
    Wants   float64
    Needs   float64
}

// DefaultAllocation - Returns default 50-30-20 budget allocation
// Tujuan: Menyediakan alokasi budget default dengan persentase 50% tabungan, 30% keinginan, 20% kebutuhan
// Parameter: Tidak ada
// Return: BudgetAllocation struct dengan default percentages
// Penjelasan: Menggunakan aturan 50-30-20 yang populer dalam financial planning
func DefaultAllocation() BudgetAllocation {
    return BudgetAllocation{
        SavingsPercent: 50.0, 
        WantsPercent:   30.0, 
        NeedsPercent:   20.0, 
    }
}

// NewCustomAllocation - Creates custom budget allocation with specified percentages
// Tujuan: Membuat alokasi budget kustom sesuai persentase yang ditentukan user
// Parameter: savings (float64) - persentase untuk tabungan, wants (float64) - persentase untuk keinginan, needs (float64) - persentase untuk kebutuhan
// Return: BudgetAllocation struct dengan custom percentages
// Penjelasan: Memungkinkan user untuk mengatur alokasi budget sesuai preferensi personal
func NewCustomAllocation(savings, wants, needs float64) BudgetAllocation {
    return BudgetAllocation{
        SavingsPercent: savings,
        WantsPercent:   wants,
        NeedsPercent:   needs,
    }
}

// NewUserBudgetService - Creates new instance of UserBudgetService
// Tujuan: Membuat instance baru dari UserBudgetService dengan dependency injection
// Parameter: userRepo (*repositories.UserRepository) - repository untuk user, userBudgetRepo (*repositories.UserBudgetRepository) - repository untuk budget
// Return: *UserBudgetService - instance baru dari service
// Penjelasan: Constructor pattern untuk inisialisasi service dengan dependencies yang dibutuhkan
func NewUserBudgetService(userRepo *repositories.UserRepository, userBudgetRepo *repositories.UserBudgetRepository) *UserBudgetService {
    return &UserBudgetService{
        UserRepo:       userRepo,
        UserBudgetRepo: userBudgetRepo,
    }
}

// CreateUserBudget - Creates new budget for user with specified allocation
// Tujuan: Membuat budget baru untuk user dengan alokasi default atau kustom
// Parameter: userID (uuid.UUID) - ID user, weeklyIncome (float64) - penghasilan mingguan, allocation (*BudgetAllocation) - alokasi budget (nil untuk default)
// Return: (*models.UserBudget, error) - budget yang dibuat atau error jika gagal
// Penjelasan: Validasi user, income, dan alokasi sebelum membuat budget baru dengan perhitungan otomatis
func (s *UserBudgetService) CreateUserBudget(userID uuid.UUID, weeklyIncome float64, allocation *BudgetAllocation) (*models.UserBudget, error) {
    // Validate user exists
    _, err := s.UserRepo.FindByID(userID)
    if err != nil {
        return nil, errors.New("user not found")
    }

    if weeklyIncome <= 0 {
        return nil, errors.New("weekly income must be greater than 0")
    }

    // Use default allocation if none provided
    if allocation == nil {
        defaultAlloc := DefaultAllocation()
        allocation = &defaultAlloc
    }

    // Validate allocation percentages
    if err := s.validateAllocation(*allocation); err != nil {
        return nil, err
    }

    // Check if user already has budget
    existingBudget, _ := s.UserBudgetRepo.GetUserBudget(userID)
    if existingBudget != nil {
        return nil, errors.New("user budget already exists, use update instead")
    }

    // Calculate budget amounts
    budgetAmounts := s.calculateBudgetAmounts(weeklyIncome, *allocation)

    userBudget := &models.UserBudget{
        UserID:        userID,
        IncomeWeekly:  int64(weeklyIncome),
        NeedsBudget:   int64(budgetAmounts.Needs),
        WantsBudget:   int64(budgetAmounts.Wants),
        SavedMoney:    int64(budgetAmounts.Savings),
        NeedsUsed:     0,
        WantsUsed:     0,
        SavingsUsed:   0,
    }

    // Validate budget distribution
    if err := userBudget.ValidateBudgetDistribution(); err != nil {
        return nil, fmt.Errorf("budget validation failed: %w", err)
    }

    if err := s.UserBudgetRepo.CreateUserBudget(userBudget); err != nil {
        return nil, fmt.Errorf("failed to create user budget: %w", err)
    }

    return userBudget, nil
}

// UpdateUserBudget - Updates existing user budget with new allocation
// Tujuan: Mengupdate budget yang sudah ada dengan alokasi baru
// Parameter: userID (uuid.UUID) - ID user, weeklyIncome (float64) - penghasilan mingguan baru, allocation (*BudgetAllocation) - alokasi budget baru (nil untuk default)
// Return: (*models.UserBudget, error) - budget yang diupdate atau error jika gagal
// Penjelasan: Update budget dengan mempertahankan usage yang sudah ada dan menyesuaikan dengan alokasi baru
func (s *UserBudgetService) UpdateUserBudget(userID uuid.UUID, weeklyIncome float64, allocation *BudgetAllocation) (*models.UserBudget, error) {
    // Validate user exists
    _, err := s.UserRepo.FindByID(userID)
    if err != nil {
        return nil, errors.New("user not found")
    }

    if weeklyIncome <= 0 {
        return nil, errors.New("weekly income must be greater than 0")
    }

    // Use default allocation if none provided
    if allocation == nil {
        defaultAlloc := DefaultAllocation()
        allocation = &defaultAlloc
    }

    // Validate allocation percentages
    if err := s.validateAllocation(*allocation); err != nil {
        return nil, err
    }

    // Check if budget exists
    existingBudget, err := s.UserBudgetRepo.GetUserBudget(userID)
    if err != nil {
        return nil, errors.New("user budget not found")
    }

    // Calculate budget amounts
    budgetAmounts := s.calculateBudgetAmounts(weeklyIncome, *allocation)

    if err := s.UserBudgetRepo.EdituserBudget(int64(budgetAmounts.Needs), int64(budgetAmounts.Savings), int64(budgetAmounts.Wants), int64(weeklyIncome), userID); err != nil {
        return nil, fmt.Errorf("failed to update user budget: %w", err)
    }

    // Return updated budget (keep existing usage)
    return &models.UserBudget{
        ID:            existingBudget.ID,
        UserID:        userID,
        IncomeWeekly:  int64(weeklyIncome),
        NeedsBudget:   int64(budgetAmounts.Needs),
        WantsBudget:   int64(budgetAmounts.Wants),
        SavedMoney:    int64(budgetAmounts.Savings),
        NeedsUsed:     existingBudget.NeedsUsed,
        WantsUsed:     existingBudget.WantsUsed,
        SavingsUsed:   existingBudget.SavingsUsed,
        CreatedAt:     existingBudget.CreatedAt,
        UpdatedAt:     existingBudget.UpdatedAt,
    }, nil
}

// GetUserBudget - Retrieves user budget by user ID
// Tujuan: Mengambil data budget user berdasarkan ID user
// Parameter: userID (uuid.UUID) - ID user yang budgetnya ingin diambil
// Return: (*models.UserBudget, error) - data budget user atau error jika tidak ditemukan
// Penjelasan: Validasi user ada kemudian mengambil budget dari repository
func (s *UserBudgetService) GetUserBudget(userID uuid.UUID) (*models.UserBudget, error) {
    _, err := s.UserRepo.FindByID(userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("user not found")
        }
        return nil, err
    }

    budget, err := s.UserBudgetRepo.GetUserBudget(userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("user budget not found")
        }
        return nil, err
    }

    return budget, nil
}

// DeleteUserBudget - Removes user budget from database
// Tujuan: Menghapus budget user dari database
// Parameter: userID (uuid.UUID) - ID user yang budgetnya ingin dihapus
// Return: error - error jika gagal menghapus atau nil jika berhasil
// Penjelasan: Validasi user dan budget ada sebelum melakukan penghapusan
func (s *UserBudgetService) DeleteUserBudget(userID uuid.UUID) error {
    // Validate user exists
    _, err := s.UserRepo.FindByID(userID)
    if err != nil {
        return errors.New("user not found")
    }

    // Check if budget exists
    _, err = s.UserBudgetRepo.GetUserBudget(userID)
    if err != nil {
        return errors.New("user budget not found")
    }

    if err := s.UserBudgetRepo.DeleteUserBudget(userID); err != nil {
        return fmt.Errorf("failed to delete user budget: %w", err)
    }

    return nil
}

// GetBudgetSummary - Gets detailed budget summary with usage statistics
// Tujuan: Mengambil ringkasan lengkap budget termasuk statistik penggunaan
// Parameter: userID (uuid.UUID) - ID user yang ringkasannya ingin diambil
// Return: (map[string]interface{}, error) - data ringkasan budget atau error jika gagal
// Penjelasan: Menghitung persentase, total allocated, dan total used untuk memberikan overview budget
func (s *UserBudgetService) GetBudgetSummary(userID uuid.UUID) (map[string]interface{}, error) {
    budget, err := s.GetUserBudget(userID)
    if err != nil {
        return nil, err
    }

    summary := map[string]interface{}{
        "weekly_income": budget.IncomeWeekly,
        "budget_allocation": map[string]interface{}{
            "needs": map[string]interface{}{
                "allocated": budget.NeedsBudget,
                "used":      budget.NeedsUsed,
                "remaining": budget.NeedsBudget - budget.NeedsUsed,
                "percentage": float64(budget.NeedsBudget) / float64(budget.IncomeWeekly) * 100,
            },
            "wants": map[string]interface{}{
                "allocated": budget.WantsBudget,
                "used":      budget.WantsUsed,
                "remaining": budget.WantsBudget - budget.WantsUsed,
                "percentage": float64(budget.WantsBudget) / float64(budget.IncomeWeekly) * 100,
            },
            "savings": map[string]interface{}{
                "allocated": budget.SavedMoney,
                "used":      budget.SavingsUsed,
                "remaining": budget.SavedMoney - budget.SavingsUsed,
            },
        },
        "total_allocated": budget.NeedsBudget + budget.WantsBudget + budget.SavedMoney,
        "total_used":      budget.NeedsUsed + budget.WantsUsed + budget.SavingsUsed,
    }

    return summary, nil
}

// SpendFromBudget - Deducts amount from specific budget category
// Tujuan: Mengurangi jumlah uang dari kategori budget tertentu saat ada pengeluaran
// Parameter: userID (uuid.UUID) - ID user, amount (int64) - jumlah yang dikeluarkan, categoryType (models.CategoryType) - tipe kategori (needs/wants/savings)
// Return: error - error jika gagal atau budget tidak mencukupi, nil jika berhasil
// Penjelasan: Menggunakan repository method untuk mengurangi budget dengan validasi ketersediaan dana
func (s *UserBudgetService) SpendFromBudget(userID uuid.UUID, amount int64, categoryType models.CategoryType) error {
    // Validate user exists
    _, err := s.UserRepo.FindByID(userID)
    if err != nil {
        return errors.New("user not found")
    }

    return s.UserBudgetRepo.Spend(userID, amount, categoryType)
}

// AddIncomeToBudget - Adds additional income to specific budget category
// Tujuan: Menambahkan penghasilan tambahan ke kategori budget tertentu
// Parameter: userID (uuid.UUID) - ID user, amount (int64) - jumlah penghasilan tambahan, categoryType (models.CategoryType) - tipe kategori tujuan
// Return: error - error jika gagal menambahkan, nil jika berhasil
// Penjelasan: Menambahkan penghasilan tambahan ke total income dan kategori budget yang dipilih
func (s *UserBudgetService) AddIncomeToBudget(userID uuid.UUID, amount int64, categoryType models.CategoryType) error {
    // Validate user exists
    _, err := s.UserRepo.FindByID(userID)
    if err != nil {
        return errors.New("user not found")
    }

    return s.UserBudgetRepo.AddIncomeToBucket(userID, amount, categoryType)
}

// validateAllocation - Validates budget allocation percentages
// Tujuan: Memvalidasi persentase alokasi budget tidak melebihi 100% dan tidak negatif
// Parameter: allocation (BudgetAllocation) - struct alokasi budget yang ingin divalidasi
// Return: error - error jika validasi gagal, nil jika valid
// Penjelasan: Memeriksa total persentase tidak lebih dari 100% dan semua nilai tidak negatif
func (s *UserBudgetService) validateAllocation(allocation BudgetAllocation) error {
    totalPercent := allocation.SavingsPercent + allocation.WantsPercent + allocation.NeedsPercent
    
    if totalPercent > 100 {
        return errors.New("total percentage cannot exceed 100%")
    }
    
    if allocation.SavingsPercent < 0 || allocation.WantsPercent < 0 || allocation.NeedsPercent < 0 {
        return errors.New("percentage values cannot be negative")
    }
    
    return nil
}

// calculateBudgetAmounts - Calculates budget amounts based on income and allocation
// Tujuan: Menghitung jumlah budget untuk setiap kategori berdasarkan penghasilan dan alokasi persentase
// Parameter: weeklyIncome (float64) - penghasilan mingguan, allocation (BudgetAllocation) - alokasi persentase
// Return: BudgetAmounts - struct berisi jumlah budget untuk setiap kategori
// Penjelasan: Mengkalikan penghasilan dengan persentase masing-masing kategori untuk mendapat nominal budget
func (s *UserBudgetService) calculateBudgetAmounts(weeklyIncome float64, allocation BudgetAllocation) BudgetAmounts {
    return BudgetAmounts{
        Savings: weeklyIncome * (allocation.SavingsPercent / 100),
        Wants:   weeklyIncome * (allocation.WantsPercent / 100),
        Needs:   weeklyIncome * (allocation.NeedsPercent / 100),
    }
}

