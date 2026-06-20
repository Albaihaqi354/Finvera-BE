package repository

import (
	"finvera-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateWithTx(tx *gorm.DB, transaction *models.Transaction) error
	DeleteWithTx(tx *gorm.DB, id uuid.UUID) error
	GetByUserID(userID uuid.UUID) ([]models.Transaction, error)
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

func (r *transactionRepository) GetByUserID(userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	// Preload Account, Category, TargetAccount, and Tags
	err := r.db.Where("user_id = ?", userID).
		Preload("Account").
		Preload("TargetAccount").
		Preload("Category").
		Preload("Tags").
		Order("date desc, created_at desc").
		Find(&transactions).Error
	return transactions, err
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
