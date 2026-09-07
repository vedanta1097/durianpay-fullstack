package transport

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteAppError(w http.ResponseWriter, appErr *entity.AppError) {
	WriteJSONError(w, appErr.Code.HTTPStatus(), appErr.Message)
}

func WriteJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Code: status, Message: message}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, err error) {
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	var aErr *entity.AppError
	if errors.As(err, &aErr) {
		WriteAppError(w, aErr)
		return
	}
	// fallback
	WriteJSONError(w, http.StatusInternalServerError, "internal server error")
}
