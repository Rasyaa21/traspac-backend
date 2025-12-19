package request

import (
	"time"
)
type PeriodRequest struct {
	StartDate      time.Time             `form:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate        time.Time             `form:"end_date" binding:"required" time_format:"2006-01-02"`
}