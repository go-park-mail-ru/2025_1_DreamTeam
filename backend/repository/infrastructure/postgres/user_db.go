package postgres

import (
	"encoding/base64"
	"errors"
	"log"
	"skillForce/backend/models"
	"skillForce/internal/hash"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var SECRET = []byte("dream_team_secret_jehpfqjbhjfkjlGUGeqJIBxcfimor")

func (d *Database) saveSession(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
	})

	secretToken, err := token.SignedString(SECRET)
	if err != nil {
		return "", err
	}

	_, err = d.conn.Exec("INSERT INTO sessions (user_id, token, expires) VALUES ($1, $2, $3)", userId, secretToken, time.Now().Add(10*time.Hour))
	if err != nil {
		return "", err
	}

	return secretToken, nil
}

// userExists - проверяет, существует ли пользователь с указанным email
func (d *Database) userExists(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM usertable WHERE email = $1)"
	err := d.conn.QueryRow(query, email).Scan(&exists)
	return exists, err
}

// RegisterUser - сохраняет нового пользователя в базе данных и создает сессию, тоже сохраняя её в базе
func (d *Database) RegisterUser(user *models.User) (string, error) {
	emailExists, err := d.userExists(user.Email)
	if err != nil {
		return "", err
	}
	if emailExists {
		return "", errors.New("email exists")
	}
	saltBase64 := base64.StdEncoding.EncodeToString(user.Salt)
	_, err = d.conn.Exec("INSERT INTO usertable (email, name, password, salt) VALUES ($1, $2, $3, $4)", user.Email, user.Name, user.Password, saltBase64)
	if err != nil {
		return "", err
	}

	log.Printf("Save user %+v in db", user)

	rows, err := d.conn.Query("SELECT id FROM usertable WHERE email = $1", user.Email)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		err = rows.Scan(&user.Id)
		if err != nil {
			return "", err
		}
	}

	cookieValue, err := d.saveSession(user.Id)
	if err != nil {
		return "", err
	}

	log.Printf("Save sessions of user %+v in db", user)

	return cookieValue, nil
}

// GetUserByCookie - получение пользователя по cookie
func (d *Database) GetUserByCookie(cookieValue string) (*models.UserProfile, error) {
	var userProfile models.UserProfile
	err := d.conn.QueryRow("SELECT u.id, u.email, u.name, COALESCE(u.bio, ''), u.avatar_src, u.hide_email FROM usertable u JOIN sessions s ON u.id = s.user_id WHERE s.token = $1 AND s.expires > NOW();",
		cookieValue).Scan(&userProfile.Id, &userProfile.Email, &userProfile.Name, &userProfile.Bio, &userProfile.AvatarSrc, &userProfile.HideEmail)
	if err != nil {
		return nil, err
	}
	return &userProfile, err
}

// AuthenticateUser - проверяет есть ли пользователь с указанным email и паролем в базе данных, елси есть - возвращает его id и сохраняет сессию в базе
func (d *Database) AuthenticateUser(email string, password string) (string, error) {
	var id int
	emailExists, err := d.userExists(email)
	if err != nil {
		return "", err
	}
	if !emailExists {
		return "", errors.New("email or password incorrect")
	}
	var passwordFromDB string
	var salt string
	err2 := d.conn.QueryRow("SELECT id, password, salt FROM usertable WHERE email = $1", email).Scan(&id, &passwordFromDB, &salt)
	if err2 != nil {
		return "", err2
	}
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}

	if !hash.CheckPassword(password, passwordFromDB, saltBytes) {
		return "", errors.New("email or password incorrect")
	}

	log.Printf("Login user with email %+v in db", email)

	cookieValue, err := d.saveSession(id)
	if err != nil {
		return "", err
	}

	log.Printf("Save sessions of user with email %+v in db", email)
	return cookieValue, nil
}

// LogoutUser - удаляет сессию пользователя из базы данных
func (d *Database) LogoutUser(userId int) error {
	_, err := d.conn.Exec("DELETE FROM sessions WHERE user_id = $1", userId)
	if err != nil {
		return err
	}
	log.Printf("Logout user with id %+v in db", userId)
	return err
}

func (d *Database) UpdateProfile(userId int, userProfile *models.UserProfile) error { // TODO: полчистить avatar_src
	_, err := d.conn.Exec("UPDATE usertable SET email = $1, name = $2, bio = $3, avatar_src = $4 WHERE id = $5",
		userProfile.Email, userProfile.Name, userProfile.Bio, userProfile.AvatarSrc, userId)
	if err != nil {
		return err
	}
	log.Printf("Update profile %+v of user with id %+v in db", userProfile, userId)
	return err
}
