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


type TokenType string

const (
	TokenTypeEmailVerification TokenType = "email_verification"
	TokenTypePasswordReset     TokenType = "password_reset"
)