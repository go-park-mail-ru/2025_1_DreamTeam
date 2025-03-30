package infrastructure

import (
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

func (i *Infrastructure) RegisterUser(user *models.User) (string, error) {
	return i.Database.RegisterUser(user)
}

func (i *Infrastructure) AuthenticateUser(email string, password string) (string, error) {
	return i.Database.AuthenticateUser(email, password)
}

func (i *Infrastructure) GetUserByCookie(cookieValue string) (*models.UserProfile, error) {
	return i.Database.GetUserByCookie(cookieValue)
}

func (i *Infrastructure) LogoutUser(userId int) error {
	return i.Database.LogoutUser(userId)
}

func (i *Infrastructure) GetBucketCourses() ([]*models.Course, error) {
	return i.Database.GetBucketCourses()
}

func (i *Infrastructure) UpdateProfile(userId int, userProfile *models.UserProfile) error {
	return i.Database.UpdateProfile(userId, userProfile)
}

func (i *Infrastructure) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	return i.Minio.UploadFileToMinIO(file, fileHeader)
}

func (i *Infrastructure) UpdateProfilePhoto(photo_url string, userId int) error {
	return i.Database.UpdateProfilePhoto(photo_url, userId)
}
