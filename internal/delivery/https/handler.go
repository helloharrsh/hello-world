package http

import (
	"encoding/json"
	"net/http"

	"mailer_application/internal/service"
)

type Handler struct {
	OTPService *service.OTPService
}

// RequestOTP handles POST /request-otp
// Expects JSON: { "email": "user@example.com" }
func (h *Handler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	// Decode and validate input
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Generate and send OTP
	err := h.OTPService.GenerateAndSendOTP(req.Email)
	if err != nil {
		if err.Error() == "email is already verified" {
			http.Error(w, "You are already subscribed. No need to verify again.", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		return
	}

	// Success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OTP sent successfully"))
}

// VerifyOTP handles POST /verify-otp
// Expects JSON: { "email": "user@example.com", "otp": "123456" }
func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	// Decode and validate input
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.OTP == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Verify OTP
	valid, err := h.OTPService.VerifyOTP(req.Email, req.OTP)
	if err != nil {
		http.Error(w, "Failed to verify OTP", http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// Success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OTP verified successfully"))
}
