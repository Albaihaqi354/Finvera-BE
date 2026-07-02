package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null;index:idx_accounts_user_deleted,composite:user" json:"userId"`
	User             User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ParentID         *uuid.UUID     `gorm:"type:uuid;index" json:"parentId"`
	Parent           *Account       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Name             string         `gorm:"type:varchar(100);not null" json:"name"`
	Type             string         `gorm:"type:varchar(50);not null" json:"type"` // e.g., 'cash', 'credit', 'debit'
	Currency         string         `gorm:"type:varchar(10);not null;default:'IDR'" json:"currency"`
	Icon             string         `gorm:"type:varchar(100)" json:"icon"`
	Color            string         `gorm:"type:varchar(20)" json:"color"`
	Balance          float64        `gorm:"type:decimal(15,2);default:0" json:"balance"`
	InitialBalance   float64        `gorm:"type:decimal(15,2);default:0" json:"initialBalance"`
	StatementDay     *int           `gorm:"type:int" json:"statementDay"`
	IsHidden         bool           `gorm:"type:boolean;default:false" json:"isHidden"`
	SortOrder        int            `gorm:"type:int;default:0" json:"sortOrder"`
	Note             string         `gorm:"type:varchar(500)" json:"note"`
	LastReconciledAt *time.Time     `json:"lastReconciledAt"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index:idx_accounts_user_deleted,composite:deleted" json:"-"`
}
