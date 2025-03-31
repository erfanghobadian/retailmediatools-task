package scheduler

import (
	"go.uber.org/zap"
	"time"

	"github.com/robfig/cron/v3"
	"sweng-task/internal/service"
)

type Scheduler struct {
	lineItemService *service.LineItemService
	log             *zap.SugaredLogger
}

func NewScheduler(lineItemService *service.LineItemService, log *zap.SugaredLogger) *Scheduler {
	return &Scheduler{
		lineItemService: lineItemService,
		log:             log,
	}
}

func (s *Scheduler) Start() {
	c := cron.New(cron.WithLocation(time.Local))

	_, err := c.AddFunc("0 0 * * *", func() {
		s.log.Info("Starting daily budget reset...")
		err := s.lineItemService.ResetDailySpending()
		if err != nil {
			s.log.Errorf("Failed to reset budgets: %v", err)
		} else {
			s.log.Info("Daily budget reset complete.")
		}
	})
	if err != nil {
		s.log.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()
	s.log.Info("Scheduler started")

	go func() {
		select {}
	}()
}
