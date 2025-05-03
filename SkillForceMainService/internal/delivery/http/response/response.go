package response

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"skillForce/internal/models/dto"
)

type ErrorResponse struct {
	ErrorStr string `json:"error"`
}

type BucketCoursesResponse struct {
	BucketCourses []*dto.CourseDTO `json:"bucket_courses"`
}

type CourseResponse struct {
	Course *dto.CourseDTO `json:"course"`
}

type UserProfileResponse struct {
	UserProfile *dto.UserProfileDTO `json:"user"`
}

type PhotoUrlResponse struct {
	Url string `json:"url"`
}

type LessonResponse struct {
	Lesson *dto.LessonDTO `json:"lesson"`
}

type LessonBodyResponse struct {
	LessonBody *dto.LessonDtoBody `json:"lesson_body"`
}

type CourseRoadmapResponse struct {
	CourseRoadmap *dto.CourseRoadmapDTO `json:"course_roadmap"`
}

type SurveyResponse struct {
	Survey *dto.SurveyDTO `json:"survey"`
}

type SurveyMetricsResponse struct {
	SurveyMetrics *dto.SurveyMetricsDTO `json:"survey_metrics"`
}

type TestResponse struct {
	Test *dto.Test `json:"test"`
}

type Result struct {
	Result bool `json:"result"`
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
func SendBucketCoursesResponse(bucketCourses []*dto.CourseDTO, w http.ResponseWriter, r *http.Request) {
	response := BucketCoursesResponse{BucketCourses: bucketCourses}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SendTestLessonResponse(test *dto.Test, w http.ResponseWriter, r *http.Request) {
	response := TestResponse{Test: test}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SendQuizResult(result bool, w http.ResponseWriter, r *http.Request) {
	response := Result{Result: result}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SendUserProfile - отправка профиля пользователя в JSON-формате
func SendUserProfile(UserProfile *dto.UserProfileDTO, w http.ResponseWriter, r *http.Request) {
	response := UserProfileResponse{UserProfile: UserProfile}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SendPhotoUrl - отправка ссылки на фото в JSON-формате
func SendPhotoUrl(url string, w http.ResponseWriter, r *http.Request) {
	response := PhotoUrlResponse{Url: url}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SendLesson - отправка урока в JSON-формате
func SendLesson(lesson *dto.LessonDTO, w http.ResponseWriter, r *http.Request) {
	response := LessonResponse{Lesson: lesson}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SendLessonBody(lessonBody *dto.LessonDtoBody, w http.ResponseWriter, r *http.Request) {
	response := LessonBodyResponse{LessonBody: lessonBody}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SendCourseRoadmap(courseRoadmap *dto.CourseRoadmapDTO, w http.ResponseWriter, r *http.Request) {
	response := CourseRoadmapResponse{CourseRoadmap: courseRoadmap}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SendCourseResponse(course *dto.CourseDTO, w http.ResponseWriter, r *http.Request) {
	response := CourseResponse{Course: course}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SendVideoRange(start, end, total int64, reader io.Reader, w http.ResponseWriter, r *http.Request) {
	length := end - start + 1

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, total))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", length))
	w.Header().Set("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusPartialContent)

	buf := make([]byte, 64*1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}
}

func SendSurveyResponse(survey *dto.SurveyDTO, w http.ResponseWriter, r *http.Request) {
	response := SurveyResponse{Survey: survey}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SendSurveyMetricsResponse(surveyMetrics *dto.SurveyMetricsDTO, w http.ResponseWriter, r *http.Request) {
	response := SurveyMetricsResponse{SurveyMetrics: surveyMetrics}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
