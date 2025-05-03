package usecase

import (
	"context"
	"io"
	"skillForce/internal/models/dto"
)

type CourseUsecase struct {
	repo CourseRepository
}

func NewCourseUsecase(repo CourseRepository) *CourseUsecase {
	return &CourseUsecase{
		repo: repo,
	}
}

func (uc *CourseUsecase) GetVideoUrl(ctx context.Context, lesson_id int) (string, error) {
	return uc.repo.GetVideoUrl(ctx, lesson_id)
}

func (uc *CourseUsecase) GetMeta(ctx context.Context, name string) (dto.VideoMeta, error) {
	return uc.repo.Stat(ctx, name)
}

func (uc *CourseUsecase) GetFragment(ctx context.Context, name string, start, end int64) (io.ReadCloser, error) {
	return uc.repo.GetVideoRange(ctx, name, start, end)
}
