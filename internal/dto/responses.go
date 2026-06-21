package dto

import (
	"time"

	"github.com/google/uuid"
)

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
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Icon      string    `json:"icon"`
	Color     string    `json:"color"`
	SortOrder int       `json:"sortOrder"`
	Note      string    `json:"note"`
	IsSystem  bool      `json:"isSystem"`
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
	Account         AccountResponse   `json:"account"`
	TargetAccount   *AccountResponse  `json:"targetAccount,omitempty"`
	Category        CategoryResponse  `json:"category"`
	Note            string            `json:"note"`
	Frequency       string            `json:"frequency"`
	NextRun         time.Time         `json:"nextRun"`
	LastRun         *time.Time        `json:"lastRun,omitempty"`
	IsActive        bool              `json:"isActive"`
}
