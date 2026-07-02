package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduledTransaction struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index:idx_scheduled_user_deleted,composite:user" json:"userId"`
	User            User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Name            string         `gorm:"type:varchar(100);not null" json:"name"`
	Type            string         `gorm:"type:varchar(50);not null" json:"type"` // e.g., 'income', 'expense', 'transfer'
	Amount          float64        `gorm:"type:decimal(15,2);not null" json:"amount"`
	AccountID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"accountId"`
	Account         Account        `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	TargetAccountID *uuid.UUID     `gorm:"type:uuid;index" json:"targetAccountId"`
	TargetAccount   *Account       `gorm:"foreignKey:TargetAccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	CategoryID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"categoryId"`
	Category        Category       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Note            string         `gorm:"type:varchar(500)" json:"note"`
	Frequency       string         `gorm:"type:varchar(20);not null" json:"frequency"` // e.g., 'daily', 'weekly', 'monthly', 'yearly'
	NextRun         time.Time      `gorm:"not null" json:"nextRun"`
	IsActive        bool           `gorm:"type:boolean;default:true" json:"isActive"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index:idx_scheduled_user_deleted,composite:deleted" json:"-"`
}
