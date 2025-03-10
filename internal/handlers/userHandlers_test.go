package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"skillForce/internal/models"
	"skillForce/internal/repository"
	"skillForce/internal/response"
	"skillForce/internal/usecase"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestOKRegisterUser(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(false)
	uc := usecase.NewUserUsecase(mockDB)
	h := &UserHandler{useCase: uc}

	bodyString := `{
		"email": "user@mail.ru",
		"name": "John Doe",
		"password": "password"
	}`
	body := bytes.NewReader([]byte(bodyString))

	r := httptest.NewRequest("POST", "/api/register", body)
	w := httptest.NewRecorder()

	h.RegisterUser(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("ожидали %d, получили %d", http.StatusOK, resp.StatusCode)
	}

	cookies := resp.Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		t.Errorf("cookie session_id не найден")
	}

	expectedSessionId := strconv.Itoa(0)
	if sessionCookie.Value != expectedSessionId {
		t.Errorf("ожидали куку session_id=%s, но получили %s", expectedSessionId, sessionCookie.Value)
	}

	if !sessionCookie.HttpOnly {
		t.Error("ожидали, что куки будут HttpOnly, но они не были установлены как HttpOnly")
	}

	expectedExpiration := time.Now().Add(10 * time.Hour)
	if sessionCookie.Expires.Before(time.Now()) || sessionCookie.Expires.After(expectedExpiration.Add(time.Minute)) {
		t.Error("время жизни куки установлено некорректно")
	}

}

func TestInvalidEmailRegisterUser(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(false)
	uc := usecase.NewUserUsecase(mockDB)
	h := &UserHandler{useCase: uc}

	bodyString := `{
		"email": "user@",
		"name": "John Doe",
		"password": "password"
	}`
	body := bytes.NewReader([]byte(bodyString))

	r := httptest.NewRequest("POST", "/api/register", body)
	w := httptest.NewRecorder()

	h.RegisterUser(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("ожидали %d, получили %d", http.StatusBadRequest, resp.StatusCode)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("ожидали Content-Type application/json, получили %s", contentType)
	}

	var errorResp response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Fatalf("не удалось распарсить JSON: %v", err)
	}

	expectedError := "invalid email"
	if strings.TrimSpace(errorResp.ErrorStr) != expectedError {
		t.Errorf("ожидали ошибку \"%s\", но получили \"%s\"", expectedError, errorResp.ErrorStr)
	}
}

func TestInvalidNameRegisterUser(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(false)
	uc := usecase.NewUserUsecase(mockDB)
	h := &UserHandler{useCase: uc}

	bodyString := `{
		"email": "user@mail.ru",
		"name": "",
		"password": "password"
	}`
	body := bytes.NewReader([]byte(bodyString))

	r := httptest.NewRequest("POST", "/api/register", body)
	w := httptest.NewRecorder()

	h.RegisterUser(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("ожидали %d, получили %d", http.StatusBadRequest, resp.StatusCode)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("ожидали Content-Type application/json, получили %s", contentType)
	}

	var errorResp response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Fatalf("не удалось распарсить JSON: %v", err)
	}

	expectedError := "missing required fields"
	if strings.TrimSpace(errorResp.ErrorStr) != expectedError {
		t.Errorf("ожидали ошибку \"%s\", но получили \"%s\"", expectedError, errorResp.ErrorStr)
	}
}

func TestInvalidPasswordRegisterUser(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(false)
	uc := usecase.NewUserUsecase(mockDB)
	h := &UserHandler{useCase: uc}

	bodyString := `{
		"email": "user@mail.ru",
		"name": "Joe",
		"password": "pass"
	}`
	body := bytes.NewReader([]byte(bodyString))

	r := httptest.NewRequest("POST", "/api/register", body)
	w := httptest.NewRecorder()

	h.RegisterUser(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("ожидали %d, получили %d", http.StatusBadRequest, resp.StatusCode)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("ожидали Content-Type application/json, получили %s", contentType)
	}

	var errorResp response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Fatalf("не удалось распарсить JSON: %v", err)
	}

	expectedError := "password too short"
	if strings.TrimSpace(errorResp.ErrorStr) != expectedError {
		t.Errorf("ожидали ошибку \"%s\", но получили \"%s\"", expectedError, errorResp.ErrorStr)
	}
}

func TestInvalidMethodRegisterUser(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(false)
	uc := usecase.NewUserUsecase(mockDB)
	h := &UserHandler{useCase: uc}

	bodyString := `{
		"email": "user@mail.ru",
		"name": "Joe",
		"password": "password"
	}`
	body := bytes.NewReader([]byte(bodyString))

	r := httptest.NewRequest("GET", "/api/register", body)
	w := httptest.NewRecorder()

	h.RegisterUser(w, r)

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

func TestOKLoginUser(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(false)
	uc := usecase.NewUserUsecase(mockDB)
	h := &UserHandler{useCase: uc}

	bodyString := `{
		"email": "user@mail.ru",
		"name": "John Doe",
		"password": "password"
	}`
	body := bytes.NewReader([]byte(bodyString))

	r := httptest.NewRequest("POST", "/api/login", body)
	w := httptest.NewRecorder()

	h.LoginUser(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("ожидали %d, получили %d", http.StatusOK, resp.StatusCode)
	}

	cookies := resp.Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		t.Errorf("cookie session_id не найден")
	}

	expectedSessionId := strconv.Itoa(0)
	if sessionCookie.Value != expectedSessionId {
		t.Errorf("ожидали куку session_id=%s, но получили %s", expectedSessionId, sessionCookie.Value)
	}

	if !sessionCookie.HttpOnly {
		t.Error("ожидали, что куки будут HttpOnly, но они не были установлены как HttpOnly")
	}

	expectedExpiration := time.Now().Add(10 * time.Hour)
	if sessionCookie.Expires.Before(time.Now()) || sessionCookie.Expires.After(expectedExpiration.Add(time.Minute)) {
		t.Error("время жизни куки установлено некорректно")
	}
}

func TestInvalidEmailLoginUser(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(false)
	uc := usecase.NewUserUsecase(mockDB)
	h := &UserHandler{useCase: uc}

	bodyString := `{
		"email": "user@",
		"password": "password"
	}`
	body := bytes.NewReader([]byte(bodyString))

	r := httptest.NewRequest("POST", "/api/login", body)
	w := httptest.NewRecorder()

	h.LoginUser(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("ожидали %d, получили %d", http.StatusBadRequest, resp.StatusCode)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("ожидали Content-Type application/json, получили %s", contentType)
	}

	var errorResp response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Fatalf("не удалось распарсить JSON: %v", err)
	}

	expectedError := "invalid email"
	if strings.TrimSpace(errorResp.ErrorStr) != expectedError {
		t.Errorf("ожидали ошибку \"%s\", но получили \"%s\"", expectedError, errorResp.ErrorStr)
	}
}

func TestInvalidPasswordLoginUser(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(false)
	uc := usecase.NewUserUsecase(mockDB)
	h := &UserHandler{useCase: uc}

	bodyString := `{
		"email": "user@mail",
		"password": "pass"
	}`
	body := bytes.NewReader([]byte(bodyString))

	r := httptest.NewRequest("POST", "/api/login", body)
	w := httptest.NewRecorder()

	h.LoginUser(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("ожидали %d, получили %d", http.StatusBadRequest, resp.StatusCode)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("ожидали Content-Type application/json, получили %s", contentType)
	}

	var errorResp response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Fatalf("не удалось распарсить JSON: %v", err)
	}

	expectedError := "password too short"
	if strings.TrimSpace(errorResp.ErrorStr) != expectedError {
		t.Errorf("ожидали ошибку \"%s\", но получили \"%s\"", expectedError, errorResp.ErrorStr)
	}
}

func TestInvalidMethodLoginUser(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(false)
	uc := usecase.NewUserUsecase(mockDB)
	h := &UserHandler{useCase: uc}

	bodyString := `{
		"email": "user@mail.ru",
		"password": "password"
	}`
	body := bytes.NewReader([]byte(bodyString))

	r := httptest.NewRequest("GET", "/api/login", body)
	w := httptest.NewRecorder()

	h.LoginUser(w, r)

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

func TestOKLogoutUnauthorizedUser(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(false)
	uc := usecase.NewUserUsecase(mockDB)
	h := &UserHandler{useCase: uc}

	r := httptest.NewRequest("GET", "/api/logout", nil)
	w := httptest.NewRecorder()

	h.LogoutUser(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("ожидали %d, получили %d", http.StatusOK, resp.StatusCode)
	}
}

func TestOKLogoutAuthorizedUser(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(true)
	uc := usecase.NewUserUsecase(mockDB)
	h := &UserHandler{useCase: uc}

	r := httptest.NewRequest("GET", "/api/logout", nil)
	w := httptest.NewRecorder()

	r.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    strconv.Itoa(1),
		HttpOnly: true,
	})

	h.LogoutUser(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("ожидали %d, получили %d", http.StatusOK, resp.StatusCode)
	}

	cookies := resp.Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		t.Errorf("cookie session_id не найден")
	}

	expectedSessionId := strconv.Itoa(1)
	if sessionCookie.Value != expectedSessionId {
		t.Errorf("ожидали куку session_id=%s, но получили %s", expectedSessionId, sessionCookie.Value)
	}

	if !sessionCookie.HttpOnly {
		t.Error("ожидали, что куки будут HttpOnly, но они не были установлены как HttpOnly")
	}

	if !sessionCookie.Expires.Before(time.Now()) {
		t.Error("время жизни куки установлено некорректно для удаления")
	}
}

func TestFalseIsAuthorized(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(false)
	uc := usecase.NewUserUsecase(mockDB)
	h := &UserHandler{useCase: uc}

	r := httptest.NewRequest("GET", "/api/isAuthorized", nil)
	w := httptest.NewRecorder()

	h.IsAuthorized(w, r)

	h.LoginUser(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("ожидали %d, получили %d", http.StatusUnauthorized, resp.StatusCode)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("ожидали Content-Type application/json, получили %s", contentType)
	}

	var errorResp response.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Fatalf("не удалось распарсить JSON: %v", err)
	}

	expectedError := "not authorized"
	if strings.TrimSpace(errorResp.ErrorStr) != expectedError {
		t.Errorf("ожидали ошибку \"%s\", но получили \"%s\"", expectedError, errorResp.ErrorStr)
	}
}

func TestTrueIsAuthorized(t *testing.T) {
	t.Parallel()

	mockDB := repository.NewmockDB(true)
	uc := usecase.NewUserUsecase(mockDB)
	h := &UserHandler{useCase: uc}

	r := httptest.NewRequest("GET", "/api/isAuthorized", nil)
	w := httptest.NewRecorder()

	r.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    strconv.Itoa(1),
		HttpOnly: true,
	})

	h.IsAuthorized(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("ожидали %d, получили %d", http.StatusOK, resp.StatusCode)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("ожидали Content-Type application/json, получили %s", contentType)
	}

	var returnedUser models.User
	err := json.NewDecoder(resp.Body).Decode(&returnedUser)
	if err != nil {
		t.Fatalf("не удалось распарсить JSON: %v", err)
	}
}

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
