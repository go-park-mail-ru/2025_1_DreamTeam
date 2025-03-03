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
	repo *repository.Database
}

// NewUserUsecase - конструктор
func NewUserUsecase(repo *repository.Database) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func hashPassword(password string, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	hashedPassword := fmt.Sprintf("%s$___$%s", base64.StdEncoding.EncodeToString(salt), base64.StdEncoding.EncodeToString(hash))
	return hashedPassword
}

// Хэширование пароля с солью
func hashPasswordAndCreateSalt(user *models.User) error {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatal(err)
		return errors.New("cannot generate salt")
	}

	hashedPassword := hashPassword(user.Password, salt)

	user.Salt = salt
	user.Password = hashedPassword

	return nil
}

// RegisterUser - регистрация пользователя
func (uc *UserUsecase) RegisterUser(user *models.User) error {
	err := hashPasswordAndCreateSalt(user)
	if err != nil {
		return err
	}

	return uc.repo.RegisterUser(user)
}

// AuthenticateUser - авторизация пользователя
func (uc *UserUsecase) AuthenticateUser(user *models.User) (int, error) {
	return uc.repo.Authenticate(user.Email, user.Password)
}

// GetUserByCookie - получение пользователя по cookie
func (uc *UserUsecase) GetUserByCookie(cookieValue string) (*models.User, error) {
	return uc.repo.GetUserByCookie(cookieValue)
}
