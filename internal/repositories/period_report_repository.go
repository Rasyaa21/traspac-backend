package repositories

import (
	"encoding/json"
	"gin-backend-app/internal/dto/response"
	"gin-backend-app/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PeriodReportRepository struct {
	DB *gorm.DB
}

func NewPeriodReportRepository(db *gorm.DB) *PeriodReportRepository {
	return &PeriodReportRepository{DB: db}
}

type categoryAgg struct {
	Bucket *models.CategoryType `gorm:"column:bucket"`
	InSum  int64                `gorm:"column:in_sum"`
	OutSum int64                `gorm:"column:out_sum"`
}

func (r *PeriodReportRepository) CreatePeriodReport(userId uuid.UUID, startDate time.Time,
	endDate time.Time) (*response.EStatement, error)  {

		var startingBalance int64
		var endingBalance int64

		err := r.DB.Model(&models.Transaction{}).
		Select(`COALESCE(SUM(
			CASE 
			   WHEN type = ? THEN amount
			   WHEN type = ? THEN -amount
			   ELSE 0
			END
		), 0)
		`,
		models.TransactionGroupIncome,
		models.TransactionGroupExpense,
	).Where("user_id = ?", userId).
	Where("date < ?", startDate).
	Scan(&startingBalance).Error

	

	if err != nil {
		return nil, err
	}

	err = r.DB.
		Model(&models.Transaction{}).
		Select(`
			COALESCE(SUM(
				CASE
					WHEN type = ? THEN amount
					WHEN type = ? THEN -amount
					ELSE 0
				END
			), 0)
		`,
			models.TransactionGroupIncome,
			models.TransactionGroupExpense,
		).
		Where("user_id = ?", userId).
		Where("date <= ?", endDate).
		Scan(&endingBalance).Error

	if err != nil {
		return nil, err
	}

	var catStart []categoryAgg
	err = r.DB.
		Model(&models.Transaction{}).
		Select(`
			budget_category AS bucket,
			COALESCE(SUM(CASE WHEN type = ? THEN amount ELSE 0 END), 0) AS in_sum,
			COALESCE(SUM(CASE WHEN type = ? THEN amount ELSE 0 END), 0) AS out_sum
		`,
			models.TransactionGroupIncome,
			models.TransactionGroupExpense,
		).
		Where("user_id = ?", userId).
		Where("date < ?", startDate).
		Group("budget_category").
		Scan(&catStart).Error
	if err != nil {
		return nil, err
	}

	var catMove []categoryAgg
	err = r.DB.
		Model(&models.Transaction{}).
		Select(`
			budget_category AS bucket,
			COALESCE(SUM(CASE WHEN type = ? THEN amount ELSE 0 END), 0) AS in_sum,
			COALESCE(SUM(CASE WHEN type = ? THEN amount ELSE 0 END), 0) AS out_sum
		`,
			models.TransactionGroupIncome,
			models.TransactionGroupExpense,
		).
		Where("user_id = ?", userId).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Group("budget_category").
		Scan(&catMove).Error
	if err != nil {
		return nil, err
	}

	categories := map[models.CategoryType]response.BucketBalance{
		models.CategoryTypeNeeds:   {},
		models.CategoryTypeWants:   {},
		models.CategoryTypeSavings: {},
	}

	for _, row := range catStart {
		if row.Bucket == nil {
			continue
		}
		b := categories[*row.Bucket]
		b.Starting = row.InSum - row.OutSum
		categories[*row.Bucket] = b
	}

	for _, row := range catMove {
		if row.Bucket == nil {
			continue
		}
		b := categories[*row.Bucket]
		b.In = row.InSum
		b.Out = row.OutSum
		b.Ending = b.Starting + b.In - b.Out
		categories[*row.Bucket] = b
	}

	var totalBuckets int64
	for _, b := range categories {
		totalBuckets += b.Ending
	}

	var txLines []response.TransactionLine
		err = r.DB.
			Model(&models.Transaction{}).
			Select(`
				id,
				date,
				COALESCE(description, '') AS description,
				type,
				budget_category AS bucket,
				amount
			`).
			Where("user_id = ?", userId).
			Where("date BETWEEN ? AND ?", startDate, endDate).
			Order("date ASC, created_at ASC").
			Scan(&txLines).Error
		if err != nil {
			return nil, err
		}

	res := &response.EStatement{
		GeneratedAt: time.Now(),
		Period: response.PeriodWindow{
			Start: startDate,
			End:   endDate,
		},
		Balances: response.BalanceSection{
			StartingBalance: startingBalance,
			EndingBalance:   endingBalance,
			Categories:      categories,
			TotalBuckets:    totalBuckets,
		},
		Transactions: txLines,
	}

	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	periodReport := *&models.PeriodReport{
		UserID: userId,
		PeriodStart: startDate,
		PeriodEnd: endDate,
		ReportData: datatypes.JSON(b),
		GeneratedAt: time.Now(),
	}

	err = r.DB.Create(&periodReport).Error

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *PeriodReportRepository) GetAllUserReports(userId uuid.UUID) (*models.PeriodReport, error) {
	var reports *models.PeriodReport
	err := r.DB.Model(&models.PeriodReport{}).Where("user_id = ?", userId).Find(&reports).Error
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (r *PeriodReportRepository) GetUserPeriodById(userId uuid.UUID, reportId uuid.UUID) (*models.PeriodReport, error) {
	var report *models.PeriodReport
	err := r.DB.Model(&models.PeriodReport{}).Where("user_id = ? AND id = ?", userId, reportId).First(&report).Error
	if err != nil {
		return report, nil
	}
	return report, nil
}