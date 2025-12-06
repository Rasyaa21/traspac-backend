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
		// format: postgres://admin:root@localhost:5432/traspac_db?sslmode=disable
		return url
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "admin"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "root"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "traspac_db"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	// DSN style URL untuk postgres driver
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbName)
}

func InitDatabase() (*gorm.DB, error) {
	dsn := buildDSN()
	env := getEnv("DB_NAME", "traspac_db")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
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

	log.Printf("✅ DB connected (%s) - MaxOpen=%d MaxIdle=%d", env, maxOpen, maxIdle)

	// Optional: monitoring pool
	setupPoolMonitoring(sqlDB)

	// ENUMs & index perlu dibuat manual sebelum AutoMigrate
	if err := setupEnums(db); err != nil {
		return nil, err
	}

	// AutoMigrate model-model keuangan
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

	// Tambah index tambahan yang GORM gak generate otomatis
	if err := addFinanceIndexes(db); err != nil {
		return nil, err
	}

	log.Println("✅ Finance schema migrated & indexes created")
	return db, nil
}

func setupEnums(db *gorm.DB) error {
	// Enable extension UUID kalau belum
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto`).Error; err != nil {
		return fmt.Errorf("failed to enable pgcrypto: %w", err)
	}

	// ENUM income/expense
	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transaction_group_enum') THEN
				CREATE TYPE transaction_group_enum AS ENUM ('income', 'expense');
			END IF;
		END$$;
	`).Error; err != nil {
		return fmt.Errorf("failed to create transaction_group_enum: %w", err)
	}

	// ENUM period_type_enum
	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'period_type_enum') THEN
				CREATE TYPE period_type_enum AS ENUM ('weekly', 'monthly');
			END IF;
		END$$;
	`).Error; err != nil {
		return fmt.Errorf("failed to create period_type_enum: %w", err)
	}

	// ENUM ai_analysis_type_enum
	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ai_analysis_type_enum') THEN
				CREATE TYPE ai_analysis_type_enum AS ENUM (
					'weekly_summary',
					'monthly_summary',
					'yearly_summary',
					'compare_period',
					'budget_evaluation'
				);
			END IF;
		END$$;
	`).Error; err != nil {
		return fmt.Errorf("failed to create ai_analysis_type_enum: %w", err)
	}

	return nil
}

func addFinanceIndexes(db *gorm.DB) error {

	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique 
		ON users(LOWER(email));
	`).Error; err != nil {
		return fmt.Errorf("failed to create idx_users_email_unique %w", err)
	}

	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name_unique
		ON users (LOWER(name));
	`).Error; err != nil {
		return fmt.Errorf("failed to create idx_users_name_unique: %w", err)
	}

	// user_budgets
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_user_budgets_user_id
		ON user_budgets(user_id);
	`).Error; err != nil {
		return err
	}

	// categories
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_user_name_type
		ON categories (user_id, LOWER(name), group_type);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_categories_user_group_type
		ON categories (user_id, group_type);
	`).Error; err != nil {
		return err
	}

	// transactions
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_transactions_user_date
		ON transactions (user_id, date DESC);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_transactions_user_category_date
		ON transactions (user_id, category_id, date DESC);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_transactions_user_type_date
		ON transactions (user_id, type, date DESC);
	`).Error; err != nil {
		return err
	}

	// period_reports
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_period_reports_user_period
		ON period_reports (user_id, period_start, period_end);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_period_reports_user_type
		ON period_reports (user_id, period_type);
	`).Error; err != nil {
		return err
	}

	// ai_logs
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_ai_logs_user_created_at
		ON ai_logs (user_id, created_at DESC);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_ai_logs_user_analysis_type
		ON ai_logs (user_id, analysis_type, created_at DESC);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_ai_logs_transaction_id
		ON ai_logs (transaction_id);
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_ai_logs_category_id
		ON ai_logs (category_id);
	`).Error; err != nil {
		return err
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