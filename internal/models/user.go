package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Email           string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash    string         `gorm:"type:varchar(255);not null" json:"-"`
	DefaultCurrency string         `gorm:"type:varchar(10);default:'IDR'" json:"defaultCurrency"`
	FirstDayOfWeek  int            `gorm:"type:int;default:1" json:"firstDayOfWeek"`
	FiscalYearStart int            `gorm:"type:int;default:1" json:"fiscalYearStart"`
	Theme           string         `gorm:"type:varchar(20);default:'dark'" json:"theme"`
	TOTPSecret      *string        `gorm:"type:varchar(255)" json:"-"`
	TOTPEnabled     bool           `gorm:"type:boolean;default:false" json:"totpEnabled"`
	EmailVerifiedAt *time.Time     `json:"emailVerifiedAt"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
