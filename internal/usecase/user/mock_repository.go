// Code generated by MockGen. DO NOT EDIT.
// Source: internal/usecase/user/user_usecase.go

// Package usecase is a generated GoMock package.
package usecase

import (
	context "context"
	multipart "mime/multipart"
	reflect "reflect"
	usermodels "skillForce/internal/models/user"

	gomock "github.com/golang/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// AuthenticateUser mocks base method.
func (m *MockUserRepository) AuthenticateUser(ctx context.Context, email, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthenticateUser", ctx, email, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthenticateUser indicates an expected call of AuthenticateUser.
func (mr *MockUserRepositoryMockRecorder) AuthenticateUser(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticateUser", reflect.TypeOf((*MockUserRepository)(nil).AuthenticateUser), ctx, email, password)
}

// DeleteProfilePhoto mocks base method.
func (m *MockUserRepository) DeleteProfilePhoto(ctx context.Context, userId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProfilePhoto", ctx, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProfilePhoto indicates an expected call of DeleteProfilePhoto.
func (mr *MockUserRepositoryMockRecorder) DeleteProfilePhoto(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProfilePhoto", reflect.TypeOf((*MockUserRepository)(nil).DeleteProfilePhoto), ctx, userId)
}

// GetUserByCookie mocks base method.
func (m *MockUserRepository) GetUserByCookie(ctx context.Context, cookieValue string) (*usermodels.UserProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByCookie", ctx, cookieValue)
	ret0, _ := ret[0].(*usermodels.UserProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByCookie indicates an expected call of GetUserByCookie.
func (mr *MockUserRepositoryMockRecorder) GetUserByCookie(ctx, cookieValue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByCookie", reflect.TypeOf((*MockUserRepository)(nil).GetUserByCookie), ctx, cookieValue)
}

// GetUserByToken mocks base method.
func (m *MockUserRepository) GetUserByToken(ctx context.Context, token string) (*usermodels.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByToken", ctx, token)
	ret0, _ := ret[0].(*usermodels.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByToken indicates an expected call of GetUserByToken.
func (mr *MockUserRepositoryMockRecorder) GetUserByToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByToken", reflect.TypeOf((*MockUserRepository)(nil).GetUserByToken), ctx, token)
}

// LogoutUser mocks base method.
func (m *MockUserRepository) LogoutUser(ctx context.Context, userId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogoutUser", ctx, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogoutUser indicates an expected call of LogoutUser.
func (mr *MockUserRepositoryMockRecorder) LogoutUser(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogoutUser", reflect.TypeOf((*MockUserRepository)(nil).LogoutUser), ctx, userId)
}

// RegisterUser mocks base method.
func (m *MockUserRepository) RegisterUser(ctx context.Context, user *usermodels.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", ctx, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockUserRepositoryMockRecorder) RegisterUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockUserRepository)(nil).RegisterUser), ctx, user)
}

// SendRegMail mocks base method.
func (m *MockUserRepository) SendRegMail(ctx context.Context, user *usermodels.User, token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendRegMail", ctx, user, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendRegMail indicates an expected call of SendRegMail.
func (mr *MockUserRepositoryMockRecorder) SendRegMail(ctx, user, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendRegMail", reflect.TypeOf((*MockUserRepository)(nil).SendRegMail), ctx, user, token)
}

// SendWelcomeMail mocks base method.
func (m *MockUserRepository) SendWelcomeMail(ctx context.Context, user *usermodels.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendWelcomeMail", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendWelcomeMail indicates an expected call of SendWelcomeMail.
func (mr *MockUserRepositoryMockRecorder) SendWelcomeMail(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendWelcomeMail", reflect.TypeOf((*MockUserRepository)(nil).SendWelcomeMail), ctx, user)
}

// UpdateProfile mocks base method.
func (m *MockUserRepository) UpdateProfile(ctx context.Context, userId int, userProfile *usermodels.UserProfile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfile", ctx, userId, userProfile)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProfile indicates an expected call of UpdateProfile.
func (mr *MockUserRepositoryMockRecorder) UpdateProfile(ctx, userId, userProfile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockUserRepository)(nil).UpdateProfile), ctx, userId, userProfile)
}

// UpdateProfilePhoto mocks base method.
func (m *MockUserRepository) UpdateProfilePhoto(ctx context.Context, url string, userId int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfilePhoto", ctx, url, userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfilePhoto indicates an expected call of UpdateProfilePhoto.
func (mr *MockUserRepositoryMockRecorder) UpdateProfilePhoto(ctx, url, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfilePhoto", reflect.TypeOf((*MockUserRepository)(nil).UpdateProfilePhoto), ctx, url, userId)
}

// UploadFile mocks base method.
func (m *MockUserRepository) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", ctx, file, fileHeader)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockUserRepositoryMockRecorder) UploadFile(ctx, file, fileHeader interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockUserRepository)(nil).UploadFile), ctx, file, fileHeader)
}

// ValidUser mocks base method.
func (m *MockUserRepository) ValidUser(ctx context.Context, user *usermodels.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidUser", ctx, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidUser indicates an expected call of ValidUser.
func (mr *MockUserRepositoryMockRecorder) ValidUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidUser", reflect.TypeOf((*MockUserRepository)(nil).ValidUser), ctx, user)
}
