package usecase

import (
	"skillForce/backend/models"
	"skillForce/backend/repository"
)

type CourseUsecaseInterface interface {
	GetBucketCourses() ([]*models.Course, error)
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
func (uc *CourseUsecase) GetBucketCourses() ([]*models.Course, error) {
	return uc.repo.GetBucketCourses()
}
