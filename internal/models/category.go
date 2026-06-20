package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	User      User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ParentID  *uuid.UUID     `gorm:"type:uuid;index" json:"parentId"`
	Parent    *Category      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Type      string         `gorm:"type:varchar(50);not null" json:"type"` // e.g., 'income', 'expense', 'transfer'
	Icon      string         `gorm:"type:varchar(100)" json:"icon"`
	Color     string         `gorm:"type:varchar(20)" json:"color"`
	SortOrder int            `gorm:"type:int;default:0" json:"sortOrder"`
	Note      string         `gorm:"type:varchar(500)" json:"note"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
