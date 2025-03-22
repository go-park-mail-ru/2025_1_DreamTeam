package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestInvalidRegisterUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Invalid Email",
			method:         "POST",
			body:           `{"email": "user@", "name": "John Doe", "password": "password"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid email",
		},
		{
			name:           "Invalid Name",
			method:         "POST",
			body:           `{"email": "user@mail.ru", "name": "", "password": "password"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "missing required fields",
		},
		{
			name:           "Invalid Password",
			method:         "POST",
			body:           `{"email": "user@mail.ru", "name": "Joe", "password": "pass"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "password too short",
		},
		{
			name:           "Invalid Method",
			method:         "GET",
			body:           `{"email": "user@mail.ru", "name": "Joe", "password": "password"}`,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "method not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := repository.NewmockDB(false)
			uc := usecase.NewUserUsecase(mockDB)
			h := &UserHandler{useCase: uc}

			body := bytes.NewReader([]byte(tt.body))
			r := httptest.NewRequest(tt.method, "/api/register", body)
			w := httptest.NewRecorder()

			h.RegisterUser(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("ожидали %d, получили %d", tt.expectedStatus, resp.StatusCode)
			}

			if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("ожидали Content-Type application/json, получили %s", contentType)
			}

			var errorResp response.ErrorResponse
			err := json.NewDecoder(resp.Body).Decode(&errorResp)
			if err != nil {
				t.Fatalf("не удалось распарсить JSON: %v", err)
			}

			if strings.TrimSpace(errorResp.ErrorStr) != tt.expectedError {
				t.Errorf("ожидали ошибку \"%s\", но получили \"%s\"", tt.expectedError, errorResp.ErrorStr)
			}
		})
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

func TestInvalidLoginUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Invalid Email",
			method:         "POST",
			body:           `{"email": "user@", "password": "password"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid email",
		},
		{
			name:           "Invalid Password",
			method:         "POST",
			body:           `{"email": "user@mail", "password": "pass"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "password too short",
		},
		{
			name:           "Invalid Method",
			method:         "GET",
			body:           `{"email": "user@mail.ru", "password": "password"}`,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "method not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := repository.NewmockDB(false)
			uc := usecase.NewUserUsecase(mockDB)
			h := &UserHandler{useCase: uc}

			body := bytes.NewReader([]byte(tt.body))
			r := httptest.NewRequest(tt.method, "/api/login", body)
			w := httptest.NewRecorder()

			h.LoginUser(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("ожидали %d, получили %d", tt.expectedStatus, resp.StatusCode)
			}

			if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("ожидали Content-Type application/json, получили %s", contentType)
			}

			var errorResp response.ErrorResponse
			err := json.NewDecoder(resp.Body).Decode(&errorResp)
			if err != nil {
				t.Fatalf("не удалось распарсить JSON: %v", err)
			}

			if strings.TrimSpace(errorResp.ErrorStr) != tt.expectedError {
				t.Errorf("ожидали ошибку \"%s\", но получили \"%s\"", tt.expectedError, errorResp.ErrorStr)
			}
		})
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
