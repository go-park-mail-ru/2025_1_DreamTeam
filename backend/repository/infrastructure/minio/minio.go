package minio

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go"
)

type Minio struct {
	MinioClient   *minio.Client
	AvatarsBucket string
}

func NewMinio(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool, bucketName string) (*Minio, error) {
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	return &Minio{MinioClient: minioClient, AvatarsBucket: bucketName}, err
}

func (mn *Minio) UploadFileToMinIO(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {

	uniqueID := uuid.New().String()
	ext := ""
	if fileHeader.Filename != "" {
		ext = fileHeader.Filename[strings.LastIndex(fileHeader.Filename, "."):]
	}
	objectName := fmt.Sprintf("%s%s", uniqueID, ext)
	contentType := fileHeader.Header.Get("Content-Type")

	_, err := mn.MinioClient.PutObject(
		mn.AvatarsBucket,
		objectName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("http://217.16.21.64:8006/%s/%s", mn.AvatarsBucket, objectName)
	return fileURL, nil
}
