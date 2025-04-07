package usecase

import (
	"context"
	"errors"
	"fmt"
	"skillForce/internal/models"
	"skillForce/internal/models/dto"
	"skillForce/pkg/logs"
)

// GetBucketCourses - извлекает список курсов из базы данных
func (uc *Usecase) GetBucketCourses(ctx context.Context) ([]*dto.CourseDTO, error) {
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

	courseTags, err := uc.repo.GetCoursesTags(ctx, bucketCoursesWithoutRating)
	if err != nil {
		logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	bucketCourses := make([]*dto.CourseDTO, 0, len(bucketCoursesWithoutRating))
	for _, course := range bucketCoursesWithoutRating {
		rating, ok := coursesRatings[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("no rating for course %d", course.Id))
			rating = models.CourseRating{
				Rating: 0,
			}
		}

		tags, ok := courseTags[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("no tags for course %d", course.Id))
			tags = []string{}
		}
		bucketCourses = append(bucketCourses, &dto.CourseDTO{
			Id:              course.Id,
			CreatorId:       course.CreatorId,
			Title:           course.Title,
			Description:     course.Description,
			ScrImage:        course.ScrImage,
			Price:           course.Price,
			TimeToPass:      course.TimeToPass,
			Rating:          rating.Rating,
			Tags:            tags,
			PurchasesAmount: course.PurchasesAmount,
		})
	}

	logs.PrintLog(ctx, "GetBucketCourses", "get bucket courses with ratings and tags from db, mapping to dto")

	return bucketCourses, nil
}

func (uc *Usecase) GetCourseLesson(ctx context.Context, userId int, courseId int) (*dto.LessonDTO, error) {
	var lessonHeader dto.LessonDtoHeader

	course, err := uc.repo.GetCourseById(ctx, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
		return nil, err
	}
	lessonHeader.CourseTitle = course.Title

	currentLessonId, currentBucketId, lessonType, err := uc.repo.FillLessonHeader(ctx, userId, courseId, &lessonHeader)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
		return nil, err
	}

	if lessonType == "text" { //TODO: сделать switch case
		lessonDto := &dto.LessonDTO{
			LessonHeader: lessonHeader,
		}

		blocks, err := uc.repo.GetLessonBlocks(ctx, currentLessonId)

		if err != nil {
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
			return nil, err
		}

		for _, block := range blocks {
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("block: %+v", block))
			lessonDto.Blocks = append(lessonDto.Blocks, struct {
				Body string `json:"body"`
			}{
				Body: block,
			})
		}

		footers, err := uc.repo.GetLessonFooters(ctx, currentLessonId, currentBucketId)

		if err != nil {
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
			return nil, err
		}

		if len(footers) != 2 {
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("lesson %d has %d footers", currentLessonId, len(footers)))
			return nil, errors.New("lesson has wrong number of footers")
		}

		lessonDto.Footer.NextLessonId = footers[1]
		lessonDto.Footer.PreviousLessonId = footers[0]

		return lessonDto, err
	}

	return nil, nil
}
