package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"skillForce/internal/models"
	"skillForce/internal/repository"

	"golang.org/x/crypto/argon2"
)

// UserUsecase - структура бизнес-логики
type UserUsecase struct {
	repo *repository.UserRepository
}

// NewUserUsecase - конструктор
func NewUserUsecase(repo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// Хэширование пароля с солью
func hashPassword(user *models.User) error {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatal(err)
		return errors.New("cannot generate salt")
	}

	hash := argon2.IDKey([]byte(user.Password), salt, 1, 64*1024, 4, 32)
	hashedPassword := fmt.Sprintf("%s$___$%s", base64.StdEncoding.EncodeToString(salt), base64.StdEncoding.EncodeToString(hash))

	user.Salt = salt
	user.Password = hashedPassword

	return nil
}

// RegisterUser - регистрация пользователя
func (uc *UserUsecase) RegisterUser(user *models.User) error {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return errors.New("missing required fields")
	}

	err := hashPassword(user)
	if err != nil {
		return err
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
