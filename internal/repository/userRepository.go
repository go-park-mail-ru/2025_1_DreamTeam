package repository

import (
	"errors"
	"log"
	"skillForce/internal/models"

	"golang.org/x/crypto/bcrypt"
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

// hashPassword - хеширование пароля
func hashPassword(user *models.User) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// checkPasswordHash - проверка соответствия хешированного пароля и пароля
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) //TODO: обработать и добавить в сваггер
	return err == nil
}

// Save - сохранение пользователя
func (r *UserRepository) Save(user *models.User) error {
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return errors.New("email exists")
		}
	}
	user.Id = len(r.users)
	err := hashPassword(user)
	if err != nil {
		return err //TODO: обработать и добавить в сваггер
	}
	r.users[user.Id] = user
	log.Printf("user saved: id: %d %+v", user.Id, user)
	return nil
}

// Authenticate - аутентификация пользователя
func (r *UserRepository) Authenticate(email string, password string) (int, error) {
	// TODO: password to hash
	// Добавить обработку случая, когда пользователя может и не быть в базе
	for _, existingUser := range r.users {
		if existingUser.Email == email && checkPasswordHash(password, existingUser.Password) {
			return existingUser.Id, nil
		}
	}
	return 0, errors.New("email or password incorrect")
}
