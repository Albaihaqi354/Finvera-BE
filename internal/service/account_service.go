package service

import (
	"errors"
	"finvera-be/internal/dto"
	"finvera-be/internal/models"
	"finvera-be/internal/repository"

	"github.com/google/uuid"
)

type AccountService interface {
	CreateAccount(userID uuid.UUID, req dto.CreateAccountRequest) (*dto.AccountResponse, error)
	GetAccounts(userID uuid.UUID, page, limit int) ([]dto.AccountResponse, int64, error)
	GetAccountByID(userID, accountID uuid.UUID) (*dto.AccountResponse, error)
	UpdateAccount(userID, accountID uuid.UUID, req dto.UpdateAccountRequest) (*dto.AccountResponse, error)
	DeleteAccount(userID, accountID uuid.UUID) error
}

type accountService struct {
	repo     repository.AccountRepository
	userRepo repository.UserRepository
}

func NewAccountService(repo repository.AccountRepository, userRepo repository.UserRepository) AccountService {
	return &accountService{repo: repo, userRepo: userRepo}
}

// Mapper
func mapAccountToResponse(account *models.Account) *dto.AccountResponse {
	if account == nil {
		return nil
	}
	return &dto.AccountResponse{
		ID:             account.ID,
		Name:           account.Name,
		Type:           account.Type,
		Currency:       account.Currency,
		Icon:           account.Icon,
		Color:          account.Color,
		Balance:        account.Balance,
		InitialBalance: account.InitialBalance,
		Note:           account.Note,
	}
}

func (s *accountService) CreateAccount(userID uuid.UUID, req dto.CreateAccountRequest) (*dto.AccountResponse, error) {
	currency := req.Currency
	if currency == "" {
		// Use user's defaultCurrency setting so new accounts match the chosen currency
		user, err := s.userRepo.GetUserByID(userID.String())
		if err == nil && user != nil && user.DefaultCurrency != "" {
			currency = user.DefaultCurrency
		} else {
			currency = "IDR"
		}
	}

	account := &models.Account{
		UserID:         userID,
		Name:           req.Name,
		Type:           req.Type,
		Currency:       currency,
		Icon:           req.Icon,
		Color:          req.Color,
		InitialBalance: req.InitialBalance,
		Balance:        req.InitialBalance, // Balance starts equal to InitialBalance
		Note:           req.Note,
	}

	if err := s.repo.Create(account); err != nil {
		return nil, err
	}

	return mapAccountToResponse(account), nil
}

func (s *accountService) GetAccounts(userID uuid.UUID, page, limit int) ([]dto.AccountResponse, int64, error) {
	accounts, total, err := s.repo.GetByUserID(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.AccountResponse
	for _, acc := range accounts {
		responses = append(responses, *mapAccountToResponse(&acc))
	}

	return responses, total, nil
}

func (s *accountService) GetAccountByID(userID, accountID uuid.UUID) (*dto.AccountResponse, error) {
	account, err := s.repo.GetByID(accountID)
	if err != nil {
		return nil, err
	}
	// Check ownership
	if account.UserID != userID {
		return nil, errors.New("unauthorized: account does not belong to user")
	}
	return mapAccountToResponse(account), nil
}

func (s *accountService) UpdateAccount(userID, accountID uuid.UUID, req dto.UpdateAccountRequest) (*dto.AccountResponse, error) {
	account, err := s.repo.GetByID(accountID)
	if err != nil {
		return nil, err
	}
	if account.UserID != userID {
		return nil, errors.New("unauthorized: account does not belong to user")
	}

	account.Name = req.Name
	account.Type = req.Type
	if req.Currency != "" {
		account.Currency = req.Currency
	}
	account.Icon = req.Icon
	account.Color = req.Color
	balanceDiff := req.InitialBalance - account.InitialBalance
	account.Balance += balanceDiff
	account.InitialBalance = req.InitialBalance
	account.Note = req.Note

	if err := s.repo.Update(account); err != nil {
		return nil, err
	}

	return mapAccountToResponse(account), nil
}

func (s *accountService) DeleteAccount(userID, accountID uuid.UUID) error {
	account, err := s.repo.GetByID(accountID)
	if err != nil {
		return err
	}
	if account.UserID != userID {
		return errors.New("unauthorized: account does not belong to user")
	}
	return s.repo.Delete(account.ID)
}
