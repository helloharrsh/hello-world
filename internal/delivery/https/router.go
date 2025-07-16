package http

import (
	"net/http"

	"mailer_application/internal/service"

	"github.com/gorilla/mux"
)

func NewRouter(otpService *service.OTPService) http.Handler {
	router := mux.NewRouter()

	handler := &Handler{
		OTPService: otpService,
	}

	router.HandleFunc("/request-otp", handler.RequestOTP).Methods("POST")
	router.HandleFunc("/verify-otp", handler.VerifyOTP).Methods("POST")

	return router
}
