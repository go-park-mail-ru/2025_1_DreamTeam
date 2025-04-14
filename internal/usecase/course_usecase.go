package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"skillForce/internal/models"
	"skillForce/internal/models/dto"
	"skillForce/pkg/logs"
	"skillForce/pkg/sanitize"
)

// GetBucketCourses - извлекает список курсов из базы данных
func (uc *Usecase) GetBucketCourses(ctx context.Context) ([]*dto.CourseDTO, error) {
	bucketCourses, err := uc.repo.GetBucketCourses(ctx)
	if err != nil {
		logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	coursesRatings, err := uc.repo.GetCoursesRaitings(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	courseTags, err := uc.repo.GetCoursesTags(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	coursePurchases, err := uc.repo.GetCoursesPurchases(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	resultBucketCourses := make([]*dto.CourseDTO, 0, len(bucketCourses))
	for _, course := range bucketCourses {
		rating, ok := coursesRatings[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("no rating for course %d", course.Id))
			rating = 0
		}

		tags, ok := courseTags[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("no tags for course %d", course.Id))
			tags = []string{}
		}

		purchases, ok := coursePurchases[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("no purchases for course %d", course.Id))
			purchases = 0
		}
		resultBucketCourses = append(resultBucketCourses, &dto.CourseDTO{
			Id:              course.Id,
			CreatorId:       course.CreatorId,
			Title:           course.Title,
			Description:     sanitize.Sanitize(course.Description),
			ScrImage:        course.ScrImage,
			Price:           course.Price,
			TimeToPass:      course.TimeToPass,
			Rating:          rating,
			Tags:            tags,
			PurchasesAmount: purchases,
		})
	}

	logs.PrintLog(ctx, "GetBucketCourses", "get bucket courses with ratings and tags from db, mapping to dto")

	return resultBucketCourses, nil
}

func (uc *Usecase) GetCourse(ctx context.Context, courseId int, userProfile *models.UserProfile) (*dto.CourseDTO, error) {
	course, err := uc.repo.GetCourseById(ctx, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("%+v", err))
		return nil, err
	}
	logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("get course %+v from db", course))

	bucketCourses := []*models.Course{course}

	coursesRatings, err := uc.repo.GetCoursesRaitings(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("%+v", err))
		return nil, err
	}

	courseTags, err := uc.repo.GetCoursesTags(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("%+v", err))
		return nil, err
	}

	coursePurchases, err := uc.repo.GetCoursesPurchases(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("%+v", err))
		return nil, err
	}

	resultBucketCourses := make([]*dto.CourseDTO, 0, len(bucketCourses))
	for _, course := range bucketCourses {
		rating, ok := coursesRatings[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("no rating for course %d", course.Id))
			rating = 0
		}

		tags, ok := courseTags[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("no tags for course %d", course.Id))
			tags = []string{}
		}

		purchases, ok := coursePurchases[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("no purchases for course %d", course.Id))
			purchases = 0
		}
		resultBucketCourses = append(resultBucketCourses, &dto.CourseDTO{
			Id:              course.Id,
			CreatorId:       course.CreatorId,
			Title:           course.Title,
			Description:     sanitize.Sanitize(course.Description),
			ScrImage:        course.ScrImage,
			Price:           course.Price,
			TimeToPass:      course.TimeToPass,
			Rating:          rating,
			Tags:            tags,
			PurchasesAmount: purchases,
		})
	}

	if userProfile != nil {
		resultBucketCourses[0].IsPurchased, err = uc.repo.IsUserPurchasedCourse(ctx, userProfile.Id, course.Id)
		if err != nil {
			logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("can't check if user purchased course: %+v", err))
			return nil, err
		}
	}

	logs.PrintLog(ctx, "GetCourse", "get course with ratings and tags from db, mapping to dto")

	return resultBucketCourses[0], nil

}

func (uc *Usecase) GetCourseLesson(ctx context.Context, userId int, courseId int) (*dto.LessonDTO, error) {
	err := uc.repo.AddUserToCourse(ctx, userId, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
		return nil, err
	}

	lessonHeader, currentLessonId, lessonType, first, err := uc.repo.GetLastLessonHeader(ctx, userId, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
		return nil, err
	}

	if lessonType == "text" { //TODO: сделать switch case
		lessonDto := &dto.LessonDTO{
			LessonHeader: *lessonHeader,
		}

		blocks, err := uc.repo.GetLessonBlocks(ctx, currentLessonId)

		if err != nil {
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
			return nil, err
		}
		var LessonBody dto.LessonDtoBody
		for _, block := range blocks {
			block = sanitize.Sanitize(block)
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("block: %+v", block))
			LessonBody.Blocks = append(LessonBody.Blocks, struct {
				Body string `json:"body"`
			}{
				Body: block,
			})
		}

		footers, err := uc.repo.GetLessonFooters(ctx, currentLessonId)

		if err != nil {
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
			return nil, err
		}

		if len(footers) != 3 {
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("lesson %d has %d footers", currentLessonId, len(footers)))
			return nil, errors.New("lesson has wrong number of footers")
		}

		LessonBody.Footer.NextLessonId = footers[2]
		LessonBody.Footer.CurrentLessonId = footers[1]
		LessonBody.Footer.PreviousLessonId = footers[0]
		lessonDto.LessonBody = LessonBody

		if first {
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("first lesson of the course of the user %+v", userId))
			user, err := uc.repo.GetUserById(ctx, userId)
			if err != nil {
				logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("can't get user by id: %+v", err))
				return nil, err
			}

			if !user.HideEmail {
				err = uc.repo.SendWelcomeCourseMail(ctx, user, courseId)
				if err != nil {
					logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("can't send welcome course mail: %+v", err))
				}
			}
		}

		return lessonDto, err
	}

	return nil, nil
}

func (uc *Usecase) GetNextLesson(ctx context.Context, userId int, courseId int, lessonId int) (*dto.LessonDTO, error) {
	blocks, err := uc.repo.GetLessonBlocks(ctx, lessonId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
		return nil, err
	}

	var LessonBody dto.LessonDtoBody
	for _, block := range blocks {
		block = sanitize.Sanitize(block)
		logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("block: %+v", block))
		LessonBody.Blocks = append(LessonBody.Blocks, struct {
			Body string `json:"body"`
		}{
			Body: block,
		})
	}

	footers, err := uc.repo.GetLessonFooters(ctx, lessonId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
		return nil, err
	}

	if len(footers) != 3 {
		logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("lesson %d has %d footers", lessonId, len(footers)))
		return nil, errors.New("lesson has wrong number of footers")
	}

	LessonBody.Footer.NextLessonId = footers[2]
	LessonBody.Footer.CurrentLessonId = footers[1]
	LessonBody.Footer.PreviousLessonId = footers[0]

	err = uc.repo.MarkLessonCompleted(ctx, userId, courseId, lessonId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
		return nil, err
	}

	lessonHeader, err := uc.repo.GetLessonHeaderByLessonId(ctx, userId, lessonId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
		return nil, err
	}

	lessonDto := &dto.LessonDTO{
		LessonHeader: *lessonHeader,
		LessonBody:   LessonBody,
	}

	return lessonDto, err
}

func (uc *Usecase) MarkLessonAsNotCompleted(ctx context.Context, userId int, lessonId int) error {
	return uc.repo.MarkLessonAsNotCompleted(ctx, userId, lessonId)
}

func (uc *Usecase) GetCourseRoadmap(ctx context.Context, userId int, courseId int) (*dto.CourseRoadmapDTO, error) {
	var roadmap dto.CourseRoadmapDTO

	var parts []*models.CoursePart
	parts, err := uc.repo.GetCourseParts(ctx, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseRoadmap", fmt.Sprintf("%+v", err))
		return nil, err
	}

	for _, part := range parts {
		buckets, err := uc.repo.GetPartBuckets(ctx, part.Id)
		if err != nil {
			logs.PrintLog(ctx, "GetCourseRoadmap", fmt.Sprintf("%+v", err))
			return nil, err
		}
		part.Buckets = buckets

		var bucketsDto []*dto.LessonBucketDTO
		for _, bucket := range buckets {
			lessonPoints, err := uc.repo.GetBucketLessons(ctx, userId, courseId, bucket.Id)
			if err != nil {
				logs.PrintLog(ctx, "GetCourseRoadmap", fmt.Sprintf("%+v", err))
				return nil, err
			}

			var lessonDtoPoints []*dto.LessonPointDTO
			for _, lessonPoint := range lessonPoints {
				var lessonDto dto.LessonPointDTO
				lessonDto.LessonId = lessonPoint.LessonId
				lessonDto.Title = lessonPoint.Title
				lessonDto.IsDone = lessonPoint.IsDone
				lessonDto.Type = lessonPoint.Type

				lessonDtoPoints = append(lessonDtoPoints, &lessonDto)
			}

			var bucketDto dto.LessonBucketDTO
			bucketDto.Id = bucket.Id
			bucketDto.Title = bucket.Title
			bucketDto.Lessons = lessonDtoPoints

			bucketsDto = append(bucketsDto, &bucketDto)
		}

		var partDto dto.CoursePartDTO
		partDto.Id = part.Id
		partDto.Title = part.Title
		partDto.Buckets = bucketsDto

		roadmap.Parts = append(roadmap.Parts, &partDto)
	}

	return &roadmap, nil
}

func (uc *Usecase) GetVideoUrl(ctx context.Context, lesson_id int) (string, error) {
	return uc.repo.GetVideoUrl(ctx, lesson_id)
}

func (uc *Usecase) GetMeta(ctx context.Context, name string) (dto.VideoMeta, error) {
	return uc.repo.Stat(ctx, name)
}

func (uc *Usecase) GetFragment(ctx context.Context, name string, start, end int64) (io.ReadCloser, error) {
	return uc.repo.GetVideoRange(ctx, name, start, end)
}
