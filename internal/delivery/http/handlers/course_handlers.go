package handlers

import (
	"fmt"
	"net/http"
	"skillForce/internal/delivery/http/response"
	"skillForce/internal/usecase"
	"skillForce/pkg/logs"
)

// CourseHandler - структура обработчика HTTP-запросов
type CourseHandler struct {
	useCase usecase.CourseUsecaseInterface
}

// NewCourseHandler - конструктор
func NewCourseHandler(uc *usecase.CourseUsecase) *CourseHandler {
	return &CourseHandler{useCase: uc}
}

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
func (h *CourseHandler) GetCourses(w http.ResponseWriter, r *http.Request) {
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
