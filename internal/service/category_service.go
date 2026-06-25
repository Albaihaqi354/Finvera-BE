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
	GetPresetCategories() ([]dto.PresetCategoryGroupResponse, error)
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
		ID:         category.ID,
		ParentID:   category.ParentID,
		Name:       category.Name,
		Type:       category.Type,
		Icon:       category.Icon,
		Color:      category.Color,
		ColorClass: category.ColorClass,
		SortOrder:  category.SortOrder,
		Note:       category.Note,
		IsSystem:   category.UserID == nil,
	}
}

func (s *categoryService) CreateCategory(userID uuid.UUID, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := &models.Category{
		UserID:    &userID,
		Name:      req.Name,
		Type:      req.Type,
		Icon:      req.Icon,
		Color:      req.Color,
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

func (s *categoryService) GetPresetCategories() ([]dto.PresetCategoryGroupResponse, error) {
	categories, err := s.repo.GetPresetCategories()
	if err != nil {
		return nil, err
	}

	parentMap := make(map[uuid.UUID]*dto.PresetCategoryGroupResponse)
	var order []*dto.PresetCategoryGroupResponse

	// First pass: create parents
	for i := range categories {
		if categories[i].ParentID == nil {
			resp := &dto.PresetCategoryGroupResponse{
				CategoryResponse: *mapCategoryToResponse(&categories[i]),
				Children:         []dto.CategoryResponse{},
			}
			parentMap[categories[i].ID] = resp
			order = append(order, resp)
		}
	}

	// Second pass: add children
	for i := range categories {
		if categories[i].ParentID != nil {
			if parent, exists := parentMap[*categories[i].ParentID]; exists {
				parent.Children = append(parent.Children, *mapCategoryToResponse(&categories[i]))
			}
		}
	}

	var result []dto.PresetCategoryGroupResponse
	for _, p := range order {
		result = append(result, *p)
	}

	return result, nil
}

func (s *categoryService) GetCategoryByID(userID, categoryID uuid.UUID) (*dto.CategoryResponse, error) {
	category, err := s.repo.GetByID(categoryID)
	if err != nil {
		return nil, err
	}

	// Ensure the user owns this category or it's a global category
	if category.UserID != nil && *category.UserID != userID {
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
	if category.UserID == nil {
		return nil, errors.New("forbidden: cannot modify default system category")
	}

	// Check ownership
	if *category.UserID != userID {
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
	if category.UserID == nil {
		return errors.New("forbidden: cannot delete default system category")
	}

	// Check ownership
	if *category.UserID != userID {
		return errors.New("unauthorized: you do not own this category")
	}

	return s.repo.Delete(categoryID)
}
