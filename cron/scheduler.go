package cron

import (
	"log"

	"github.com/robfig/cron/v3"

	"mailer_application/internal/infrastructure/db"
	"mailer_application/internal/infrastructure/mail"
)

type Scheduler struct {
	Repo   *db.Repository
	Mailer *mail.Mailer
}

func NewScheduler(repo *db.Repository, mailer *mail.Mailer) *Scheduler {
	return &Scheduler{
		Repo:   repo,
		Mailer: mailer,
	}
}

func (s *Scheduler) Start() {
	c := cron.New()

	// Run every day at 9:00 AM
	_, err := c.AddFunc("0 9 * * *", func() {
		log.Println("⏰ Running daily interview question job...")
		s.sendInterviewQuestions()
	})

	if err != nil {
		log.Fatalf("Failed to schedule job: %v", err)
	}

	c.Start()
	log.Println("✅ Scheduler started.")
}

func (s *Scheduler) sendInterviewQuestions() {
	// Sample static question (you can fetch from a DB or API in the future)
	question := "What is the difference between a goroutine and a thread in Go?"

	users, err := s.Repo.GetAllVerifiedUsers()
	if err != nil {
		log.Printf("Failed to fetch users: %v", err)
		return
	}

	for _, user := range users {
		err := s.Mailer.Send(user.Email, "Daily Interview Question", question)
		if err != nil {
			log.Printf("Failed to send email to %s: %v", user.Email, err)
		}
	}
}
