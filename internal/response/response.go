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

type UserProfileResponse struct {
	UserProfile *models.UserProfile `json:"user"`
}

// SendErrorResponse - отправка ошибки в JSON-формате
func SendErrorResponse(textError string, headerStatus int, w http.ResponseWriter, r *http.Request) {
	response := ErrorResponse{ErrorStr: textError}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(headerStatus)
	json.NewEncoder(w).Encode(response)
}

// SendOKResponse - отправка пустого ответа со статусом 200 OK
func SendOKResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("200 OK")
}

// SendBucketCoursesResponse - отправка списка курсов в JSON-формате
func SendBucketCoursesResponse(bucketCourses []*models.Course, w http.ResponseWriter, r *http.Request) {
	response := BucketCoursesResponse{BucketCourses: bucketCourses}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SendUserProfile - отправка профиля пользователя в JSON-формате
func SendUserProfile(UserProfile *models.UserProfile, w http.ResponseWriter, r *http.Request) {
	response := UserProfileResponse{UserProfile: UserProfile}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
