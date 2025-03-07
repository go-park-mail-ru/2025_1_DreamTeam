package handlers

import (
	"log"
	"net/http"
	"skillForce/internal/response"
	"skillForce/internal/usecase"
)

// CourseHandler - структура обработчика HTTP-запросов
type CourseHandler struct {
	useCase *usecase.CourseUsecase
}

// NewCourseHandler - конструктор
func NewCourseHandler(uc *usecase.CourseUsecase) *CourseHandler {
	return &CourseHandler{useCase: uc}
}

// GetCourses - обработчик получения списка курсов
func (h *CourseHandler) GetCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("from getCourses: method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w)
		return
	}

	bucketCourses, err := h.useCase.GetBucketCourses()
	if err != nil {
		log.Printf("from getCourses: %v", err)
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w)
		return
	}

	log.Print("send bucket courses")
	response.SendBucketCoursesResponse(bucketCourses, w)
}
