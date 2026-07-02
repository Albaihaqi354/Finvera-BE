package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index:idx_transactions_user_deleted,composite:user" json:"userId"`
	User            User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	AccountID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"accountId"`
	Account         Account        `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	TargetAccountID *uuid.UUID     `gorm:"type:uuid;index" json:"targetAccountId"`
	TargetAccount   *Account       `gorm:"foreignKey:TargetAccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	CategoryID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"categoryId"`
	Category        Category       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Type            string         `gorm:"type:varchar(50);not null" json:"type"` // income, expense, transfer
	Amount          float64        `gorm:"type:decimal(15,2);not null" json:"amount"`
	Currency        string         `gorm:"type:varchar(10);not null;default:'IDR'" json:"currency"`
	Date            time.Time      `gorm:"type:timestamp with time zone;not null;index" json:"date"`
	Note            string         `gorm:"type:varchar(500)" json:"note"`
	GeoLat          *float64       `gorm:"type:decimal(10,7)" json:"geoLat"`
	GeoLng          *float64       `gorm:"type:decimal(10,7)" json:"geoLng"`
	GeoName         string         `gorm:"type:varchar(255)" json:"geoName"`
	Tags            []Tag          `gorm:"many2many:transaction_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tags"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index:idx_transactions_user_deleted,composite:deleted" json:"-"`
}
