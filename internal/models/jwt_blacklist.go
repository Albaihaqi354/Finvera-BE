package models

import (
	"time"
)

type JWTBlacklist struct {
	Token     string    `gorm:"type:varchar(512);primaryKey" json:"token"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expiresAt"`
}
