package service

import (
	"errors"
	"finvera-be/internal/models"
	"finvera-be/internal/repository"

	"github.com/google/uuid"
)

type CategoryService interface {
	CreateCategory(userID uuid.UUID, req CreateCategoryRequest) (*models.Category, error)
	GetCategories(userID uuid.UUID) ([]models.Category, error)
	GetCategoryByID(userID, categoryID uuid.UUID) (*models.Category, error)
	UpdateCategory(userID, categoryID uuid.UUID, req UpdateCategoryRequest) (*models.Category, error)
	DeleteCategory(userID, categoryID uuid.UUID) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

// Request DTOs
type CreateCategoryRequest struct {
	Name      string  `json:"name" binding:"required"`
	Type      string  `json:"type" binding:"required,oneof=income expense transfer"`
	Icon      string  `json:"icon"`
	Color     string  `json:"color"`
	SortOrder int     `json:"sortOrder"`
	Note      string  `json:"note"`
}

type UpdateCategoryRequest struct {
	Name      string  `json:"name" binding:"required"`
	Type      string  `json:"type" binding:"required,oneof=income expense transfer"`
	Icon      string  `json:"icon"`
	Color     string  `json:"color"`
	SortOrder int     `json:"sortOrder"`
	Note      string  `json:"note"`
}

func (s *categoryService) CreateCategory(userID uuid.UUID, req CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		Icon:      req.Icon,
		Color:     req.Color,
		SortOrder: req.SortOrder,
		Note:      req.Note,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetCategories(userID uuid.UUID) ([]models.Category, error) {
	return s.repo.GetAvailableForUser(userID)
}

func (s *categoryService) GetCategoryByID(userID, categoryID uuid.UUID) (*models.Category, error) {
	category, err := s.repo.GetByID(categoryID)
	if err != nil {
		return nil, err
	}

	// Ensure the user owns this category or it's a global category
	if category.UserID != uuid.Nil && category.UserID != userID {
		return nil, errors.New("unauthorized: you do not have permission to access this category")
	}

	return category, nil
}

func (s *categoryService) UpdateCategory(userID, categoryID uuid.UUID, req UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.repo.GetByID(categoryID)
	if err != nil {
		return nil, err
	}

	// Check if this is a default system category
	if category.UserID == uuid.Nil {
		return nil, errors.New("forbidden: cannot modify default system category")
	}

	// Check ownership
	if category.UserID != userID {
		return nil, errors.New("unauthorized: you do not own this category")
	}

	category.Name = req.Name
	category.Type = req.Type
	category.Icon = req.Icon
	category.Color = req.Color
	category.SortOrder = req.SortOrder
	category.Note = req.Note

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(userID, categoryID uuid.UUID) error {
	category, err := s.repo.GetByID(categoryID)
	if err != nil {
		return err
	}

	// Check if this is a default system category
	if category.UserID == uuid.Nil {
		return errors.New("forbidden: cannot delete default system category")
	}

	// Check ownership
	if category.UserID != userID {
		return errors.New("unauthorized: you do not own this category")
	}

	return s.repo.Delete(categoryID)
}
