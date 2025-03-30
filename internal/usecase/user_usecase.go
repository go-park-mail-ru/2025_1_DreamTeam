package usecase

import (
	"mime/multipart"
	"skillForce/internal/models"
	"skillForce/internal/repository"
	"skillForce/pkg/hash"
)

type UserUsecaseInterface interface {
	RegisterUser(user *models.User) (string, error)
	AuthenticateUser(user *models.User) (string, error)
	GetUserByCookie(cookieValue string) (*models.UserProfile, error)
	LogoutUser(userId int) error
	UpdateProfile(userId int, userProfile *models.UserProfile) error
	UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
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
func (uc *UserUsecase) RegisterUser(user *models.User) (string, error) {
	err := hash.HashPasswordAndCreateSalt(user)
	if err != nil {
		return "", err
	}

	return uc.repo.RegisterUser(user)
}

// AuthenticateUser - авторизация пользователя
func (uc *UserUsecase) AuthenticateUser(user *models.User) (string, error) {
	return uc.repo.AuthenticateUser(user.Email, user.Password)
}

// GetUserByCookie - получение пользователя по cookie
func (uc *UserUsecase) GetUserByCookie(cookieValue string) (*models.UserProfile, error) {
	return uc.repo.GetUserByCookie(cookieValue)
}

func (uc *UserUsecase) LogoutUser(userId int) error {
	return uc.repo.LogoutUser(userId)
}

func (uc *UserUsecase) UpdateProfile(userId int, userProfile *models.UserProfile) error {
	return uc.repo.UpdateProfile(userId, userProfile)
}

func (uc *UserUsecase) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	return uc.repo.UploadFile(file, fileHeader)
}
