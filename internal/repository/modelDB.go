package repository

import (
	"database/sql"
	"skillForce/internal/models"

	_ "github.com/lib/pq"
)

type Repository interface {
	RegisterUser(user *models.User) (string, error)
	AuthenticateUser(email string, password string) (string, error)
	GetUserByCookie(cookieValue string) (*models.UserProfile, error)
	LogoutUser(userId int) error
	GetBucketCourses() ([]*models.Course, error)
	UpdateProfile(userId int, userProfile *models.UserProfile) error
}

type Database struct {
	conn *sql.DB
}

// NewDatabase - конструктор
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

// Close - закрытие соединения с базой данных
func (d *Database) Close() {
	d.conn.Close()
}
