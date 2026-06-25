package cron

import (
	"finvera-be/internal/dto"
	"finvera-be/internal/models"
	"finvera-be/internal/service"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronService struct {
	scheduler *cron.Cron
	db        *gorm.DB
	txService service.TransactionService
}

func NewCronService(db *gorm.DB, txService service.TransactionService) *CronService {
	// create scheduler with seconds precision for easier local testing, or just standard
	// Standard parser: cron.New() runs on minute precision.
	scheduler := cron.New()
	return &CronService{
		scheduler: scheduler,
		db:        db,
		txService: txService,
	}
}

func (s *CronService) Start() {
	// Jalankan pengecekan setiap jam 00:01
	// Untuk local testing yang cepat, bisa pakai "* * * * *" (setiap menit)
	_, err := s.scheduler.AddFunc("1 0 * * *", s.ProcessScheduledTransactions)
	if err != nil {
		log.Printf("Error scheduling ProcessScheduledTransactions: %v", err)
	}

	// Buat testing lokal (jalan setiap menit):
	_, err = s.scheduler.AddFunc("* * * * *", s.ProcessScheduledTransactions)
	if err != nil {
		log.Printf("Error scheduling test cron: %v", err)
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
	// Find all active scheduled transactions where NextRun is <= time.Now()
	now := time.Now()
	var scheduled []models.ScheduledTransaction

	if err := s.db.Where("is_active = ? AND next_run <= ?", true, now).Find(&scheduled).Error; err != nil {
		log.Printf("Error finding scheduled transactions: %v", err)
		return
	}

	for _, st := range scheduled {
		// Logika: Execute this transaction!
		// But remember to make it inside a transaction or handle properly
		s.executeScheduled(st, now)
	}
}

func (s *CronService) executeScheduled(st models.ScheduledTransaction, runTime time.Time) {
	// Create the transaction using existing service (which handles account balance updates)
	req := dto.CreateTransactionRequest{
		Type:            st.Type,
		Amount:          st.Amount,
		AccountID:       st.AccountID,
		TargetAccountID: st.TargetAccountID,
		CategoryID:      st.CategoryID,
		Date:            runTime,
		Note:            st.Note + " (Auto-generated)",
		TagIDs:          []uuid.UUID{}, // Currently ScheduledTransaction doesn't have TagIDs explicitly mapped, assume empty
	}

	_, err := s.txService.CreateTransaction(st.UserID, req)
	if err != nil {
		log.Printf("Failed to execute scheduled transaction %s: %v", st.ID, err)
		return
	}

	// Update NextRun based on Frequency
	var nextRun time.Time
	switch st.Frequency {
	case "daily":
		nextRun = st.NextRun.AddDate(0, 0, 1)
	case "weekly":
		nextRun = st.NextRun.AddDate(0, 0, 7)
	case "monthly":
		nextRun = st.NextRun.AddDate(0, 1, 0)
	case "yearly":
		nextRun = st.NextRun.AddDate(1, 0, 0)
	default:
		// Fallback safe daily
		nextRun = st.NextRun.AddDate(0, 0, 1)
	}

	st.NextRun = nextRun
	if err := s.db.Save(&st).Error; err != nil {
		log.Printf("Failed to update NextRun for scheduled transaction %s: %v", st.ID, err)
	} else {
		log.Printf("Successfully executed and updated scheduled transaction %s", st.ID)
	}
}
