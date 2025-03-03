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

// SendErrorResponse - отправка ошибки в JSON-формате
func SendErrorResponse(textError string, headerStatus int, w http.ResponseWriter) {
	response := ErrorResponse{ErrorStr: textError}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(headerStatus)
	json.NewEncoder(w).Encode(response)
}

// SendOKResponse - отправка пустого ответа со статусом 200 OK
func SendOKResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// SendBucketCoursesResponse - отправка списка курсов в JSON-формате
func SendBucketCoursesResponse(bucketCourses []*models.Course, w http.ResponseWriter) {
	response := BucketCoursesResponse{BucketCourses: bucketCourses}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
