package courseCourseInfrastructure

import (
	"context"
	"fmt"
	"io"
	"log"
	"skillForce/config"
	"skillForce/internal/models/dto"
	"skillForce/internal/repository/course/minio"
	"skillForce/internal/repository/course/postgres"
)

type CourseInfrastructure struct {
	Database *postgres.Database
	Minio    *minio.Minio
}

func NewCourseInfrastructure(conf *config.Config) *CourseInfrastructure {
	mn, err := minio.NewMinio(conf.Minio.Endpoint, conf.Minio.AccessKey, conf.Minio.SecretAccessKey, conf.Minio.UseSSL, conf.Minio.BucketName, conf.Minio.VideoBucket)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Name)
	database, err := postgres.NewDatabase(dsn, conf.Secrets.JwtSessionSecret)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return &CourseInfrastructure{
		Database: database,
		Minio:    mn,
	}
}

func (ci *CourseInfrastructure) Close() {
	ci.Database.Close()
}

func (ci *CourseInfrastructure) Stat(ctx context.Context, name string) (dto.VideoMeta, error) {
	return ci.Minio.Stat(ctx, name)
}

func (ci *CourseInfrastructure) GetVideoRange(ctx context.Context, name string, start, end int64) (io.ReadCloser, error) {
	return ci.Minio.GetVideoRange(ctx, name, start, end)
}

func (ci *CourseInfrastructure) GetVideoUrl(ctx context.Context, lessonId int) (string, error) {
	return ci.Database.GetVideoUrl(ctx, lessonId)
}
