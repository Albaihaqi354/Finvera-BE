package repository

import (
	"finvera-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionFilter holds validated, sanitized filter values from the handler.
// All string fields have been checked before reaching the repository.
type TransactionFilter struct {
	StartDate string
	EndDate   string
	// Type must be one of: income, expense, transfer — validated in handler
	Type string
	// AccountID must be a valid UUID string — validated in handler
	AccountID string
	Search    string
}

type TransactionRepository interface {
	CreateWithTx(tx *gorm.DB, transaction *models.Transaction) error
	DeleteWithTx(tx *gorm.DB, id uuid.UUID) error
	GetByUserID(userID uuid.UUID, page, limit int, filter TransactionFilter) ([]models.Transaction, int64, error)
	GetByID(id uuid.UUID) (*models.Transaction, error)
	GetDB() *gorm.DB
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *transactionRepository) CreateWithTx(tx *gorm.DB, transaction *models.Transaction) error {
	return tx.Create(transaction).Error
}

func (r *transactionRepository) DeleteWithTx(tx *gorm.DB, id uuid.UUID) error {
	return tx.Delete(&models.Transaction{}, "id = ?", id).Error
}

func (r *transactionRepository) GetByUserID(userID uuid.UUID, page, limit int, filter TransactionFilter) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	// Always scope by the authenticated user's ID first
	query := r.db.Where("transactions.user_id = ?", userID)

	if filter.StartDate != "" {
		query = query.Where("transactions.date >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("transactions.date <= ?", filter.EndDate)
	}
	// Type is already validated as oneof=income,expense,transfer by the handler binding
	if filter.Type != "" {
		query = query.Where("transactions.type = ?", filter.Type)
	}
	// AccountID is a valid UUID string — safe to use as a parameterized query
	if filter.AccountID != "" {
		query = query.Where("transactions.account_id = ? OR transactions.target_account_id = ?", filter.AccountID, filter.AccountID)
	}
	// Search uses ILIKE with parameterized value — safe against injection
	if filter.Search != "" {
		query = query.Where("transactions.note ILIKE ?", "%"+filter.Search+"%")
	}

	if err := query.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.
		Preload("Account").
		Preload("TargetAccount").
		Preload("Category").
		Preload("Tags").
		Order("transactions.date desc, transactions.created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *transactionRepository) GetByID(id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.
		Preload("Account").
		Preload("TargetAccount").
		Preload("Category").
		Preload("Tags").
		First(&transaction, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
