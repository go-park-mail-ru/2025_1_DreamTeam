package usecase

import (
	"context"
	"io"
	"skillForce/internal/models/dto"
)

type CourseRepository interface {
	GetVideoUrl(ctx context.Context, lessonId int) (string, error)
	GetVideoRange(ctx context.Context, name string, start, end int64) (io.ReadCloser, error)
	Stat(ctx context.Context, name string) (dto.VideoMeta, error)
}
