package service

import (
	"errors"
	"finvera-be/internal/dto"
	"finvera-be/internal/models"
	"finvera-be/internal/repository"

	"github.com/google/uuid"
)

type ScheduledTransactionService interface {
	CreateScheduled(userID uuid.UUID, req dto.CreateScheduledRequest) (*dto.ScheduledTransactionResponse, error)
	GetScheduleds(userID uuid.UUID, page, limit int) ([]dto.ScheduledTransactionResponse, int64, error)
	GetScheduledByID(userID, scheduledID uuid.UUID) (*dto.ScheduledTransactionResponse, error)
	UpdateScheduled(userID, scheduledID uuid.UUID, req dto.UpdateScheduledRequest) (*dto.ScheduledTransactionResponse, error)
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

// Mapper
func mapScheduledToResponse(scheduled *models.ScheduledTransaction) *dto.ScheduledTransactionResponse {
	if scheduled == nil {
		return nil
	}

	var targetAcc *dto.AccountResponse
	if scheduled.TargetAccount != nil {
		targetAcc = mapAccountToResponse(scheduled.TargetAccount)
	}

	return &dto.ScheduledTransactionResponse{
		ID:            scheduled.ID,
		Name:          scheduled.Name,
		Type:          scheduled.Type,
		Amount:        scheduled.Amount,
		Account:       *mapAccountToResponse(&scheduled.Account),
		TargetAccount: targetAcc,
		Category:      *mapCategoryToResponse(&scheduled.Category),
		Note:          scheduled.Note,
		Frequency:     scheduled.Frequency,
		NextRun:       scheduled.NextRun,
		LastRun:       nil,
		IsActive:      scheduled.IsActive,
	}
}

func (s *scheduledTransactionService) CreateScheduled(userID uuid.UUID, req dto.CreateScheduledRequest) (*dto.ScheduledTransactionResponse, error) {
	// Validations
	account, err := s.accountRepo.GetByID(req.AccountID)
	if err != nil || account.UserID != userID {
		return nil, errors.New("invalid or unauthorized accountId")
	}

	category, err := s.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("invalid categoryId")
	}
	if category.UserID != nil && *category.UserID != userID {
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

	reloaded, _ := s.repo.GetByID(scheduled.ID)
	return mapScheduledToResponse(reloaded), nil
}

func (s *scheduledTransactionService) GetScheduleds(userID uuid.UUID, page, limit int) ([]dto.ScheduledTransactionResponse, int64, error) {
	scheduleds, total, err := s.repo.GetByUserID(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.ScheduledTransactionResponse
	for _, sched := range scheduleds {
		responses = append(responses, *mapScheduledToResponse(&sched))
	}

	return responses, total, nil
}

func (s *scheduledTransactionService) GetScheduledByID(userID, scheduledID uuid.UUID) (*dto.ScheduledTransactionResponse, error) {
	scheduled, err := s.repo.GetByID(scheduledID)
	if err != nil {
		return nil, err
	}
	if scheduled.UserID != userID {
		return nil, errors.New("unauthorized: scheduled transaction does not belong to user")
	}
	return mapScheduledToResponse(scheduled), nil
}

func (s *scheduledTransactionService) UpdateScheduled(userID, scheduledID uuid.UUID, req dto.UpdateScheduledRequest) (*dto.ScheduledTransactionResponse, error) {
	scheduled, err := s.repo.GetByID(scheduledID)
	if err != nil {
		return nil, err
	}
	if scheduled.UserID != userID {
		return nil, errors.New("unauthorized: scheduled transaction does not belong to user")
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
	if category.UserID != nil && *category.UserID != userID {
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

	reloaded, _ := s.repo.GetByID(scheduledID)
	return mapScheduledToResponse(reloaded), nil
}

func (s *scheduledTransactionService) DeleteScheduled(userID, scheduledID uuid.UUID) error {
	scheduled, err := s.repo.GetByID(scheduledID)
	if err != nil {
		return err
	}
	if scheduled.UserID != userID {
		return errors.New("unauthorized: scheduled transaction does not belong to user")
	}
	return s.repo.Delete(scheduled.ID)
}
