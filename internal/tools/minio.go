package tools

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go"
)

var MinioClient *minio.Client
var AvatarsBucket string

func InitMinio(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool, bucketName string) error {
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	MinioClient = minioClient
	AvatarsBucket = bucketName
	return err
}

func UploadToMinIO(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {

	uniqueID := uuid.New().String()
	ext := ""
	if fileHeader.Filename != "" {
		ext = fileHeader.Filename[strings.LastIndex(fileHeader.Filename, "."):]
	}
	objectName := fmt.Sprintf("%s%s", uniqueID, ext)
	contentType := fileHeader.Header.Get("Content-Type")

	_, err := MinioClient.PutObject(
		AvatarsBucket,
		objectName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("http://217.16.21.64:8006/%s/%s", AvatarsBucket, objectName)
	return fileURL, nil
}
