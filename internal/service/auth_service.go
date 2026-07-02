package service

import (
	"errors"
	"strings"

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
	// Sanitize non-password fields only
	username = strings.TrimSpace(username)
	email = strings.ToLower(strings.TrimSpace(email))

	// Validate username length after trimming
	if len(username) < 3 {
		return nil, errors.New("username must be at least 3 characters")
	}

	// Check if user exists
	existingUser, _ := s.userRepo.GetUserByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}
	existingEmail, _ := s.userRepo.GetUserByEmail(email)
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	// Password is hashed as-is — do NOT trim, preserving intentional leading/trailing spaces
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err = s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(identifier, password string, cfg *config.Config) (string, error) {
	// Trim identifier (username/email) but NOT the password
	identifier = strings.TrimSpace(identifier)

	if identifier == "" {
		return "", errors.New("identifier must not be empty")
	}
	if password == "" {
		return "", errors.New("password must not be empty")
	}

	var user *models.User
	var err error

	if strings.Contains(identifier, "@") {
		user, err = s.userRepo.GetUserByEmail(strings.ToLower(identifier))
	} else {
		user, err = s.userRepo.GetUserByUsername(identifier)
	}

	if err != nil || user == nil {
		// Return a generic error to prevent user enumeration
		return "", errors.New("invalid credentials")
	}

	// Compare password as received — no trimming, exact match
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, cfg.JWTSecret, cfg.JWTIssuer)
	if err != nil {
		return "", err
	}

	return token, nil
}
