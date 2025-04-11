package infrastructure

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"skillForce/config"
	"skillForce/internal/models"
	"skillForce/internal/models/dto"
	"skillForce/internal/repository/infrastructure/minio"
	"skillForce/internal/repository/infrastructure/postgres"
)

type Infrastructure struct {
	Database *postgres.Database
	Minio    *minio.Minio
}

func NewInfrastructure(conf *config.Config) *Infrastructure {
	mn, err := minio.NewMinio(conf.Minio.Endpoint, conf.Minio.AccessKey, conf.Minio.SecretAccessKey, conf.Minio.UseSSL, conf.Minio.BucketName)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Name)
	database, err := postgres.NewDatabase(dsn, conf.Secrets.JwtSessionSecret)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return &Infrastructure{
		Database: database,
		Minio:    mn,
	}
}

func (i *Infrastructure) Close() {
	i.Database.Close()
}

func (i *Infrastructure) RegisterUser(ctx context.Context, user *models.User) (string, error) {
	return i.Database.RegisterUser(ctx, user)
}

func (i *Infrastructure) AuthenticateUser(ctx context.Context, email string, password string) (string, error) {
	return i.Database.AuthenticateUser(ctx, email, password)
}

func (i *Infrastructure) GetUserByCookie(ctx context.Context, cookieValue string) (*models.UserProfile, error) {
	return i.Database.GetUserByCookie(ctx, cookieValue)
}

func (i *Infrastructure) LogoutUser(ctx context.Context, userId int) error {
	return i.Database.LogoutUser(ctx, userId)
}

func (i *Infrastructure) GetBucketCourses(ctx context.Context) ([]*models.Course, error) {
	return i.Database.GetBucketCourses(ctx)
}

func (i *Infrastructure) UpdateProfile(ctx context.Context, userId int, userProfile *models.UserProfile) error {
	return i.Database.UpdateProfile(ctx, userId, userProfile)
}

func (i *Infrastructure) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	return i.Minio.UploadFileToMinIO(ctx, file, fileHeader)
}

func (i *Infrastructure) UpdateProfilePhoto(ctx context.Context, photo_url string, userId int) (string, error) {
	return i.Database.UpdateProfilePhoto(ctx, photo_url, userId)
}

func (i *Infrastructure) GetCoursesRaitings(ctx context.Context, bucketCoursesWithoutRating []*models.Course) (map[int]float32, error) {
	return i.Database.GetCoursesRaitings(ctx, bucketCoursesWithoutRating)
}

func (i *Infrastructure) GetCoursesTags(ctx context.Context, bucketCoursesWithoutTags []*models.Course) (map[int][]string, error) {
	return i.Database.GetCoursesTags(ctx, bucketCoursesWithoutTags)
}

func (i *Infrastructure) GetCourseById(ctx context.Context, courseId int) (*models.Course, error) {
	return i.Database.GetCourseById(ctx, courseId)
}

func (i *Infrastructure) FillLastLessonHeader(ctx context.Context, userId int, courseId int, LessonHeader *dto.LessonDtoHeader) (int, int, string, error) {
	return i.Database.FillLastLessonHeader(ctx, userId, courseId, LessonHeader)
}

func (i *Infrastructure) GetLessonBlocks(ctx context.Context, currentLessonId int) ([]string, error) {
	return i.Database.GetLessonBlocks(ctx, currentLessonId)
}

func (i *Infrastructure) GetLessonFooters(ctx context.Context, currentLessonId int) ([]int, error) {
	return i.Database.GetLessonFooters(ctx, currentLessonId)
}

func (i *Infrastructure) MarkLessonCompleted(ctx context.Context, userId int, courseId int, lessonId int) error {
	return i.Database.MarkLessonCompleted(ctx, userId, courseId, lessonId)
}

func (i *Infrastructure) MarkLessonAsNotCompleted(ctx context.Context, userId int, lessonId int) error {
	return i.Database.MarkLessonAsNotCompleted(ctx, userId, lessonId)
}

func (i *Infrastructure) GetCourseParts(ctx context.Context, courseId int) ([]*models.CoursePart, error) {
	return i.Database.GetCourseParts(ctx, courseId)
}

func (i *Infrastructure) GetPartBuckets(ctx context.Context, partId int) ([]*models.LessonBucket, error) {
	return i.Database.GetPartBuckets(ctx, partId)
}

func (i *Infrastructure) GetBucketLessons(ctx context.Context, userId int, courseId int, bucketId int) ([]*models.LessonPoint, error) {
	return i.Database.GetBucketLessons(ctx, userId, courseId, bucketId)
}

func (i *Infrastructure) AddUserToCourse(ctx context.Context, userId int, courseId int) error {
	return i.Database.AddUserToCourse(ctx, userId, courseId)
}

func (i *Infrastructure) GetCoursesPurchases(ctx context.Context, bucketCoursesWithoutPurchases []*models.Course) (map[int]int, error) {
	return i.Database.GetCoursesPurchases(ctx, bucketCoursesWithoutPurchases)
}

func (i *Infrastructure) GetBucketByLessonId(ctx context.Context, lessonId int) (*models.LessonBucket, error) {
	return i.Database.GetBucketByLessonId(ctx, lessonId)
}

func (i *Infrastructure) FillLessonHeaderByLessonId(ctx context.Context, userId int, courseId int, currentLessonId int, LessonHeader *dto.LessonDtoHeader) error {
	return i.Database.FillLessonHeaderByLessonId(ctx, userId, courseId, currentLessonId, LessonHeader)
}

func (i *Infrastructure) DeleteProfilePhoto(ctx context.Context, userId int) error {
	return i.Database.DeleteProfilePhoto(ctx, userId)
}
