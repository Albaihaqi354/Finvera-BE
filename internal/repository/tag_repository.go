package repository

import (
	"finvera-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagRepository interface {
	Create(tag *models.Tag) error
	GetByUserID(userID uuid.UUID) ([]models.Tag, error)
	GetByID(id uuid.UUID) (*models.Tag, error)
	Update(tag *models.Tag) error
	Delete(id uuid.UUID) error
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) GetByUserID(userID uuid.UUID) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Where("user_id = ?", userID).Order("name asc").Find(&tags).Error
	return tags, err
}

func (r *tagRepository) GetByID(id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) Update(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

func (r *tagRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Tag{}, "id = ?", id).Error
}
