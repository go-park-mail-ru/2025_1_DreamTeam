package repository

import (
	"context"
	"mime/multipart"
	"skillForce/internal/models"
)

type Repository interface {
	RegisterUser(ctx context.Context, user *models.User) (string, error)
	AuthenticateUser(ctx context.Context, email string, password string) (string, error)
	GetUserByCookie(ctx context.Context, cookieValue string) (*models.UserProfile, error)
	LogoutUser(ctx context.Context, userId int) error
	GetBucketCourses(ctx context.Context) ([]*models.Course, error)
	UpdateProfile(ctx context.Context, userId int, userProfile *models.UserProfile) error
	UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	UpdateProfilePhoto(ctx context.Context, url string, userId int) error
}
