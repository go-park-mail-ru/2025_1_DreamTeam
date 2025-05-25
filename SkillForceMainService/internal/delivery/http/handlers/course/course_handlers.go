package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	coursepb "skillForce/internal/delivery/grpc/proto/course"
	"skillForce/internal/delivery/http/response"
	"skillForce/internal/models/dto"
	models "skillForce/internal/models/user"
	"skillForce/pkg/logs"
	"strconv"

	"strings"

	"github.com/mailru/easyjson"
	"google.golang.org/grpc"
)

type CookieManagerInterface interface {
	CheckCookie(r *http.Request) *models.UserProfile
}

type VideoManagerInterface interface {
	GetVideoUrl(ctx context.Context, lesson_id int) (string, error)
	GetMeta(ctx context.Context, name string) (dto.VideoMeta, error)
	GetFragment(ctx context.Context, name string, start, end int64) (io.ReadCloser, error)
}

type Handler struct {
	courseClient  coursepb.CourseServiceClient
	cookieManager CookieManagerInterface
	videoManager  VideoManagerInterface
}

func NewHandler(cookieManager CookieManagerInterface, videoManager VideoManagerInterface) *Handler {
	conn, err := grpc.Dial("course-service:8082", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	courseClient := coursepb.NewCourseServiceClient(conn)
	return &Handler{
		courseClient:  courseClient,
		cookieManager: cookieManager,
	}
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
func (h *Handler) GetCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetCourses", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)

	var grpcUserProfile *coursepb.UserProfile
	if userProfile != nil {
		grpcUserProfile = &coursepb.UserProfile{
			Id:        int32(userProfile.Id),
			Email:     userProfile.Email,
			Bio:       userProfile.Bio,
			Name:      userProfile.Name,
			AvatarSrc: userProfile.AvatarSrc,
			HideEmail: userProfile.HideEmail,
			IsAdmin:   userProfile.IsAdmin,
		}
	}

	grpcGetBucketcourses := coursepb.GetBucketCoursesRequest{
		UserProfile: grpcUserProfile,
	}

	grpcBucketCoursesResponse, err := h.courseClient.GetBucketCourses(r.Context(), &grpcGetBucketcourses)
	if err != nil {
		logs.PrintLog(r.Context(), "GetCourses", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	bucketCourses := make([]*dto.CourseDTO, len(grpcBucketCoursesResponse.Courses))
	for i, grpcBucketCourse := range grpcBucketCoursesResponse.Courses {
		bucketCourses[i] = &dto.CourseDTO{
			Id:              int(grpcBucketCourse.Id),
			Title:           grpcBucketCourse.Title,
			ScrImage:        grpcBucketCourse.ScrImage,
			Tags:            grpcBucketCourse.Tags,
			Rating:          float32(grpcBucketCourse.Rating),
			TimeToPass:      int(grpcBucketCourse.TimeToPass),
			PurchasesAmount: int(grpcBucketCourse.PurchasesAmount),
			IsPurchased:     grpcBucketCourse.IsPurchased,
			IsFavorite:      grpcBucketCourse.IsFavorite,
			CreatorId:       int(grpcBucketCourse.CreatorId),
			Description:     grpcBucketCourse.Description,
			Price:           int(grpcBucketCourse.Price),
		}
	}

	logs.PrintLog(r.Context(), "GetCourses", "send bucket courses")
	response.SendBucketCoursesResponse(bucketCourses, w, r)
}

func (h *Handler) GetPurchasedCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetPurchasedCourses", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "GetPurchasedCourses", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	grpcUserProfile := &coursepb.UserProfile{
		Id:        int32(userProfile.Id),
		Email:     userProfile.Email,
		Bio:       userProfile.Bio,
		Name:      userProfile.Name,
		AvatarSrc: userProfile.AvatarSrc,
		HideEmail: userProfile.HideEmail,
		IsAdmin:   userProfile.IsAdmin,
	}

	grpcGetBucketcourses := coursepb.GetBucketCoursesRequest{
		UserProfile: grpcUserProfile,
	}

	grpcBucketCoursesResponse, err := h.courseClient.GetPurchasedBucketCourses(r.Context(), &grpcGetBucketcourses)
	if err != nil {
		logs.PrintLog(r.Context(), "GetCourses", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	if len(grpcBucketCoursesResponse.Courses) == 0 {
		logs.PrintLog(r.Context(), "GetPurchasedCourses", "send purchased bucket courses")
		response.SendNoContentOKResponse(w, r)
		return
	}

	bucketCourses := make([]*dto.CourseDTO, len(grpcBucketCoursesResponse.Courses))
	for i, grpcBucketCourse := range grpcBucketCoursesResponse.Courses {
		bucketCourses[i] = &dto.CourseDTO{
			Id:              int(grpcBucketCourse.Id),
			Title:           grpcBucketCourse.Title,
			ScrImage:        grpcBucketCourse.ScrImage,
			Tags:            grpcBucketCourse.Tags,
			Rating:          float32(grpcBucketCourse.Rating),
			TimeToPass:      int(grpcBucketCourse.TimeToPass),
			PurchasesAmount: int(grpcBucketCourse.PurchasesAmount),
			IsPurchased:     grpcBucketCourse.IsPurchased,
			IsFavorite:      grpcBucketCourse.IsFavorite,
			CreatorId:       int(grpcBucketCourse.CreatorId),
			Description:     grpcBucketCourse.Description,
			Price:           int(grpcBucketCourse.Price),
		}
	}

	logs.PrintLog(r.Context(), "GetPurchasedCourses", "send purchased bucket courses")
	response.SendBucketCoursesResponse(bucketCourses, w, r)
}

func (h *Handler) GetCompletedCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetCompletedCourses", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "GetCompletedCourses", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	grpcUserProfile := &coursepb.UserProfile{
		Id:        int32(userProfile.Id),
		Email:     userProfile.Email,
		Bio:       userProfile.Bio,
		Name:      userProfile.Name,
		AvatarSrc: userProfile.AvatarSrc,
		HideEmail: userProfile.HideEmail,
		IsAdmin:   userProfile.IsAdmin,
	}

	grpcGetBucketcourses := coursepb.GetBucketCoursesRequest{
		UserProfile: grpcUserProfile,
	}

	grpcBucketCoursesResponse, err := h.courseClient.GetCompletedBucketCourses(r.Context(), &grpcGetBucketcourses)
	if err != nil {
		logs.PrintLog(r.Context(), "GetCompletedCourses", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	if len(grpcBucketCoursesResponse.Courses) == 0 {
		logs.PrintLog(r.Context(), "GetCompletedCourses", "send NoContentOKResponse")
		response.SendNoContentOKResponse(w, r)
		return
	}

	bucketCourses := make([]*dto.CourseDTO, len(grpcBucketCoursesResponse.Courses))
	for i, grpcBucketCourse := range grpcBucketCoursesResponse.Courses {
		bucketCourses[i] = &dto.CourseDTO{
			Id:              int(grpcBucketCourse.Id),
			Title:           grpcBucketCourse.Title,
			ScrImage:        grpcBucketCourse.ScrImage,
			Tags:            grpcBucketCourse.Tags,
			Rating:          float32(grpcBucketCourse.Rating),
			TimeToPass:      int(grpcBucketCourse.TimeToPass),
			PurchasesAmount: int(grpcBucketCourse.PurchasesAmount),
			IsPurchased:     grpcBucketCourse.IsPurchased,
			IsFavorite:      grpcBucketCourse.IsFavorite,
			CreatorId:       int(grpcBucketCourse.CreatorId),
			Description:     grpcBucketCourse.Description,
			Price:           int(grpcBucketCourse.Price),
		}
	}

	logs.PrintLog(r.Context(), "GetCompletedCourses", "send completed bucket courses")
	response.SendBucketCoursesResponse(bucketCourses, w, r)
}

// SearchCourses godoc
// @Summary Search courses by title
// @Description Retrieves a list of available courses by title
// @Tags courses
// @Accept json
// @Produce json
// @Param keywords query string true "Keywords to search courses"
// @Success 200 {object} response.BucketCoursesResponse "List of courses"
// @Failure 405 {object} response.ErrorResponse "method not allowed"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /api/searchCourses [get]
func (h *Handler) SearchCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetCourses", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)

	var grpcUserProfile *coursepb.UserProfile
	if userProfile != nil {
		grpcUserProfile = &coursepb.UserProfile{
			Id:        int32(userProfile.Id),
			Email:     userProfile.Email,
			Bio:       userProfile.Bio,
			Name:      userProfile.Name,
			AvatarSrc: userProfile.AvatarSrc,
			HideEmail: userProfile.HideEmail,
			IsAdmin:   userProfile.IsAdmin,
		}
	}

	keywords := r.URL.Query().Get("keywords")

	grpcSearchCourses := coursepb.SearchCoursesByTitleRequest{
		UserProfile: grpcUserProfile,
		Keywords:    keywords,
	}

	grpcSearchCoursesResponse, err := h.courseClient.SearchCoursesByTitle(r.Context(), &grpcSearchCourses)
	if err != nil {
		logs.PrintLog(r.Context(), "GetCourses", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	bucketCourses := make([]*dto.CourseDTO, len(grpcSearchCoursesResponse.Courses))
	for i, grpcBucketCourse := range grpcSearchCoursesResponse.Courses {
		bucketCourses[i] = &dto.CourseDTO{
			Id:              int(grpcBucketCourse.Id),
			Title:           grpcBucketCourse.Title,
			ScrImage:        grpcBucketCourse.ScrImage,
			Tags:            grpcBucketCourse.Tags,
			Rating:          float32(grpcBucketCourse.Rating),
			TimeToPass:      int(grpcBucketCourse.TimeToPass),
			PurchasesAmount: int(grpcBucketCourse.PurchasesAmount),
			IsPurchased:     grpcBucketCourse.IsPurchased,
			IsFavorite:      grpcBucketCourse.IsFavorite,
			CreatorId:       int(grpcBucketCourse.CreatorId),
			Description:     grpcBucketCourse.Description,
			Price:           int(grpcBucketCourse.Price),
		}
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

	userProfile := h.cookieManager.CheckCookie(r)

	var grpcUserProfile *coursepb.UserProfile
	if userProfile != nil {
		grpcUserProfile = &coursepb.UserProfile{
			Id:        int32(userProfile.Id),
			Email:     userProfile.Email,
			Bio:       userProfile.Bio,
			Name:      userProfile.Name,
			AvatarSrc: userProfile.AvatarSrc,
			HideEmail: userProfile.HideEmail,
			IsAdmin:   userProfile.IsAdmin,
		}
	}

	grpcGetCourseRequest := coursepb.GetCourseRequest{
		CourseId:    int32(courseId),
		UserProfile: grpcUserProfile,
	}

	grpcGetCourseResponse, err := h.courseClient.GetCourse(r.Context(), &grpcGetCourseRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "GetCourse", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	course := &dto.CourseDTO{
		Id:              int(grpcGetCourseResponse.Course.Id),
		Title:           grpcGetCourseResponse.Course.Title,
		ScrImage:        grpcGetCourseResponse.Course.ScrImage,
		Tags:            grpcGetCourseResponse.Course.Tags,
		Rating:          float32(grpcGetCourseResponse.Course.Rating),
		TimeToPass:      int(grpcGetCourseResponse.Course.TimeToPass),
		PurchasesAmount: int(grpcGetCourseResponse.Course.PurchasesAmount),
		IsPurchased:     grpcGetCourseResponse.Course.IsPurchased,
		IsCompleted:     grpcGetCourseResponse.Course.IsCompleted,
		IsFavorite:      grpcGetCourseResponse.Course.IsFavorite,
		CreatorId:       int(grpcGetCourseResponse.Course.CreatorId),
		Description:     grpcGetCourseResponse.Course.Description,
		Price:           int(grpcGetCourseResponse.Course.Price),
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

	userProfile := h.cookieManager.CheckCookie(r)
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

	grpcGetCourseLessonRequest := coursepb.GetCourseLessonRequest{
		CourseId: int32(courseId),
		UserId:   int32(userProfile.Id),
	}
	grpcGetCourseLessonResponse, err := h.courseClient.GetCourseLesson(r.Context(), &grpcGetCourseLessonRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "GetCourseLesson", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	lesson := dto.LessonDTO{
		LessonHeader: dto.LessonDtoHeader{
			CourseTitle: grpcGetCourseLessonResponse.Lesson.Header.CourseTitle,
			CourseId:    int(grpcGetCourseLessonResponse.Lesson.Header.CourseId),
			Part: struct {
				Order int    `json:"order"`
				Title string `json:"title"`
			}{
				Order: int(grpcGetCourseLessonResponse.Lesson.Header.Part.Order),
				Title: grpcGetCourseLessonResponse.Lesson.Header.Part.Title,
			},
			Bucket: struct {
				Order int    `json:"order"`
				Title string `json:"title"`
			}{
				Order: int(grpcGetCourseLessonResponse.Lesson.Header.Bucket.Order),
				Title: grpcGetCourseLessonResponse.Lesson.Header.Bucket.Title,
			},
		},
	}

	for _, grpcPoint := range grpcGetCourseLessonResponse.Lesson.Header.Points {
		point := struct {
			LessonId int    `json:"lesson_id"`
			Type     string `json:"type"`
			IsDone   bool   `json:"is_done"`
		}{
			LessonId: int(grpcPoint.LessonId),
			Type:     grpcPoint.Type,
			IsDone:   grpcPoint.IsDone,
		}
		lesson.LessonHeader.Points = append(lesson.LessonHeader.Points, point)
	}

	for _, grpcBlock := range grpcGetCourseLessonResponse.Lesson.Body.Blocks {
		block := struct {
			Body string `json:"body"`
		}{
			Body: grpcBlock.Body,
		}
		lesson.LessonBody.Blocks = append(lesson.LessonBody.Blocks, block)
	}

	lesson.LessonBody.Footer = struct {
		NextLessonId     int `json:"next_lesson_id"`
		CurrentLessonId  int `json:"current_lesson_id"`
		PreviousLessonId int `json:"previous_lesson_id"`
	}{
		NextLessonId:     int(grpcGetCourseLessonResponse.Lesson.Body.Footer.NextLessonId),
		CurrentLessonId:  int(grpcGetCourseLessonResponse.Lesson.Body.Footer.CurrentLessonId),
		PreviousLessonId: int(grpcGetCourseLessonResponse.Lesson.Body.Footer.PreviousLessonId),
	}

	logs.PrintLog(r.Context(), "GetCourseLesson", "send course lesson")
	response.SendLesson(&lesson, w, r)
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

	userProfile := h.cookieManager.CheckCookie(r)
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

	grpcGetNextLessonRequest := coursepb.GetNextLessonRequest{
		UserId:   int32(userProfile.Id),
		CourseId: int32(courseId),
		LessonId: int32(lessonId),
	}

	grpcGetNextLessonResponse, err := h.courseClient.GetNextLesson(r.Context(), &grpcGetNextLessonRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "GetNextLesson", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	lesson := dto.LessonDTO{
		LessonHeader: dto.LessonDtoHeader{
			CourseTitle: grpcGetNextLessonResponse.Lesson.Header.CourseTitle,
			CourseId:    int(grpcGetNextLessonResponse.Lesson.Header.CourseId),
			Part: struct {
				Order int    `json:"order"`
				Title string `json:"title"`
			}{
				Order: int(grpcGetNextLessonResponse.Lesson.Header.Part.Order),
				Title: grpcGetNextLessonResponse.Lesson.Header.Part.Title,
			},
			Bucket: struct {
				Order int    `json:"order"`
				Title string `json:"title"`
			}{
				Order: int(grpcGetNextLessonResponse.Lesson.Header.Bucket.Order),
				Title: grpcGetNextLessonResponse.Lesson.Header.Bucket.Title,
			},
		},
	}

	for _, grpcPoint := range grpcGetNextLessonResponse.Lesson.Header.Points {
		point := struct {
			LessonId int    `json:"lesson_id"`
			Type     string `json:"type"`
			IsDone   bool   `json:"is_done"`
		}{
			LessonId: int(grpcPoint.LessonId),
			Type:     grpcPoint.Type,
			IsDone:   grpcPoint.IsDone,
		}
		lesson.LessonHeader.Points = append(lesson.LessonHeader.Points, point)
	}

	for _, grpcBlock := range grpcGetNextLessonResponse.Lesson.Body.Blocks {
		block := struct {
			Body string `json:"body"`
		}{
			Body: grpcBlock.Body,
		}
		lesson.LessonBody.Blocks = append(lesson.LessonBody.Blocks, block)
	}

	lesson.LessonBody.Footer = struct {
		NextLessonId     int `json:"next_lesson_id"`
		CurrentLessonId  int `json:"current_lesson_id"`
		PreviousLessonId int `json:"previous_lesson_id"`
	}{
		NextLessonId:     int(grpcGetNextLessonResponse.Lesson.Body.Footer.NextLessonId),
		CurrentLessonId:  int(grpcGetNextLessonResponse.Lesson.Body.Footer.CurrentLessonId),
		PreviousLessonId: int(grpcGetNextLessonResponse.Lesson.Body.Footer.PreviousLessonId),
	}

	logs.PrintLog(r.Context(), "GetNextLesson", "send lesson body to user")
	response.SendLesson(&lesson, w, r)
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
// @Failure      405 {object} response.ErrorResponse "method not allowed"
// @Failure      500 {object} response.ErrorResponse "internal server error"
// @Router       /api/markLessonAsNotCompleted [post]
func (h *Handler) MarkLessonAsNotCompleted(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.PrintLog(r.Context(), "MarkLessonAsNotCompleted", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "MarkLessonAsNotCompleted", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	logs.PrintLog(r.Context(), "MarkLessonAsNotCompleted", fmt.Sprintf("user %+v is authorized", userProfile))

	lessonId := dto.LessonIDRequest{}
	if err := easyjson.UnmarshalFromReader(r.Body, &lessonId); err != nil {
		logs.PrintLog(r.Context(), "MarkLessonAsNotCompleted", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	grpcMarkLessonAsNotCompletedRequest := &coursepb.MarkLessonAsNotCompletedRequest{
		UserId:   int32(userProfile.Id),
		LessonId: int32(lessonId.Id),
	}
	_, err := h.courseClient.MarkLessonAsNotCompleted(r.Context(), grpcMarkLessonAsNotCompletedRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "MarkLessonAsNotCompleted", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	response.SendOKResponse(w, r)
}

// MarkLessonAsCompleted godoc
// @Summary      Mark a lesson as not completed
// @Description  Marks the specified lesson as not completed for the authenticated user
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        lessonId body dto.LessonIDRequest true "Lesson ID"
// @Success      200 {object} string "OK"
// @Failure      400 {object} response.ErrorResponse "ivalid lesson ID"
// @Failure      401 {object} response.ErrorResponse "unauthorized"
// @Failure      405 {object} response.ErrorResponse "method not allowed"
// @Failure      500 {object} response.ErrorResponse "internal server error"
// @Router       /api/markLessonAsCompleted [post]
func (h *Handler) MarkLessonAsCompleted(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.PrintLog(r.Context(), "MarkLessonAsCompleted", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "MarkLessonAsCompleted", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	logs.PrintLog(r.Context(), "MarkLessonAsCompleted", fmt.Sprintf("user %+v is authorized", userProfile))

	lessonId := dto.LessonIDRequest{}
	if err := easyjson.UnmarshalFromReader(r.Body, &lessonId); err != nil {
		logs.PrintLog(r.Context(), "MarkLessonAsCompleted", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	grpcMarkLessonAsCompletedRequest := &coursepb.MarkLessonAsCompletedRequest{
		UserId:   int32(userProfile.Id),
		LessonId: int32(lessonId.Id),
	}
	_, err := h.courseClient.MarkLessonAsCompleted(r.Context(), grpcMarkLessonAsCompletedRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "MarkLessonAsCompleted", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	response.SendOKResponse(w, r)
}

// MarkCourseAsCompleted godoc
// @Summary Mark a course as completed
// @Description Marks the specified course as completed for the authenticated user
// @Tags courses
// @Accept json
// @Produce json
// @Param courseId body dto.LessonIDRequest true "Course ID"
// @Success 200 {object} string "OK"
// @Failure 400 {object} response.ErrorResponse "invalid course ID"
// @Failure 401 {object} response.ErrorResponse "unauthorized"
// @Failure 405 {object} response.ErrorResponse "method not allowed"
// @Failure 500 {object} response.ErrorResponse "internal server error"
// @Router /api/markCourseAsCompleted [post]
func (h *Handler) MarkCourseAsCompleted(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.PrintLog(r.Context(), "MarkCourseAsCompleted", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "MarkCourseAsCompleted", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	logs.PrintLog(r.Context(), "MarkCourseAsCompleted", fmt.Sprintf("user %+v is authorized", userProfile))

	courseId := dto.CourseIDRequest{}
	if err := easyjson.UnmarshalFromReader(r.Body, &courseId); err != nil {
		logs.PrintLog(r.Context(), "MarkCourseAsCompleted", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	logs.PrintLog(r.Context(), "MarkCourseAsCompleted", fmt.Sprintf("user %+v wants to mark course %+v as completed", userProfile, courseId))

	grpcMarkCourseAsCompletedRequest := &coursepb.MarkCourseAsCompletedRequest{
		UserId:   int32(userProfile.Id),
		CourseId: int32(courseId.Id),
	}
	_, err := h.courseClient.MarkCourseAsCompleted(r.Context(), grpcMarkCourseAsCompletedRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "MarkCourseAsCompleted", fmt.Sprintf("%+v", err))
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

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "GetCourseRoadmap", "user not logged in")
		userProfile = &models.UserProfile{Id: -1}
	}

	logs.PrintLog(r.Context(), "GetCourseRoadmap", fmt.Sprintf("user %+v is authorized", userProfile))

	courseIdStr := r.URL.Query().Get("courseId")
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		logs.PrintLog(r.Context(), "GetCourseRoadmap", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	grpcGetCourseRoadmapRequest := &coursepb.GetCourseRoadmapRequest{
		UserId:   int32(userProfile.Id),
		CourseId: int32(courseId),
	}
	grpcGetCourseRoadmapResponse, err := h.courseClient.GetCourseRoadmap(r.Context(), grpcGetCourseRoadmapRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "GetCourseRoadmap", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	courseRoadmap := dto.CourseRoadmapDTO{}
	for _, grpcPart := range grpcGetCourseRoadmapResponse.Roadmap.Parts {
		part := &dto.CoursePartDTO{
			Id:    int(grpcPart.Id),
			Title: grpcPart.Title,
		}

		for _, grpcBucket := range grpcPart.Buckets {
			bucket := &dto.LessonBucketDTO{
				Id:    int(grpcBucket.Id),
				Title: grpcBucket.Title,
			}

			for _, grpcLesson := range grpcBucket.Lessons {
				lesson := &dto.LessonPointDTO{
					LessonId: int(grpcLesson.LessonId),
					Type:     grpcLesson.Type,
					Title:    grpcLesson.Title,
					Value:    grpcLesson.Value,
					IsDone:   grpcLesson.IsDone,
				}
				bucket.Lessons = append(bucket.Lessons, lesson)
			}

			part.Buckets = append(part.Buckets, bucket)
		}

		courseRoadmap.Parts = append(courseRoadmap.Parts, part)
	}

	logs.PrintLog(r.Context(), "GetCourseRoadmap", "send course roadmap to user")
	response.SendCourseRoadmap(&courseRoadmap, w, r)

}

func (h *Handler) GetRating(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetRating", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "GetRating", "user not logged in")
		userProfile = &models.UserProfile{Id: -1}
	}

	logs.PrintLog(r.Context(), "GetRating", fmt.Sprintf("user %+v is authorized", userProfile))

	courseIdStr := r.URL.Query().Get("courseId")
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		logs.PrintLog(r.Context(), "GetRating", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	grpcGetRatingsRequest := &coursepb.GetRatingRequest{
		UserId:   int32(userProfile.Id),
		CourseId: int32(courseId),
	}

	grpcGetRatingsResponse, err := h.courseClient.GetRating(r.Context(), grpcGetRatingsRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "GetRaiting", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	ratingItems := make([]dto.RaitingItem, 0, len(grpcGetRatingsResponse.GetRating()))
	for _, item := range grpcGetRatingsResponse.GetRating() {
		user := item.GetUser()
		if user == nil {
			continue
		}

		ratingItems = append(ratingItems, dto.RaitingItem{
			User: dto.UserProfileDTO{
				Name:      user.GetName(),
				Email:     user.GetEmail(),
				Bio:       user.GetBio(),
				AvatarSrc: user.GetAvatarSrc(),
				HideEmail: user.GetHideEmail(),
				IsAdmin:   user.GetIsAdmin(),
			},
			Rating: int(item.GetRating()),
		})
	}

	rating := dto.Raiting{Rating: ratingItems}
	logs.PrintLog(r.Context(), "GetRaiting", "get rating from grpc")
	response.SendRatingResponse(&rating, w, r)
}

func (h *Handler) GetSertificate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetSertificate", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "GetSertificate", "user not logged in")
		userProfile = &models.UserProfile{Id: -1}
	}

	logs.PrintLog(r.Context(), "GetSertificate", fmt.Sprintf("user %+v is authorized", userProfile))

	courseIdStr := r.URL.Query().Get("courseId")
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		logs.PrintLog(r.Context(), "GetSertificate", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	var grpcUserProfile *coursepb.UserProfile
	if userProfile != nil {
		grpcUserProfile = &coursepb.UserProfile{
			Id:        int32(userProfile.Id),
			Email:     userProfile.Email,
			Bio:       userProfile.Bio,
			Name:      userProfile.Name,
			AvatarSrc: userProfile.AvatarSrc,
			HideEmail: userProfile.HideEmail,
			IsAdmin:   userProfile.IsAdmin,
		}
	}

	grpcGetSertificateRequest := &coursepb.GetSertificateRequest{
		User:     grpcUserProfile,
		CourseId: int32(courseId),
	}

	grpcGetSertificateResponse, err := h.courseClient.GetSertificate(r.Context(), grpcGetSertificateRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "GetSertificate", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	sertificateUrl := grpcGetSertificateResponse.SertificateUrl
	logs.PrintLog(r.Context(), "GetSertificate", fmt.Sprintf("get sertificate url: %s from grpc", sertificateUrl))
	response.SendSertificateUrl(sertificateUrl, w, r)
}

func (h *Handler) GetGeneratedSertificate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetSertificate", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "GetSertificate", "user not logged in")
		userProfile = &models.UserProfile{Id: -1}
	}

	logs.PrintLog(r.Context(), "GetSertificate", fmt.Sprintf("user %+v is authorized", userProfile))

	courseIdStr := r.URL.Query().Get("courseId")
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		logs.PrintLog(r.Context(), "GetSertificate", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	var grpcUserProfile *coursepb.UserProfile
	if userProfile != nil {
		grpcUserProfile = &coursepb.UserProfile{
			Id:        int32(userProfile.Id),
			Email:     userProfile.Email,
			Bio:       userProfile.Bio,
			Name:      userProfile.Name,
			AvatarSrc: userProfile.AvatarSrc,
			HideEmail: userProfile.HideEmail,
			IsAdmin:   userProfile.IsAdmin,
		}
	}

	grpcGetSertificateRequest := &coursepb.GetSertificateRequest{
		User:     grpcUserProfile,
		CourseId: int32(courseId),
	}

	grpcGetSertificateResponse, err := h.courseClient.GetGeneratedSertificate(r.Context(), grpcGetSertificateRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "GetSertificate", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	sertificateUrl := grpcGetSertificateResponse.SertificateUrl
	logs.PrintLog(r.Context(), "GetSertificate", fmt.Sprintf("get sertificate url: %s from grpc", sertificateUrl))
	response.SendSertificateUrl(sertificateUrl, w, r)
}

func (h *Handler) GetStatistic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetStatistic", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "GetStatistic", "user not logged in")
		userProfile = &models.UserProfile{Id: -1}
	}

	logs.PrintLog(r.Context(), "GetSertificate", fmt.Sprintf("user %+v is authorized", userProfile))

	courseIdStr := r.URL.Query().Get("courseId")
	courseId, err := strconv.Atoi(courseIdStr)
	if err != nil {
		logs.PrintLog(r.Context(), "GetStatistic", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	grpcGetStatisticRequest := &coursepb.GetStatisticRequest{
		CourseId: int32(courseId),
		UserId:   int32(userProfile.Id),
	}

	grpcGetStatisticResponse, err := h.courseClient.GetStatistic(r.Context(), grpcGetStatisticRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "GetStatistic", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	statistic := &dto.UserStats{
		Percentage:            int(grpcGetStatisticResponse.Percentage),
		CompletedTextLessons:  int(grpcGetStatisticResponse.CompletedTextLessons),
		AmountTextLessons:     int(grpcGetStatisticResponse.AmountTextLessons),
		CompletedVideoLessons: int(grpcGetStatisticResponse.CompletedVideoLessons),
		AmountVideoLessons:    int(grpcGetStatisticResponse.AmountVideoLessons),
		RecievedPoints:        int(grpcGetStatisticResponse.ReceivedPoints),
		AmountPoints:          int(grpcGetStatisticResponse.AmountPoints),
		CompletedTests:        int(grpcGetStatisticResponse.CompletedTests),
		AmountTests:           int(grpcGetStatisticResponse.AmountTests),
		CompletedQuestions:    int(grpcGetStatisticResponse.CompletedQuestions),
		AmountQuestions:       int(grpcGetStatisticResponse.AmountQuestions),
	}

	logs.PrintLog(r.Context(), "GetStatistic", "get statistic from grpc")
	response.SendStatistic(statistic, w, r)
}

// ServeVideo godoc
// @Summary Serve video content
// @Description Streams video content for a lesson based on the lesson ID provided in the query parameters.
//
//	If a "Range" header is present, it streams the requested byte range; otherwise, it streams the entire video.
//
// @Tags videos
// @Accept */*
// @Produce video/mp4
// @Param lesson_id query int true "Lesson ID"
// @Success 206 {file} video/mp4 "Partial Content"
// @Failure 400 {object} response.ErrorResponse "Invalid lesson ID parameter"
// @Failure 404 {object} response.ErrorResponse "Video not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/serveVideo [get]
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

	videoSrc, err := h.videoManager.GetVideoUrl(ctx, lesson_id_int)

	if err != nil {
		response.SendErrorResponse("video not found", http.StatusNotFound, w, r)
		return
	}

	name := strings.Split(videoSrc, "/")[len(strings.Split(videoSrc, "/"))-1]

	meta, err := h.videoManager.GetMeta(ctx, name)
	if err != nil {
		response.SendErrorResponse("video not found", http.StatusNotFound, w, r)
		return
	}

	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		reader, err := h.videoManager.GetFragment(ctx, name, 0, meta.Size-1)
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

	reader, err := h.videoManager.GetFragment(ctx, name, start, end)
	if err != nil {
		response.SendErrorResponse("reading frame error"+err.Error(), http.StatusInternalServerError, w, r)
		return
	}
	defer reader.Close()

	response.SendVideoRange(start, end, meta.Size, reader, w, r)
}

func (h *Handler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.PrintLog(r.Context(), "CreateCourse", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "CreateCourse", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	var CourseInput dto.CourseDTO
	if err := easyjson.UnmarshalFromReader(r.Body, &CourseInput); err != nil {
		logs.PrintLog(r.Context(), "CreateCourse", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	var parts []*coursepb.CoursePartDTO
	for _, part := range CourseInput.Parts {
		var buckets []*coursepb.LessonBucketDTO
		for _, bucket := range part.Buckets {
			var lessons []*coursepb.LessonPointDTO
			for _, lesson := range bucket.Lessons {
				lessons = append(lessons, &coursepb.LessonPointDTO{
					LessonId: int32(lesson.LessonId),
					Type:     lesson.Type,
					Title:    lesson.Title,
					Value:    lesson.Value,
					IsDone:   lesson.IsDone,
				})
			}

			buckets = append(buckets, &coursepb.LessonBucketDTO{
				Id:      int32(bucket.Id),
				Title:   bucket.Title,
				Lessons: lessons,
			})
		}

		parts = append(parts, &coursepb.CoursePartDTO{
			Id:      int32(part.Id),
			Title:   part.Title,
			Buckets: buckets,
		})
	}

	grpcCreateCourseRequest := &coursepb.CreateCourseRequest{
		Course: &coursepb.CourseDTO{
			Id:              int32(CourseInput.Id),
			Price:           int32(CourseInput.Price),
			PurchasesAmount: int32(CourseInput.PurchasesAmount),
			CreatorId:       int32(CourseInput.CreatorId),
			TimeToPass:      int32(CourseInput.TimeToPass),
			Rating:          CourseInput.Rating,
			Tags:            CourseInput.Tags,
			Title:           CourseInput.Title,
			Description:     CourseInput.Description,
			IsPurchased:     CourseInput.IsPurchased,
			IsFavorite:      CourseInput.IsFavorite,
			Parts:           parts,
		},
		UserProfile: &coursepb.UserProfile{
			Id:        int32(userProfile.Id),
			Name:      userProfile.Name,
			Email:     userProfile.Email,
			Bio:       userProfile.Bio,
			AvatarSrc: userProfile.AvatarSrc,
			HideEmail: userProfile.HideEmail,
			IsAdmin:   userProfile.IsAdmin,
		},
	}

	_, err := h.courseClient.CreateCourse(r.Context(), grpcCreateCourseRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "CreateCourse", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	response.SendOKResponse(w, r)
}

// AddCourseToFavourites godoc
// @Summary      Add course to user's favorites
// @Description  Adds a course to the authenticated user's favorites list
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        course body dto.CourseDTO true "Course data to add to favorites"
// @Success      200 {object} string "OK"
// @Failure      400 {object} response.ErrorResponse "invalid request body"
// @Failure      401 {object} response.ErrorResponse "user not authorized"
// @Failure      405 {object} response.ErrorResponse "method not allowed"
// @Failure      500 {object} response.ErrorResponse "internal server error"
// @Router       /api/addCourseToFavourites [post]
func (h *Handler) AddCourseToFavourites(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.PrintLog(r.Context(), "AddCourseToFavourites", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "AddCourseToFavourites", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	var CourseInput dto.CourseDTO
	if err := easyjson.UnmarshalFromReader(r.Body, &CourseInput); err != nil {
		logs.PrintLog(r.Context(), "AddCourseToFavourites", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	var parts []*coursepb.CoursePartDTO
	for _, part := range CourseInput.Parts {
		var buckets []*coursepb.LessonBucketDTO
		for _, bucket := range part.Buckets {
			var lessons []*coursepb.LessonPointDTO
			for _, lesson := range bucket.Lessons {
				lessons = append(lessons, &coursepb.LessonPointDTO{
					LessonId: int32(lesson.LessonId),
					Type:     lesson.Type,
					Title:    lesson.Title,
					Value:    lesson.Value,
					IsDone:   lesson.IsDone,
				})
			}

			buckets = append(buckets, &coursepb.LessonBucketDTO{
				Id:      int32(bucket.Id),
				Title:   bucket.Title,
				Lessons: lessons,
			})
		}

		parts = append(parts, &coursepb.CoursePartDTO{
			Id:      int32(part.Id),
			Title:   part.Title,
			Buckets: buckets,
		})
	}

	grpcAddCourseToFavourites := &coursepb.AddToFavouritesRequest{
		Course: &coursepb.CourseDTO{
			Id:              int32(CourseInput.Id),
			Price:           int32(CourseInput.Price),
			PurchasesAmount: int32(CourseInput.PurchasesAmount),
			CreatorId:       int32(CourseInput.CreatorId),
			TimeToPass:      int32(CourseInput.TimeToPass),
			Rating:          CourseInput.Rating,
			Tags:            CourseInput.Tags,
			Title:           CourseInput.Title,
			Description:     CourseInput.Description,
			IsPurchased:     CourseInput.IsPurchased,
			IsFavorite:      CourseInput.IsFavorite,
			Parts:           parts,
		},
		UserProfile: &coursepb.UserProfile{
			Id:        int32(userProfile.Id),
			Name:      userProfile.Name,
			Email:     userProfile.Email,
			Bio:       userProfile.Bio,
			AvatarSrc: userProfile.AvatarSrc,
			HideEmail: userProfile.HideEmail,
			IsAdmin:   userProfile.IsAdmin,
		},
	}

	_, err := h.courseClient.AddCourseToFavourites(r.Context(), grpcAddCourseToFavourites)
	if err != nil {
		logs.PrintLog(r.Context(), "AddCourseToFavourites", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	response.SendOKResponse(w, r)
}

// DeleteCourseFromFavourites godoc
// @Summary      Remove course from favorites
// @Description  Removes a course from the authenticated user's favorites list
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        course body dto.CourseDTO true "Course data to remove from favorites"
// @Success      200 {object} string "OK"
// @Failure      400 {object} response.ErrorResponse "invalid request body"
// @Failure      401 {object} response.ErrorResponse "user not authorized"
// @Failure      405 {object} response.ErrorResponse "method not allowed"
// @Failure      500 {object} response.ErrorResponse "internal server error"
// @Router       /api/deleteCourseFromFavourites [post]
func (h *Handler) DeleteCourseFromFavourites(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.PrintLog(r.Context(), "DeleteCourseFromFavourites", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "DeleteCourseFromFavourites", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	var CourseInput dto.CourseDTO
	if err := easyjson.UnmarshalFromReader(r.Body, &CourseInput); err != nil {
		logs.PrintLog(r.Context(), "DeleteCourseFromFavourites", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	var parts []*coursepb.CoursePartDTO
	for _, part := range CourseInput.Parts {
		var buckets []*coursepb.LessonBucketDTO
		for _, bucket := range part.Buckets {
			var lessons []*coursepb.LessonPointDTO
			for _, lesson := range bucket.Lessons {
				lessons = append(lessons, &coursepb.LessonPointDTO{
					LessonId: int32(lesson.LessonId),
					Type:     lesson.Type,
					Title:    lesson.Title,
					Value:    lesson.Value,
					IsDone:   lesson.IsDone,
				})
			}

			buckets = append(buckets, &coursepb.LessonBucketDTO{
				Id:      int32(bucket.Id),
				Title:   bucket.Title,
				Lessons: lessons,
			})
		}

		parts = append(parts, &coursepb.CoursePartDTO{
			Id:      int32(part.Id),
			Title:   part.Title,
			Buckets: buckets,
		})
	}

	grpcDeleteCourseFromFavourites := &coursepb.DeleteCourseFromFavouritesRequest{
		Course: &coursepb.CourseDTO{
			Id:              int32(CourseInput.Id),
			Price:           int32(CourseInput.Price),
			PurchasesAmount: int32(CourseInput.PurchasesAmount),
			CreatorId:       int32(CourseInput.CreatorId),
			TimeToPass:      int32(CourseInput.TimeToPass),
			Rating:          CourseInput.Rating,
			Tags:            CourseInput.Tags,
			Title:           CourseInput.Title,
			Description:     CourseInput.Description,
			IsPurchased:     CourseInput.IsPurchased,
			IsFavorite:      CourseInput.IsFavorite,
			Parts:           parts,
		},
		UserProfile: &coursepb.UserProfile{
			Id:        int32(userProfile.Id),
			Name:      userProfile.Name,
			Email:     userProfile.Email,
			Bio:       userProfile.Bio,
			AvatarSrc: userProfile.AvatarSrc,
			HideEmail: userProfile.HideEmail,
			IsAdmin:   userProfile.IsAdmin,
		},
	}

	_, err := h.courseClient.DeleteCourseFromFavourites(r.Context(), grpcDeleteCourseFromFavourites)
	if err != nil {
		logs.PrintLog(r.Context(), "DeleteCourseFromFavourites", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	response.SendOKResponse(w, r)
}

// GetFavouriteCourses godoc
// @Summary Get user's favorite courses
// @Description Retrieves all courses marked as favorite by the authenticated user
// @Tags courses
// @Produce json
// @Success      200 {object} string "OK"
// @Failure 401 {object} response.ErrorResponse "User not authorized"
// @Failure 405 {object} response.ErrorResponse "Method not allowed"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/getFavouriteCourses [get]
func (h *Handler) GetFavouriteCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetFavouriteCourses", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "GetFavouriteCourses", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	grpcGetFavouriteCourses := &coursepb.GetFavouritesRequest{
		UserProfile: &coursepb.UserProfile{
			Id:        int32(userProfile.Id),
			Name:      userProfile.Name,
			Email:     userProfile.Email,
			Bio:       userProfile.Bio,
			AvatarSrc: userProfile.AvatarSrc,
			HideEmail: userProfile.HideEmail,
			IsAdmin:   userProfile.IsAdmin,
		},
	}

	grpcGetFavouritesResponse, err := h.courseClient.GetFavouriteCourses(r.Context(), grpcGetFavouriteCourses)
	if err != nil {
		logs.PrintLog(r.Context(), "GetFavouriteCourses", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	bucketCourses := make([]*dto.CourseDTO, len(grpcGetFavouritesResponse.Courses))
	for i, grpcBucketCourse := range grpcGetFavouritesResponse.Courses {
		bucketCourses[i] = &dto.CourseDTO{
			Id:              int(grpcBucketCourse.Id),
			Title:           grpcBucketCourse.Title,
			ScrImage:        grpcBucketCourse.ScrImage,
			Tags:            grpcBucketCourse.Tags,
			Rating:          float32(grpcBucketCourse.Rating),
			TimeToPass:      int(grpcBucketCourse.TimeToPass),
			PurchasesAmount: int(grpcBucketCourse.PurchasesAmount),
			IsPurchased:     grpcBucketCourse.IsPurchased,
			IsFavorite:      grpcBucketCourse.IsFavorite,
			CreatorId:       int(grpcBucketCourse.CreatorId),
			Description:     grpcBucketCourse.Description,
			Price:           int(grpcBucketCourse.Price),
		}
	}

	response.SendBucketCoursesResponse(bucketCourses, w, r)
}

// GetTestLesson godoc
// @Summary Get test lesson details
// @Description Retrieves test lesson information including questions and answers
// @Tags lessons
// @Produce json
// @Param lessonId query int true "Lesson ID to retrieve test for"
// @Success 200 {object} dto.Test "Test lesson details including question and answers"
// @Failure 400 {object} response.ErrorResponse "Invalid lesson ID format"
// @Failure 401 {object} response.ErrorResponse "User not authorized"
// @Failure 404 {object} response.ErrorResponse "Lesson not found"
// @Failure 405 {object} response.ErrorResponse "Method not allowed"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/GetTestLesson [get]
func (h *Handler) GetTestLesson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetTestLesson", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "GetTestLesson", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	lessonIdStr := r.URL.Query().Get("lessonId")
	lessonId, err := strconv.Atoi(lessonIdStr)

	if err != nil {
		logs.PrintLog(r.Context(), "GetTestLesson", "not found")
		response.SendErrorResponse("not found", http.StatusNotFound, w, r)
		return
	}

	grpcGetTestLesson := &coursepb.GetTestLessonRequest{
		LessonId: int32(lessonId),
		UserId:   int32(userProfile.Id),
	}

	grpcGetTestLessonResponse, err := h.courseClient.GetTestLesson(r.Context(), grpcGetTestLesson)
	if err != nil {
		logs.PrintLog(r.Context(), "GetTestLesson", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	answers := make([]*dto.QuizAnswer, len(grpcGetTestLessonResponse.TestDTO.Answers))
	for i, grpcAnswer := range grpcGetTestLessonResponse.TestDTO.Answers {
		answers[i] = &dto.QuizAnswer{
			AnswerID: int64(grpcAnswer.AnswerId),
			Answer:   grpcAnswer.Answer,
			IsRight:  grpcAnswer.IsRight,
		}
	}

	test := &dto.Test{
		QuestionID: int64(grpcGetTestLessonResponse.TestDTO.QuestionId),
		Question:   grpcGetTestLessonResponse.TestDTO.Question,
		Answers:    answers,
		UserAnswer: dto.UserAnswer{
			IsRight:    grpcGetTestLessonResponse.UserAnswer.IsRight,
			QuestionID: int64(grpcGetTestLessonResponse.TestDTO.QuestionId),
			AnswerID:   int64(grpcGetTestLessonResponse.UserAnswer.AnswerId),
		},
	}

	response.SendTestLessonResponse(test, w, r)
}

// AnswerQuiz godoc
// @Summary Submit quiz answer
// @Description Processes user's answer to a quiz question and returns the result
// @Tags lessons
// @Accept json
// @Produce json
// @Param answer body dto.Answer true "User's quiz answer"
// @Success 200 {object} bool "Result of the quiz answer (true/false)"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "User not authorized"
// @Failure 405 {object} response.ErrorResponse "Method not allowed"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/AnswerQuiz [post]
func (h *Handler) AnswerQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.PrintLog(r.Context(), "AnswerQuiz", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "AnswerQuiz", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	var AnswerInput dto.Answer
	if err := easyjson.UnmarshalFromReader(r.Body, &AnswerInput); err != nil {
		logs.PrintLog(r.Context(), "AnswerQuiz", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	grpcAnswerQuiz := &coursepb.AnswerQuizRequest{
		QuestionId: int32(AnswerInput.QuestionID),
		AnswerId:   int32(AnswerInput.Answer_ID),
		UserId:     int32(userProfile.Id),
		CourseId:   int32(AnswerInput.Course_ID),
	}

	grpcAnswerQuizResponse, err := h.courseClient.AnswerQuiz(r.Context(), grpcAnswerQuiz)
	if err != nil {
		logs.PrintLog(r.Context(), "AnswerQuiz", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	result := grpcAnswerQuizResponse.IsRight

	response.SendQuizResult(result, w, r)
}

// GetQuestionTestLesson godoc
// @Summary Get test question for lesson
// @Description Retrieves a test question and user's previous answer (if exists) for a specific lesson
// @Tags lessons
// @Produce json
// @Param lessonId query int true "Lesson ID to retrieve test question for"
// @Success 200 {object} dto.QuestionTest "Test question details including user's previous answer status"
// @Failure 400 {object} response.ErrorResponse "Invalid lesson ID format"
// @Failure 401 {object} response.ErrorResponse "User not authorized"
// @Failure 404 {object} response.ErrorResponse "Lesson not found"
// @Failure 405 {object} response.ErrorResponse "Method not allowed"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/GetQuestionTestLesson [get]
func (h *Handler) GetQuestionTestLesson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "GetQuestionTestLesson", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "GetQuestionTestLesson", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	lessonIdStr := r.URL.Query().Get("lessonId")
	lessonId, err := strconv.Atoi(lessonIdStr)

	if err != nil {
		logs.PrintLog(r.Context(), "GetQuestionTestLesson", "not found")
		response.SendErrorResponse("not found", http.StatusNotFound, w, r)
		return
	}

	grpcGetQuestionTestLesson := &coursepb.GetQuestionTestLessonRequest{
		LessonId: int32(lessonId),
		UserId:   int32(userProfile.Id),
	}

	grpcGetQuestionTestLessonResponse, err := h.courseClient.GetQuestionTestLesson(r.Context(), grpcGetQuestionTestLesson)
	if err != nil {
		logs.PrintLog(r.Context(), "GetQuestionTestLesson", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	question := &dto.QuestionTest{
		QuestionID: int64(grpcGetQuestionTestLessonResponse.QuestionId),
		Question:   grpcGetQuestionTestLessonResponse.Question,
		UserAnswer: dto.UserQuestionAnswer{
			Status: grpcGetQuestionTestLessonResponse.UserAnswer.Status,
			Answer: grpcGetQuestionTestLessonResponse.UserAnswer.Answer,
		},
	}

	response.SendQuestionTestLessonResponse(question, w, r)
}

// AnswerQuestion godoc
// @Summary Submit answer to a question
// @Description Processes and stores a user's answer to a specific test question
// @Tags lessons
// @Accept json
// @Produce json
// @Param answer body dto.AnswerQuestion true "User's answer to the question"
// @Success      200 {object} string "OK"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "User not authorized"
// @Failure 405 {object} response.ErrorResponse "Method not allowed"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/AnswerQuestion [post]
func (h *Handler) AnswerQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.PrintLog(r.Context(), "AnswerQuestion", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "AnswerQuestion", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	var AnswerInput dto.AnswerQuestion
	if err := easyjson.UnmarshalFromReader(r.Body, &AnswerInput); err != nil {
		logs.PrintLog(r.Context(), "AnswerQuestion", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	grpcAnswerQuestion := &coursepb.AnswerQuestionRequest{
		QuestionId: int32(AnswerInput.QuestionID),
		UserId:     int32(userProfile.Id),
		Answer:     AnswerInput.Answer,
	}

	_, err := h.courseClient.AnswerQuestion(r.Context(), grpcAnswerQuestion)
	if err != nil {
		logs.PrintLog(r.Context(), "AnswerQuestion", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	response.SendOKResponse(w, r)
}
