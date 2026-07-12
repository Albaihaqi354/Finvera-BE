package database

import (
	"log"
	"os"
	"time"

	"finvera-be/internal/config"
	"finvera-be/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the shared database instance.
// Prefer injecting this via function arguments rather than using as a global.
var DB *gorm.DB

// ConnectDB opens a postgres connection, configures the connection pool,
// and returns the *gorm.DB so callers don't have to read the global.
func ConnectDB(cfg *config.Config) *gorm.DB {
	if cfg.DBDSN == "" {
		log.Fatal("Database DSN is empty")
	}

	// Set log level based on environment
	gormLogLevel := logger.Warn
	if cfg.AppEnv == "development" {
		gormLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to the database")

	// AutoMigrate is safe for development.
	// In production, AUTO_MIGRATE should not be set (or set to false) to disable it and use a proper migration tool.
	autoMigrate := os.Getenv("AUTO_MIGRATE")
	if autoMigrate == "true" {
		log.Println("Running AutoMigrate...")
		err = db.AutoMigrate(
			&models.User{},
			&models.Account{},
			&models.Category{},
			&models.TagGroup{},
			&models.Tag{},
			&models.Transaction{},
			&models.ScheduledTransaction{},
			&models.JWTBlacklist{},
		)
		if err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("AutoMigrate completed")
	} else {
		log.Println("AutoMigrate skipped (AUTO_MIGRATE != true)")
	}

	DB = db
	return db
}
