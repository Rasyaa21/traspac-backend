package utils

import (
	"errors"
	"gin-backend-app/internal/models"
)

func BuildPeriodQuery(period models.PeriodType) (string, error) {
	switch period {
	case models.PeriodDaily:
		return "DATE(t.date)", nil
	case models.PeriodWeekly:
		return "date_trunc('week', t.date)", nil
	case models.PeriodMonthly:
		return "date_trunc('month', t.date)", nil
	default:
		return "", errors.New("invalid period type")
	}
}