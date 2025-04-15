package minio

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"skillForce/internal/models/dto"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go"
)

type Minio struct {
	MinioClient   *minio.Client
	AvatarsBucket string
	VideoBucket   string
}

func NewMinio(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool, bucketName string, videoBucket string) (*Minio, error) {
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	return &Minio{MinioClient: minioClient, AvatarsBucket: bucketName, VideoBucket: videoBucket}, err
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

func (mn *Minio) GetVideoRange(ctx context.Context, name string, start, end int64) (io.ReadCloser, error) {
	opts := minio.GetObjectOptions{}
	opts.SetRange(start, end)
	return mn.MinioClient.GetObject(mn.VideoBucket, name, opts)
}

func (mn *Minio) Stat(ctx context.Context, name string) (dto.VideoMeta, error) {
	info, err := mn.MinioClient.StatObject(mn.VideoBucket, name, minio.StatObjectOptions{})
	if err != nil {
		return dto.VideoMeta{}, err
	}
	return dto.VideoMeta{Name: name, Size: info.Size}, nil
}
