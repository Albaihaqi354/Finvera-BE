package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	DefaultCurrency string    `json:"defaultCurrency"`
	Theme           string    `json:"theme"`
}

type AccountResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	Currency       string    `json:"currency"`
	Icon           string    `json:"icon"`
	Color          string    `json:"color"`
	Balance        float64   `json:"balance"`
	InitialBalance float64   `json:"initialBalance"`
	Note           string    `json:"note"`
}

type CategoryResponse struct {
	ID         uuid.UUID  `json:"id"`
	ParentID   *uuid.UUID `json:"parentId,omitempty"`
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	Icon       string     `json:"icon"`
	Color      string     `json:"color"`
	ColorClass string     `json:"colorClass"`
	SortOrder  int        `json:"sortOrder"`
	Note       string     `json:"note"`
	IsSystem   bool       `json:"isSystem"`
}

type PresetCategoryGroupResponse struct {
	CategoryResponse
	Children []CategoryResponse `json:"children"`
}

type TagResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Color string    `json:"color"`
}

type TransactionResponse struct {
	ID              uuid.UUID         `json:"id"`
	Type            string            `json:"type"`
	Amount          float64           `json:"amount"`
	Currency        string            `json:"currency"`
	Account         AccountResponse   `json:"account"`
	TargetAccount   *AccountResponse  `json:"targetAccount,omitempty"`
	Category        CategoryResponse  `json:"category"`
	Date            time.Time         `json:"date"`
	Note            string            `json:"note"`
	Tags            []TagResponse     `json:"tags"`
	CreatedAt       time.Time         `json:"createdAt"`
}

type ScheduledTransactionResponse struct {
	ID              uuid.UUID         `json:"id"`
	Name            string            `json:"name"`
	Type            string            `json:"type"`
	Amount          float64           `json:"amount"`
	Currency        string            `json:"currency"`
	Account         AccountResponse   `json:"account"`
	TargetAccount   *AccountResponse  `json:"targetAccount,omitempty"`
	Category        CategoryResponse  `json:"category"`
	Note            string            `json:"note"`
	Frequency       string            `json:"frequency"`
	NextRun         time.Time         `json:"nextRun"`
	LastRun         *time.Time        `json:"lastRun,omitempty"`
	IsActive        bool              `json:"isActive"`
}
