package repository

import (
	"context"
	"io"
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
	GetLastLessonHeader(ctx context.Context, userId int, courseId int) (*dto.LessonDtoHeader, int, string, bool, error)
	GetLessonHeaderByLessonId(ctx context.Context, userId int, currentLessonId int) (*dto.LessonDtoHeader, error)
	GetLessonBlocks(ctx context.Context, currentLessonId int) ([]string, error)
	GetLessonVideo(ctx context.Context, currentLessonId int) ([]string, error)
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
	GetVideoUrl(ctx context.Context, lesson_id int) (string, error)
	GetVideoRange(ctx context.Context, name string, start, end int64) (io.ReadCloser, error)
	Stat(ctx context.Context, name string) (dto.VideoMeta, error)
	ValidUser(ctx context.Context, user *models.User) (string, error)
	SendRegMail(ctx context.Context, user *models.User, token string) error
	SendWelcomeMail(ctx context.Context, user *models.User) error
	GetUserByToken(ctx context.Context, token string) (*models.User, error)
	SendWelcomeCourseMail(ctx context.Context, user *models.User, courseId int) error
	GetUserById(ctx context.Context, userId int) (*models.User, error)
	IsUserPurchasedCourse(ctx context.Context, userId int, courseId int) (bool, error)
	GetLessonById(ctx context.Context, lessonId int) (*models.LessonPoint, error)
}
