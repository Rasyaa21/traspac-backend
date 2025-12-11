package cron

import (
	"log"

	"gin-backend-app/internal/services"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron      *cron.Cron
	EmailVerificationService  *services.EmailVerficationService
}

func NewScheduler(emailVerificationService *services.EmailVerficationService) *Scheduler {
	c := cron.New(cron.WithSeconds()) 

	s := &Scheduler{
		cron:     c,
		EmailVerificationService: emailVerificationService,
	}

	s.registerJobs()

	return s
}

func (s *Scheduler) registerJobs() {
	s.CleanTokenJobs()
}

func (s *Scheduler) Start() {
	log.Println("[CRON] start scheduler")
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	log.Println("[CRON] stop scheduler")
	s.cron.Stop()
}