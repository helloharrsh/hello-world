package db

import (
	"database/sql"
	"time"

	"mailer_application/internal/domain/model"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// SaveOTP inserts a new OTP for a given email and creates the user if not exists
func (r *Repository) SaveOTP(email, code string, expiry time.Time) error {
	_, err := r.DB.Exec(`
		INSERT INTO users (email, verified, created_at)
		VALUES (?, false, NOW())
		ON DUPLICATE KEY UPDATE email = email`, email)
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(`
		INSERT INTO otps (email, code, expires_at, created_at)
		VALUES (?, ?, ?, NOW())`, email, code, expiry)
	return err
}

// GetOTP retrieves the most recent OTP for the given email
func (r *Repository) GetOTP(email string) (*model.OTP, error) {
	row := r.DB.QueryRow(`
		SELECT id, email, code, expires_at, created_at
		FROM otps
		WHERE email = ?
		ORDER BY created_at DESC
		LIMIT 1`, email)

	var otp model.OTP
	err := row.Scan(&otp.ID, &otp.Email, &otp.Code, &otp.ExpiresAt, &otp.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// MarkUserAsVerified sets verified = true for a given email
func (r *Repository) MarkUserAsVerified(email string) error {
	_, err := r.DB.Exec(`UPDATE users SET verified = true WHERE email = ?`, email)
	return err
}

// GetAllVerifiedUsers fetches all users who have verified their OTP
func (r *Repository) GetAllVerifiedUsers() ([]*model.User, error) {
	rows, err := r.DB.Query(`
		SELECT id, email, verified, created_at
		FROM users
		WHERE verified = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Verified, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *Repository) IsUserVerified(email string) (bool, error) {
	var verified bool
	err := r.DB.QueryRow(`SELECT verified FROM users WHERE email = ?`, email).Scan(&verified)
	if err == sql.ErrNoRows {
		return false, nil // not found → not verified
	}
	return verified, err
}
