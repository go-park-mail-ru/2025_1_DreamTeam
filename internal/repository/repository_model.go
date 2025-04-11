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
	GetCoursesRaitings(ctx context.Context, bucketCoursesWithoutRating []*models.Course) (map[int]float32, error)
	GetCoursesTags(ctx context.Context, bucketCoursesWithoutTags []*models.Course) (map[int][]string, error)
	GetCourseById(ctx context.Context, courseId int) (*models.Course, error)
	FillLastLessonHeader(ctx context.Context, userId int, courseId int, LessonHeader *dto.LessonDtoHeader) (int, int, string, error)
	FillLessonHeaderByLessonId(ctx context.Context, userId int, courseId int, currentLessonId int, LessonHeader *dto.LessonDtoHeader) error
	GetLessonBlocks(ctx context.Context, currentLessonId int) ([]string, error)
	GetLessonFooters(ctx context.Context, currentLessonId int) ([]int, error)
	MarkLessonCompleted(ctx context.Context, userId int, courseId int, lessonId int) error
	MarkLessonAsNotCompleted(ctx context.Context, userId int, lessonId int) error
	GetCourseParts(ctx context.Context, courseId int) ([]*models.CoursePart, error)
	GetPartBuckets(ctx context.Context, partId int) ([]*models.LessonBucket, error)
	GetBucketLessons(ctx context.Context, userId int, courseId int, bucketId int) ([]*models.LessonPoint, error)
	AddUserToCourse(ctx context.Context, userId int, courseId int) error
	GetCoursesPurchases(ctx context.Context, bucketCoursesWithoutPurchases []*models.Course) (map[int]int, error)
	GetBucketByLessonId(ctx context.Context, lessonId int) (*models.LessonBucket, error)
	DeleteProfilePhoto(ctx context.Context, userId int) error
}
