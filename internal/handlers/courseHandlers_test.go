package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"skillForce/internal/models"
	"skillForce/internal/repository"
	"skillForce/internal/response"
	"skillForce/internal/usecase"
	"testing"
)

func TestOKGetCourses(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(false)
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
