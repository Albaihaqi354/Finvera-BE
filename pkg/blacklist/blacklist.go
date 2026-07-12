package blacklist

import (
	"log"
	"time"

	"finvera-be/internal/database"
	"finvera-be/internal/models"
)

func init() {
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			Cleanup()
		}
	}()
}

// Add adds a token to the blacklist with its expiration time
func Add(token string, exp time.Time) {
	if database.DB == nil {
		return
	}
	item := models.JWTBlacklist{
		Token:     token,
		ExpiresAt: exp,
	}
	if err := database.DB.Create(&item).Error; err != nil {
		log.Printf("Failed to add token to blacklist: %v", err)
	}
}

// IsBlacklisted checks if a token is in the blacklist and not expired
func IsBlacklisted(token string) bool {
	if database.DB == nil {
		return false
	}
	var item models.JWTBlacklist
	if err := database.DB.First(&item, "token = ?", token).Error; err != nil {
		return false
	}

	if time.Now().After(item.ExpiresAt) {
		// Clean up expired token
		database.DB.Delete(&item)
		return false
	}

	return true
}

// Cleanup removes expired tokens from the DB
func Cleanup() {
	if database.DB == nil {
		return
	}
	database.DB.Where("expires_at < ?", time.Now()).Delete(&models.JWTBlacklist{})
}
