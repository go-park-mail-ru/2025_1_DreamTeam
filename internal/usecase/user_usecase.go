package usecase

import (
	"context"
	"mime/multipart"
	"skillForce/internal/models"
	"skillForce/internal/repository"
	"skillForce/pkg/hash"
)

type UserUsecaseInterface interface {
	RegisterUser(ctx context.Context, user *models.User) (string, error)
	AuthenticateUser(ctx context.Context, user *models.User) (string, error)
	GetUserByCookie(ctx context.Context, cookieValue string) (*models.UserProfile, error)
	LogoutUser(ctx context.Context, userId int) error
	UpdateProfile(ctx context.Context, userId int, userProfile *models.UserProfile) error
	UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	SaveProfilePhoto(ctx context.Context, url string, userId int) error
}

// UserUsecase - структура бизнес-логики, которая ожидает интерфейс репозитория
type UserUsecase struct {
	repo repository.Repository
}

// NewUserUsecase - конструктор
func NewUserUsecase(repo repository.Repository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// RegisterUser - регистрация пользователя
func (uc *UserUsecase) RegisterUser(ctx context.Context, user *models.User) (string, error) {
	err := hash.HashPasswordAndCreateSalt(user)
	if err != nil {
		return "", err
	}

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

func (uc *UserUsecase) SaveProfilePhoto(ctx context.Context, url string, userId int) error {
	return uc.repo.UpdateProfilePhoto(ctx, url, userId)
}
