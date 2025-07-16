package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"mailer_application/internal/config"
	httpDelivery "mailer_application/internal/delivery/https"
	"mailer_application/internal/infrastructure/db"
	"mailer_application/internal/infrastructure/mail"
	"mailer_application/internal/service"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using system env.")
	}

	cfg := config.Load()

	// DB connection
	conn, err := db.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer conn.Close()

	// Initialize mailer
	mailer := mail.NewMailer(cfg)

	// Initialize repository
	repo := db.NewRepository(conn)

	// Initialize service
	otpService := service.NewOTPService(repo, mailer)

	// Setup routes
	router := httpDelivery.NewRouter(otpService)

	log.Printf("Server running on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
