package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var ValidPaths = map[string]bool{
	"/api/register":                   true,
	"/api/login":                      true,
	"/api/logout":                     true,
	"/api/isAuthorized":               true,
	"/api/updateProfile":              true,
	"/api/updateProfilePhoto":         true,
	"/api/deleteProfilePhoto":         true,
	"/api/validEmail":                 true,
	"/api/getCourses":                 true,
	"/api/searchCourses":              true,
	"/api/getCourse":                  true,
	"/api/getCourseLesson":            true,
	"/api/getNextLesson":              true,
	"/api/markLessonAsNotCompleted":   true,
	"/api/getCourseRoadmap":           true,
	"/api/video":                      true,
	"/api/createCourse":               true,
	"/api/addCourseToFavourites":      true,
	"/api/deleteCourseFromFavourites": true,
	"/api/getFavouriteCourses":        true,
	"/api/GetTestLesson":              true,
	"/api/AnswerQuiz":                 true,
	"/api/GetQuestionTestLesson":      true,
	"/api/AnswerQuestion":             true,
	"/api/docs/":                      true,
	"/metrics":                        true,
}

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Количество HTTP запросов",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Длительность HTTP запросов",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)

		duration := time.Since(start).Seconds()
		path := r.URL.Path

		if _, ok := ValidPaths[path]; !ok {
			return
		}

		method := r.Method
		status := strconv.Itoa(recorder.statusCode)

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}
