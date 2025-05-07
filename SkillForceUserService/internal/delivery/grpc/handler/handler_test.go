package handler

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"testing"

	userpb "skillForce/internal/delivery/grpc/proto"
	models "skillForce/internal/models/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock для UserUsecaseInterface
type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) RegisterUser(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) ValidUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUsecase) AuthenticateUser(ctx context.Context, user *models.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) LogoutUser(ctx context.Context, userId int) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockUserUsecase) UpdateProfile(ctx context.Context, userId int, userProfile *models.UserProfile) error {
	args := m.Called(ctx, userId, userProfile)
	return args.Error(0)
}

func (m *MockUserUsecase) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	args := m.Called(ctx, file, fileHeader)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) SaveProfilePhoto(ctx context.Context, url string, userId int) (string, error) {
	args := m.Called(ctx, url, userId)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) DeleteProfilePhoto(ctx context.Context, userId int) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

// Тест для RegisterUser
func TestRegisterUser(t *testing.T) {
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	mockUsecase.On("RegisterUser", mock.Anything, "test_token").Return("cookie_value", nil)

	req := &userpb.RegisterRequest{Token: "test_token"}
	resp, err := handler.RegisterUser(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "cookie_value", resp.CookieVal)
	mockUsecase.AssertExpectations(t)
}

// Тест для ValidUser
func TestValidUser(t *testing.T) {
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	user := &models.User{Name: "John", Email: "john@example.com", Password: "password"}
	mockUsecase.On("ValidUser", mock.Anything, user).Return(nil)

	req := &userpb.User{Name: "John", Email: "john@example.com", Password: "password"}
	resp, err := handler.ValidUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockUsecase.AssertExpectations(t)
}

// Тест для AuthenticateUser
func TestAuthenticateUser(t *testing.T) {
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	user := &models.User{Name: "John", Email: "john@example.com", Password: "password"}
	mockUsecase.On("AuthenticateUser", mock.Anything, user).Return("cookie_value", nil)

	req := &userpb.User{Name: "John", Email: "john@example.com", Password: "password"}
	resp, err := handler.AuthenticateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "cookie_value", resp.CookieVal)
	mockUsecase.AssertExpectations(t)
}

// Тест для ConvertToMultipart
func TestConvertToMultipart(t *testing.T) {
	fileData := []byte("test data")
	fileName := "test.txt"
	contentType := "text/plain"

	file, fileHeader, err := ConvertToMultipart(fileData, fileName, contentType)

	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.NotNil(t, fileHeader)

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	assert.NoError(t, err)
	assert.Equal(t, "test data", buf.String())
}
