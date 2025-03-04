package repository

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"skillForce/internal/models"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/argon2"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(connStr string) (*Database, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &Database{conn: db}, nil
}

func (d *Database) Close() {
	d.conn.Close()
}

func (d *Database) GetBucketCourses() ([]*models.Course, error) {
	//TODO: можно заморочиться и сделать самописную пагинацию через LIMIT OFFSET
	var bucketCourses []*models.Course
	rows, err := d.conn.Query("SELECT id, creator_user_id, title, description, avatar_src, price, time_to_pass FROM course LIMIT 16")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.Id, &course.CreatorId, &course.Title, &course.Description, &course.ScrImage, &course.Price, &course.TimeToPass); err != nil {
			return nil, err
		}
		bucketCourses = append(bucketCourses, &course)
	}

	return bucketCourses, nil
}

func (d *Database) userExists(email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM usertable WHERE email = $1)"
	err := d.conn.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func (d *Database) RegisterUser(user *models.User) error {
	emailExists, err := d.userExists(user.Email)
	if err != nil {
		return err
	}
	if emailExists {
		return errors.New("email exists")
	}
	saltBase64 := base64.StdEncoding.EncodeToString(user.Salt)
	_, err2 := d.conn.Exec("INSERT INTO usertable (email, password, salt) VALUES ($1, $2, $3)", user.Email, user.Password, saltBase64)
	if err2 != nil {
		return err2
	}
	return nil
}

func (d *Database) GetUserByCookie(cookieValue string) (*models.User, error) {
	var user models.User
	err := d.conn.QueryRow("SELECT u.email, u.password, u.salt FROM usertable u JOIN sessions s ON u.id = s.user_id WHERE s.token = $1 AND s.expires > NOW();",
		cookieValue).Scan(&user.Email, &user.Password, &user.Salt)
	return &user, err
}

func hashPassword(password string, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	hashedPassword := fmt.Sprintf("%s$___$%s", base64.StdEncoding.EncodeToString(salt), base64.StdEncoding.EncodeToString(hash))
	return hashedPassword
}

func (d *Database) Authenticate(email string, password string) (int, error) {
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
	hashedInputPassword := hashPassword(password, saltBytes)
	if hashedInputPassword != passwordFromDB {
		return -1, errors.New("email or password incorrect")
	}
	return id, nil
}
