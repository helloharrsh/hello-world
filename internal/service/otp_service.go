package service

import (
	"fmt"
	"math/rand"
	"time"

	"mailer_application/internal/infrastructure/db"
	"mailer_application/internal/infrastructure/mail"
)

type OTPService struct {
	Repo   *db.Repository
	Mailer *mail.Mailer
}

func NewOTPService(repo *db.Repository, mailer *mail.Mailer) *OTPService {
	return &OTPService{
		Repo:   repo,
		Mailer: mailer,
	}
}

func (s *OTPService) GenerateAndSendOTP(email string) error {
	code := generateOTP()
	expiry := time.Now().Add(5 * time.Minute)

	err := s.Repo.SaveOTP(email, code, expiry)
	if err != nil {
		return err
	}

	body := fmt.Sprintf("Your OTP is: %s\nIt expires in 5 minutes.", code)
	return s.Mailer.Send(email, "Your OTP Code", body)
}

func (s *OTPService) VerifyOTP(email, code string) (bool, error) {
	storedOTP, err := s.Repo.GetOTP(email)
	if err != nil {
		return false, err
	}

	if storedOTP.Code != code || storedOTP.ExpiresAt.Before(time.Now()) {
		return false, nil
	}

	// OTP is valid; mark user as verified
	err = s.Repo.MarkUserAsVerified(email)
	if err != nil {
		return false, err
	}

	return true, nil
}

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
