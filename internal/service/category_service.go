package service

import (
	"errors"
	"finvera-be/internal/dto"
	"finvera-be/internal/models"
	"finvera-be/internal/repository"

	"github.com/google/uuid"
)

type CategoryService interface {
	CreateCategory(userID uuid.UUID, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategories(userID uuid.UUID, page, limit int) ([]dto.CategoryResponse, int64, error)
	GetCategoryByID(userID, categoryID uuid.UUID) (*dto.CategoryResponse, error)
	UpdateCategory(userID, categoryID uuid.UUID, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(userID, categoryID uuid.UUID) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

// Mapper
func mapCategoryToResponse(category *models.Category) *dto.CategoryResponse {
	if category == nil {
		return nil
	}
	return &dto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Type:      category.Type,
		Icon:      category.Icon,
		Color:     category.Color,
		SortOrder: category.SortOrder,
		Note:      category.Note,
		IsSystem:  category.UserID == uuid.Nil,
	}
}

func (s *categoryService) CreateCategory(userID uuid.UUID, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
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

	return mapCategoryToResponse(category), nil
}

func (s *categoryService) GetCategories(userID uuid.UUID, page, limit int) ([]dto.CategoryResponse, int64, error) {
	categories, total, err := s.repo.GetAvailableForUser(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.CategoryResponse
	for _, cat := range categories {
		responses = append(responses, *mapCategoryToResponse(&cat))
	}

	return responses, total, nil
}

func (s *categoryService) GetCategoryByID(userID, categoryID uuid.UUID) (*dto.CategoryResponse, error) {
	category, err := s.repo.GetByID(categoryID)
	if err != nil {
		return nil, err
	}

	// Ensure the user owns this category or it's a global category
	if category.UserID != uuid.Nil && category.UserID != userID {
		return nil, errors.New("unauthorized: you do not have permission to access this category")
	}

	return mapCategoryToResponse(category), nil
}

func (s *categoryService) UpdateCategory(userID, categoryID uuid.UUID, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
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

	return mapCategoryToResponse(category), nil
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
