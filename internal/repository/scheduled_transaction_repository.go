package repository

import (
	"finvera-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduledTransactionRepository interface {
	Create(scheduled *models.ScheduledTransaction) error
	GetByUserID(userID uuid.UUID) ([]models.ScheduledTransaction, error)
	GetByID(id uuid.UUID) (*models.ScheduledTransaction, error)
	Update(scheduled *models.ScheduledTransaction) error
	Delete(id uuid.UUID) error
}

type scheduledTransactionRepository struct {
	db *gorm.DB
}

func NewScheduledTransactionRepository(db *gorm.DB) ScheduledTransactionRepository {
	return &scheduledTransactionRepository{db: db}
}

func (r *scheduledTransactionRepository) Create(scheduled *models.ScheduledTransaction) error {
	return r.db.Create(scheduled).Error
}

func (r *scheduledTransactionRepository) GetByUserID(userID uuid.UUID) ([]models.ScheduledTransaction, error) {
	var scheduleds []models.ScheduledTransaction
	err := r.db.Where("user_id = ?", userID).
		Preload("Account").
		Preload("TargetAccount").
		Preload("Category").
		Order("next_run asc").
		Find(&scheduleds).Error
	return scheduleds, err
}

func (r *scheduledTransactionRepository) GetByID(id uuid.UUID) (*models.ScheduledTransaction, error) {
	var scheduled models.ScheduledTransaction
	err := r.db.Preload("Account").
		Preload("TargetAccount").
		Preload("Category").
		First(&scheduled, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &scheduled, nil
}

func (r *scheduledTransactionRepository) Update(scheduled *models.ScheduledTransaction) error {
	return r.db.Save(scheduled).Error
}

func (r *scheduledTransactionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ScheduledTransaction{}, "id = ?", id).Error
}
