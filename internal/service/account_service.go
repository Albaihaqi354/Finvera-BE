package service

import (
	"errors"
	"finvera-be/internal/models"
	"finvera-be/internal/repository"

	"github.com/google/uuid"
)

type AccountService interface {
	CreateAccount(userID uuid.UUID, req CreateAccountRequest) (*models.Account, error)
	GetAccounts(userID uuid.UUID) ([]models.Account, error)
	GetAccountByID(userID, accountID uuid.UUID) (*models.Account, error)
	UpdateAccount(userID, accountID uuid.UUID, req UpdateAccountRequest) (*models.Account, error)
	DeleteAccount(userID, accountID uuid.UUID) error
}

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

// Request DTOs
type CreateAccountRequest struct {
	Name           string  `json:"name" binding:"required"`
	Type           string  `json:"type" binding:"required,oneof=asset liability"`
	Currency       string  `json:"currency"`
	Icon           string  `json:"icon"`
	Color          string  `json:"color"`
	InitialBalance float64 `json:"initialBalance"`
	Note           string  `json:"note"`
}

type UpdateAccountRequest struct {
	Name           string  `json:"name" binding:"required"`
	Type           string  `json:"type" binding:"required,oneof=asset liability"`
	Currency       string  `json:"currency"`
	Icon           string  `json:"icon"`
	Color          string  `json:"color"`
	InitialBalance float64 `json:"initialBalance"`
	Note           string  `json:"note"`
}

func (s *accountService) CreateAccount(userID uuid.UUID, req CreateAccountRequest) (*models.Account, error) {
	currency := req.Currency
	if currency == "" {
		currency = "IDR"
	}

	account := &models.Account{
		UserID:         userID,
		Name:           req.Name,
		Type:           req.Type,
		Currency:       currency,
		Icon:           req.Icon,
		Color:          req.Color,
		InitialBalance: req.InitialBalance,
		Note:           req.Note,
	}

	if err := s.repo.Create(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountService) GetAccounts(userID uuid.UUID) ([]models.Account, error) {
	return s.repo.GetByUserID(userID)
}

func (s *accountService) GetAccountByID(userID, accountID uuid.UUID) (*models.Account, error) {
	account, err := s.repo.GetByID(accountID)
	if err != nil {
		return nil, err
	}
	// Check ownership
	if account.UserID != userID {
		return nil, errors.New("unauthorized: account does not belong to user")
	}
	return account, nil
}

func (s *accountService) UpdateAccount(userID, accountID uuid.UUID, req UpdateAccountRequest) (*models.Account, error) {
	account, err := s.GetAccountByID(userID, accountID)
	if err != nil {
		return nil, err
	}

	account.Name = req.Name
	account.Type = req.Type
	if req.Currency != "" {
		account.Currency = req.Currency
	}
	account.Icon = req.Icon
	account.Color = req.Color
	account.InitialBalance = req.InitialBalance
	account.Note = req.Note

	if err := s.repo.Update(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountService) DeleteAccount(userID, accountID uuid.UUID) error {
	account, err := s.GetAccountByID(userID, accountID)
	if err != nil {
		return err
	}
	return s.repo.Delete(account.ID)
}
