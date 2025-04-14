package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"skillForce/internal/models"
	"skillForce/internal/repository"
	"skillForce/pkg/hash"
	"skillForce/pkg/logs"
)

type UserUsecase struct {
	repo repository.Repository
}

func NewUserUsecase(repo repository.Repository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (uc *UserUsecase) ValidUser(ctx context.Context, user *models.User) error {
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
func (uc *UserUsecase) AuthenticateUser(ctx context.Context, user *models.User) (string, error) {
	return uc.repo.AuthenticateUser(ctx, user.Email, user.Password)
}

// GetUserByCookie - получение пользователя по cookie
func (uc *UserUsecase) GetUserByCookie(ctx context.Context, cookieValue string) (*models.UserProfile, error) {
	return uc.repo.GetUserByCookie(ctx, cookieValue)
}

func (uc *UserUsecase) LogoutUser(ctx context.Context, userId int) error {
	return uc.repo.LogoutUser(ctx, userId)
}

func (uc *UserUsecase) UpdateProfile(ctx context.Context, userId int, userProfile *models.UserProfile) error {
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
