package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/hash"
	"skillForce/pkg/logs"
)

type userRepo interface {
	RegisterUser(ctx context.Context, user *usermodels.User) (string, error)
	AuthenticateUser(ctx context.Context, email string, password string) (string, error)
	GetUserByCookie(ctx context.Context, cookieValue string) (*usermodels.UserProfile, error)
	LogoutUser(ctx context.Context, userId int) error
	UpdateProfile(ctx context.Context, userId int, userProfile *usermodels.UserProfile) error
	UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	UpdateProfilePhoto(ctx context.Context, url string, userId int) (string, error)
	DeleteProfilePhoto(ctx context.Context, userId int) error
	ValidUser(ctx context.Context, user *usermodels.User) (string, error)
	SendRegMail(ctx context.Context, user *usermodels.User, token string) error
	SendWelcomeMail(ctx context.Context, user *usermodels.User) error
	GetUserByToken(ctx context.Context, token string) (*usermodels.User, error)
}

type UserUsecase struct {
	repo userRepo
}

func NewUserUsecase(repo userRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (uc *UserUsecase) ValidUser(ctx context.Context, user *usermodels.User) error {
	token, err := uc.repo.ValidUser(ctx, user)
	if err != nil {
		logs.PrintLog(ctx, "ValidUser", fmt.Sprintf("%+v", err))
		return err
	}

	err = uc.repo.SendRegMail(ctx, user, token)
	if err != nil {
		logs.PrintLog(ctx, "SendMail", fmt.Sprintf("%+v", err))
		return err
	}
	return nil
}

func (uc *UserUsecase) RegisterUser(ctx context.Context, token string) (string, error) {
	user, err := uc.repo.GetUserByToken(ctx, token)
	if err != nil {
		return "", err
	}

	err = hash.HashPasswordAndCreateSalt(user)
	if err != nil {
		return "", err
	}

	_ = uc.repo.SendWelcomeMail(ctx, user)

	return uc.repo.RegisterUser(ctx, user)
}

// AuthenticateUser - авторизация пользователя
func (uc *UserUsecase) AuthenticateUser(ctx context.Context, user *usermodels.User) (string, error) {
	return uc.repo.AuthenticateUser(ctx, user.Email, user.Password)
}

// GetUserByCookie - получение пользователя по cookie
func (uc *UserUsecase) GetUserByCookie(ctx context.Context, cookieValue string) (*usermodels.UserProfile, error) {
	return uc.repo.GetUserByCookie(ctx, cookieValue)
}

func (uc *UserUsecase) LogoutUser(ctx context.Context, userId int) error {
	return uc.repo.LogoutUser(ctx, userId)
}

func (uc *UserUsecase) UpdateProfile(ctx context.Context, userId int, userProfile *usermodels.UserProfile) error {
	return uc.repo.UpdateProfile(ctx, userId, userProfile)
}

func (uc *UserUsecase) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	return uc.repo.UploadFile(ctx, file, fileHeader)
}

func (uc *UserUsecase) SaveProfilePhoto(ctx context.Context, url string, userId int) (string, error) {
	return uc.repo.UpdateProfilePhoto(ctx, url, userId)
}

func (uc *UserUsecase) DeleteProfilePhoto(ctx context.Context, userId int) error {
	return uc.repo.DeleteProfilePhoto(ctx, userId)
}
