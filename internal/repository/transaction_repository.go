package repository

import (
	"finvera-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionFilter struct {
	StartDate string
	EndDate   string
	Type      string
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

	query := r.db.Where("user_id = ?", userID)

	if filter.StartDate != "" {
		query = query.Where("date >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("date <= ?", filter.EndDate)
	}
	if filter.Type != "" && filter.Type != "All Types" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.AccountID != "" && filter.AccountID != "All Accounts" {
		query = query.Where("account_id = ? OR target_account_id = ?", filter.AccountID, filter.AccountID)
	}
	if filter.Search != "" {
		// Search in note. Later can join category and account.
		query = query.Where("note ILIKE ?", "%"+filter.Search+"%")
	}

	err := query.Model(&models.Transaction{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Preload("Account").
		Preload("TargetAccount").
		Preload("Category").
		Preload("Tags").
		Order("date desc, created_at desc").
		Offset(offset).Limit(limit).
		Find(&transactions).Error
		
	return transactions, total, err
}

func (r *transactionRepository) GetByID(id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Account").
		Preload("TargetAccount").
		Preload("Category").
		Preload("Tags").
		First(&transaction, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
