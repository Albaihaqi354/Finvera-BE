package repository

import (
	"finvera-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	GetAvailableForUser(userID uuid.UUID) ([]models.Category, error)
	GetByID(id uuid.UUID) (*models.Category, error)
	Update(category *models.Category) error
	Delete(id uuid.UUID) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetAvailableForUser(userID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	// Get categories where user_id matches OR user_id is null (system defaults)
	// Order by type and sort_order
	err := r.db.Where("user_id = ? OR user_id IS NULL", userID).
		Order("type asc, sort_order asc, created_at asc").
		Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Category{}, "id = ?", id).Error
}
