package repository

import (
	"encoding/base64"
	"errors"
	"log"
	"skillForce/internal/models"
	"strings"

	"golang.org/x/crypto/argon2"
)

// UserRepository - структура хранилища пользователей
type UserRepository struct {
	users map[int]*models.User
}

// NewUserRepository - создание нового репозитория пользователей
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[int]*models.User),
	}
}

// Проверка пароля
func checkPassword(password, storedHash string) (bool, error) {
	parts := strings.Split(storedHash, "$___$")
	cleanStoredHash := parts[1]
	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return true, errors.New("cannot decode salt")
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	expectedHash := base64.StdEncoding.EncodeToString(hash)
	return expectedHash == cleanStoredHash, nil
}

// Save - сохранение пользователя
func (r *UserRepository) Save(user *models.User) error {
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return errors.New("email exists")
		}
	}
	user.Id = len(r.users)
	r.users[user.Id] = user
	log.Printf("user saved: %+v", user)
	return nil
}

// Authenticate - аутентификация пользователя
func (r *UserRepository) Authenticate(email string, password string) (int, error) {
	for _, existingUser := range r.users {
		isPasswordMatch, err := checkPassword(password, existingUser.Password)
		if err != nil {
			return 0, err
		}
		if existingUser.Email == email && isPasswordMatch {
			log.Printf("user authenticated: %+v", existingUser)
			return existingUser.Id, nil
		}
	}
	return 0, errors.New("email or password incorrect")
}
