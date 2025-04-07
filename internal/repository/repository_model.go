package repository

import (
	"context"
	"mime/multipart"
	"skillForce/internal/models"
	"skillForce/internal/models/dto"
)

type Repository interface {
	RegisterUser(ctx context.Context, user *models.User) (string, error)
	AuthenticateUser(ctx context.Context, email string, password string) (string, error)
	GetUserByCookie(ctx context.Context, cookieValue string) (*models.UserProfile, error)
	LogoutUser(ctx context.Context, userId int) error
	GetBucketCourses(ctx context.Context) ([]*models.Course, error)
	UpdateProfile(ctx context.Context, userId int, userProfile *models.UserProfile) error
	UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	UpdateProfilePhoto(ctx context.Context, url string, userId int) (string, error)
	GetCoursesRaitings(ctx context.Context, bucketCoursesWithoutRating []*models.Course) (map[int]models.CourseRating, error)
	GetCoursesTags(ctx context.Context, bucketCoursesWithoutTags []*models.Course) (map[int][]string, error)
	GetCourseById(ctx context.Context, courseId int) (*models.Course, error)
	FillLessonHeader(ctx context.Context, userId int, courseId int, LessonHeader *dto.LessonDtoHeader) (int, int, string, error)
	GetLessonBlocks(ctx context.Context, currentLessonId int) ([]string, error)
	GetLessonFooters(ctx context.Context, currentLessonId int, currentBucketId int) ([]int, error)
}
