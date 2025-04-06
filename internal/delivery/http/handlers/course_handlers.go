package handlers

import (
	"fmt"
	"net/http"
	"skillForce/internal/delivery/http/response"
	"skillForce/pkg/logs"
	"strconv"
)

// GetCourses godoc
// @Summary Get list of courses
// @Description Retrieves a list of available courses
// @Tags courses
// @Accept json
// @Produce json
// @Success 200 {object} response.BucketCoursesResponse "List of courses"
// @Failure 405 {object} response.ErrorResponse "method not allowed"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /api/getCourses [get]
func (h *Handler) GetCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetCourses", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	bucketCourses, err := h.useCase.GetBucketCourses(r.Context())
	if err != nil {
		logs.PrintLog(r.Context(), "GetCourses", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	logs.PrintLog(r.Context(), "GetCourses", "send bucket courses")
	response.SendBucketCoursesResponse(bucketCourses, w, r)
}

// GetCourseLesson godoc
// @Summary Get next lesson of a course for the user
// @Description Returns the next lesson the user should take in the specified course
// @Tags courses
// @Accept json
// @Produce json
// @Param courseId query int true "Course ID"
// @Success 200 {object} response.LessonResponse "next lesson of the course"
// @Failure 400 {object} response.ErrorResponse "invalid course ID"
// @Failure 401 {object} response.ErrorResponse "not authorized"
// @Failure 405 {object} response.ErrorResponse "method not allowed"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /api/getCourseLesson [get]

func (h *Handler) GetCourseLesson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetCourseLesson", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.checkCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "GetCourseLesson", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	courseIdStr := r.URL.Query().Get("courseId")
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		logs.PrintLog(r.Context(), "GetCourseLesson", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	lesson, err := h.useCase.GetCourseLesson(r.Context(), userProfile.Id, courseId)
	if err != nil {
		logs.PrintLog(r.Context(), "GetCourseLesson", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	logs.PrintLog(r.Context(), "GetCourseLesson", "send course lesson")
	response.SendLesson(lesson, w, r)
}
