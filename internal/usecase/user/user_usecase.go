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
	authRepo    UserAuthRepository
	profileRepo UserProfileRepository
	mailRepo    UserMailRepository
}

func NewUserUsecase(
	authRepo UserAuthRepository,
	profileRepo UserProfileRepository,
	mailRepo UserMailRepository,
) *UserUsecase {
	return &UserUsecase{
		authRepo:    authRepo,
		profileRepo: profileRepo,
		mailRepo:    mailRepo,
	}
}

func (uc *UserUsecase) ValidUser(ctx context.Context, user *usermodels.User) error {
	token, err := uc.authRepo.ValidUser(ctx, user)
	if err != nil {
		logs.PrintLog(ctx, "ValidUser", fmt.Sprintf("%+v", err))
		return err
	}

	err = uc.mailRepo.SendRegMail(ctx, user, token)
	if err != nil {
		logs.PrintLog(ctx, "SendMail", fmt.Sprintf("%+v", err))
		return err
	}
	return nil
}

func (uc *UserUsecase) RegisterUser(ctx context.Context, token string) (string, error) {
	user, err := uc.authRepo.GetUserByToken(ctx, token)
	if err != nil {
		return "", err
	}

	err = hash.HashPasswordAndCreateSalt(user)
	if err != nil {
		return "", err
	}

	_ = uc.mailRepo.SendWelcomeMail(ctx, user)

	return uc.authRepo.RegisterUser(ctx, user)
}

func (uc *UserUsecase) AuthenticateUser(ctx context.Context, user *usermodels.User) (string, error) {
	return uc.authRepo.AuthenticateUser(ctx, user.Email, user.Password)
}

func (uc *UserUsecase) GetUserByCookie(ctx context.Context, cookieValue string) (*usermodels.UserProfile, error) {
	return uc.authRepo.GetUserByCookie(ctx, cookieValue)
}

func (uc *UserUsecase) LogoutUser(ctx context.Context, userId int) error {
	return uc.authRepo.LogoutUser(ctx, userId)
}

func (uc *UserUsecase) UpdateProfile(ctx context.Context, userId int, userProfile *usermodels.UserProfile) error {
	return uc.profileRepo.UpdateProfile(ctx, userId, userProfile)
}

func (uc *UserUsecase) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	return uc.profileRepo.UploadFile(ctx, file, fileHeader)
}

func (uc *UserUsecase) SaveProfilePhoto(ctx context.Context, url string, userId int) (string, error) {
	return uc.profileRepo.UpdateProfilePhoto(ctx, url, userId)
}

func (uc *UserUsecase) DeleteProfilePhoto(ctx context.Context, userId int) error {
	return uc.profileRepo.DeleteProfilePhoto(ctx, userId)
}
