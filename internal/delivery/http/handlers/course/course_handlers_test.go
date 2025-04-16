package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handlers "skillForce/internal/delivery/http/handlers/course"
	"skillForce/pkg/logs"

	// "skillForce/internal/logs"
	"skillForce/internal/models/dto"
	usermodels "skillForce/internal/models/user"

	"github.com/golang/mock/gomock"
)

func TestGetCourses(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourses", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().GetBucketCourses(gomock.Any()).Return([]*dto.CourseDTO{
		{Id: 1, Title: "Test Course"},
	}, nil)

	handler.GetCourses(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetCourses_MethodNotAllowed(t *testing.T) {
	handler := handlers.NewHandler(nil, nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodPost, "/api/getCourses", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetCourses(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestGetCourse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourses?courseId=123", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 1}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)
	mockUsecase.EXPECT().GetCourse(gomock.Any(), 123, profile).Return(&dto.CourseDTO{Id: 123}, nil)

	handler.GetCourse(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetCourseLesson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourseLesson?courseId=5", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 42}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)
	mockUsecase.EXPECT().GetCourseLesson(gomock.Any(), 42, 5).Return(&dto.LessonDTO{}, nil)

	handler.GetCourseLesson(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetNextLesson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getNextLesson?courseId=1&lessonId=2", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 10}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)
	mockUsecase.EXPECT().GetNextLesson(gomock.Any(), 10, 1, 2).Return(&dto.LessonDTO{}, nil)

	handler.GetNextLesson(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestMarkLessonAsNotCompleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	body := dto.LessonIDRequest{Id: 7}
	bodyBytes, _ := json.Marshal(body)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodPost, "/api/markLessonAsNotCompleted", bytes.NewReader(bodyBytes)).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 123}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)
	mockUsecase.EXPECT().MarkLessonAsNotCompleted(gomock.Any(), 123, 7).Return(nil)

	handler.MarkLessonAsNotCompleted(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetCourseRoadmap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourseRoadmap?courseId=99", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 999}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)
	mockUsecase.EXPECT().GetCourseRoadmap(gomock.Any(), 999, 99).Return(&dto.CourseRoadmapDTO{}, nil)

	handler.GetCourseRoadmap(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestServeVideo_NoRange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/serveVideo?lesson_id=1", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	videoURL := "https://cdn.com/videos/video.mp4"
	meta := dto.VideoMeta{Name: "video.mp4", Size: 1000}
	reader := ioutil.NopCloser(strings.NewReader("video data"))

	mockUsecase.EXPECT().GetVideoUrl(gomock.Any(), 1).Return(videoURL, nil)
	mockUsecase.EXPECT().GetMeta(gomock.Any(), "video.mp4").Return(meta, nil)
	mockUsecase.EXPECT().GetFragment(gomock.Any(), "video.mp4", int64(0), int64(999)).Return(reader, nil)

	handler.ServeVideo(w, req)

	if w.Code != http.StatusPartialContent && w.Code != http.StatusOK {
		t.Errorf("expected 206 or 200, got %d", w.Code)
	}
}

func TestGetCourse_InvalidCourseID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourses?courseId=invalid", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetCourse(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetCourse_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourses?courseId=123", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 1}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)
	mockUsecase.EXPECT().GetCourse(gomock.Any(), 123, profile).Return(nil, fmt.Errorf("course not found"))

	handler.GetCourse(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetCourseLesson_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourseLesson?courseId=5", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockCookie.EXPECT().CheckCookie(req).Return(nil)

	handler.GetCourseLesson(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetCourseLesson_InvalidCourseID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourseLesson?courseId=invalid", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 42}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)

	handler.GetCourseLesson(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetCourseLesson_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourseLesson?courseId=5", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 42}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)
	mockUsecase.EXPECT().GetCourseLesson(gomock.Any(), 42, 5).Return(nil, fmt.Errorf("lesson not found"))

	handler.GetCourseLesson(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetNextLesson_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getNextLesson?courseId=1&lessonId=2", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockCookie.EXPECT().CheckCookie(req).Return(nil)

	handler.GetNextLesson(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetNextLesson_InvalidCourseID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getNextLesson?courseId=invalid&lessonId=2", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 10}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)

	handler.GetNextLesson(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetNextLesson_InvalidLessonID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getNextLesson?courseId=1&lessonId=invalid", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 10}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)

	handler.GetNextLesson(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestMarkLessonAsNotCompleted_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	body := dto.LessonIDRequest{Id: 7}
	bodyBytes, _ := json.Marshal(body)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodPost, "/api/markLessonAsNotCompleted", bytes.NewReader(bodyBytes)).WithContext(ctx)
	w := httptest.NewRecorder()

	mockCookie.EXPECT().CheckCookie(req).Return(nil) // User not logged in

	handler.MarkLessonAsNotCompleted(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestMarkLessonAsNotCompleted_MethodNotAllowed(t *testing.T) {
	handler := handlers.NewHandler(nil, nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/markLessonAsNotCompleted", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.MarkLessonAsNotCompleted(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestGetCourseRoadmap_InvalidCourseID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourseRoadmap?courseId=invalid", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 999}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)

	handler.GetCourseRoadmap(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetCourseRoadmap_MethodNotAllowed(t *testing.T) {
	handler := handlers.NewHandler(nil, nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodPost, "/api/getCourseRoadmap", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetCourseRoadmap(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestServeVideo_VideoNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/serveVideo?lesson_id=1", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().GetVideoUrl(gomock.Any(), 1).Return("", fmt.Errorf("video not found"))

	handler.ServeVideo(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestServeVideo_InvalidLessonID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/serveVideo?lesson_id=invalid", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeVideo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestServeVideo_MissingLessonID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/serveVideo", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeVideo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetBucketCourses_EmptyResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourses", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().GetBucketCourses(gomock.Any()).Return(nil, nil)

	handler.GetCourses(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestMarkLessonAsNotCompleted_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	body := dto.LessonIDRequest{Id: 7}
	bodyBytes, _ := json.Marshal(body)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodPost, "/api/markLessonAsNotCompleted", bytes.NewReader(bodyBytes)).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 123}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)
	mockUsecase.EXPECT().MarkLessonAsNotCompleted(gomock.Any(), 123, 7).Return(fmt.Errorf("some error"))

	handler.MarkLessonAsNotCompleted(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetCourseRoadmap_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourseRoadmap?courseId=99", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 999}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)
	mockUsecase.EXPECT().GetCourseRoadmap(gomock.Any(), 999, 99).Return(nil, fmt.Errorf("roadmap error"))

	handler.GetCourseRoadmap(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetNextLesson_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getNextLesson?courseId=1&lessonId=2", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 10}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)
	mockUsecase.EXPECT().GetNextLesson(gomock.Any(), 10, 1, 2).Return(nil, fmt.Errorf("next lesson error"))

	handler.GetNextLesson(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
func TestGetCourse_UnexpectedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourses?courseId=123", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 1}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)
	mockUsecase.EXPECT().GetCourse(gomock.Any(), 123, profile).Return(nil, fmt.Errorf("unexpected error"))

	handler.GetCourse(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
func TestServeVideo_MetaError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/serveVideo?lesson_id=1", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	videoURL := "https://cdn.com/videos/video.mp4"
	mockUsecase.EXPECT().GetVideoUrl(gomock.Any(), 1).Return(videoURL, nil)
	mockUsecase.EXPECT().GetMeta(gomock.Any(), "video.mp4").Return(dto.VideoMeta{}, fmt.Errorf("meta error"))

	handler.ServeVideo(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestServeVideo_FragmentError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/serveVideo?lesson_id=1", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	videoURL := "https://cdn.com/videos/video.mp4"
	meta := dto.VideoMeta{Name: "video.mp4", Size: 1000}
	mockUsecase.EXPECT().GetVideoUrl(gomock.Any(), 1).Return(videoURL, nil)
	mockUsecase.EXPECT().GetMeta(gomock.Any(), "video.mp4").Return(meta, nil)
	mockUsecase.EXPECT().GetFragment(gomock.Any(), "video.mp4", int64(0), int64(999)).Return(nil, fmt.Errorf("fragment error"))

	handler.ServeVideo(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
func TestGetCourses_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	handler := handlers.NewHandler(mockUsecase, nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourses", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockUsecase.EXPECT().GetBucketCourses(gomock.Any()).Return(nil, fmt.Errorf("database error"))

	handler.GetCourses(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetCourse_MethodNotAllowed(t *testing.T) {
	handler := handlers.NewHandler(nil, nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodPost, "/api/getCourses?courseId=1", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetCourse(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestGetCourseLesson_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodPost, "/api/getCourses", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetCourseLesson(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestGetNextLesson_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)

	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodPost, "/api/getNextLesson?courseId=1&lessonId=2", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetNextLesson(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

func TestMarkLessonAsNotCompleted_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)
	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodPost, "/api/markLessonAsNotCompleted",
		strings.NewReader("invalid json")).WithContext(ctx)
	w := httptest.NewRecorder()

	profile := &usermodels.UserProfile{Id: 123}
	mockCookie.EXPECT().CheckCookie(req).Return(profile)

	handler.MarkLessonAsNotCompleted(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetCourseRoadmap_UnauthorizedUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	mockCookie := handlers.NewMockCookieManagerInterface(ctrl)
	handler := handlers.NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/getCourseRoadmap?courseId=99", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	mockCookie.EXPECT().CheckCookie(req).Return(nil)

	mockUsecase.EXPECT().GetCourseRoadmap(gomock.Any(), -1, 99).Return(&dto.CourseRoadmapDTO{}, nil)

	handler.GetCourseRoadmap(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestServeVideo_WithRangeHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := handlers.NewMockCourseUsecaseInterface(ctrl)
	handler := handlers.NewHandler(mockUsecase, nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/serveVideo?lesson_id=1", nil).WithContext(ctx)
	req.Header.Set("Range", "bytes=0-499")
	w := httptest.NewRecorder()

	videoURL := "https://cdn.com/videos/video.mp4"
	meta := dto.VideoMeta{Name: "video.mp4", Size: 1000}
	reader := ioutil.NopCloser(strings.NewReader("partial video data"))

	mockUsecase.EXPECT().GetVideoUrl(gomock.Any(), 1).Return(videoURL, nil)
	mockUsecase.EXPECT().GetMeta(gomock.Any(), "video.mp4").Return(meta, nil)
	mockUsecase.EXPECT().GetFragment(gomock.Any(), "video.mp4", int64(0), int64(499)).Return(reader, nil)

	handler.ServeVideo(w, req)

	if w.Code != http.StatusPartialContent {
		t.Errorf("expected 206, got %d", w.Code)
	}
}
