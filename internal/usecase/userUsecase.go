package usecase

import (
	"skillForce/internal/hash"
	"skillForce/internal/models"
	"skillForce/internal/repository"
)

// UserUsecase - структура бизнес-логики
type UserUsecase struct {
	repo *repository.Database
}

// NewUserUsecase - конструктор
func NewUserUsecase(repo *repository.Database) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// RegisterUser - регистрация пользователя
func (uc *UserUsecase) RegisterUser(user *models.User) error {
	err := hash.HashPasswordAndCreateSalt(user)
	if err != nil {
		return err
	}

	return uc.repo.RegisterUser(user)
}

// AuthenticateUser - авторизация пользователя
func (uc *UserUsecase) AuthenticateUser(user *models.User) (int, error) {
	return uc.repo.AuthenticateUser(user.Email, user.Password)
}

// GetUserByCookie - получение пользователя по cookie
func (uc *UserUsecase) GetUserByCookie(cookieValue string) (*models.User, error) {
	return uc.repo.GetUserByCookie(cookieValue)
}

func (uc *UserUsecase) LogoutUser(userId int) error {
	return uc.repo.LogoutUser(userId)
}
