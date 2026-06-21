package service

import (
	"errors"
	"finvera-be/internal/dto"
	"finvera-be/internal/models"
	"finvera-be/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService interface {
	CreateTransaction(userID uuid.UUID, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	GetTransactions(userID uuid.UUID, page, limit int) ([]dto.TransactionResponse, int64, error)
	GetTransactionByID(userID, transactionID uuid.UUID) (*dto.TransactionResponse, error)
	UpdateTransaction(userID, transactionID uuid.UUID, req dto.UpdateTransactionRequest) (*dto.TransactionResponse, error)
	DeleteTransaction(userID, transactionID uuid.UUID) error
}

type transactionService struct {
	repo         repository.TransactionRepository
	accountRepo  repository.AccountRepository
	categoryRepo repository.CategoryRepository
	tagRepo      repository.TagRepository
}

func NewTransactionService(
	repo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
	tagRepo repository.TagRepository,
) TransactionService {
	return &transactionService{
		repo:         repo,
		accountRepo:  accountRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}
}

// Mapper
func mapTransactionToResponse(transaction *models.Transaction) *dto.TransactionResponse {
	if transaction == nil {
		return nil
	}

	var targetAcc *dto.AccountResponse
	if transaction.TargetAccount != nil {
		targetAcc = mapAccountToResponse(transaction.TargetAccount)
	}

	var tagResponses []dto.TagResponse
	for _, tag := range transaction.Tags {
		tagResponses = append(tagResponses, *mapTagToResponse(&tag))
	}

	return &dto.TransactionResponse{
		ID:            transaction.ID,
		Type:          transaction.Type,
		Amount:        transaction.Amount,
		Account:       *mapAccountToResponse(&transaction.Account),
		TargetAccount: targetAcc,
		Category:      *mapCategoryToResponse(&transaction.Category),
		Date:          transaction.Date,
		Note:          transaction.Note,
		Tags:          tagResponses,
		CreatedAt:     transaction.CreatedAt,
	}
}

func (s *transactionService) CreateTransaction(userID uuid.UUID, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	// Validations
	account, err := s.accountRepo.GetByID(req.AccountID)
	if err != nil || account.UserID != userID {
		return nil, errors.New("invalid or unauthorized accountId")
	}

	category, err := s.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("invalid categoryId")
	}
	// Category belongs to user or is a global category
	if category.UserID != uuid.Nil && category.UserID != userID {
		return nil, errors.New("unauthorized categoryId")
	}

	var targetAccount *models.Account
	if req.Type == "transfer" {
		if req.TargetAccountID == nil {
			return nil, errors.New("targetAccountId is required for transfer transactions")
		}
		if *req.TargetAccountID == req.AccountID {
			return nil, errors.New("target account must be different from source account")
		}
		targetAccount, err = s.accountRepo.GetByID(*req.TargetAccountID)
		if err != nil || targetAccount.UserID != userID {
			return nil, errors.New("invalid or unauthorized targetAccountId")
		}
	} else {
		if req.TargetAccountID != nil {
			return nil, errors.New("targetAccountId should only be provided for transfer transactions")
		}
	}

	var tags []models.Tag
	for _, tagID := range req.TagIDs {
		tag, err := s.tagRepo.GetByID(tagID)
		if err != nil || tag.UserID != userID {
			return nil, errors.New("invalid or unauthorized tagId")
		}
		tags = append(tags, *tag)
	}

	transaction := &models.Transaction{
		UserID:          userID,
		Type:            req.Type,
		Amount:          req.Amount,
		AccountID:       req.AccountID,
		TargetAccountID: req.TargetAccountID,
		CategoryID:      req.CategoryID,
		Date:            req.Date,
		Note:            req.Note,
		Tags:            tags,
	}

	db := s.repo.GetDB()
	
	// Start DB Transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		// 1. Save Transaction
		if err := s.repo.CreateWithTx(tx, transaction); err != nil {
			return err
		}

		// 2. Update Account Balance
		if req.Type == "income" {
			account.Balance += req.Amount
			if err := tx.Save(account).Error; err != nil {
				return err
			}
		} else if req.Type == "expense" {
			account.Balance -= req.Amount
			if err := tx.Save(account).Error; err != nil {
				return err
			}
		} else if req.Type == "transfer" {
			account.Balance -= req.Amount
			targetAccount.Balance += req.Amount
			if err := tx.Save(account).Error; err != nil {
				return err
			}
			if err := tx.Save(targetAccount).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload to get preloaded relations properly
	reloaded, _ := s.repo.GetByID(transaction.ID)
	return mapTransactionToResponse(reloaded), nil
}

func (s *transactionService) GetTransactions(userID uuid.UUID, page, limit int) ([]dto.TransactionResponse, int64, error) {
	transactions, total, err := s.repo.GetByUserID(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.TransactionResponse
	for _, t := range transactions {
		responses = append(responses, *mapTransactionToResponse(&t))
	}

	return responses, total, nil
}

func (s *transactionService) GetTransactionByID(userID, transactionID uuid.UUID) (*dto.TransactionResponse, error) {
	transaction, err := s.repo.GetByID(transactionID)
	if err != nil {
		return nil, err
	}
	if transaction.UserID != userID {
		return nil, errors.New("unauthorized: transaction does not belong to user")
	}
	return mapTransactionToResponse(transaction), nil
}

// RevertBalance is a helper function to reverse the effect of a transaction on balances
func (s *transactionService) revertBalance(tx *gorm.DB, transaction *models.Transaction) error {
	account := transaction.Account
	
	if transaction.Type == "income" {
		account.Balance -= transaction.Amount
		if err := tx.Save(&account).Error; err != nil {
			return err
		}
	} else if transaction.Type == "expense" {
		account.Balance += transaction.Amount
		if err := tx.Save(&account).Error; err != nil {
			return err
		}
	} else if transaction.Type == "transfer" {
		account.Balance += transaction.Amount
		if err := tx.Save(&account).Error; err != nil {
			return err
		}
		
		if transaction.TargetAccount != nil {
			targetAccount := *transaction.TargetAccount
			targetAccount.Balance -= transaction.Amount
			if err := tx.Save(&targetAccount).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *transactionService) UpdateTransaction(userID, transactionID uuid.UUID, req dto.UpdateTransactionRequest) (*dto.TransactionResponse, error) {
	// Updating a transaction requires reverting its old impact and applying the new impact
	
	// We'll delete the old one and create a new one to ensure everything runs perfectly in tx, 
	// or we can manually update fields and balances.
	// For full safety, let's revert the old balances, update fields, and apply new balances.

	// Since GetTransactionByID returns DTO now, we have to get the raw model directly from Repo.
	transaction, err := s.repo.GetByID(transactionID)
	if err != nil {
		return nil, err
	}
	if transaction.UserID != userID {
		return nil, errors.New("unauthorized: transaction does not belong to user")
	}

	// For simplicity in a robust way, let's just reuse Delete & Create logic by making a DB transaction here
	// This ensures we don't duplicate the balance logic.
	db := s.repo.GetDB()
	
	err = db.Transaction(func(tx *gorm.DB) error {
		// 1. Revert Old Balance
		if err := s.revertBalance(tx, transaction); err != nil {
			return err
		}

		// 2. Clear old tags and update fields
		if err := tx.Model(transaction).Association("Tags").Clear(); err != nil {
			return err
		}

		// Validations for new fields
		account, err := s.accountRepo.GetByID(req.AccountID)
		if err != nil || account.UserID != userID {
			return errors.New("invalid or unauthorized accountId")
		}

		category, err := s.categoryRepo.GetByID(req.CategoryID)
		if err != nil {
			return errors.New("invalid categoryId")
		}
		if category.UserID != uuid.Nil && category.UserID != userID {
			return errors.New("unauthorized categoryId")
		}

		var targetAccount *models.Account
		if req.Type == "transfer" {
			if req.TargetAccountID == nil {
				return errors.New("targetAccountId is required for transfer transactions")
			}
			if *req.TargetAccountID == req.AccountID {
				return errors.New("target account must be different from source account")
			}
			targetAccount, err = s.accountRepo.GetByID(*req.TargetAccountID)
			if err != nil || targetAccount.UserID != userID {
				return errors.New("invalid or unauthorized targetAccountId")
			}
		}

		var tags []models.Tag
		for _, tagID := range req.TagIDs {
			tag, err := s.tagRepo.GetByID(tagID)
			if err != nil || tag.UserID != userID {
				return errors.New("invalid or unauthorized tagId")
			}
			tags = append(tags, *tag)
		}

		transaction.Type = req.Type
		transaction.Amount = req.Amount
		transaction.AccountID = req.AccountID
		transaction.TargetAccountID = req.TargetAccountID
		transaction.CategoryID = req.CategoryID
		transaction.Date = req.Date
		transaction.Note = req.Note
		transaction.Tags = tags

		if err := tx.Save(transaction).Error; err != nil {
			return err
		}

		// 3. Apply New Balance
		if req.Type == "income" {
			account.Balance += req.Amount
			if err := tx.Save(account).Error; err != nil {
				return err
			}
		} else if req.Type == "expense" {
			account.Balance -= req.Amount
			if err := tx.Save(account).Error; err != nil {
				return err
			}
		} else if req.Type == "transfer" {
			account.Balance -= req.Amount
			targetAccount.Balance += req.Amount
			if err := tx.Save(account).Error; err != nil {
				return err
			}
			if err := tx.Save(targetAccount).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	reloaded, _ := s.repo.GetByID(transactionID)
	return mapTransactionToResponse(reloaded), nil
}

func (s *transactionService) DeleteTransaction(userID, transactionID uuid.UUID) error {
	transaction, err := s.repo.GetByID(transactionID)
	if err != nil {
		return err
	}
	if transaction.UserID != userID {
		return errors.New("unauthorized: transaction does not belong to user")
	}

	db := s.repo.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Revert balance
		if err := s.revertBalance(tx, transaction); err != nil {
			return err
		}

		// 2. Delete the transaction (along with association if soft delete handles it, GORM will delete relations)
		// Clear Many-to-Many associations first to be safe
		if err := tx.Model(transaction).Association("Tags").Clear(); err != nil {
			return err
		}

		return s.repo.DeleteWithTx(tx, transactionID)
	})
}
