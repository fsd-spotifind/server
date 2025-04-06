package server

import (
	"context"
	db "fsd-backend/prisma/db"
	"log"
	"strings"

	"github.com/robfig/cron/v3"
)

func (s *Server) GenerateStatsForAllUsers(period string) error {
	accounts, err := s.db.GetAllUserAccounts(context.Background())
	if err != nil {
		return err
	}
	for _, account := range accounts {
		go func(account db.AccountModel) {
			accessToken, err := s.RefreshAccessTokenIfNeeded(context.Background(), &account)
			if err != nil {
				log.Printf("User %s (%s): %v", account.UserID, period, err)
				return
			}
			statsPeriod := db.StatisticPeriod(strings.ToUpper(period))
			if err := s.GenerateUserStatistic(context.Background(), account.UserID, accessToken, statsPeriod); err != nil {
				log.Printf("User %s (%s): %v", account.UserID, period, err)
			}
		}(account)
	}
	return nil
}

func (s *Server) StartCronJobs() {
	c := cron.New()

	// Weekly stats: Every Monday at 1am
	c.AddFunc("0 1 * * 1", func() {
		log.Println("[CRON] Weekly stats generation started")
		if err := s.GenerateStatsForAllUsers("weekly"); err != nil {
			log.Println("Error:", err)
		}
	})

	// Monthly stats: First day of month at 2am
	c.AddFunc("0 2 1 * *", func() {
		log.Println("[CRON] Monthly stats generation started")
		if err := s.GenerateStatsForAllUsers("monthly"); err != nil {
			log.Println("Error:", err)
		}
	})

	c.Start()
}
