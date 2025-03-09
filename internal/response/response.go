package response

import (
	"encoding/json"
	"net/http"
	"skillForce/internal/models"
)

type ErrorResponse struct {
	ErrorStr string `json:"error"`
}

type BucketCoursesResponse struct {
	BucketCourses []*models.Course `json:"bucket_courses"`
}

type UserResponse struct {
	User *models.User `json:"user"`
}

func AddCorsHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8001")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Credentials")
}

// SendErrorResponse - отправка ошибки в JSON-формате
func SendErrorResponse(textError string, headerStatus int, w http.ResponseWriter, r *http.Request) {
	response := ErrorResponse{ErrorStr: textError}
	w.Header().Set("Content-Type", "application/json")
	AddCorsHeaders(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(headerStatus)
	json.NewEncoder(w).Encode(response)
}

// SendOKResponse - отправка пустого ответа со статусом 200 OK
func SendOKResponse(w http.ResponseWriter, r *http.Request) {
	AddCorsHeaders(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Пользователь авторизирован!")
}

// SendBucketCoursesResponse - отправка списка курсов в JSON-формате
func SendBucketCoursesResponse(bucketCourses []*models.Course, w http.ResponseWriter, r *http.Request) {
	response := BucketCoursesResponse{BucketCourses: bucketCourses}
	AddCorsHeaders(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SendUser(user *models.User, w http.ResponseWriter, r *http.Request) {
	response := UserResponse{User: user}
	AddCorsHeaders(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SendCors(w http.ResponseWriter, r *http.Request) {
	AddCorsHeaders(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
}
