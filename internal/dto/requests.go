package dto

import (
	"time"

	"github.com/google/uuid"
)

// ── Auth ──────────────────────────────────────────────────────────────────────

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email"    binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email"    binding:"required,email,max=255"`
}

// ── Account ───────────────────────────────────────────────────────────────────

type CreateAccountRequest struct {
	Name           string  `json:"name"           binding:"required,max=100"`
	Type           string  `json:"type"           binding:"required,oneof=asset liability"`
	Currency       string  `json:"currency"       binding:"max=10"`
	Icon           string  `json:"icon"           binding:"max=100"`
	Color          string  `json:"color"          binding:"max=20"`
	InitialBalance float64 `json:"initialBalance"`
	Note           string  `json:"note"           binding:"max=500"`
}

type UpdateAccountRequest struct {
	Name           string  `json:"name"           binding:"required,max=100"`
	Type           string  `json:"type"           binding:"required,oneof=asset liability"`
	Currency       string  `json:"currency"       binding:"max=10"`
	Icon           string  `json:"icon"           binding:"max=100"`
	Color          string  `json:"color"          binding:"max=20"`
	InitialBalance float64 `json:"initialBalance"`
	Note           string  `json:"note"           binding:"max=500"`
}

// ── Category ──────────────────────────────────────────────────────────────────

type CreateCategoryRequest struct {
	Name      string `json:"name"      binding:"required,max=100"`
	Type      string `json:"type"      binding:"required,oneof=income expense transfer"`
	Icon      string `json:"icon"      binding:"max=100"`
	Color     string `json:"color"     binding:"max=20"`
	SortOrder int    `json:"sortOrder"`
	Note      string `json:"note"      binding:"max=500"`
}

type UpdateCategoryRequest struct {
	Name      string `json:"name"      binding:"required,max=100"`
	Type      string `json:"type"      binding:"required,oneof=income expense transfer"`
	Icon      string `json:"icon"      binding:"max=100"`
	Color     string `json:"color"     binding:"max=20"`
	SortOrder int    `json:"sortOrder"`
	Note      string `json:"note"      binding:"max=500"`
}

// ── Tag ───────────────────────────────────────────────────────────────────────

type CreateTagRequest struct {
	Name  string `json:"name"  binding:"required,max=100"`
	Color string `json:"color" binding:"max=20"`
}

type UpdateTagRequest struct {
	Name  string `json:"name"  binding:"required,max=100"`
	Color string `json:"color" binding:"max=20"`
}

// ── Transaction ───────────────────────────────────────────────────────────────
// TransactionRequest is shared by both Create and Update to eliminate duplication (DRY).
// The handler maps this to the service layer.
type TransactionRequest struct {
	Type            string      `json:"type"            binding:"required,oneof=income expense transfer"`
	Amount          float64     `json:"amount"          binding:"required,gt=0"`
	AccountID       uuid.UUID   `json:"accountId"       binding:"required"`
	TargetAccountID *uuid.UUID  `json:"targetAccountId"`
	CategoryID      uuid.UUID   `json:"categoryId"      binding:"required"`
	Date            time.Time   `json:"date"            binding:"required"`
	Note            string      `json:"note"            binding:"max=500"`
	TagIDs          []uuid.UUID `json:"tagIds"`
}

// CreateTransactionRequest and UpdateTransactionRequest are type aliases for
// TransactionRequest so existing swagger docs and callers are unaffected.
type CreateTransactionRequest = TransactionRequest
type UpdateTransactionRequest = TransactionRequest

// ── Scheduled Transaction ─────────────────────────────────────────────────────

type CreateScheduledRequest struct {
	Name            string     `json:"name"            binding:"required,max=100"`
	Type            string     `json:"type"            binding:"required,oneof=income expense transfer"`
	Amount          float64    `json:"amount"          binding:"required,gt=0"`
	AccountID       uuid.UUID  `json:"accountId"       binding:"required"`
	TargetAccountID *uuid.UUID `json:"targetAccountId"`
	CategoryID      uuid.UUID  `json:"categoryId"      binding:"required"`
	Note            string     `json:"note"            binding:"max=500"`
	Frequency       string     `json:"frequency"       binding:"required,oneof=daily weekly monthly yearly"`
	NextRun         time.Time  `json:"nextRun"         binding:"required"`
	IsActive        bool       `json:"isActive"`
}

type UpdateScheduledRequest struct {
	Name            string     `json:"name"            binding:"required,max=100"`
	Type            string     `json:"type"            binding:"required,oneof=income expense transfer"`
	Amount          float64    `json:"amount"          binding:"required,gt=0"`
	AccountID       uuid.UUID  `json:"accountId"       binding:"required"`
	TargetAccountID *uuid.UUID `json:"targetAccountId"`
	CategoryID      uuid.UUID  `json:"categoryId"      binding:"required"`
	Note            string     `json:"note"            binding:"max=500"`
	Frequency       string     `json:"frequency"       binding:"required,oneof=daily weekly monthly yearly"`
	NextRun         time.Time  `json:"nextRun"         binding:"required"`
	IsActive        bool       `json:"isActive"`
}
