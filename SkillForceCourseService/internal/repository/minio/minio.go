package minio

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go"
)

type Minio struct {
	MinioClient          *minio.Client
	SertificatesBucket   string
	CoursesAvatarsBucket string
}

func NewMinio(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool, bucketName string, courseBucketName string) (*Minio, error) {
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	return &Minio{MinioClient: minioClient, SertificatesBucket: bucketName, CoursesAvatarsBucket: courseBucketName}, err
}

func (mn *Minio) UploadFileToMinIO(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {

	uniqueID := uuid.New().String()
	ext := ""
	if fileHeader.Filename != "" {
		ext = fileHeader.Filename[strings.LastIndex(fileHeader.Filename, "."):]
	}
	objectName := fmt.Sprintf("%s%s", uniqueID, ext)
	contentType := fileHeader.Header.Get("Content-Type")

	_, err := mn.MinioClient.PutObject(
		mn.SertificatesBucket,
		objectName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("https://skill-force.ru/%s/%s", mn.SertificatesBucket, objectName)
	return fileURL, nil
}

func (mn *Minio) UploadCourseAvatarToMinIO(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {

	uniqueID := uuid.New().String()
	ext := ""
	if fileHeader.Filename != "" {
		ext = fileHeader.Filename[strings.LastIndex(fileHeader.Filename, "."):]
	}
	objectName := fmt.Sprintf("%s%s", uniqueID, ext)
	contentType := fileHeader.Header.Get("Content-Type")

	_, err := mn.MinioClient.PutObject(
		mn.CoursesAvatarsBucket,
		objectName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("https://skill-force.ru/%s/%s", mn.CoursesAvatarsBucket, objectName)
	return fileURL, nil
}
