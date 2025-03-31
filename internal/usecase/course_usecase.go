package usecase

import (
	"context"
	"skillForce/internal/models"
	"skillForce/internal/repository"
)

type CourseUsecaseInterface interface {
	GetBucketCourses(ctx context.Context) ([]*models.Course, error)
}

// CourseUsecase - структура бизнес-логики
type CourseUsecase struct {
	repo repository.Repository
}

// NewCourseUsecase - конструктор
func NewCourseUsecase(repo repository.Repository) *CourseUsecase {
	return &CourseUsecase{repo: repo}
}

// GetBucketCourses - извлекает список курсов из базы данных
func (uc *CourseUsecase) GetBucketCourses(ctx context.Context) ([]*models.Course, error) {
	return uc.repo.GetBucketCourses(ctx)
}
