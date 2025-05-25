//go:generate easyjson -all response.go

package response

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"skillForce/internal/models/dto"

	jwriter "github.com/mailru/easyjson/jwriter"
)

//easyjson:json
type ErrorResponse struct {
	ErrorStr string `json:"error"`
}

//easyjson:json
type BucketCoursesResponse struct {
	BucketCourses []*dto.CourseDTO `json:"bucket_courses"`
}

//easyjson:json
type CourseResponse struct {
	Course *dto.CourseDTO `json:"course"`
}

//easyjson:json
type UserProfileResponse struct {
	UserProfile *dto.UserProfileDTO `json:"user"`
}

//easyjson:json
type PhotoUrlResponse struct {
	Url string `json:"url"`
}

//easyjson:json
type SertificateUrlResponse struct {
	Url string `json:"url"`
}

type LessonResponse struct {
	Lesson *dto.LessonDTO `json:"lesson"`
}

//easyjson:json
type LessonBodyResponse struct {
	LessonBody *dto.LessonDtoBody `json:"lesson_body"`
}

//easyjson:json
type CourseRoadmapResponse struct {
	CourseRoadmap *dto.CourseRoadmapDTO `json:"course_roadmap"`
}

//easyjson:json
type SurveyResponse struct {
	Survey *dto.SurveyDTO `json:"survey"`
}

//easyjson:json
type SurveyMetricsResponse struct {
	SurveyMetrics *dto.SurveyMetricsDTO `json:"survey_metrics"`
}

//easyjson:json
type TestResponse struct {
	Test *dto.Test `json:"test"`
}

//easyjson:json
type Result struct {
	Result bool `json:"result"`
}

//easyjson:json
type Billing struct {
	Continue_url string `json:"continue_url"`
}

//easyjson:json
type QuestionTestResponse struct {
	Question *dto.QuestionTest `json:"question"`
}

//easyjson:json
type RaitingResponse struct {
	Raiting *dto.Raiting `json:"course_raiting"`
}

//easyjson:json
type StatisticResponse struct {
	Statistic *dto.UserStats `json:"statistic"`
}

func marshaling(w http.ResponseWriter, response interface{ MarshalEasyJSON(*jwriter.Writer) }) {
	jw := jwriter.Writer{}
	response.MarshalEasyJSON(&jw)

	if _, err := jw.DumpTo(w); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// SendErrorResponse - отправка ошибки в JSON-формате
func SendErrorResponse(textError string, headerStatus int, w http.ResponseWriter, r *http.Request) {
	response := ErrorResponse{ErrorStr: textError}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(headerStatus)
	marshaling(w, response)
}

// SendOKResponse - отправка пустого ответа со статусом 200 OK
func SendOKResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("200 OK")
}

func SendBillingRedirect(w http.ResponseWriter, r *http.Request, continue_url string) {
	response := Billing{Continue_url: continue_url}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}

func SendNoContentOKResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("204 OK")
}

// SendBucketCoursesResponse - отправка списка курсов в JSON-формате
func SendBucketCoursesResponse(bucketCourses []*dto.CourseDTO, w http.ResponseWriter, r *http.Request) {
	response := BucketCoursesResponse{BucketCourses: bucketCourses}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}

func SendTestLessonResponse(test *dto.Test, w http.ResponseWriter, r *http.Request) {
	response := TestResponse{Test: test}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}

func SendQuizResult(result bool, w http.ResponseWriter, r *http.Request) {
	response := Result{Result: result}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}

func SendQuestionTestLessonResponse(question *dto.QuestionTest, w http.ResponseWriter, r *http.Request) {
	response := QuestionTestResponse{Question: question}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}

// SendUserProfile - отправка профиля пользователя в JSON-формате
func SendUserProfile(UserProfile *dto.UserProfileDTO, w http.ResponseWriter, r *http.Request) {
	response := UserProfileResponse{UserProfile: UserProfile}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}

// SendPhotoUrl - отправка ссылки на фото в JSON-формате
func SendPhotoUrl(url string, w http.ResponseWriter, r *http.Request) {
	response := PhotoUrlResponse{Url: url}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}

func SendSertificateUrl(url string, w http.ResponseWriter, r *http.Request) {
	response := SertificateUrlResponse{Url: url}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}

// SendLesson - отправка урока в JSON-формате
func SendLesson(lesson *dto.LessonDTO, w http.ResponseWriter, r *http.Request) {
	response := LessonResponse{Lesson: lesson}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}

func SendLessonBody(lessonBody *dto.LessonDtoBody, w http.ResponseWriter, r *http.Request) {
	response := LessonBodyResponse{LessonBody: lessonBody}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}

func SendCourseRoadmap(courseRoadmap *dto.CourseRoadmapDTO, w http.ResponseWriter, r *http.Request) {
	response := CourseRoadmapResponse{CourseRoadmap: courseRoadmap}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}

func SendCourseResponse(course *dto.CourseDTO, w http.ResponseWriter, r *http.Request) {
	response := CourseResponse{Course: course}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
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
	marshaling(w, response)
}

func SendSurveyMetricsResponse(surveyMetrics *dto.SurveyMetricsDTO, w http.ResponseWriter, r *http.Request) {
	response := SurveyMetricsResponse{SurveyMetrics: surveyMetrics}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}

func SendRatingResponse(raiting *dto.Raiting, w http.ResponseWriter, r *http.Request) {
	response := RaitingResponse{Raiting: raiting}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}

func SendStatistic(statistic *dto.UserStats, w http.ResponseWriter, r *http.Request) {
	response := StatisticResponse{Statistic: statistic}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshaling(w, response)
}
