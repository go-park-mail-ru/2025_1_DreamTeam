package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"skillForce/internal/models/dto"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/logs"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockUserUsecaseInterface(ctrl)
	mockCookie := NewMockCookieManagerInterface(ctrl)
	handler := NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	t.Run("Method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/register", nil).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.RegisterUser(rr, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Invalid JSON body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer([]byte("invalid_json"))).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.RegisterUser(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Missing required fields", func(t *testing.T) {
		body := `{"email":"test@mail.com","password":"12345"}`
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer([]byte(body))).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.RegisterUser(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "missing required fields")
	})

	t.Run("Invalid email", func(t *testing.T) {
		body := `{"name":"John","email":"invalid_email","password":"12345"}`
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer([]byte(body))).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.RegisterUser(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid email")
	})

	t.Run("Password too short", func(t *testing.T) {
		body := `{"name":"John","email":"test@mail.com","password":"123"}`
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer([]byte(body))).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.RegisterUser(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "password too short")
	})

	t.Run("Successful registration", func(t *testing.T) {
		body := `{"name":"John","email":"john@mail.com","password":"12345"}`
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer([]byte(body))).WithContext(ctx)
		rr := httptest.NewRecorder()

		mockUsecase.EXPECT().
			ValidUser(gomock.Any(), &usermodels.User{
				Name:     "John",
				Email:    "john@mail.com",
				Password: "12345",
			}).
			Return(nil)

		handler.RegisterUser(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestConfirmUserEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockUserUsecaseInterface(ctrl)
	mockCookie := NewMockCookieManagerInterface(ctrl)
	handler := NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	t.Run("Method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/validEmail", nil).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ConfirmUserEmail(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/validEmail", nil).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ConfirmUserEmail(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid token")
	})

	t.Run("Unexpected server error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/validEmail?token=abc123", nil).WithContext(ctx)
		rr := httptest.NewRecorder()

		mockUsecase.EXPECT().
			RegisterUser(gomock.Any(), "abc123").
			Return("", errors.New("db error"))

		handler.ConfirmUserEmail(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "server error")
	})

	t.Run("Successful confirmation", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/validEmail?token=abc123", nil).WithContext(ctx)
		rr := httptest.NewRecorder()

		mockUsecase.EXPECT().
			RegisterUser(gomock.Any(), "abc123").
			Return("cookie-value", nil)

		mockCookie.EXPECT().
			SetCookie(rr, "cookie-value")

		handler.ConfirmUserEmail(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockUserUsecaseInterface(ctrl)
	mockCookie := NewMockCookieManagerInterface(ctrl)
	handler := NewHandler(mockUsecase, mockCookie)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	t.Run("Method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/login", nil).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.LoginUser(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		body := bytes.NewBufferString("invalid json")
		req := httptest.NewRequest(http.MethodPost, "/api/login", body).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.LoginUser(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid request")
	})

	t.Run("Invalid email format", func(t *testing.T) {
		input := dto.UserDTO{
			Email:    "invalid-email",
			Password: "123456",
		}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body)).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.LoginUser(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid email")
	})

	t.Run("Password too short", func(t *testing.T) {
		input := dto.UserDTO{
			Email:    "user@example.com",
			Password: "123",
		}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body)).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.LoginUser(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "password too short")
	})

	t.Run("Internal server error", func(t *testing.T) {
		input := dto.UserDTO{
			Email:    "user@example.com",
			Password: "password",
		}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body)).WithContext(ctx)
		rr := httptest.NewRecorder()

		mockUsecase.EXPECT().
			AuthenticateUser(gomock.Any(), gomock.Any()).
			Return("", errors.New("db is down"))

		handler.LoginUser(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "server error")
	})

	t.Run("Successful login", func(t *testing.T) {
		input := dto.UserDTO{
			Email:    "user@example.com",
			Password: "password",
		}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body)).WithContext(ctx)
		rr := httptest.NewRecorder()

		mockUsecase.EXPECT().
			AuthenticateUser(gomock.Any(), gomock.Any()).
			Return("cookie-value", nil)

		mockCookie.EXPECT().
			SetCookie(rr, "cookie-value")

		handler.LoginUser(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
