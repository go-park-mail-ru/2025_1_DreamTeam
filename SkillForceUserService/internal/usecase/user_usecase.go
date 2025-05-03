package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/hash"
	"skillForce/pkg/logs"
)

type UserUsecase struct {
	repo UserRepository
}

func NewUserUsecase(repo UserRepository) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (uc *UserUsecase) ValidUser(ctx context.Context, user *usermodels.User) error {
	token, err := uc.repo.ValidUser(ctx, user)
	if err != nil {
		logs.PrintLog(ctx, "ValidUser", fmt.Sprintf("%+v", err))
		return err
	}

	go uc.repo.SendRegMail(ctx, user, token)
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

	go uc.repo.SendWelcomeMail(ctx, user)

	return uc.repo.RegisterUser(ctx, user)
}

func (uc *UserUsecase) AuthenticateUser(ctx context.Context, user *usermodels.User) (string, error) {
	return uc.repo.AuthenticateUser(ctx, user.Email, user.Password)
}

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
