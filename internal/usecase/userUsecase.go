package usecase

import (
	"errors"
	"skillForce/internal/models"
	"skillForce/internal/repository"
)

// UserUsecase - структура бизнес-логики
type UserUsecase struct {
	repo *repository.UserRepository
}

// NewUserUsecase - конструктор
func NewUserUsecase(repo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// RegisterUser - регистрация пользователя
func (uc *UserUsecase) RegisterUser(user *models.User) error {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return errors.New("missing required fields")
	}
	return uc.repo.Save(user)
}

// AuthenticateUser - авторизация пользователя
func (uc *UserUsecase) AuthenticateUser(user *models.User) (int, error) {
	if user.Email == "" || user.Password == "" {
		return 0, errors.New("missing required fields")
	}
	return uc.repo.Authenticate(user.Email, user.Password)
}
