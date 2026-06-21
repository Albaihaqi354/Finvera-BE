package service

import (
	"errors"
	"finvera-be/internal/dto"
	"finvera-be/internal/models"
	"finvera-be/internal/repository"

	"github.com/google/uuid"
)

type TagService interface {
	CreateTag(userID uuid.UUID, req dto.CreateTagRequest) (*dto.TagResponse, error)
	GetTags(userID uuid.UUID, page, limit int) ([]dto.TagResponse, int64, error)
	GetTagByID(userID, tagID uuid.UUID) (*dto.TagResponse, error)
	UpdateTag(userID, tagID uuid.UUID, req dto.UpdateTagRequest) (*dto.TagResponse, error)
	DeleteTag(userID, tagID uuid.UUID) error
}

type tagService struct {
	repo repository.TagRepository
}

func NewTagService(repo repository.TagRepository) TagService {
	return &tagService{repo: repo}
}

// Mapper
func mapTagToResponse(tag *models.Tag) *dto.TagResponse {
	if tag == nil {
		return nil
	}
	return &dto.TagResponse{
		ID:    tag.ID,
		Name:  tag.Name,
		Color: tag.Color,
	}
}

func (s *tagService) CreateTag(userID uuid.UUID, req dto.CreateTagRequest) (*dto.TagResponse, error) {
	tag := &models.Tag{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	}

	if err := s.repo.Create(tag); err != nil {
		return nil, err
	}

	return mapTagToResponse(tag), nil
}

func (s *tagService) GetTags(userID uuid.UUID, page, limit int) ([]dto.TagResponse, int64, error) {
	tags, total, err := s.repo.GetByUserID(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.TagResponse
	for _, t := range tags {
		responses = append(responses, *mapTagToResponse(&t))
	}

	return responses, total, nil
}

func (s *tagService) GetTagByID(userID, tagID uuid.UUID) (*dto.TagResponse, error) {
	tag, err := s.repo.GetByID(tagID)
	if err != nil {
		return nil, err
	}
	if tag.UserID != userID {
		return nil, errors.New("unauthorized: tag does not belong to user")
	}
	return mapTagToResponse(tag), nil
}

func (s *tagService) UpdateTag(userID, tagID uuid.UUID, req dto.UpdateTagRequest) (*dto.TagResponse, error) {
	tag, err := s.repo.GetByID(tagID)
	if err != nil {
		return nil, err
	}
	if tag.UserID != userID {
		return nil, errors.New("unauthorized: tag does not belong to user")
	}

	tag.Name = req.Name
	tag.Color = req.Color

	if err := s.repo.Update(tag); err != nil {
		return nil, err
	}

	return mapTagToResponse(tag), nil
}

func (s *tagService) DeleteTag(userID, tagID uuid.UUID) error {
	tag, err := s.repo.GetByID(tagID)
	if err != nil {
		return err
	}
	if tag.UserID != userID {
		return errors.New("unauthorized: tag does not belong to user")
	}
	return s.repo.Delete(tag.ID)
}
