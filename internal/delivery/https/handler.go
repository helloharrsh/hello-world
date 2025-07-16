package http

import (
	"encoding/json"
	"net/http"

	"mailer_application/internal/service"
)

type Handler struct {
	OTPService *service.OTPService
}

func (h *Handler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := h.OTPService.GenerateAndSendOTP(req.Email)
	if err != nil {
		http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OTP sent successfully"))
}

func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.OTP == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	valid, err := h.OTPService.VerifyOTP(req.Email, req.OTP)
	if err != nil || !valid {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OTP verified successfully"))
}
