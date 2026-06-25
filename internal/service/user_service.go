package service

import (
	"errors"
	"finvera-be/internal/dto"
	"finvera-be/internal/repository"

	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(userID uuid.UUID) (*dto.UserResponse, error)
	UpdateProfile(userID uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.repo.GetUserByID(userID.String())
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *userService) UpdateProfile(userID uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetUserByID(userID.String())
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if username is already taken by someone else
	if req.Username != user.Username {
		existingUser, err := s.repo.GetUserByUsername(req.Username)
		if err != nil {
			return nil, err
		}
		if existingUser != nil {
			return nil, errors.New("username already exists")
		}
	}

	// Check if email is already taken by someone else
	if req.Email != user.Email {
		existingEmail, err := s.repo.GetUserByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		if existingEmail != nil {
			return nil, errors.New("email already exists")
		}
	}

	user.Username = req.Username
	user.Email = req.Email

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
