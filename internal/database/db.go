package database

import (
	"log"

	"finvera-be/internal/config"
	"finvera-be/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) {
	if cfg.DBDSN == "" {
		log.Fatal("Database DSN is empty")
	}

	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Successfully connected to the database")

	log.Println("Running AutoMigrate...")
	err = db.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Category{},
		&models.TagGroup{},
		&models.Tag{},
		&models.Transaction{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	DB = db
}
