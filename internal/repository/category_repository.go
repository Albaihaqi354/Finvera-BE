package repository

import (
	"finvera-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	GetAvailableForUser(userID uuid.UUID, page, limit int) ([]models.Category, int64, error)
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

func (r *categoryRepository) GetAvailableForUser(userID uuid.UUID, page, limit int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	query := r.db.Where("user_id = ? OR user_id IS NULL", userID)

	// Get total count
	err := query.Model(&models.Category{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	offset := (page - 1) * limit
	err = query.Order("type asc, sort_order asc, created_at asc").Offset(offset).Limit(limit).Find(&categories).Error
	
	return categories, total, err
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
