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

// DefaultAllocation returns the default 50-30-20 allocation
func DefaultAllocation() BudgetAllocation {
    return BudgetAllocation{
        SavingsPercent: 50.0, // 50% for savings
        WantsPercent:   30.0, // 30% for wants
        NeedsPercent:   20.0, // 20% for needs
    }
}

// NewCustomAllocation creates a custom allocation
func NewCustomAllocation(savings, wants, needs float64) BudgetAllocation {
    return BudgetAllocation{
        SavingsPercent: savings,
        WantsPercent:   wants,
        NeedsPercent:   needs,
    }
}

func NewUserBudgetService(userRepo *repositories.UserRepository, userBudgetRepo *repositories.UserBudgetRepository) *UserBudgetService {
    return &UserBudgetService{
        UserRepo:       userRepo,
        UserBudgetRepo: userBudgetRepo,
    }
}

// CreateUserBudget creates budget with specified allocation (default or custom)
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
        IncomeWeekly:  weeklyIncome,
        NeedsBudget:   budgetAmounts.Needs,
        WantsBudget:   budgetAmounts.Wants,
        SavingsBudget: budgetAmounts.Savings,
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

// UpdateUserBudget updates budget with specified allocation (default or custom)
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

    if err := s.UserBudgetRepo.EdituserBudget(budgetAmounts.Needs, budgetAmounts.Savings, budgetAmounts.Wants, weeklyIncome, userID); err != nil {
        return nil, fmt.Errorf("failed to update user budget: %w", err)
    }

    // Return updated budget
    return &models.UserBudget{
        ID:            existingBudget.ID,
        UserID:        userID,
        IncomeWeekly:  weeklyIncome,
        NeedsBudget:   budgetAmounts.Needs,
        WantsBudget:   budgetAmounts.Wants,
        SavingsBudget: budgetAmounts.Savings,
        CreatedAt:     existingBudget.CreatedAt,
        UpdatedAt:     existingBudget.UpdatedAt,
    }, nil
}

// GetUserBudget retrieves user budget
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

// DeleteUserBudget removes user budget
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

// CalculateBudgetPreview calculates budget preview without saving
func (s *UserBudgetService) CalculateBudgetPreview(weeklyIncome float64, allocation *BudgetAllocation) (map[string]interface{}, error) {
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

    // Calculate budget amounts
    budgetAmounts := s.calculateBudgetAmounts(weeklyIncome, *allocation)

    preview := map[string]interface{}{
        "preview": map[string]interface{}{
            "weekly_income":    weeklyIncome,
            "savings_budget":   budgetAmounts.Savings,
            "wants_budget":     budgetAmounts.Wants,
            "needs_budget":     budgetAmounts.Needs,
            "total_allocated":  budgetAmounts.Savings + budgetAmounts.Wants + budgetAmounts.Needs,
            "remaining":        weeklyIncome - (budgetAmounts.Savings + budgetAmounts.Wants + budgetAmounts.Needs),
        },
        "allocation": map[string]interface{}{
            "savings_percent": fmt.Sprintf("%.1f%%", allocation.SavingsPercent),
            "wants_percent":   fmt.Sprintf("%.1f%%", allocation.WantsPercent),
            "needs_percent":   fmt.Sprintf("%.1f%%", allocation.NeedsPercent),
        },
        "monthly_projection": map[string]interface{}{
            "monthly_income":  weeklyIncome * 4.33,
            "monthly_savings": budgetAmounts.Savings * 4.33,
            "monthly_wants":   budgetAmounts.Wants * 4.33,
            "monthly_needs":   budgetAmounts.Needs * 4.33,
        },
        "breakdown": map[string]interface{}{
            "savings_amount": budgetAmounts.Savings,
            "wants_amount":   budgetAmounts.Wants,
            "needs_amount":   budgetAmounts.Needs,
        },
    }

    return preview, nil
}

type BudgetAmounts struct {
    Savings float64
    Wants   float64
    Needs   float64
}

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

func (s *UserBudgetService) calculateBudgetAmounts(weeklyIncome float64, allocation BudgetAllocation) BudgetAmounts {
    return BudgetAmounts{
        Savings: weeklyIncome * (allocation.SavingsPercent / 100),
        Wants:   weeklyIncome * (allocation.WantsPercent / 100),
        Needs:   weeklyIncome * (allocation.NeedsPercent / 100),
    }
}

