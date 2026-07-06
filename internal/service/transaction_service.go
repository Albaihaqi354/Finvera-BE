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
	CreateTransaction(userID uuid.UUID, req dto.TransactionRequest) (*dto.TransactionResponse, error)
	GetTransactions(userID uuid.UUID, page, limit int, filter repository.TransactionFilter) ([]dto.TransactionResponse, int64, error)
	GetTransactionByID(userID, transactionID uuid.UUID) (*dto.TransactionResponse, error)
	UpdateTransaction(userID, transactionID uuid.UUID, req dto.TransactionRequest) (*dto.TransactionResponse, error)
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

// ── Mapper helpers ────────────────────────────────────────────────────────────

func mapTransactionToResponse(t *models.Transaction) *dto.TransactionResponse {
	if t == nil {
		return nil
	}
	var targetAcc *dto.AccountResponse
	if t.TargetAccount != nil {
		targetAcc = mapAccountToResponse(t.TargetAccount)
	}
	var tags []dto.TagResponse
	for _, tag := range t.Tags {
		tags = append(tags, *mapTagToResponse(&tag))
	}
	return &dto.TransactionResponse{
		ID:            t.ID,
		Type:          t.Type,
		Amount:        t.Amount,
		Account:       *mapAccountToResponse(&t.Account),
		TargetAccount: targetAcc,
		Category:      *mapCategoryToResponse(&t.Category),
		Date:          t.Date,
		Note:          t.Note,
		Tags:          tags,
		CreatedAt:     t.CreatedAt,
	}
}

// ── Validation helpers ────────────────────────────────────────────────────────

func (s *transactionService) validateTransactionRequest(
	userID uuid.UUID,
	req dto.TransactionRequest,
) (account *models.Account, targetAccount *models.Account, tags []models.Tag, err error) {
	// Validate account ownership
	account, err = s.accountRepo.GetByID(req.AccountID)
	if err != nil || account.UserID != userID {
		return nil, nil, nil, errors.New("invalid or unauthorized accountId")
	}

	// Validate category
	category, err := s.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		return nil, nil, nil, errors.New("invalid categoryId")
	}
	if category.UserID != nil && *category.UserID != userID {
		return nil, nil, nil, errors.New("unauthorized categoryId")
	}

	// Transfer-specific validation
	if req.Type == "transfer" {
		if req.TargetAccountID == nil {
			return nil, nil, nil, errors.New("targetAccountId is required for transfer transactions")
		}
		if *req.TargetAccountID == req.AccountID {
			return nil, nil, nil, errors.New("target account must be different from source account")
		}
		targetAccount, err = s.accountRepo.GetByID(*req.TargetAccountID)
		if err != nil || targetAccount.UserID != userID {
			return nil, nil, nil, errors.New("invalid or unauthorized targetAccountId")
		}
	} else if req.TargetAccountID != nil {
		return nil, nil, nil, errors.New("targetAccountId should only be provided for transfer transactions")
	}

	// Validate tags ownership
	if len(req.TagIDs) > 0 {
		foundTags, err := s.tagRepo.GetByIDs(req.TagIDs)
		if err != nil {
			return nil, nil, nil, err
		}
		if len(foundTags) != len(req.TagIDs) {
			return nil, nil, nil, errors.New("invalid tagId found")
		}
		for _, tag := range foundTags {
			if tag.UserID != userID {
				return nil, nil, nil, errors.New("invalid or unauthorized tagId")
			}
			tags = append(tags, tag)
		}
	}

	return account, targetAccount, tags, nil
}

// ── Balance helpers ───────────────────────────────────────────────────────────

func applyBalance(tx *gorm.DB, txType string, amount float64, account, targetAccount *models.Account) error {
	switch txType {
	case "income":
		account.Balance += amount
		return tx.Save(account).Error
	case "expense":
		account.Balance -= amount
		return tx.Save(account).Error
	case "transfer":
		account.Balance -= amount
		targetAccount.Balance += amount
		if err := tx.Save(account).Error; err != nil {
			return err
		}
		return tx.Save(targetAccount).Error
	}
	return nil
}

func (s *transactionService) revertBalance(tx *gorm.DB, transaction *models.Transaction) error {
	account := transaction.Account
	switch transaction.Type {
	case "income":
		account.Balance -= transaction.Amount
		return tx.Save(&account).Error
	case "expense":
		account.Balance += transaction.Amount
		return tx.Save(&account).Error
	case "transfer":
		account.Balance += transaction.Amount
		if err := tx.Save(&account).Error; err != nil {
			return err
		}
		if transaction.TargetAccount != nil {
			ta := *transaction.TargetAccount
			ta.Balance -= transaction.Amount
			return tx.Save(&ta).Error
		} else if transaction.TargetAccountID != nil {
			var ta models.Account
			if err := tx.First(&ta, "id = ?", *transaction.TargetAccountID).Error; err != nil {
				return err
			}
			ta.Balance -= transaction.Amount
			return tx.Save(&ta).Error
		}
	}
	return nil
}

// ── Service methods ───────────────────────────────────────────────────────────

func (s *transactionService) CreateTransaction(userID uuid.UUID, req dto.TransactionRequest) (*dto.TransactionResponse, error) {
	account, targetAccount, tags, err := s.validateTransactionRequest(userID, req)
	if err != nil {
		return nil, err
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
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.CreateWithTx(tx, transaction); err != nil {
			return err
		}
		return applyBalance(tx, req.Type, req.Amount, account, targetAccount)
	}); err != nil {
		return nil, err
	}

	reloaded, err := s.repo.GetByID(transaction.ID)
	if err != nil {
		return nil, err
	}
	return mapTransactionToResponse(reloaded), nil
}

func (s *transactionService) GetTransactions(userID uuid.UUID, page, limit int, filter repository.TransactionFilter) ([]dto.TransactionResponse, int64, error) {
	transactions, total, err := s.repo.GetByUserID(userID, page, limit, filter)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]dto.TransactionResponse, 0, len(transactions))
	for i := range transactions {
		responses = append(responses, *mapTransactionToResponse(&transactions[i]))
	}
	return responses, total, nil
}

func (s *transactionService) GetTransactionByID(userID, transactionID uuid.UUID) (*dto.TransactionResponse, error) {
	transaction, err := s.repo.GetByID(transactionID)
	if err != nil {
		return nil, err
	}
	if transaction.UserID != userID {
		return nil, errors.New("transaction not found")
	}
	return mapTransactionToResponse(transaction), nil
}

func (s *transactionService) UpdateTransaction(userID, transactionID uuid.UUID, req dto.TransactionRequest) (*dto.TransactionResponse, error) {
	transaction, err := s.repo.GetByID(transactionID)
	if err != nil {
		return nil, err
	}
	if transaction.UserID != userID {
		return nil, errors.New("transaction not found")
	}

	account, targetAccount, tags, err := s.validateTransactionRequest(userID, req)
	if err != nil {
		return nil, err
	}

	db := s.repo.GetDB()
	if err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Revert old balances
		if err := s.revertBalance(tx, transaction); err != nil {
			return err
		}
		// 2. Clear old many2many tags
		if err := tx.Model(transaction).Association("Tags").Clear(); err != nil {
			return err
		}
		// 3. Update fields
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
		// 4. Apply new balances
		return applyBalance(tx, req.Type, req.Amount, account, targetAccount)
	}); err != nil {
		return nil, err
	}

	reloaded, err := s.repo.GetByID(transactionID)
	if err != nil {
		return nil, err
	}
	return mapTransactionToResponse(reloaded), nil
}

func (s *transactionService) DeleteTransaction(userID, transactionID uuid.UUID) error {
	transaction, err := s.repo.GetByID(transactionID)
	if err != nil {
		return err
	}
	if transaction.UserID != userID {
		return errors.New("transaction not found")
	}

	db := s.repo.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		if err := s.revertBalance(tx, transaction); err != nil {
			return err
		}
		if err := tx.Model(transaction).Association("Tags").Clear(); err != nil {
			return err
		}
		return s.repo.DeleteWithTx(tx, transactionID)
	})
}
