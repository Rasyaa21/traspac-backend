package cron

import "log"

func (s *Scheduler) ClearWeeklyBudget() {
	spec := "0 0 2 * * 1"
	_, err := s.cron.AddFunc(spec, func ()  {
		if err := s.UserBudgetService.UserBudgetRepo.ResetWeeklyBudget(); err != nil {
			log.Printf("[CRON] weekly budget error: %v\n", err)
			return
		}
		log.Println("[CRON] weekly budget reset success")
	})
		if err != nil {
		log.Fatalf("[CRON] gagal daftar token job: %v", err)
	}
}