package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"skillForce/internal/delivery/http/response"
	"skillForce/internal/models"
	"skillForce/internal/models/dto"
	"skillForce/pkg/logs"
	"strconv"
	"strings"
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

// GetCourse godoc
// @Summary Get course
// @Description Retrieves a course by ID
// @Tags courses
// @Accept json
// @Produce json
// @Param courseId query int true "Course ID"
// @Success 200 {object} response.CourseResponse "course"
// @Failure 405 {object} response.ErrorResponse "method not allowed"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /api/getCourses [get]
func (h *Handler) GetCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetCourse", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	courseId, err := strconv.Atoi(r.URL.Query().Get("courseId"))
	if err != nil {
		logs.PrintLog(r.Context(), "GetCourse", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid course ID", http.StatusBadRequest, w, r)
		return
	}

	course, err := h.useCase.GetCourse(r.Context(), courseId)
	if err != nil {
		logs.PrintLog(r.Context(), "GetCourse", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	logs.PrintLog(r.Context(), "GetCourse", "send course")
	response.SendCourseResponse(course, w, r)
}

// GetCourseLesson godoc
// @Summary Get lesson of a course for the user
// @Description Returns the lesson the user should take in the course
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

	logs.PrintLog(r.Context(), "GetCourseLesson", fmt.Sprintf("user %+v is authorized", userProfile))

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

// GetNextLesson godoc
// @Summary Get next lesson in a course
// @Description Returns the next lesson the user should take based on current lesson and course
// @Tags courses
// @Accept json
// @Produce json
// @Param courseId query int true "Course ID"
// @Param lessonId query int true "Current Lesson ID"
// @Success 200 {object} response.LessonResponse "next lesson content"
// @Failure 400 {object} response.ErrorResponse "invalid course or lesson ID"
// @Failure 401 {object} response.ErrorResponse "not authorized"
// @Failure 405 {object} response.ErrorResponse "method not allowed"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /api/getNextLesson [get]
func (h *Handler) GetNextLesson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetNextLesson", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.checkCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "GetNextLesson", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	logs.PrintLog(r.Context(), "GetNextLesson", fmt.Sprintf("user %+v is authorized", userProfile))

	lessonIdStr := r.URL.Query().Get("lessonId")
	lessonId, err := strconv.Atoi(lessonIdStr)
	if err != nil {
		logs.PrintLog(r.Context(), "GetNextLesson", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	courseIdStr := r.URL.Query().Get("courseId")
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		logs.PrintLog(r.Context(), "GetNextLesson", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	lesson, err := h.useCase.GetNextLesson(r.Context(), userProfile.Id, courseId, lessonId)
	if err != nil {
		logs.PrintLog(r.Context(), "GetNextLesson", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	logs.PrintLog(r.Context(), "GetNextLesson", "send lesson body to user")
	response.SendLesson(lesson, w, r)
}

// MarkLessonAsNotCompleted godoc
// @Summary      Mark a lesson as not completed
// @Description  Marks the specified lesson as not completed for the authenticated user
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        lessonId body dto.LessonIDRequest true "Lesson ID"
// @Success      200 {object} string "OK"
// @Failure      400 {object} response.ErrorResponse "ivalid lesson ID"
// @Failure      401 {object} response.ErrorResponse "unauthorized"
// @Failure      405 {object} response.ErrorResponse "uethod not allowed"
// @Failure      500 {object} response.ErrorResponse "internal server error"
// @Router       /api/markLessonAsNotCompleted [post]
func (h *Handler) MarkLessonAsNotCompleted(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.PrintLog(r.Context(), "MarkLessonAsNotCompleted", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.checkCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "MarkLessonAsNotCompleted", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	logs.PrintLog(r.Context(), "MarkLessonAsNotCompleted", fmt.Sprintf("user %+v is authorized", userProfile))

	lessonId := dto.LessonIDRequest{}
	err := json.NewDecoder(r.Body).Decode(&lessonId)
	if err != nil {
		logs.PrintLog(r.Context(), "MarkLessonAsNotCompleted", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	err = h.useCase.MarkLessonAsNotCompleted(r.Context(), userProfile.Id, lessonId.Id)
	if err != nil {
		logs.PrintLog(r.Context(), "MarkLessonAsNotCompleted", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	response.SendOKResponse(w, r)
}

// GetCourseRoadmap godoc
// @Summary      Get course roadmap
// @Description  Returns the roadmap of a course for the authenticated user
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        courseId query int true "Course ID"
// @Success      200 {object} response.CourseRoadmapResponse "Course roadmap"
// @Failure      400 {object} response.ErrorResponse "invalid course ID"
// @Failure      405 {object} response.ErrorResponse "method not allowed"
// @Failure      500 {object} response.ErrorResponse "internal server error"
// @Router       /api/getCourseRoadmap [get]
func (h *Handler) GetCourseRoadmap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetCourseRoadmap", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.checkCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "GetCourseRoadmap", "user not logged in")
		userProfile = &models.UserProfile{Id: -1}
	}

	logs.PrintLog(r.Context(), "GetCourseRoadmap", fmt.Sprintf("user %+v is authorized", userProfile))

	courseIdStr := r.URL.Query().Get("courseId")
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		logs.PrintLog(r.Context(), "GetNextLesson", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	courseRoadmap, err := h.useCase.GetCourseRoadmap(r.Context(), userProfile.Id, courseId)
	if err != nil {
		logs.PrintLog(r.Context(), "GetCourseRoadmap", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	logs.PrintLog(r.Context(), "GetCourseRoadmap", "send course roadmap to user")
	response.SendCourseRoadmap(courseRoadmap, w, r)

}

func (h *Handler) ServeVideo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	lesson_id := r.URL.Query().Get("lesson_id")

	if lesson_id == "" {
		response.SendErrorResponse("not found lesson_id parameter", http.StatusBadRequest, w, r)
		return
	}

	lesson_id_int, err := strconv.Atoi(lesson_id)
	if err != nil {
		response.SendErrorResponse("invalid lesson_id parameter", http.StatusBadRequest, w, r)
		return
	}

	videoSrc, err := h.useCase.GetVideoUrl(ctx, lesson_id_int)

	if err != nil {
		response.SendErrorResponse("video not found", http.StatusNotFound, w, r)
		return
	}

	name := strings.Split(videoSrc, "/")[len(strings.Split(videoSrc, "/"))-1]

	meta, err := h.useCase.GetMeta(ctx, name)
	if err != nil {
		response.SendErrorResponse("video not found", http.StatusNotFound, w, r)
		return
	}

	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		reader, err := h.useCase.GetFragment(ctx, name, 0, meta.Size-1)
		if err != nil {
			response.SendErrorResponse("video getting error"+err.Error(), http.StatusInternalServerError, w, r)
			return
		}
		defer reader.Close()

		response.SendVideoRange(0, meta.Size-1, meta.Size, reader, w, r)
		return
	}

	var start, end int64
	rangeParts := strings.Split(strings.Replace(rangeHeader, "bytes=", "", 1), "-")
	start, _ = strconv.ParseInt(rangeParts[0], 10, 64)
	if rangeParts[1] != "" {
		end, _ = strconv.ParseInt(rangeParts[1], 10, 64)
	} else {
		end = meta.Size - 1
	}
	if end >= meta.Size {
		end = meta.Size - 1
	}

	reader, err := h.useCase.GetFragment(ctx, name, start, end)
	if err != nil {
		response.SendErrorResponse("reading frame error"+err.Error(), http.StatusInternalServerError, w, r)
		return
	}
	defer reader.Close()

	response.SendVideoRange(start, end, meta.Size, reader, w, r)
}
