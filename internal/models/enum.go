package models

type TransactionGroupType string

const (
	TransactionGroupIncome  TransactionGroupType = "income"
	TransactionGroupExpense TransactionGroupType = "expense"
)

type CategoryType string

const (
	CategoryTypeNeeds   CategoryType = "needs"
	CategoryTypeWants   CategoryType = "wants" 
	CategoryTypeSavings CategoryType = "savings"
)

type PeriodType string

const (
	PeriodWeekly  PeriodType = "weekly"
	PeriodMonthly PeriodType = "monthly"
	PeriodDaily PeriodType = "daily"
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
	TokenTypePasswordReset     TokenType = "password_reset"
)