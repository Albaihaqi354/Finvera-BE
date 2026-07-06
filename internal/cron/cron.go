package cron

import (
	"finvera-be/internal/dto"
	"finvera-be/internal/models"
	"finvera-be/internal/service"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronService struct {
	scheduler *cron.Cron
	db        *gorm.DB
	txService service.TransactionService
	mutex     sync.Mutex
}

func NewCronService(db *gorm.DB, txService service.TransactionService) *CronService {
	scheduler := cron.New()
	return &CronService{
		scheduler: scheduler,
		db:        db,
		txService: txService,
	}
}

func (s *CronService) Start() {
	_, err := s.scheduler.AddFunc("1 0 * * *", s.ProcessScheduledTransactions)
	if err != nil {
		log.Printf("Error scheduling ProcessScheduledTransactions: %v", err)
	}

	s.scheduler.Start()
	log.Println("Cron Service started")
}

func (s *CronService) Stop() {
	ctx := s.scheduler.Stop()
	<-ctx.Done()
	log.Println("Cron Service stopped")
}

func (s *CronService) ProcessScheduledTransactions() {
	// Prevent overlapping runs on the same instance
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	var scheduled []models.ScheduledTransaction

	if err := s.db.Where("is_active = ? AND next_run <= ?", true, now).Find(&scheduled).Error; err != nil {
		log.Printf("Error finding scheduled transactions: %v", err)
		return
	}

	for _, st := range scheduled {
		s.executeScheduled(st, now)
	}
}

func (s *CronService) executeScheduled(st models.ScheduledTransaction, runTime time.Time) {
	// Execute within a database transaction to ensure atomicity
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Ensure we lock this row so other instances don't process it concurrently
		var lockedSt models.ScheduledTransaction
		if err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&lockedSt, "id = ?", st.ID).Error; err != nil {
			return err
		}
		
		if !lockedSt.IsActive || lockedSt.NextRun.After(runTime) {
			// Already processed by another instance or deactivated
			return nil
		}

		// Calculate NextRun
		var nextRun time.Time
		switch lockedSt.Frequency {
		case "daily":
			nextRun = lockedSt.NextRun.AddDate(0, 0, 1)
		case "weekly":
			nextRun = lockedSt.NextRun.AddDate(0, 0, 7)
		case "monthly":
			nextRun = lockedSt.NextRun.AddDate(0, 1, 0)
		case "yearly":
			nextRun = lockedSt.NextRun.AddDate(1, 0, 0)
		default:
			nextRun = lockedSt.NextRun.AddDate(0, 0, 1)
		}

		// Update NextRun FIRST. If this fails, transaction rolls back.
		lockedSt.NextRun = nextRun
		if err := tx.Save(&lockedSt).Error; err != nil {
			return err
		}

		// Because txService.CreateTransaction creates its own transaction,
		// there's a risk of nested transaction issues or connection locks.
		// However, we rely on it here for simplicity. A better long-term fix
		// is injecting `tx` into CreateTransaction.
		req := dto.CreateTransactionRequest{
			Type:            lockedSt.Type,
			Amount:          lockedSt.Amount,
			AccountID:       lockedSt.AccountID,
			TargetAccountID: lockedSt.TargetAccountID,
			CategoryID:      lockedSt.CategoryID,
			Date:            runTime,
			Note:            lockedSt.Note + " (Auto-generated)",
			TagIDs:          []uuid.UUID{},
		}

		_, err := s.txService.CreateTransaction(lockedSt.UserID, req)
		if err != nil {
			return err
		}
		
		return nil
	})

	if err != nil {
		log.Printf("Failed to execute scheduled transaction %s: %v", st.ID, err)
	} else {
		log.Printf("Successfully executed and updated scheduled transaction %s", st.ID)
	}
}
