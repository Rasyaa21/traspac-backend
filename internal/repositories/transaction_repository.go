package repositories

import (
	"errors"
	"gin-backend-app/internal/dto/request"
	"gin-backend-app/internal/dto/response"
	"gin-backend-app/internal/models"
	"gin-backend-app/pkg/utils"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type txRow struct {
	models.Transaction
	Period time.Time `json:"period"`
}

type TransactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) CreateTransaction(transaction *models.Transaction) (*models.Transaction, error) {
	err := r.DB.Create(transaction).Error
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactionByID - Get transaction by ID with ownership validation
func (r *TransactionRepository) GetTransactionByID(userID uuid.UUID, transactionID uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.DB.Where("id = ? AND user_id = ?", transactionID, userID).
		First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) EditTransaction(userID uuid.UUID, transactionID uuid.UUID, req *request.TransactionUpdateRequest) (*models.Transaction, error) {
	updates := map[string]any{}

	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.Date != nil {
		updates["date"] = *req.Date
	}
	if req.BudgetCategory != nil {
		updates["budget_category"] = *req.BudgetCategory
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.PhotoURL != nil {
		updates["photo_url"] = *req.PhotoURL
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	updates["updated_at"] = time.Now()

	err := r.DB.Model(&models.Transaction{}).
		Where("id = ? AND user_id = ?", transactionID, userID).
		Updates(updates).Error
	if err != nil {
		return nil, err
	}

	// Get and return updated transaction (no preload)
	var tx models.Transaction
	err = r.DB.
		Where("id = ? AND user_id = ?", transactionID, userID).
		First(&tx).Error
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func (r *TransactionRepository) DeleteTransaction(userId uuid.UUID, transactionId uuid.UUID) error {
	err := r.DB.Model(&models.Transaction{}).Where("user_id = ?", userId).Where("id = ?", transactionId).Delete(&models.Transaction{}).Error
	return err
}

func (r *TransactionRepository) GetAllUserTransaction(userId uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.DB.Model(&models.Transaction{}).Where("user_id = ?", userId).Order("date DESC, created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetAllByPeriod(userId uuid.UUID, periodType models.PeriodType) (*response.GetAllTransactionByPeriodResponse, error) {
	periodQuery, err := utils.BuildPeriodQuery(periodType)
	if err != nil {
		return nil, err
	}

	var rows []txRow
	err = r.DB.Table("transactions t").
		Select("t.*, "+periodQuery+"AS period").
		Where("t.user_id = ?", userId).
		Order("period DESC, t.date DESC, t.created_at DESC").Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	groupMap := map[time.Time]*response.PeriodTransactionGroup{}
	var total response.TransactionSummary

	for _, row := range rows {
		g, ok := groupMap[row.Period]
		if !ok {
			groupMap[row.Period] = &response.PeriodTransactionGroup{
				Period:       row.Period,
				Transactions: []models.Transaction{},
			}
			g = groupMap[row.Period]
		}
		tx := row.Transaction
		g.Transactions = append(g.Transactions, tx)

		if tx.Type == models.TransactionGroupType("income") {
			g.Summary.Income += tx.Amount
			total.Income += tx.Amount
		} else if tx.Type == models.TransactionGroupType("expense") {
			g.Summary.Expense += tx.Amount
			total.Expense += tx.Amount
		}
	}

	total.Total = total.Income - total.Expense

	periods := make([]time.Time, 0, len(groupMap))
	for p := range groupMap {
		periods = append(periods, p)
	}

	sort.Slice(periods, func(i, j int) bool { return periods[i].After(periods[j]) })

	groups := make([]response.PeriodTransactionGroup, 0, len(periods))
	for _, p := range periods {
		g := groupMap[p]
		g.Summary.Total = g.Summary.Income - g.Summary.Expense
		groups = append(groups, *g)
	}

	return &response.GetAllTransactionByPeriodResponse{
		Groups: groups,
		Total:  total,
	}, nil
}

func (r *TransactionRepository) GetByCategoryByPeriod(
	userID uuid.UUID,
	periodType models.PeriodType,
	startDate time.Time,
	endDate time.Time,
) (*response.GetByCategoryResponse, error) {

	periodQuery, err := utils.BuildPeriodQuery(periodType)
	if err != nil {
		return nil, err
	}

	var items []response.CategoryPeriodSummary
	err = r.DB.Table("transactions t").
		Select(periodQuery+" AS period, t.type, COALESCE(SUM(t.amount),0) AS total").
		Where("t.user_id = ?", userID).
		Where("t.date >= ? AND t.date <= ?", startDate, endDate).
		Group("period, t.type").
		Order("period ASC").
		Scan(&items).Error
	if err != nil {
		return nil, err
	}

	type sumRow struct {
		Type  string
		Total int64
	}

	var sums []sumRow
	err = r.DB.Table("transactions t").
		Select("t.type, COALESCE(SUM(t.amount),0) AS total").
		Where("t.user_id = ?", userID).
		Where("t.date >= ? AND t.date <= ?", startDate, endDate).
		Group("t.type").
		Scan(&sums).Error
	if err != nil {
		return nil, err
	}

	var total response.TransactionSummary
	for _, s := range sums {
		if s.Type == "income" {
			total.Income = s.Total
		} else if s.Type == "expense" {
			total.Expense = s.Total
		}
	}
	total.Total = total.Income - total.Expense

	return &response.GetByCategoryResponse{
		Items: items,
		Total: total,
	}, nil
}