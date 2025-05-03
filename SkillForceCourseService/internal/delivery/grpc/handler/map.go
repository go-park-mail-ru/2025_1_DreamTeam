package grpc

import (
	coursepb "skillForce/internal/delivery/grpc/proto"
	"skillForce/internal/models/dto"
	usermodels "skillForce/internal/models/user"
)

func mapToCourseDTO(dto *dto.CourseDTO) *coursepb.CourseDTO {
	if dto == nil {
		return nil
	}
	return &coursepb.CourseDTO{
		Id:              int32(dto.Id),
		Price:           int32(dto.Price),
		PurchasesAmount: int32(dto.PurchasesAmount),
		CreatorId:       int32(dto.CreatorId),
		TimeToPass:      int32(dto.TimeToPass),
		Rating:          dto.Rating,
		Tags:            dto.Tags,
		Title:           dto.Title,
		Description:     dto.Description,
		ScrImage:        dto.ScrImage,
		IsPurchased:     dto.IsPurchased,
		IsFavorite:      dto.IsFavorite,
		Parts:           mapCourseParts(dto.Parts),
	}
}

func mapCourseParts(parts []*dto.CoursePartDTO) []*coursepb.CoursePartDTO {
	res := make([]*coursepb.CoursePartDTO, 0, len(parts))
	for _, part := range parts {
		res = append(res, &coursepb.CoursePartDTO{
			Id:      int32(part.Id),
			Title:   part.Title,
			Buckets: mapLessonBuckets(part.Buckets),
		})
	}
	return res
}

func mapLessonBuckets(buckets []*dto.LessonBucketDTO) []*coursepb.LessonBucketDTO {
	res := make([]*coursepb.LessonBucketDTO, 0, len(buckets))
	for _, b := range buckets {
		res = append(res, &coursepb.LessonBucketDTO{
			Id:      int32(b.Id),
			Title:   b.Title,
			Lessons: mapLessonPoints(b.Lessons),
		})
	}
	return res
}

func mapLessonPoints(points []*dto.LessonPointDTO) []*coursepb.LessonPointDTO {
	res := make([]*coursepb.LessonPointDTO, 0, len(points))
	for _, p := range points {
		res = append(res, &coursepb.LessonPointDTO{
			LessonId: int32(p.LessonId),
			Type:     p.Type,
			Title:    p.Title,
			Value:    p.Value,
			IsDone:   p.IsDone,
		})
	}
	return res
}

func mapToLessonDTO(dto *dto.LessonDTO) *coursepb.LessonDTO {
	if dto == nil {
		return nil
	}
	return &coursepb.LessonDTO{
		Header: &coursepb.LessonDtoHeader{
			CourseTitle: dto.LessonHeader.CourseTitle,
			CourseId:    int32(dto.LessonHeader.CourseId),
			Part: &coursepb.Part{
				Order: int32(dto.LessonHeader.Part.Order),
				Title: dto.LessonHeader.Part.Title,
			},
			Bucket: &coursepb.Bucket{
				Order: int32(dto.LessonHeader.Bucket.Order),
				Title: dto.LessonHeader.Bucket.Title,
			},
			Points: mapLessonHeaderPoints([]struct {
				LessonId int
				Type     string
				IsDone   bool
			}(dto.LessonHeader.Points)),
		},
		Body: &coursepb.LessonDtoBody{
			Blocks: mapBlocks([]struct{ Body string }(dto.LessonBody.Blocks)),
			Footer: &coursepb.Footer{
				NextLessonId:     int32(dto.LessonBody.Footer.NextLessonId),
				CurrentLessonId:  int32(dto.LessonBody.Footer.CurrentLessonId),
				PreviousLessonId: int32(dto.LessonBody.Footer.PreviousLessonId),
			},
		},
	}
}

func mapLessonHeaderPoints(points []struct {
	LessonId int
	Type     string
	IsDone   bool
}) []*coursepb.Point {
	res := make([]*coursepb.Point, 0, len(points))
	for _, p := range points {
		res = append(res, &coursepb.Point{
			LessonId: int32(p.LessonId),
			Type:     p.Type,
			IsDone:   p.IsDone,
		})
	}
	return res
}

func mapBlocks(blocks []struct{ Body string }) []*coursepb.Block {
	res := make([]*coursepb.Block, 0, len(blocks))
	for _, b := range blocks {
		res = append(res, &coursepb.Block{Body: b.Body})
	}
	return res
}

func mapToCourseRoadmapResponse(dto *dto.CourseRoadmapDTO) *coursepb.GetCourseRoadmapResponse {
	return &coursepb.GetCourseRoadmapResponse{
		Roadmap: &coursepb.CourseRoadmapDTO{
			Parts: mapCourseParts(dto.Parts),
		},
	}
}

func mapToGetBucketCoursesResponse(courses []*dto.CourseDTO) *coursepb.GetBucketCoursesResponse {
	res := make([]*coursepb.CourseDTO, 0, len(courses))
	for _, c := range courses {
		res = append(res, mapToCourseDTO(c))
	}
	return &coursepb.GetBucketCoursesResponse{Courses: res}
}

func mapToCourseLessonResponse(lesson *dto.LessonDTO) *coursepb.GetCourseLessonResponse {
	return &coursepb.GetCourseLessonResponse{
		Lesson: mapToLessonDTO(lesson),
	}
}

func mapToNextLessonResponse(lesson *dto.LessonDTO) *coursepb.GetNextLessonResponse {
	return &coursepb.GetNextLessonResponse{
		Lesson: mapToLessonDTO(lesson),
	}
}

func mapToCourseResponse(course *dto.CourseDTO) *coursepb.GetCourseResponse {
	return &coursepb.GetCourseResponse{
		Course: mapToCourseDTO(course),
	}
}

func mapToGetFavouritesResponse(courses []*dto.CourseDTO) *coursepb.GetFavouritesResponse {
	res := make([]*coursepb.CourseDTO, 0, len(courses))
	for _, c := range courses {
		res = append(res, mapToCourseDTO(c))
	}
	return &coursepb.GetFavouritesResponse{Courses: res}
}

func mapPbCourseDTOToDTO(coursepb *coursepb.CourseDTO) *dto.CourseDTO {
	var parts []*dto.CoursePartDTO
	for _, part := range coursepb.Parts {
		var buckets []*dto.LessonBucketDTO
		for _, bucket := range part.Buckets {
			var lessons []*dto.LessonPointDTO
			for _, lesson := range bucket.Lessons {
				lessons = append(lessons, &dto.LessonPointDTO{
					LessonId: int(lesson.LessonId),
					Type:     lesson.Type,
					Title:    lesson.Title,
					Value:    lesson.Value,
					IsDone:   lesson.IsDone,
				})
			}
			buckets = append(buckets, &dto.LessonBucketDTO{
				Id:      int(bucket.Id),
				Title:   bucket.Title,
				Lessons: lessons,
			})
		}
		parts = append(parts, &dto.CoursePartDTO{
			Id:      int(part.Id),
			Title:   part.Title,
			Buckets: buckets,
		})
	}

	return &dto.CourseDTO{
		Id:              int(coursepb.Id),
		Price:           int(coursepb.Price),
		PurchasesAmount: int(coursepb.PurchasesAmount),
		CreatorId:       int(coursepb.CreatorId),
		TimeToPass:      int(coursepb.TimeToPass),
		Rating:          coursepb.Rating,
		Tags:            coursepb.Tags,
		Title:           coursepb.Title,
		Description:     coursepb.Description,
		ScrImage:        coursepb.ScrImage,
		IsPurchased:     coursepb.IsPurchased,
		IsFavorite:      coursepb.IsFavorite,
		Parts:           parts,
	}
}

func mapToGetUserProfile(user *coursepb.UserProfile) *usermodels.UserProfile {
	if user == nil {
		return nil
	}

	return &usermodels.UserProfile{
		Id:        int(user.Id),
		Email:     user.Email,
		Name:      user.Name,
		Bio:       user.Bio,
		AvatarSrc: user.AvatarSrc,
		HideEmail: user.HideEmail,
		IsAdmin:   user.IsAdmin,
	}
}

func mapToTestDTO(test *dto.Test) *coursepb.TestDTO {
	if test == nil {
		return nil
	}
	return &coursepb.TestDTO{
		QuestionId: int32(test.QuestionID),
		Question:   test.Question,
		Answers:    mapAnswers(test.Answers),
	}
}

func mapToUserAnswerDTO(test *dto.Test) *coursepb.UserAnswer {
	if test == nil {
		return nil
	}
	return &coursepb.UserAnswer{
		IsRight:    test.UserAnswer.IsRight,
		QuestionId: int32(test.UserAnswer.QuestionID),
		AnswerId:   int32(test.UserAnswer.AnswerID),
	}
}

func mapAnswers(parts []*dto.QuizAnswer) []*coursepb.AnswerTestDTO {
	res := make([]*coursepb.AnswerTestDTO, 0, len(parts))
	for _, part := range parts {
		res = append(res, &coursepb.AnswerTestDTO{
			AnswerId: int32(part.AnswerID),
			Answer:   part.Answer,
			IsRight:  part.IsRight,
		})
	}
	return res
}

func mapToGetTestLessonResponse(test *dto.Test) *coursepb.GetTestLessonResponse {
	return &coursepb.GetTestLessonResponse{
		TestDTO:    mapToTestDTO(test),
		UserAnswer: mapToUserAnswerDTO(test),
	}
}

func mapToAnswerQuizResponse(test *dto.QuizResult) *coursepb.AnswerQuizResponse {
	return &coursepb.AnswerQuizResponse{
		IsRight: test.Result,
	}
}
