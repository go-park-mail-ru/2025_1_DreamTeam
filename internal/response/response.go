package response

import (
	"encoding/json"
	"net/http"
	"skillForce/internal/errors"
)

// SendErrorResponse - отправка ошибки в JSON-формате
func SendErrorResponse(textError string, headerStatus int, w http.ResponseWriter) {
	response := errors.ErrorResponse{Error: textError}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(headerStatus)
	json.NewEncoder(w).Encode(response)
	// Вопрос: не лучше ли использовать http.Error()???
}

// SendOKResponse - отправка пустого ответа со статусом 200 OK
func SendOKResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
