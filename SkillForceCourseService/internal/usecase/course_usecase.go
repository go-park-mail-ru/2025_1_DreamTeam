package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"os"
	coursemodels "skillForce/internal/models/course"
	"skillForce/internal/models/dto"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/logs"
	"skillForce/pkg/sanitize"
	"skillForce/pkg/sertificate"
	"time"
)

type CourseUsecase struct {
	repo CourseRepository
}

func NewCourseUsecase(repo CourseRepository) *CourseUsecase {
	return &CourseUsecase{
		repo: repo,
	}
}

func (uc *CourseUsecase) GetBucketCourses(ctx context.Context, userProfile *usermodels.UserProfile) ([]*dto.CourseDTO, error) {
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

	courseFavourites := make(map[int]bool, 0)
	if userProfile != nil {
		courseFavourites, err = uc.repo.GetCoursesFavouriteStatus(ctx, bucketCourses, userProfile.Id)
		if err != nil {
			logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
			return nil, err
		}
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
			IsFavorite:      courseFavourites[course.Id],
		})
	}

	logs.PrintLog(ctx, "GetBucketCourses", "get bucket courses with ratings and tags from db, mapping to dto")

	return resultBucketCourses, nil
}

func (uc *CourseUsecase) GetPurchasedBucketCourses(ctx context.Context, userProfile *usermodels.UserProfile) ([]*dto.CourseDTO, error) {
	bucketCourses, err := uc.repo.GetPurchasedBucketCourses(ctx, userProfile.Id)
	if err != nil {
		logs.PrintLog(ctx, "GetPurchasedBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	coursesRatings, err := uc.repo.GetCoursesRaitings(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetPurchasedBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	courseTags, err := uc.repo.GetCoursesTags(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetPurchasedBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	coursePurchases, err := uc.repo.GetCoursesPurchases(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetPurchasedBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	courseFavourites := make(map[int]bool, 0)
	if userProfile != nil {
		courseFavourites, err = uc.repo.GetCoursesFavouriteStatus(ctx, bucketCourses, userProfile.Id)
		if err != nil {
			logs.PrintLog(ctx, "GetPurchasedBucketCourses", fmt.Sprintf("%+v", err))
			return nil, err
		}
	}

	resultBucketCourses := make([]*dto.CourseDTO, 0, len(bucketCourses))
	for _, course := range bucketCourses {
		rating, ok := coursesRatings[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetPurchasedBucketCourses", fmt.Sprintf("no rating for course %d", course.Id))
			rating = 0
		}

		tags, ok := courseTags[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetPurchasedBucketCourses", fmt.Sprintf("no tags for course %d", course.Id))
			tags = []string{}
		}

		purchases, ok := coursePurchases[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetPurchasedBucketCourses", fmt.Sprintf("no purchases for course %d", course.Id))
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
			IsFavorite:      courseFavourites[course.Id],
		})
	}

	logs.PrintLog(ctx, "GetPurchasedBucketCourses", "get bucket courses with ratings and tags from db, mapping to dto")

	return resultBucketCourses, nil
}

func (uc *CourseUsecase) GetCompletedBucketCourses(ctx context.Context, userProfile *usermodels.UserProfile) ([]*dto.CourseDTO, error) {
	bucketCourses, err := uc.repo.GetCompletedBucketCourses(ctx, userProfile.Id)
	if err != nil {
		logs.PrintLog(ctx, "GetCompletedBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	coursesRatings, err := uc.repo.GetCoursesRaitings(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetCompletedBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	courseTags, err := uc.repo.GetCoursesTags(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetCompletedBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	coursePurchases, err := uc.repo.GetCoursesPurchases(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetCompletedBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	courseFavourites := make(map[int]bool, 0)
	if userProfile != nil {
		courseFavourites, err = uc.repo.GetCoursesFavouriteStatus(ctx, bucketCourses, userProfile.Id)
		if err != nil {
			logs.PrintLog(ctx, "GetCompletedBucketCourses", fmt.Sprintf("%+v", err))
			return nil, err
		}
	}

	resultBucketCourses := make([]*dto.CourseDTO, 0, len(bucketCourses))
	for _, course := range bucketCourses {
		rating, ok := coursesRatings[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetCompletedBucketCourses", fmt.Sprintf("no rating for course %d", course.Id))
			rating = 0
		}

		tags, ok := courseTags[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetCompletedBucketCourses", fmt.Sprintf("no tags for course %d", course.Id))
			tags = []string{}
		}

		purchases, ok := coursePurchases[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetCompletedBucketCourses", fmt.Sprintf("no purchases for course %d", course.Id))
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
			IsFavorite:      courseFavourites[course.Id],
		})
	}

	logs.PrintLog(ctx, "GetCompletedBucketCourses", "get bucket courses with ratings and tags from db, mapping to dto")

	return resultBucketCourses, nil
}

func (uc *CourseUsecase) GetCourse(ctx context.Context, courseId int, userProfile *usermodels.UserProfile) (*dto.CourseDTO, error) {
	course, err := uc.repo.GetCourseById(ctx, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("%+v", err))
		return nil, err
	}
	logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("get course %+v from db", course))

	bucketCourses := []*coursemodels.Course{course}

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

		resultBucketCourses[0].IsCompleted, err = uc.repo.IsUserCompletedCourse(ctx, userProfile.Id, course.Id)
		if err != nil {
			logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("can't check if user completed course: %+v", err))
			return nil, err
		}

		courseFavourits, err := uc.repo.GetCoursesFavouriteStatus(ctx, bucketCourses, userProfile.Id)
		if err != nil {
			logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("%+v", err))
			return nil, err
		}
		isFavourite, ok := courseFavourits[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetCourse", fmt.Sprintf("no favourite status for course %d", course.Id))
			isFavourite = false
		}
		resultBucketCourses[0].IsFavorite = isFavourite
	}

	logs.PrintLog(ctx, "GetCourse", "get course with ratings and tags from db, mapping to dto")

	return resultBucketCourses[0], nil

}

func (uc *CourseUsecase) GetCourseLesson(ctx context.Context, userId int, courseId int) (*dto.LessonDTO, error) {
	err := uc.repo.AddUserToCourse(ctx, userId, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
		return nil, err
	}

	lessonHeader, currentLessonId, lessonType, _, err := uc.repo.GetLastLessonHeader(ctx, userId, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
		return nil, err
	}

	lessonDto := &dto.LessonDTO{
		LessonHeader: *lessonHeader,
	}

	if lessonType == "text" {
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

		/*
			first, err := uc.repo.IsWelcomeCourseMailSended(ctx, userId, courseId)
			if first {
				logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("first lesson of the course of the user %+v", userId))
				user, err := uc.repo.GetUserById(ctx, userId)
				if err != nil {
					logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("can't get user by id: %+v", err))
					return nil, err
				}
				course, err := uc.repo.GetCourseById(ctx, courseId)
				if err != nil {
					logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("can't get course by id: %+v", err))
					return nil, err
				}
				if !user.HideEmail {
					go uc.repo.SendWelcomeCourseMail(ctx, user, course)
				}
			}
		*/

		return lessonDto, err
	}

	if lessonType == "video" {
		blocks, err := uc.repo.GetLessonVideo(ctx, currentLessonId)

		if err != nil {
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
			return nil, err
		}
		var LessonBody dto.LessonDtoBody
		for _, block := range blocks {
			//block = sanitize.Sanitize(block)
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

		/*
			first, err := uc.repo.IsWelcomeCourseMailSended(ctx, userId, courseId)
			if first {
				logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("first lesson of the course of the user %+v", userId))
				user, err := uc.repo.GetUserById(ctx, userId)
				if err != nil {
					logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("can't get user by id: %+v", err))
					return nil, err
				}
				course, err := uc.repo.GetCourseById(ctx, courseId)
				if err != nil {
					logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("can't get course by id: %+v", err))
					return nil, err
				}
				if !user.HideEmail {
					go uc.repo.SendWelcomeCourseMail(ctx, user, course)
				}
			}
		*/

		return lessonDto, err

	}

	if lessonType == "quiz" {
		var LessonBody dto.LessonDtoBody
		LessonBody.Blocks = append(LessonBody.Blocks, struct {
			Body string `json:"body"`
		}{
			Body: "quiz",
		})

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

		return lessonDto, err
	}

	if lessonType == "question" {
		var LessonBody dto.LessonDtoBody
		LessonBody.Blocks = append(LessonBody.Blocks, struct {
			Body string `json:"body"`
		}{
			Body: "question",
		})

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

		return lessonDto, err
	}

	return nil, nil
}

func (uc *CourseUsecase) GetNextLesson(ctx context.Context, userId int, courseId int, lessonId int) (*dto.LessonDTO, error) {
	lesson, err := uc.repo.GetLessonById(ctx, lessonId)
	if err != nil {
		logs.PrintLog(ctx, "GetNextLesson", fmt.Sprintf("%+v", err))
		return nil, err
	}

	if lesson.Type == "text" {
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

		lessonHeader, err := uc.repo.GetLessonHeaderByLessonId(ctx, userId, lessonId)
		if err != nil {
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
			return nil, err
		}

		lessonDto := &dto.LessonDTO{
			LessonHeader: *lessonHeader,
			LessonBody:   LessonBody,
		}
		/*
			user, _ := uc.repo.GetUserById(ctx, userId)
			if user.HideEmail {
				isMiddle, _ := uc.repo.IsMiddle(ctx, userId, courseId)
				user, _ := uc.repo.GetUserById(ctx, userId)
				if isMiddle {
					go uc.repo.SendMiddleCourseMail(ctx, user, courseId)
				}
			}
		*/
		return lessonDto, err
	}

	if lesson.Type == "video" {
		blocks, err := uc.repo.GetLessonVideo(ctx, lessonId)
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

		lessonHeader, err := uc.repo.GetLessonHeaderByLessonId(ctx, userId, lessonId)
		if err != nil {
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
			return nil, err
		}

		lessonDto := &dto.LessonDTO{
			LessonHeader: *lessonHeader,
			LessonBody:   LessonBody,
		}
		/*
			user, _ := uc.repo.GetUserById(ctx, userId)
			if user.HideEmail {
				isMiddle, _ := uc.repo.IsMiddle(ctx, userId, courseId)
				user, _ := uc.repo.GetUserById(ctx, userId)
				if isMiddle {
					go uc.repo.SendMiddleCourseMail(ctx, user, courseId)
				}
			}
		*/
		return lessonDto, err
	}

	if lesson.Type == "question" {

		var LessonBody dto.LessonDtoBody
		LessonBody.Blocks = append(LessonBody.Blocks, struct {
			Body string `json:"body"`
		}{
			Body: "question",
		})

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

		// err = uc.repo.MarkLessonCompleted(ctx, userId, courseId, lessonId)
		// if err != nil {
		// 	logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
		// 	return nil, err
		// }

		lessonHeader, err := uc.repo.GetLessonHeaderByLessonId(ctx, userId, lessonId)
		if err != nil {
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
			return nil, err
		}

		lessonDto := &dto.LessonDTO{
			LessonHeader: *lessonHeader,
			LessonBody:   LessonBody,
		}
		/*
			user, _ := uc.repo.GetUserById(ctx, userId)
			if user.HideEmail {
				isMiddle, _ := uc.repo.IsMiddle(ctx, userId, courseId)
				user, _ := uc.repo.GetUserById(ctx, userId)
				if isMiddle {
					go uc.repo.SendMiddleCourseMail(ctx, user, courseId)
				}
			}
		*/
		return lessonDto, err
	}

	if lesson.Type == "quiz" {

		var LessonBody dto.LessonDtoBody
		LessonBody.Blocks = append(LessonBody.Blocks, struct {
			Body string `json:"body"`
		}{
			Body: "quiz",
		})

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

		// err = uc.repo.MarkLessonCompleted(ctx, userId, courseId, lessonId)
		// if err != nil {
		// 	logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
		// 	return nil, err
		// }

		lessonHeader, err := uc.repo.GetLessonHeaderByLessonId(ctx, userId, lessonId)
		if err != nil {
			logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
			return nil, err
		}

		lessonDto := &dto.LessonDTO{
			LessonHeader: *lessonHeader,
			LessonBody:   LessonBody,
		}
		/*
			user, _ := uc.repo.GetUserById(ctx, userId)
			if user.HideEmail {
				isMiddle, _ := uc.repo.IsMiddle(ctx, userId, courseId)
				user, _ := uc.repo.GetUserById(ctx, userId)
				if isMiddle {
					go uc.repo.SendMiddleCourseMail(ctx, user, courseId)
				}
			}
		*/
		return lessonDto, err
	}

	return nil, errors.New("next lesson has wrong type")
}

func (uc *CourseUsecase) MarkLessonAsNotCompleted(ctx context.Context, userId int, lessonId int) error {
	return uc.repo.MarkLessonAsNotCompleted(ctx, userId, lessonId)
}

func (uc *CourseUsecase) MarkCourseAsCompleted(ctx context.Context, userId int, courseId int) error {
	return uc.repo.MarkCourseAsCompleted(ctx, userId, courseId)
}

func (uc *CourseUsecase) MarkLessonAsCompleted(ctx context.Context, userId int, lessonId int) error {
	return uc.repo.MarkLessonCompleted(ctx, userId, lessonId)
}

func (uc *CourseUsecase) GetCourseRoadmap(ctx context.Context, userId int, courseId int) (*dto.CourseRoadmapDTO, error) {
	var roadmap dto.CourseRoadmapDTO

	var parts []*coursemodels.CoursePart
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

func (uc *CourseUsecase) GetRating(ctx context.Context, userId int, courseId int) (*dto.Raiting, error) {
	return uc.repo.GetRating(ctx, userId, courseId)
}

func (uc *CourseUsecase) GetStatistic(ctx context.Context, userId int, courseId int) (*dto.UserStats, error) {
	return uc.repo.GetStatistic(ctx, userId, courseId)
}

func (uc *CourseUsecase) GetSertificate(ctx context.Context, userProfile *usermodels.UserProfile, courseId int) (string, error) {
	exists, err := uc.repo.IsSertificateExists(ctx, userProfile.Id, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetSertificate", fmt.Sprintf("%+v", err))
		return "", err
	}

	if exists {
		return "", errors.New("certificate already exists")
	}

	course, err := uc.repo.GetCourseById(ctx, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetSertificate", fmt.Sprintf("%+v", err))
		return "", err
	}
	date := time.Now().Format("02.01.2006")

	stats, err := uc.repo.GetStatistic(ctx, userProfile.Id, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetSertificate", fmt.Sprintf("%+v", err))
		return "", err
	}

	tempFileName := fmt.Sprintf("certificate_%v_%v.pdf", userProfile.Name, course.Id)

	if (stats.RecievedPoints)*100/stats.AmountPoints >= 85 {
		err = sertificate.GenerateGoodCertificate(userProfile.Name, course.Title, date, tempFileName)
	} else {
		err = sertificate.GenerateNormalCertificate(userProfile.Name, course.Title, date, tempFileName)
	}

	if err != nil {
		logs.PrintLog(ctx, "GetSertificate", fmt.Sprintf("can't generate certificate: %+v", err))
		return "", err
	}
	defer func() {
		if err := os.Remove(tempFileName); err != nil {
			logs.PrintLog(ctx, "GetSertificate", fmt.Sprintf("%+v", err))
		}
	}()

	logs.PrintLog(ctx, "GetSertificate", fmt.Sprintf("certificate for user (id: %v) and course (id: %v) was generated", userProfile.Id, course.Id))

	// Открываем файл для чтения
	file, err := os.Open(tempFileName)
	if err != nil {
		logs.PrintLog(ctx, "GetSertificate", fmt.Sprintf("failed to open certificate file: %+v", err))
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logs.PrintLog(ctx, "GetSertificate", fmt.Sprintf("%+v", err))
		}
	}()

	// Получаем информацию о файле
	fileInfo, err := file.Stat()
	if err != nil {
		logs.PrintLog(ctx, "GetSertificate", fmt.Sprintf("failed to get file info: %+v", err))
		return "", err
	}

	fileHeader := &multipart.FileHeader{
		Filename: tempFileName,
		Size:     fileInfo.Size(),
		Header:   make(textproto.MIMEHeader),
	}
	fileHeader.Header.Set("Content-Type", "application/pdf")

	// Загружаем файл в MinIO
	url, err := uc.repo.UploadFileToMinIO(ctx, file, fileHeader)
	if err != nil {
		logs.PrintLog(ctx, "GetSertificate", fmt.Sprintf("failed to upload certificate: %+v", err))
		return "", err
	}

	err = uc.repo.SaveSertificate(ctx, userProfile.Id, course.Id, url)
	if err != nil {
		logs.PrintLog(ctx, "GetSertificate", fmt.Sprintf("failed to save sertificate: %+v", err))
		return "", err
	}

	return url, nil
}

func (uc *CourseUsecase) GetGeneratedSertificate(ctx context.Context, userProfile *usermodels.UserProfile, courseId int) (string, error) {
	return uc.repo.GetGeneratedSertificate(ctx, userProfile, courseId)
}

func (uc *CourseUsecase) CreateCourse(ctx context.Context, courseDto *dto.CourseDTO, userProfile *usermodels.UserProfile) error {
	course := coursemodels.Course{
		CreatorId:   userProfile.Id,
		Description: courseDto.Description,
		Title:       courseDto.Title,
		Price:       courseDto.Price,
		TimeToPass:  courseDto.TimeToPass,
	}
	for _, part := range courseDto.Parts {
		coursePart := coursemodels.CoursePart{
			Title: part.Title,
		}
		course.Parts = append(course.Parts, &coursePart)
		for _, bucket := range part.Buckets {
			courseBucket := coursemodels.LessonBucket{
				Title: bucket.Title,
			}
			coursePart.Buckets = append(coursePart.Buckets, &courseBucket)
			for _, lesson := range bucket.Lessons {
				courseLesson := coursemodels.LessonPoint{
					Value: lesson.Value,
					Type:  lesson.Type,
					Title: lesson.Title,
				}
				courseBucket.Lessons = append(courseBucket.Lessons, &courseLesson)
			}
		}
	}
	courseId, err := uc.repo.CreateCourse(ctx, &course, userProfile)
	if err != nil {
		logs.PrintLog(ctx, "CreateCourse", fmt.Sprintf("%+v", err))
		return err
	}
	course.Id = courseId
	for partOrder, part := range course.Parts {
		part.Order = partOrder + 1
		partId, err := uc.repo.CreatePart(ctx, part, course.Id)
		if err != nil {
			logs.PrintLog(ctx, "CreateCourse", fmt.Sprintf("%+v", err))
			return err
		}
		for bucketOrder, bucket := range part.Buckets {
			bucket.Order = bucketOrder + 1
			bucket.PartId = partId
			bucketId, err := uc.repo.CreateBucket(ctx, bucket, partId)
			if err != nil {
				logs.PrintLog(ctx, "CreateCourse", fmt.Sprintf("%+v", err))
				return err
			}
			for lessonOrder, lesson := range bucket.Lessons {
				lesson.Order = lessonOrder + 1
				lesson.BucketId = bucketId

				switch lesson.Type {
				case "video":
					err = uc.repo.CreateVideoLesson(ctx, lesson, bucketId)
					if err != nil {
						logs.PrintLog(ctx, "CreateCourse", fmt.Sprintf("%+v", err))
						return err
					}
				case "text":
					err = uc.repo.CreateTextLesson(ctx, lesson, bucketId)
					if err != nil {
						logs.PrintLog(ctx, "CreateCourse", fmt.Sprintf("%+v", err))
						return err
					}
				}
			}
		}
	}

	return nil
}

func (uc *CourseUsecase) AddCourseToFavourites(ctx context.Context, course *dto.CourseDTO, userProfile *usermodels.UserProfile) error {
	logs.PrintLog(ctx, "AddCourseToFavourites", fmt.Sprintf("add course with id: %v to favourites of user with id: %v", course.Id, userProfile.Id))
	return uc.repo.AddCourseToFavourites(ctx, course.Id, userProfile.Id)
}

func (uc *CourseUsecase) DeleteCourseFromFavourites(ctx context.Context, course *dto.CourseDTO, userProfile *usermodels.UserProfile) error {
	logs.PrintLog(ctx, "DeleteCourseFromFavourites", fmt.Sprintf("delete course with id: %v from favourites of user with id: %v", course.Id, userProfile.Id))
	return uc.repo.DeleteCourseFromFavourites(ctx, course.Id, userProfile.Id)
}

func (uc *CourseUsecase) GetFavouriteCourses(ctx context.Context, userProfile *usermodels.UserProfile) ([]*dto.CourseDTO, error) {
	logs.PrintLog(ctx, "GetFavouriteCourses", fmt.Sprintf("get favourite courses of user with id: %v", userProfile.Id))
	bucketCourses, err := uc.repo.GetFavouriteCourses(ctx, userProfile.Id)
	if err != nil {
		logs.PrintLog(ctx, "GetFavouriteCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	coursesRatings, err := uc.repo.GetCoursesRaitings(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetFavouriteCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	courseTags, err := uc.repo.GetCoursesTags(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetFavouriteCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	coursePurchases, err := uc.repo.GetCoursesPurchases(ctx, bucketCourses)
	if err != nil {
		logs.PrintLog(ctx, "GetFavouriteCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}

	resultBucketCourses := make([]*dto.CourseDTO, 0, len(bucketCourses))
	for _, course := range bucketCourses {
		rating, ok := coursesRatings[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetFavouriteCourses", fmt.Sprintf("no rating for course %d", course.Id))
			rating = 0
		}

		tags, ok := courseTags[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetFavouriteCourses", fmt.Sprintf("no tags for course %d", course.Id))
			tags = []string{}
		}

		purchases, ok := coursePurchases[course.Id]
		if !ok {
			logs.PrintLog(ctx, "GetFavouriteCourses", fmt.Sprintf("no purchases for course %d", course.Id))
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
			IsFavorite:      true,
		})
	}

	logs.PrintLog(ctx, "GetFavouriteCourses", "get bucket courses with ratings and tags from db, mapping to dto")

	return resultBucketCourses, nil
}

func (uc *CourseUsecase) GetTestLesson(ctx context.Context, lesson_id int, user_id int) (*dto.Test, error) {
	logs.PrintLog(ctx, "GetTestLesson", fmt.Sprintf("get test lesson owith id: %v", lesson_id))

	test, err := uc.repo.GetLessonTest(ctx, lesson_id, user_id)

	if err != nil {
		logs.PrintLog(ctx, "GetTestLesson", fmt.Sprintf("%+v", err))
		return nil, err
	}

	return test, nil
}

func (uc *CourseUsecase) AnswerQuiz(ctx context.Context, question_id int, answer_id int, user_id int, course_id int) (*dto.QuizResult, error) {
	logs.PrintLog(ctx, "AnswerQuiz", fmt.Sprintf("get quiz lesson (%v) result", question_id))

	result, err := uc.repo.AnswerQuiz(ctx, question_id, answer_id, user_id, course_id)

	if err != nil {
		logs.PrintLog(ctx, "AnswerQuiz", fmt.Sprintf("%+v", err))
		return nil, err
	}

	return result, nil
}

func (uc *CourseUsecase) GetQuestionTestLesson(ctx context.Context, lesson_id int, user_id int) (*dto.QuestionTest, error) {
	test, err := uc.repo.GetQuestionTestLesson(ctx, lesson_id, user_id)

	if err != nil {
		logs.PrintLog(ctx, "GetQuestionTestLesson", fmt.Sprintf("%+v", err))
		return nil, err
	}

	return test, nil
}

func (uc *CourseUsecase) AnswerQuestion(ctx context.Context, question_id int, user_id int, answer string) error {
	logs.PrintLog(ctx, "AnswerQuestion", fmt.Sprintf("set question lesson (%v) result", question_id))

	err := uc.repo.AnswerQuestion(ctx, question_id, user_id, answer)

	if err != nil {
		logs.PrintLog(ctx, "AnswerQuestion", fmt.Sprintf("%+v", err))
		return err
	}

	return nil
}

func (uc *CourseUsecase) SearchCoursesByTitle(ctx context.Context, userProfile *usermodels.UserProfile, keywords string) ([]*dto.CourseDTO, error) {
	bucketCourses, err := uc.repo.SearchCoursesByTitle(ctx, keywords)
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

	courseFavourites := make(map[int]bool, 0)
	if userProfile != nil {
		courseFavourites, err = uc.repo.GetCoursesFavouriteStatus(ctx, bucketCourses, userProfile.Id)
		if err != nil {
			logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
			return nil, err
		}
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
			IsFavorite:      courseFavourites[course.Id],
		})
	}

	logs.PrintLog(ctx, "GetBucketCourses", "get bucket courses with ratings and tags from db, mapping to dto")

	return resultBucketCourses, nil
}

func (uc *CourseUsecase) AddRating(ctx context.Context, course_id int, user_id int, rating int) error {
	logs.PrintLog(ctx, "AddRating", fmt.Sprintf("set rating for course (%v) by user (%v)", course_id, user_id))

	err := uc.repo.AddRaiting(ctx, user_id, course_id, rating)

	if err != nil {
		logs.PrintLog(ctx, "AddRating", fmt.Sprintf("%+v", err))
		return err
	}

	return nil
}
