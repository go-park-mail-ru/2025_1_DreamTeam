package hash

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"skillForce/internal/models"

	"golang.org/x/crypto/argon2"
)

func hashPassword(password string, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	hashedPassword := fmt.Sprintf("%s$___$%s", base64.StdEncoding.EncodeToString(salt), base64.StdEncoding.EncodeToString(hash))
	return hashedPassword
}

// Хэширование пароля с солью
func HashPasswordAndCreateSalt(user *models.User) error {
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

func CheckPassword(password string, passwordFromDB string, saltBytes []byte) bool {
	hashedInputPassword := hashPassword(password, saltBytes)
	return hashedInputPassword == passwordFromDB
}
