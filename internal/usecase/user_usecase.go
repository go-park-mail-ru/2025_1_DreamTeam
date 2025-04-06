package usecase

import (
	"context"
	"mime/multipart"
	"skillForce/internal/models"
	"skillForce/pkg/hash"
)

// RegisterUser - регистрация пользователя
func (uc *Usecase) RegisterUser(ctx context.Context, user *models.User) (string, error) {
	err := hash.HashPasswordAndCreateSalt(user)
	if err != nil {
		return "", err
	}

	return uc.repo.RegisterUser(ctx, user)
}

// AuthenticateUser - авторизация пользователя
func (uc *Usecase) AuthenticateUser(ctx context.Context, user *models.User) (string, error) {
	return uc.repo.AuthenticateUser(ctx, user.Email, user.Password)
}

// GetUserByCookie - получение пользователя по cookie
func (uc *Usecase) GetUserByCookie(ctx context.Context, cookieValue string) (*models.UserProfile, error) {
	return uc.repo.GetUserByCookie(ctx, cookieValue)
}

func (uc *Usecase) LogoutUser(ctx context.Context, userId int) error {
	return uc.repo.LogoutUser(ctx, userId)
}

func (uc *Usecase) UpdateProfile(ctx context.Context, userId int, userProfile *models.UserProfile) error {
	return uc.repo.UpdateProfile(ctx, userId, userProfile)
}

func (uc *Usecase) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	return uc.repo.UploadFile(ctx, file, fileHeader)
}

func (uc *Usecase) SaveProfilePhoto(ctx context.Context, url string, userId int) (string, error) {
	return uc.repo.UpdateProfilePhoto(ctx, url, userId)
}
