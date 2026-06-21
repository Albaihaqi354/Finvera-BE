package repository

import (
	"finvera-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account *models.Account) error
	GetByUserID(userID uuid.UUID, page, limit int) ([]models.Account, int64, error)
	GetByID(id uuid.UUID) (*models.Account, error)
	Update(account *models.Account) error
	Delete(id uuid.UUID) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepository) GetByUserID(userID uuid.UUID, page, limit int) ([]models.Account, int64, error) {
	var accounts []models.Account
	var total int64

	query := r.db.Where("user_id = ?", userID)
	
	// Get total count
	err := query.Model(&models.Account{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	offset := (page - 1) * limit
	err = query.Order("sort_order asc, created_at asc").Offset(offset).Limit(limit).Find(&accounts).Error
	
	return accounts, total, err
}

func (r *accountRepository) GetByID(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	err := r.db.First(&account, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) Update(account *models.Account) error {
	return r.db.Save(account).Error
}

func (r *accountRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Account{}, "id = ?", id).Error
}
