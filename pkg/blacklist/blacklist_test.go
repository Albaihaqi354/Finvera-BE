package blacklist_test

import (
	"testing"
	"time"

	"finvera-be/internal/database"
	"finvera-be/internal/models"
	"finvera-be/pkg/blacklist"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.JWTBlacklist{})
	database.DB = db
}

func TestAddAndIsBlacklisted(t *testing.T) {
	setupTestDB()

	token := "test-token"
	exp := time.Now().Add(1 * time.Hour)

	if blacklist.IsBlacklisted(token) {
		t.Errorf("expected token not to be blacklisted initially")
	}

	blacklist.Add(token, exp)

	if !blacklist.IsBlacklisted(token) {
		t.Errorf("expected token to be blacklisted after Add")
	}
}

func TestExpiredTokenIsCleaned(t *testing.T) {
	setupTestDB()

	token := "expired-token"
	exp := time.Now().Add(-1 * time.Hour)

	blacklist.Add(token, exp)

	if blacklist.IsBlacklisted(token) {
		t.Errorf("expected expired token not to be considered blacklisted and should be cleaned up")
	}

	// Verify it was deleted from db
	var count int64
	database.DB.Model(&models.JWTBlacklist{}).Where("token = ?", token).Count(&count)
	if count != 0 {
		t.Errorf("expected expired token to be deleted from db")
	}
}
