package usecase

import (
	"context"
	"mime/multipart"
	"skillForce/internal/models"
	"skillForce/internal/models/dto"
	"skillForce/internal/repository"
)

type UsecaseInterface interface {
	RegisterUser(ctx context.Context, user *models.User) (string, error)
	AuthenticateUser(ctx context.Context, user *models.User) (string, error)
	GetUserByCookie(ctx context.Context, cookieValue string) (*models.UserProfile, error)
	LogoutUser(ctx context.Context, userId int) error
	UpdateProfile(ctx context.Context, userId int, userProfile *models.UserProfile) error
	UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	SaveProfilePhoto(ctx context.Context, url string, userId int) (string, error)
	GetBucketCourses(ctx context.Context) ([]*dto.CourseDTO, error)
}

type Usecase struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) *Usecase {
	return &Usecase{repo: repo}
}
