package repository

import (
	"finvera-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account *models.Account) error
	GetByUserID(userID uuid.UUID) ([]models.Account, error)
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

func (r *accountRepository) GetByUserID(userID uuid.UUID) ([]models.Account, error) {
	var accounts []models.Account
	// Use order by sort_order ascending
	err := r.db.Where("user_id = ?", userID).Order("sort_order asc, created_at asc").Find(&accounts).Error
	return accounts, err
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
