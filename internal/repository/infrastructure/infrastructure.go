package infrastructure

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"skillForce/config"
	"skillForce/internal/models"
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
	database, err := postgres.NewDatabase(dsn)
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

func (i *Infrastructure) UpdateProfilePhoto(ctx context.Context, photo_url string, userId int) error {
	return i.Database.UpdateProfilePhoto(ctx, photo_url, userId)
}

func (i *Infrastructure) GetCoursesRaitings(ctx context.Context, bucketCoursesWithoutRating []*models.Course) (map[int]models.CourseRating, error) {
	return i.Database.GetCoursesRaitings(ctx, bucketCoursesWithoutRating)
}

func (i *Infrastructure) GetCoursesTags(ctx context.Context, bucketCoursesWithoutTags []*models.Course) (map[int][]string, error) {
	return i.Database.GetCoursesTags(ctx, bucketCoursesWithoutTags)
}
