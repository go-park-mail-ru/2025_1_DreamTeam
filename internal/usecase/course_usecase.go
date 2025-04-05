package usecase

import (
	"context"
	"errors"
	"fmt"
	"skillForce/internal/models/dto"
	"skillForce/internal/repository"
	"skillForce/pkg/logs"
)

type CourseUsecaseInterface interface {
	GetBucketCourses(ctx context.Context) ([]*dto.CourseDTO, error)
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
func (uc *CourseUsecase) GetBucketCourses(ctx context.Context) ([]*dto.CourseDTO, error) {
	bucketCoursesWithoutRating, err := uc.repo.GetBucketCourses(ctx)
	if err != nil {
		logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	coursesRatings, err := uc.repo.GetCoursesRaitings(ctx, bucketCoursesWithoutRating)
	if err != nil {
		logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	bucketCourses := make([]*dto.CourseDTO, 0, len(bucketCoursesWithoutRating))
	for _, course := range bucketCoursesWithoutRating {
		rating, ok := coursesRatings[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("no rating for course %d", course.Id))
			return nil, errors.New("not enough ratings for courses")
		}
		bucketCourses = append(bucketCourses, &dto.CourseDTO{
			Id:          course.Id,
			CreatorId:   course.CreatorId,
			Title:       course.Title,
			Description: course.Description,
			ScrImage:    course.ScrImage,
			Price:       course.Price,
			TimeToPass:  course.TimeToPass,
			Rating:      rating.Rating,
		})
	}

	logs.PrintLog(ctx, "GetBucketCourses", "get bucket courses and ratings from db, mapping to dto")

	return bucketCourses, nil
}
