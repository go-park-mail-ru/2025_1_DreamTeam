package repository

import (
	"encoding/base64"
	"errors"
	"log"
	"skillForce/internal/hash"
	"skillForce/internal/models"
	"time"
)

// userExists - проверяет, существует ли пользователь с указанным email
func (d *Database) userExists(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM usertable WHERE email = $1)"
	err := d.conn.QueryRow(query, email).Scan(&exists)
	return exists, err
}

// RegisterUser - сохраняет нового пользователя в базе данных и создает сессию, тоже сохраняя её в базе
func (d *Database) RegisterUser(user *models.User) error {
	emailExists, err := d.userExists(user.Email)
	if err != nil {
		return err
	}
	if emailExists {
		return errors.New("email exists")
	}
	saltBase64 := base64.StdEncoding.EncodeToString(user.Salt)
	_, err = d.conn.Exec("INSERT INTO usertable (email, name, password, salt) VALUES ($1, $2, $3, $4)", user.Email, user.Name, user.Password, saltBase64)
	if err != nil {
		return err
	}

	log.Printf("Save user %+v in db", user)

	rows, err := d.conn.Query("SELECT id FROM usertable WHERE email = $1", user.Email)
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rows.Scan(&user.Id)
		if err != nil {
			return err
		}
	}

	_, err = d.conn.Exec("INSERT INTO sessions (user_id, token, expires) VALUES ($1, $2, $3)", user.Id, user.Id, time.Now().Add(10*time.Hour))
	if err != nil {
		return err
	}

	log.Printf("Save sessions of user %+v in db", user)

	return nil
}

// GetUserByCookie - получение пользователя по cookie
func (d *Database) GetUserByCookie(cookieValue string) (*models.User, error) {
	var user models.User
	err := d.conn.QueryRow("SELECT u.id, u.email, u.name, u.password, u.salt FROM usertable u JOIN sessions s ON u.id = s.user_id WHERE s.token = $1 AND s.expires > NOW();",
		cookieValue).Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.Salt)
	return &user, err
}

// AuthenticateUser - проверяет есть ли пользователь с указанным email и паролем в базе данных, елси есть - возвращает его id и сохраняет сессию в базе
func (d *Database) AuthenticateUser(email string, password string) (int, error) {
	var id int
	emailExists, err := d.userExists(email)
	if err != nil {
		return -1, err
	}
	if !emailExists {
		return -1, errors.New("email or password incorrect")
	}
	var passwordFromDB string
	var salt string
	err2 := d.conn.QueryRow("SELECT id, password, salt FROM usertable WHERE email = $1", email).Scan(&id, &passwordFromDB, &salt)
	if err2 != nil {
		return -1, err2
	}
	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return -1, err
	}

	if !hash.CheckPassword(password, passwordFromDB, saltBytes) {
		return -1, errors.New("email or password incorrect")
	}

	log.Printf("Login user with email %+v in db", email)

	_, err = d.conn.Exec("INSERT INTO sessions (user_id, token, expires) VALUES ($1, $2, $3)", id, id, time.Now().Add(10*time.Hour))
	if err != nil {
		return -1, err
	}

	log.Printf("Save sessions of user with email %+v in db", email)
	return id, nil
}

// LogoutUser - удаляет сессию пользователя из базы данных
func (d *Database) LogoutUser(userId int) error {
	_, err := d.conn.Exec("DELETE FROM sessions WHERE user_id = $1", userId)
	return err
}
