package cron

import (
	"log"
)

func (s *Scheduler) CleanTokenJobs() {
	// tiap hari jam 02:00
	spec := "0 0 2 * * *"
	_, err := s.cron.AddFunc(spec, func() {
		if err := s.EmailVerificationService.CleanupExpiredTokens(); err != nil {
			log.Printf("[CRON] token cleanup error: %v\n", err)
			return
		}
		log.Println("[CRON] token cleanup success")
	})

	if err != nil {
		log.Fatalf("[CRON] gagal daftar token job: %v", err)
	}
}