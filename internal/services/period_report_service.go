package services

import (
	"encoding/json"
	"gin-backend-app/internal/dto/response"
	"gin-backend-app/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type PeriodReportService struct {
	PeriodReportRepo *repositories.PeriodReportRepository
}

func NewPeriodReportService(periodReportRepo *repositories.PeriodReportRepository) *PeriodReportService {
	return &PeriodReportService{
		PeriodReportRepo: periodReportRepo,
	}
}

func (s *PeriodReportService) CreatePeriodReport(userId uuid.UUID, startDate time.Time, endDate time.Time) (*response.EStatement, error) {
	return s.PeriodReportRepo.CreatePeriodReport(userId, startDate, endDate)
}

func (s *PeriodReportService) GetAllUserReports(userId uuid.UUID) (*response.PeriodReportsListResponse, error) {
	periodReports, err := s.PeriodReportRepo.GetAllUserReports(userId)
	if err != nil {
		return nil, err
	}

	var reports []response.PeriodReportResponse
	if periodReports != nil {
		var reportData response.EStatement
		if err := json.Unmarshal([]byte(periodReports.ReportData), &reportData); err != nil {
			return nil, err
		}

		report := response.PeriodReportResponse{
			ID:          periodReports.ID,
			UserID:      periodReports.UserID,
			PeriodStart: periodReports.PeriodStart,
			PeriodEnd:   periodReports.PeriodEnd,
			PDFReport:   periodReports.PdfReport,
			ReportData:  reportData,
			GeneratedAt: periodReports.GeneratedAt,
			CreatedAt:   periodReports.CreatedAt,
			UpdatedAt:   periodReports.UpdatedAt,
		}
		reports = append(reports, report)
	}

	return &response.PeriodReportsListResponse{
		Reports: reports,
	}, nil
}

func (s *PeriodReportService) GetUserReportById(userId uuid.UUID, reportId uuid.UUID) (*response.PeriodReportResponse, error) {
	periodReport, err := s.PeriodReportRepo.GetUserPeriodById(userId, reportId)
	if err != nil {
		return nil, err
	}

	var reportData response.EStatement
	if err := json.Unmarshal([]byte(periodReport.ReportData), &reportData); err != nil {
		return nil, err
	}

	report := &response.PeriodReportResponse{
		ID:          periodReport.ID,
		UserID:      periodReport.UserID,
		PeriodStart: periodReport.PeriodStart,
		PeriodEnd:   periodReport.PeriodEnd,
		PDFReport:   periodReport.PdfReport,
		ReportData:  reportData,
		GeneratedAt: periodReport.GeneratedAt,
		CreatedAt:   periodReport.CreatedAt,
		UpdatedAt:   periodReport.UpdatedAt,
	}

	return report, nil
}