package models

type TransactionGroupType string

const (
	TransactionGroupIncome TransactionGroupType = "income"
	TransactionGroupExpense TransactionGroupType = "expense"
)

type PeriodType string

const (
	PeriodWeekly PeriodType = "weekly"
	PeriodMonthly PeriodType = "monthly"
)

type AiAnalysisType string

const (
	AiWeeklySummary    AiAnalysisType = "weekly_summary"
	AiMonthlySummary   AiAnalysisType = "monthly_summary"
	AiYearlySummary    AiAnalysisType = "yearly_summary"
	AiComparePeriod    AiAnalysisType = "compare_period"
	AiBudgetEvaluation AiAnalysisType = "budget_evaluation"
)

type TokenType string

const (
	TokenTypeEmailVerification TokenType = "email_verification"
	TokenTypePasswordReset TokenType = "password_reset"
)