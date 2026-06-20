package service

import (
	"errors"

	"finvera-be/internal/config"
	"finvera-be/internal/models"
	"finvera-be/internal/repository"
	"finvera-be/pkg/utils"
)

type AuthService interface {
	Register(username, email, password string) (*models.User, error)
	Login(username, password string, cfg *config.Config) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

func (s *authService) Register(username, email, password string) (*models.User, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.GetUserByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}
	existingEmail, _ := s.userRepo.GetUserByEmail(email)
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(username, password string, cfg *config.Config) (string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, cfg.JWTSecret, cfg.JWTIssuer)
	if err != nil {
		return "", err
	}

	return token, nil
}
