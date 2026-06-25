package dto

import (
	"time"

	"github.com/google/uuid"
)

// Auth
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3" example:"johndoe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"secret123"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
}

// Account
type CreateAccountRequest struct {
	Name           string  `json:"name" binding:"required"`
	Type           string  `json:"type" binding:"required,oneof=asset liability"`
	Currency       string  `json:"currency"`
	Icon           string  `json:"icon"`
	Color          string  `json:"color"`
	InitialBalance float64 `json:"initialBalance"`
	Note           string  `json:"note"`
}

type UpdateAccountRequest struct {
	Name           string  `json:"name" binding:"required"`
	Type           string  `json:"type" binding:"required,oneof=asset liability"`
	Currency       string  `json:"currency"`
	Icon           string  `json:"icon"`
	Color          string  `json:"color"`
	InitialBalance float64 `json:"initialBalance"`
	Note           string  `json:"note"`
}

// Category
type CreateCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=income expense transfer"`
	Icon      string `json:"icon"`
	Color     string `json:"color"`
	SortOrder int    `json:"sortOrder"`
	Note      string `json:"note"`
}

type UpdateCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=income expense transfer"`
	Icon      string `json:"icon"`
	Color     string `json:"color"`
	SortOrder int    `json:"sortOrder"`
	Note      string `json:"note"`
}

// Tag
type CreateTagRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

type UpdateTagRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

// Transaction
type CreateTransactionRequest struct {
	Type            string      `json:"type" binding:"required,oneof=income expense transfer"`
	Amount          float64     `json:"amount" binding:"required,gt=0"`
	AccountID       uuid.UUID   `json:"accountId" binding:"required"`
	TargetAccountID *uuid.UUID  `json:"targetAccountId"`
	CategoryID      uuid.UUID   `json:"categoryId" binding:"required"`
	Date            time.Time   `json:"date" binding:"required"`
	Note            string      `json:"note"`
	TagIDs          []uuid.UUID `json:"tagIds"`
}

type UpdateTransactionRequest struct {
	Type            string      `json:"type" binding:"required,oneof=income expense transfer"`
	Amount          float64     `json:"amount" binding:"required,gt=0"`
	AccountID       uuid.UUID   `json:"accountId" binding:"required"`
	TargetAccountID *uuid.UUID  `json:"targetAccountId"`
	CategoryID      uuid.UUID   `json:"categoryId" binding:"required"`
	Date            time.Time   `json:"date" binding:"required"`
	Note            string      `json:"note"`
	TagIDs          []uuid.UUID `json:"tagIds"`
}

// Scheduled Transaction
type CreateScheduledRequest struct {
	Name            string     `json:"name" binding:"required"`
	Type            string     `json:"type" binding:"required,oneof=income expense transfer"`
	Amount          float64    `json:"amount" binding:"required,gt=0"`
	AccountID       uuid.UUID  `json:"accountId" binding:"required"`
	TargetAccountID *uuid.UUID `json:"targetAccountId"`
	CategoryID      uuid.UUID  `json:"categoryId" binding:"required"`
	Note            string     `json:"note"`
	Frequency       string     `json:"frequency" binding:"required,oneof=daily weekly monthly yearly"`
	NextRun         time.Time  `json:"nextRun" binding:"required"`
	IsActive        bool       `json:"isActive"`
}

type UpdateScheduledRequest struct {
	Name            string     `json:"name" binding:"required"`
	Type            string     `json:"type" binding:"required,oneof=income expense transfer"`
	Amount          float64    `json:"amount" binding:"required,gt=0"`
	AccountID       uuid.UUID  `json:"accountId" binding:"required"`
	TargetAccountID *uuid.UUID `json:"targetAccountId"`
	CategoryID      uuid.UUID  `json:"categoryId" binding:"required"`
	Note            string     `json:"note"`
	Frequency       string     `json:"frequency" binding:"required,oneof=daily weekly monthly yearly"`
	NextRun         time.Time  `json:"nextRun" binding:"required"`
	IsActive        bool       `json:"isActive"`
}
