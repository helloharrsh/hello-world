package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"mailer_application/internal/config"
	httpDelivery "mailer_application/internal/delivery/https"
	"mailer_application/internal/infrastructure/db"
	"mailer_application/internal/infrastructure/mail"
	"mailer_application/internal/service"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using system env.")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to MySQL
	conn, err := db.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer conn.Close()

	// Initialize dependencies
	mailer := mail.NewMailer(cfg)
	repo := db.NewRepository(conn)
	otpService := service.NewOTPService(repo, mailer)

	// Setup main router
	r := mux.NewRouter()

	// Mount API routes under /api/
	apiRouter := httpDelivery.NewRouter(otpService)
	r.PathPrefix("/api/").Handler(http.StripPrefix("/api", apiRouter))

	// Serve static frontend files
	fs := http.FileServer(http.Dir("./frontend"))
	r.PathPrefix("/").Handler(fs)

	// Start server
	log.Printf("✅ Server running at http://localhost:%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
