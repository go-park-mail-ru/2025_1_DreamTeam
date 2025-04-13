package usecase

import (
	"context"
	"io"
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
	GetCourseLesson(ctx context.Context, userId int, courseId int) (*dto.LessonDTO, error)
	GetNextLesson(ctx context.Context, userId int, cousreId int, lessonId int) (*dto.LessonDTO, error)
	MarkLessonAsNotCompleted(ctx context.Context, userId int, lessonId int) error
	GetCourseRoadmap(ctx context.Context, userId int, courseId int) (*dto.CourseRoadmapDTO, error)
	GetCourse(ctx context.Context, courseId int) (*dto.CourseDTO, error)
	DeleteProfilePhoto(ctx context.Context, userId int) error
	GetVideoUrl(ctx context.Context, lesson_id int) (string, error)
	GetMeta(ctx context.Context, name string) (dto.VideoMeta, error)
	GetFragment(ctx context.Context, name string, start, end int64) (io.ReadCloser, error)
}

type Usecase struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) *Usecase {
	return &Usecase{repo: repo}
}
