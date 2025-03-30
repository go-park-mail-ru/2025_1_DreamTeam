package repository

import (
	"mime/multipart"
	"skillForce/internal/models"
)

type Repository interface {
	RegisterUser(user *models.User) (string, error)
	AuthenticateUser(email string, password string) (string, error)
	GetUserByCookie(cookieValue string) (*models.UserProfile, error)
	LogoutUser(userId int) error
	GetBucketCourses() ([]*models.Course, error)
	UpdateProfile(userId int, userProfile *models.UserProfile) error
	UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}
