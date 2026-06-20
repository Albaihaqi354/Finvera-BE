package service

import (
	"errors"
	"finvera-be/internal/models"
	"finvera-be/internal/repository"

	"github.com/google/uuid"
)

type TagService interface {
	CreateTag(userID uuid.UUID, req CreateTagRequest) (*models.Tag, error)
	GetTags(userID uuid.UUID) ([]models.Tag, error)
	GetTagByID(userID, tagID uuid.UUID) (*models.Tag, error)
	UpdateTag(userID, tagID uuid.UUID, req UpdateTagRequest) (*models.Tag, error)
	DeleteTag(userID, tagID uuid.UUID) error
}

type tagService struct {
	repo repository.TagRepository
}

func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{repo: repo}
}

// Request DTOs
type CreateTagRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

type UpdateTagRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

func (s *tagService) CreateTag(userID uuid.UUID, req CreateTagRequest) (*models.Tag, error) {
	tag := &models.Tag{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	}

	if err := s.repo.Create(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *tagService) GetTags(userID uuid.UUID) ([]models.Tag, error) {
	return s.repo.GetByUserID(userID)
}

func (s *tagService) GetTagByID(userID, tagID uuid.UUID) (*models.Tag, error) {
	tag, err := s.repo.GetByID(tagID)
	if err != nil {
		return nil, err
	}
	if tag.UserID != userID {
		return nil, errors.New("unauthorized: tag does not belong to user")
	}
	return tag, nil
}

func (s *tagService) UpdateTag(userID, tagID uuid.UUID, req UpdateTagRequest) (*models.Tag, error) {
	tag, err := s.GetTagByID(userID, tagID)
	if err != nil {
		return nil, err
	}

	tag.Name = req.Name
	tag.Color = req.Color

	if err := s.repo.Update(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *tagService) DeleteTag(userID, tagID uuid.UUID) error {
	tag, err := s.GetTagByID(userID, tagID)
	if err != nil {
		return err
	}
	return s.repo.Delete(tag.ID)
}
