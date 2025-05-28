package usecase

import (
	"context"
	"errors"
	multipart "mime/multipart"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/logs"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// func TestValidUser_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := NewMockUserRepository(ctrl)
// 	uc := NewUserUsecase(mockRepo)

// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
// 		Data: make([]*logs.LogString, 0),
// 	})
// 	user := &usermodels.User{}
// 	expectedToken := "token123"

// 	mockRepo.EXPECT().
// 		ValidUser(ctx, user).
// 		Return(expectedToken, nil)
// 	mockRepo.EXPECT().SendRegMail(ctx, user, expectedToken).Return(nil)

// 	err := uc.ValidUser(ctx, user)
// 	require.NoError(t, err)
// }

func TestValidUser_ValidFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	uc := NewUserUsecase(mockRepo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	user := &usermodels.User{}

	mockRepo.EXPECT().
		ValidUser(ctx, user).
		Return("", errors.New("validation error"))

	err := uc.ValidUser(ctx, user)
	require.Error(t, err)
	require.Equal(t, "validation error", err.Error())
}

// func TestValidUser_SendMailFail(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := NewMockUserRepository(ctrl)
// 	uc := NewUserUsecase(mockRepo)

// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
// 		Data: make([]*logs.LogString, 0),
// 	})
// 	user := &usermodels.User{}
// 	expectedToken := "token123"

// 	mockRepo.EXPECT().
// 		ValidUser(ctx, user).
// 		Return(expectedToken, nil)
// 	mockRepo.EXPECT().
// 		SendRegMail(ctx, user, expectedToken).
// 		Return(errors.New("mail error"))

// 	err := uc.ValidUser(ctx, user)
// 	require.Error(t, err)
// 	require.Equal(t, "mail error", err.Error())
// }

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	uc := NewUserUsecase(mockRepo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	token := "token123"
	user := &usermodels.User{
		Password: "plainpassword", // Тестовый пароль
	}

	// Настраиваем ожидания
	mockRepo.EXPECT().
		GetUserByToken(ctx, token).
		Return(user, nil)

	mockRepo.EXPECT().
		RegisterUser(ctx, gomock.Any()).
		Return("jwtToken", nil)

	// Вызываем тестируемый метод
	result, err := uc.RegisterUser(ctx, token)

	// Проверяем результаты
	require.NoError(t, err)
	require.Equal(t, "jwtToken", result)

	// Проверяем, что пароль был захэширован
	require.NotEqual(t, "plainpassword", user.Password)
	require.NotEmpty(t, user.Salt)
}

func TestAuthenticateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	uc := NewUserUsecase(mockRepo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	user := &usermodels.User{}
	expectedToken := "jwtToken"

	mockRepo.EXPECT().
		AuthenticateUser(ctx, user.Email, user.Password).
		Return(expectedToken, nil)

	token, err := uc.AuthenticateUser(ctx, user)
	require.NoError(t, err)
	require.Equal(t, expectedToken, token)
}

func TestUpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	uc := NewUserUsecase(mockRepo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	profile := &usermodels.UserProfile{}

	mockRepo.EXPECT().UpdateProfile(ctx, 1, profile).Return(nil)

	err := uc.UpdateProfile(ctx, 1, profile)
	require.NoError(t, err)
}

func TestUploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	uc := NewUserUsecase(mockRepo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	var file multipart.File
	var header *multipart.FileHeader

	mockRepo.EXPECT().
		UploadFile(ctx, file, header).
		Return("https://example.com/file.png", nil)

	url, err := uc.UploadFile(ctx, file, header)
	require.NoError(t, err)
	require.Equal(t, "https://example.com/file.png", url)
}

func TestGetUserByCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	uc := NewUserUsecase(mockRepo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	cookie := "test_cookie"
	expectedProfile := &usermodels.UserProfile{Id: 1}

	mockRepo.EXPECT().
		GetUserByCookie(ctx, cookie).
		Return(expectedProfile, nil)

	profile, err := uc.GetUserByCookie(ctx, cookie)
	require.NoError(t, err)
	require.Equal(t, expectedProfile, profile)
}

func TestLogoutUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	uc := NewUserUsecase(mockRepo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	userId := 123

	mockRepo.EXPECT().
		LogoutUser(ctx, userId).
		Return(nil)

	err := uc.LogoutUser(ctx, userId)
	require.NoError(t, err)
}

func TestSaveProfilePhoto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	uc := NewUserUsecase(mockRepo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	userId := 456
	url := "https://example.com/photo.jpg"
	expectedUrl := "https://example.com/photo.jpg"

	mockRepo.EXPECT().
		UpdateProfilePhoto(ctx, url, userId).
		Return(expectedUrl, nil)

	resultUrl, err := uc.SaveProfilePhoto(ctx, url, userId)
	require.NoError(t, err)
	require.Equal(t, expectedUrl, resultUrl)
}

func TestDeleteProfilePhoto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)
	uc := NewUserUsecase(mockRepo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	userId := 789

	mockRepo.EXPECT().
		DeleteProfilePhoto(ctx, userId).
		Return(nil)

	err := uc.DeleteProfilePhoto(ctx, userId)
	require.NoError(t, err)
}
