package usecase

import (
	"skillForce/internal/models"
	"skillForce/internal/repository"
)

// CourseUsecase - структура бизнес-логики
type CourseUsecase struct {
	repo *repository.Database
}

// NewCourseUsecase - конструктор
func NewCourseUsecase(repo *repository.Database) *CourseUsecase {
	return &CourseUsecase{repo: repo}
}

// GetBucketCourses - извлекает список курсов из базы данных
func (uc *CourseUsecase) GetBucketCourses() ([]*models.Course, error) {
	return uc.repo.GetBucketCourses()
}
