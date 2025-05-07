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

// Тест для UpdateProfile
func TestUpdateProfile(t *testing.T) {
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	userProfile := &models.UserProfile{Id: 1, Email: "john@example.com", Name: "John", Bio: "Hello", HideEmail: false}
	mockUsecase.On("UpdateProfile", mock.Anything, 1, userProfile).Return(nil)

	req := &userpb.UpdateProfileRequest{
		UserId: 1,
		Profile: &userpb.UserProfile{
			Email:     "john@example.com",
			Name:      "John",
			Bio:       "Hello",
			HideEmail: false,
		},
	}
	resp, err := handler.UpdateProfile(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockUsecase.AssertExpectations(t)
}

// Тест для UploadFile
func TestUploadFile(t *testing.T) {
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	fileData := []byte("test data")
	fileName := "test.txt"
	contentType := "text/plain"
	mockUsecase.On("UploadFile", mock.Anything, mock.Anything, mock.Anything).Return("http://example.com/test.txt", nil)

	req := &userpb.UploadFileRequest{
		FileData:    fileData,
		FileName:    fileName,
		ContentType: contentType,
	}
	resp, err := handler.UploadFile(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "http://example.com/test.txt", resp.Url)
	mockUsecase.AssertExpectations(t)
}

// Тест для SaveProfilePhoto
func TestSaveProfilePhoto(t *testing.T) {
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	mockUsecase.On("SaveProfilePhoto", mock.Anything, "http://example.com/photo.jpg", 1).Return("http://example.com/photo.jpg", nil)

	req := &userpb.SaveProfilePhotoRequest{
		Url:    "http://example.com/photo.jpg",
		UserId: 1,
	}
	resp, err := handler.SaveProfilePhoto(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "http://example.com/photo.jpg", resp.NewPhtotoUrl)
	mockUsecase.AssertExpectations(t)
}

// Тест для DeleteProfilePhoto
func TestDeleteProfilePhoto(t *testing.T) {
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)

	mockUsecase.On("DeleteProfilePhoto", mock.Anything, 1).Return(nil)

	req := &userpb.DeleteProfilePhotoRequest{UserId: 1}
	resp, err := handler.DeleteProfilePhoto(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockUsecase.AssertExpectations(t)
}
