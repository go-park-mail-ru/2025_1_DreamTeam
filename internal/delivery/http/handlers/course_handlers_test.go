package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"skillForce/internal/delivery/http/response"
	"skillForce/internal/models"
	"skillForce/internal/repository/mock"
	"skillForce/internal/usecase"
	"strings"
	"testing"
)

func TestOKGetCourses(t *testing.T) {
	t.Parallel()

	mockDB := mock.NewmockDB(false)
	uc := usecase.NewCourseUsecase(mockDB)
	h := &CourseHandler{useCase: uc}

	r := httptest.NewRequest("GET", "/api/getCourses", nil)
	w := httptest.NewRecorder()

	h.GetCourses(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("ожидали %d, получили %d", http.StatusOK, resp.StatusCode)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("ожидали Content-Type application/json, получили %s", contentType)
	}

	expectedCourses := []*models.Course{
		{Id: 1, Price: 1, PurchasesAmount: 1, CreatorId: 1, TimeToPass: 1, Title: "Курс #1", Description: "Описание курса #1", ScrImage: "image_1.jpg"},
		{Id: 2, Price: 2, PurchasesAmount: 2, CreatorId: 2, TimeToPass: 2, Title: "Курс #2", Description: "Описание курса #2", ScrImage: "image_2.jpg"},
	}

	var actualResponse response.BucketCoursesResponse
	err := json.NewDecoder(resp.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Ошибка декодирования JSON: %v", err)
	}

	if !reflect.DeepEqual(actualResponse.BucketCourses, expectedCourses) {
		t.Errorf("ожидали %v, получили %v", expectedCourses, actualResponse.BucketCourses)
	}
}

func TestInvalidMethodGetCourses(t *testing.T) {
	t.Parallel()

	mockDB := mock.NewmockDB(false)
	uc := usecase.NewCourseUsecase(mockDB)
	h := &CourseHandler{useCase: uc}

	r := httptest.NewRequest("POST", "/api/getCourses", nil)
	w := httptest.NewRecorder()

	h.GetCourses(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("ожидали %d, получили %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("ожидали Content-Type application/json, получили %s", contentType)
	}

	var errorResp response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Fatalf("не удалось распарсить JSON: %v", err)
	}

	expectedError := "method not allowed"
	if strings.TrimSpace(errorResp.ErrorStr) != expectedError {
		t.Errorf("ожидали ошибку \"%s\", но получили \"%s\"", expectedError, errorResp.ErrorStr)
	}
}
