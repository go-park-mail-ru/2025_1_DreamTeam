package grpc

import (
	"context"
	coursepb "skillForce/internal/delivery/grpc/proto"
	"skillForce/internal/usecase"

	"google.golang.org/protobuf/types/known/emptypb"
)

type CourseHandler struct {
	coursepb.UnimplementedCourseServiceServer
	usecase *usecase.CourseUsecase
}

func NewCourseHandler(uc *usecase.CourseUsecase) *CourseHandler {
	return &CourseHandler{
		usecase: uc,
	}
}

func (h *CourseHandler) GetBucketCourses(ctx context.Context, req *coursepb.GetBucketCoursesRequest) (*coursepb.GetBucketCoursesResponse, error) {
	courses, err := h.usecase.GetBucketCourses(ctx, mapToGetUserProfile(req.UserProfile))
	if err != nil {
		return nil, err
	}
	return mapToGetBucketCoursesResponse(courses), nil
}

func (h *CourseHandler) GetPurchasedBucketCourses(ctx context.Context, req *coursepb.GetBucketCoursesRequest) (*coursepb.GetBucketCoursesResponse, error) {
	courses, err := h.usecase.GetPurchasedBucketCourses(ctx, mapToGetUserProfile(req.UserProfile))
	if err != nil {
		return nil, err
	}
	return mapToGetBucketCoursesResponse(courses), nil
}

func (h *CourseHandler) GetCourseLesson(ctx context.Context, req *coursepb.GetCourseLessonRequest) (*coursepb.GetCourseLessonResponse, error) {
	lesson, err := h.usecase.GetCourseLesson(ctx, int(req.UserId), int(req.CourseId))
	if err != nil {
		return nil, err
	}
	return mapToCourseLessonResponse(lesson), nil
}

func (h *CourseHandler) GetNextLesson(ctx context.Context, req *coursepb.GetNextLessonRequest) (*coursepb.GetNextLessonResponse, error) {
	lesson, err := h.usecase.GetNextLesson(ctx, int(req.UserId), int(req.CourseId), int(req.LessonId))
	if err != nil {
		return nil, err
	}
	return mapToNextLessonResponse(lesson), nil
}

func (h *CourseHandler) MarkLessonAsNotCompleted(ctx context.Context, req *coursepb.MarkLessonAsNotCompletedRequest) (*emptypb.Empty, error) {
	err := h.usecase.MarkLessonAsNotCompleted(ctx, int(req.UserId), int(req.LessonId))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *CourseHandler) MarkLessonAsCompleted(ctx context.Context, req *coursepb.MarkLessonAsCompletedRequest) (*emptypb.Empty, error) {
	err := h.usecase.MarkLessonAsCompleted(ctx, int(req.UserId), int(req.LessonId))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *CourseHandler) GetCourseRoadmap(ctx context.Context, req *coursepb.GetCourseRoadmapRequest) (*coursepb.GetCourseRoadmapResponse, error) {
	roadmap, err := h.usecase.GetCourseRoadmap(ctx, int(req.UserId), int(req.CourseId))
	if err != nil {
		return nil, err
	}
	return mapToCourseRoadmapResponse(roadmap), nil
}

func (h *CourseHandler) GetCourse(ctx context.Context, req *coursepb.GetCourseRequest) (*coursepb.GetCourseResponse, error) {
	course, err := h.usecase.GetCourse(ctx, int(req.CourseId), mapToGetUserProfile(req.UserProfile))
	if err != nil {
		return nil, err
	}
	return mapToCourseResponse(course), nil
}

func (h *CourseHandler) CreateCourse(ctx context.Context, req *coursepb.CreateCourseRequest) (*emptypb.Empty, error) {
	err := h.usecase.CreateCourse(ctx, mapPbCourseDTOToDTO(req.Course), mapToGetUserProfile(req.UserProfile))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *CourseHandler) AddCourseToFavourites(ctx context.Context, req *coursepb.AddToFavouritesRequest) (*emptypb.Empty, error) {
	err := h.usecase.AddCourseToFavourites(ctx, mapPbCourseDTOToDTO(req.Course), mapToGetUserProfile(req.UserProfile))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *CourseHandler) DeleteCourseFromFavourites(ctx context.Context, req *coursepb.DeleteCourseFromFavouritesRequest) (*emptypb.Empty, error) {
	err := h.usecase.DeleteCourseFromFavourites(ctx, mapPbCourseDTOToDTO(req.Course), mapToGetUserProfile(req.UserProfile))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *CourseHandler) GetFavouriteCourses(ctx context.Context, req *coursepb.GetFavouritesRequest) (*coursepb.GetFavouritesResponse, error) {
	courses, err := h.usecase.GetFavouriteCourses(ctx, mapToGetUserProfile(req.UserProfile))
	if err != nil {
		return nil, err
	}
	return mapToGetFavouritesResponse(courses), nil
}

func (h *CourseHandler) GetTestLesson(ctx context.Context, req *coursepb.GetTestLessonRequest) (*coursepb.GetTestLessonResponse, error) {
	tests, err := h.usecase.GetTestLesson(ctx, int(req.LessonId), int(req.UserId))
	if err != nil {
		return nil, err
	}
	return mapToGetTestLessonResponse(tests), nil
}

func (h *CourseHandler) AnswerQuiz(ctx context.Context, req *coursepb.AnswerQuizRequest) (*coursepb.AnswerQuizResponse, error) {
	res, err := h.usecase.AnswerQuiz(ctx, int(req.QuestionId), int(req.AnswerId), int(req.UserId), int(req.CourseId))
	if err != nil {
		return nil, err
	}
	return mapToAnswerQuizResponse(res), nil
}

func (h *CourseHandler) GetQuestionTestLesson(ctx context.Context, req *coursepb.GetQuestionTestLessonRequest) (*coursepb.GetQuestionTestLessonResponse, error) {
	question, err := h.usecase.GetQuestionTestLesson(ctx, int(req.LessonId), int(req.UserId))
	if err != nil {
		return nil, err
	}
	return mapToGetQuestionTestLessonResponse(question), nil
}

func (h *CourseHandler) AnswerQuestion(ctx context.Context, req *coursepb.AnswerQuestionRequest) (*emptypb.Empty, error) {
	err := h.usecase.AnswerQuestion(ctx, int(req.QuestionId), int(req.UserId), req.Answer)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *CourseHandler) SearchCoursesByTitle(ctx context.Context, req *coursepb.SearchCoursesByTitleRequest) (*coursepb.GetBucketCoursesResponse, error) {
	courses, err := h.usecase.SearchCoursesByTitle(ctx, mapToGetUserProfile(req.UserProfile), req.Keywords)
	if err != nil {
		return nil, err
	}
	return mapToGetBucketCoursesResponse(courses), nil
}
