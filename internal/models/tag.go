package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagGroup struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	SortOrder int       `gorm:"type:int;default:0" json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
}

type Tag struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	User      User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	GroupID   *uuid.UUID     `gorm:"type:uuid;index" json:"groupId"`
	Group     *TagGroup      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Color     string         `gorm:"type:varchar(20)" json:"color"`
	SortOrder int            `gorm:"type:int;default:0" json:"sortOrder"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
