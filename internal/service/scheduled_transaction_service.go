package service

import (
	"errors"
	"finvera-be/internal/models"
	"finvera-be/internal/repository"
	"time"

	"github.com/google/uuid"
)

type ScheduledTransactionService interface {
	CreateScheduled(userID uuid.UUID, req CreateScheduledRequest) (*models.ScheduledTransaction, error)
	GetScheduleds(userID uuid.UUID) ([]models.ScheduledTransaction, error)
	GetScheduledByID(userID, scheduledID uuid.UUID) (*models.ScheduledTransaction, error)
	UpdateScheduled(userID, scheduledID uuid.UUID, req UpdateScheduledRequest) (*models.ScheduledTransaction, error)
	DeleteScheduled(userID, scheduledID uuid.UUID) error
}

type scheduledTransactionService struct {
	repo         repository.ScheduledTransactionRepository
	accountRepo  repository.AccountRepository
	categoryRepo repository.CategoryRepository
}

func NewScheduledTransactionService(
	repo repository.ScheduledTransactionRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
) ScheduledTransactionService {
	return &scheduledTransactionService{
		repo:         repo,
		accountRepo:  accountRepo,
		categoryRepo: categoryRepo,
	}
}

// Request DTOs
type CreateScheduledRequest struct {
	Name            string     `json:"name" binding:"required"`
	Type            string     `json:"type" binding:"required,oneof=income expense transfer"`
	Amount          float64    `json:"amount" binding:"required,gt=0"`
	AccountID       uuid.UUID  `json:"accountId" binding:"required"`
	TargetAccountID *uuid.UUID `json:"targetAccountId"`
	CategoryID      uuid.UUID  `json:"categoryId" binding:"required"`
	Note            string     `json:"note"`
	Frequency       string     `json:"frequency" binding:"required,oneof=daily weekly monthly yearly"`
	NextRun         time.Time  `json:"nextRun" binding:"required"`
	IsActive        bool       `json:"isActive"`
}

type UpdateScheduledRequest struct {
	Name            string     `json:"name" binding:"required"`
	Type            string     `json:"type" binding:"required,oneof=income expense transfer"`
	Amount          float64    `json:"amount" binding:"required,gt=0"`
	AccountID       uuid.UUID  `json:"accountId" binding:"required"`
	TargetAccountID *uuid.UUID `json:"targetAccountId"`
	CategoryID      uuid.UUID  `json:"categoryId" binding:"required"`
	Note            string     `json:"note"`
	Frequency       string     `json:"frequency" binding:"required,oneof=daily weekly monthly yearly"`
	NextRun         time.Time  `json:"nextRun" binding:"required"`
	IsActive        bool       `json:"isActive"`
}

func (s *scheduledTransactionService) CreateScheduled(userID uuid.UUID, req CreateScheduledRequest) (*models.ScheduledTransaction, error) {
	// Validations
	account, err := s.accountRepo.GetByID(req.AccountID)
	if err != nil || account.UserID != userID {
		return nil, errors.New("invalid or unauthorized accountId")
	}

	category, err := s.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("invalid categoryId")
	}
	if category.UserID != uuid.Nil && category.UserID != userID {
		return nil, errors.New("unauthorized categoryId")
	}

	if req.Type == "transfer" {
		if req.TargetAccountID == nil {
			return nil, errors.New("targetAccountId is required for transfer transactions")
		}
		targetAccount, err := s.accountRepo.GetByID(*req.TargetAccountID)
		if err != nil || targetAccount.UserID != userID {
			return nil, errors.New("invalid or unauthorized targetAccountId")
		}
	}

	scheduled := &models.ScheduledTransaction{
		UserID:          userID,
		Name:            req.Name,
		Type:            req.Type,
		Amount:          req.Amount,
		AccountID:       req.AccountID,
		TargetAccountID: req.TargetAccountID,
		CategoryID:      req.CategoryID,
		Note:            req.Note,
		Frequency:       req.Frequency,
		NextRun:         req.NextRun,
		IsActive:        req.IsActive,
	}

	if err := s.repo.Create(scheduled); err != nil {
		return nil, err
	}

	return s.repo.GetByID(scheduled.ID)
}

func (s *scheduledTransactionService) GetScheduleds(userID uuid.UUID) ([]models.ScheduledTransaction, error) {
	return s.repo.GetByUserID(userID)
}

func (s *scheduledTransactionService) GetScheduledByID(userID, scheduledID uuid.UUID) (*models.ScheduledTransaction, error) {
	scheduled, err := s.repo.GetByID(scheduledID)
	if err != nil {
		return nil, err
	}
	if scheduled.UserID != userID {
		return nil, errors.New("unauthorized: scheduled transaction does not belong to user")
	}
	return scheduled, nil
}

func (s *scheduledTransactionService) UpdateScheduled(userID, scheduledID uuid.UUID, req UpdateScheduledRequest) (*models.ScheduledTransaction, error) {
	scheduled, err := s.GetScheduledByID(userID, scheduledID)
	if err != nil {
		return nil, err
	}

	// Validations
	account, err := s.accountRepo.GetByID(req.AccountID)
	if err != nil || account.UserID != userID {
		return nil, errors.New("invalid or unauthorized accountId")
	}

	category, err := s.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("invalid categoryId")
	}
	if category.UserID != uuid.Nil && category.UserID != userID {
		return nil, errors.New("unauthorized categoryId")
	}

	if req.Type == "transfer" {
		if req.TargetAccountID == nil {
			return nil, errors.New("targetAccountId is required for transfer transactions")
		}
		targetAccount, err := s.accountRepo.GetByID(*req.TargetAccountID)
		if err != nil || targetAccount.UserID != userID {
			return nil, errors.New("invalid or unauthorized targetAccountId")
		}
	}

	scheduled.Name = req.Name
	scheduled.Type = req.Type
	scheduled.Amount = req.Amount
	scheduled.AccountID = req.AccountID
	scheduled.TargetAccountID = req.TargetAccountID
	scheduled.CategoryID = req.CategoryID
	scheduled.Note = req.Note
	scheduled.Frequency = req.Frequency
	scheduled.NextRun = req.NextRun
	scheduled.IsActive = req.IsActive

	if err := s.repo.Update(scheduled); err != nil {
		return nil, err
	}

	return s.repo.GetByID(scheduledID)
}

func (s *scheduledTransactionService) DeleteScheduled(userID, scheduledID uuid.UUID) error {
	scheduled, err := s.GetScheduledByID(userID, scheduledID)
	if err != nil {
		return err
	}
	return s.repo.Delete(scheduled.ID)
}
