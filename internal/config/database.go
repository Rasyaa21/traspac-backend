package config

import (
	"database/sql"
	"fmt"
	"gin-backend-app/internal/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getPoolConfig() (maxOpen, maxIdle int, maxLifetime, maxIdleTime time.Duration) {
	env := getEnv("ENV", "development")

	switch env {
	case "production":
		return 100, 25, 2 * time.Hour, time.Hour
	case "staging":
		return 50, 15, time.Hour, 30 * time.Minute
	case "testing":
		return 150, 50, 15 * time.Minute, 5 * time.Minute
	case "development":
		return 30, 10, 30 * time.Minute, 15 * time.Minute
	default:
		return 20, 5, time.Hour, 30 * time.Minute
	}
}

func buildDSN() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "admin")
	password := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "traspac_db")
	port := getEnv("DB_PORT", "5432")

	// DSN style URL untuk postgres driver
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbName)
}

func InitDatabase() (*gorm.DB, error) {
	dsn := buildDSN()
	env := getEnv("ENV", "development")

	// Configure logger based on environment
	var logLevel logger.LogLevel
	switch env {
	case "production":
		logLevel = logger.Error
	case "testing":
		logLevel = logger.Silent
	default:
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logLevel),
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	maxOpen, maxIdle, maxLifetime, maxIdleTime := getPoolConfig()
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(maxLifetime)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if env != "testing" {
		log.Printf("✅ DB connected (%s) - MaxOpen=%d MaxIdle=%d", env, maxOpen, maxIdle)
	}

	// Setup monitoring only for non-testing environments
	if env != "testing" {
		setupPoolMonitoring(sqlDB)
	}

	// ENUMs & extensions setup
	if err := setupExtensionsAndEnums(db); err != nil {
		return nil, err
	}

	// AutoMigrate all models
	if err := db.AutoMigrate(
		&models.User{},
		&models.UserBudget{},
		&models.Category{},
		&models.Transaction{},
		&models.PeriodReport{},
		&models.AILog{},
		&models.UserToken{},
	); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	// Create custom indexes
	if err := createCustomIndexes(db); err != nil {
		return nil, err
	}

	if env != "testing" {
		log.Println("✅ Database schema migrated & indexes created")
	}

	return db, nil
}

func setupExtensionsAndEnums(db *gorm.DB) error {
	// Enable UUID extension
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return fmt.Errorf("failed to enable uuid-ossp: %w", err)
	}

	// Enable pgcrypto for gen_random_uuid()
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto`).Error; err != nil {
		return fmt.Errorf("failed to enable pgcrypto: %w", err)
	}

	// Create ENUM types
	enums := []struct {
		name   string
		values []string
	}{
		{
			name:   "transaction_group_enum",
			values: []string{"income", "expense"},
		},
		{
			name:   "category_type_enum",
			values: []string{"needs", "wants", "savings"},
		},
		{
			name:   "period_type_enum",
			values: []string{"weekly", "monthly"},
		},
		{
			name: "ai_analysis_type_enum",
			values: []string{
				"weekly_summary",
				"monthly_summary",
				"yearly_summary",
				"compare_period",
				"budget_evaluation",
			},
		},
		{
			name:   "token_type_enum",
			values: []string{"email_verification", "password_reset", "two_factor"},
		},
	}

	for _, enum := range enums {
		if err := createEnumType(db, enum.name, enum.values); err != nil {
			return fmt.Errorf("failed to create %s: %w", enum.name, err)
		}
	}

	return nil
}

func createEnumType(db *gorm.DB, typeName string, values []string) error {
	// Check if enum exists
	var exists bool
	err := db.Raw("SELECT EXISTS(SELECT 1 FROM pg_type WHERE typname = ?)", typeName).Scan(&exists).Error
	if err != nil {
		return err
	}

	if !exists {
		// Build values string
		valuesStr := ""
		for i, value := range values {
			if i > 0 {
				valuesStr += ", "
			}
			valuesStr += fmt.Sprintf("'%s'", value)
		}

		query := fmt.Sprintf("CREATE TYPE %s AS ENUM (%s)", typeName, valuesStr)
		if err := db.Exec(query).Error; err != nil {
			return err
		}
	}

	return nil
}

func createCustomIndexes(db *gorm.DB) error {
	indexes := []struct {
		name  string
		table string
		query string
	}{
		// Users indexes
		{
			name:  "idx_users_email_unique",
			table: "users",
			query: "CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique ON users(LOWER(email))",
		},
		{
			name:  "idx_users_name_unique",
			table: "users",
			query: "CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name_unique ON users(LOWER(name))",
		},

		// User Budgets indexes
		{
			name:  "idx_user_budgets_user_id",
			table: "user_budgets",
			query: "CREATE INDEX IF NOT EXISTS idx_user_budgets_user_id ON user_budgets(user_id)",
		},
		{
			name:  "idx_user_budgets_income_weekly",
			table: "user_budgets",
			query: "CREATE INDEX IF NOT EXISTS idx_user_budgets_income_weekly ON user_budgets(income_weekly)",
		},
		{
			name:  "idx_user_budgets_usage",
			table: "user_budgets",
			query: "CREATE INDEX IF NOT EXISTS idx_user_budgets_usage ON user_budgets(needs_used, wants_used, savings_used)",
		},

		// Categories indexes - Removed category_type references
		{
			name:  "idx_categories_user_name",
			table: "categories",
			query: "CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_user_name ON categories(user_id, LOWER(name))",
		},
		{
			name:  "idx_categories_user_id",
			table: "categories",
			query: "CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories(user_id)",
		},
		{
			name:  "idx_categories_is_default",
			table: "categories",
			query: "CREATE INDEX IF NOT EXISTS idx_categories_is_default ON categories(user_id, is_default)",
		},

		// Transactions indexes
		{
			name:  "idx_transactions_user_date",
			table: "transactions",
			query: "CREATE INDEX IF NOT EXISTS idx_transactions_user_date ON transactions(user_id, date DESC)",
		},
		{
			name:  "idx_transactions_user_category_date",
			table: "transactions",
			query: "CREATE INDEX IF NOT EXISTS idx_transactions_user_category_date ON transactions(user_id, category_id, date DESC)",
		},
		{
			name:  "idx_transactions_user_type_date",
			table: "transactions",
			query: "CREATE INDEX IF NOT EXISTS idx_transactions_user_type_date ON transactions(user_id, type, date DESC)",
		},
		{
			name:  "idx_transactions_budget_category",
			table: "transactions",
			query: "CREATE INDEX IF NOT EXISTS idx_transactions_budget_category ON transactions(user_id, budget_category, date DESC)",
		},

		// Period Reports indexes
		{
			name:  "idx_period_reports_user_period",
			table: "period_reports",
			query: "CREATE INDEX IF NOT EXISTS idx_period_reports_user_period ON period_reports(user_id, period_start, period_end)",
		},
		{
			name:  "idx_period_reports_user_type",
			table: "period_reports",
			query: "CREATE INDEX IF NOT EXISTS idx_period_reports_user_type ON period_reports(user_id, period_type)",
		},

		// AI Logs indexes
		{
			name:  "idx_ai_logs_user_created_at",
			table: "ai_logs",
			query: "CREATE INDEX IF NOT EXISTS idx_ai_logs_user_created_at ON ai_logs(user_id, created_at DESC)",
		},
		{
			name:  "idx_ai_logs_user_analysis_type",
			table: "ai_logs",
			query: "CREATE INDEX IF NOT EXISTS idx_ai_logs_user_analysis_type ON ai_logs(user_id, analysis_type, created_at DESC)",
		},
		{
			name:  "idx_ai_logs_transaction_id",
			table: "ai_logs",
			query: "CREATE INDEX IF NOT EXISTS idx_ai_logs_transaction_id ON ai_logs(transaction_id)",
		},
		{
			name:  "idx_ai_logs_category_id",
			table: "ai_logs",
			query: "CREATE INDEX IF NOT EXISTS idx_ai_logs_category_id ON ai_logs(category_id)",
		},

		// User Tokens indexes
		{
			name:  "idx_user_tokens_token_otp",
			table: "user_tokens",
			query: "CREATE UNIQUE INDEX IF NOT EXISTS idx_user_tokens_token_otp ON user_tokens(token_otp)",
		},
		{
			name:  "idx_user_tokens_verify_token",
			table: "user_tokens",
			query: "CREATE UNIQUE INDEX IF NOT EXISTS idx_user_tokens_verify_token ON user_tokens(verify_token) WHERE verify_token IS NOT NULL",
		},
		{
			name:  "idx_user_tokens_user_type",
			table: "user_tokens",
			query: "CREATE INDEX IF NOT EXISTS idx_user_tokens_user_type ON user_tokens(user_id, token_type)",
		},
		{
			name:  "idx_user_tokens_expires_at",
			table: "user_tokens",
			query: "CREATE INDEX IF NOT EXISTS idx_user_tokens_expires_at ON user_tokens(expires_at)",
		},
		{
			name:  "idx_user_tokens_used_at",
			table: "user_tokens",
			query: "CREATE INDEX IF NOT EXISTS idx_user_tokens_used_at ON user_tokens(used_at)",
		},
		{
			name:  "idx_user_tokens_user_id",
			table: "user_tokens",
			query: "CREATE INDEX IF NOT EXISTS idx_user_tokens_user_id ON user_tokens(user_id)",
		},
	}

	for _, idx := range indexes {
		if err := db.Exec(idx.query).Error; err != nil {
			return fmt.Errorf("failed to create index %s on table %s: %w", idx.name, idx.table, err)
		}
	}

	return nil
}

func setupPoolMonitoring(sqlDB *sql.DB) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		defer ticker.Stop()
		var lastWaitCount int64
		for range ticker.C {
			stats := sqlDB.Stats()
			newWaits := stats.WaitCount - lastWaitCount
			lastWaitCount = stats.WaitCount
			if newWaits == 0 {
				continue
			}
			if newWaits > 10 {
				log.Printf("[DB Pool] High waits: New=%d InUse=%d Idle=%d Open=%d WaitDuration=%s",
					newWaits, stats.InUse, stats.Idle, stats.OpenConnections, stats.WaitDuration)
			}
		}
	}()
}