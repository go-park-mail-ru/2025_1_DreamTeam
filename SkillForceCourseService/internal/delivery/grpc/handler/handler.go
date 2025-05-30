package grpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	coursepb "skillForce/internal/delivery/grpc/proto"
	"skillForce/internal/usecase"
	"skillForce/pkg/logs"

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

func (h *CourseHandler) GetCompletedBucketCourses(ctx context.Context, req *coursepb.GetBucketCoursesRequest) (*coursepb.GetBucketCoursesResponse, error) {
	courses, err := h.usecase.GetCompletedBucketCourses(ctx, mapToGetUserProfile(req.UserProfile))
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

func (h *CourseHandler) MarkCourseAsCompleted(ctx context.Context, req *coursepb.MarkCourseAsCompletedRequest) (*emptypb.Empty, error) {
	err := h.usecase.MarkCourseAsCompleted(ctx, int(req.UserId), int(req.CourseId))
	logs.PrintLog(ctx, "MarkCourseAsCompleted", fmt.Sprintf("course id: %v, user id: %v", int(req.CourseId), int(req.UserId)))
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

func (h *CourseHandler) GetRating(ctx context.Context, req *coursepb.GetRatingRequest) (*coursepb.GetRatingResponse, error) {
	rating, err := h.usecase.GetRating(ctx, int(req.UserId), int(req.CourseId))
	if err != nil {
		return nil, err
	}
	return mapToRatingResponse(rating), nil
}

func (h *CourseHandler) GetSertificate(ctx context.Context, req *coursepb.GetSertificateRequest) (*coursepb.GetSertificateResponse, error) {
	sertificateUrl, err := h.usecase.GetSertificate(ctx, mapToGetUserProfile(req.User), int(req.CourseId))
	if err != nil {
		return nil, err
	}
	return &coursepb.GetSertificateResponse{SertificateUrl: sertificateUrl}, nil
}

func (h *CourseHandler) GetGeneratedSertificate(ctx context.Context, req *coursepb.GetSertificateRequest) (*coursepb.GetSertificateResponse, error) {
	sertificateUrl, err := h.usecase.GetGeneratedSertificate(ctx, mapToGetUserProfile(req.User), int(req.CourseId))
	if err != nil {
		return nil, err
	}
	return &coursepb.GetSertificateResponse{SertificateUrl: sertificateUrl}, nil
}

func (h *CourseHandler) GetStatistic(ctx context.Context, req *coursepb.GetStatisticRequest) (*coursepb.GetStatisticResponse, error) {
	stats, err := h.usecase.GetStatistic(ctx, int(req.UserId), int(req.CourseId))
	if err != nil {
		return nil, err
	}
	return mapToStatisticResponse(stats), nil
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

func (h *CourseHandler) AddRaiting(ctx context.Context, req *coursepb.AddRaitingRequest) (*emptypb.Empty, error) {
	err := h.usecase.AddRating(ctx, int(req.CourseId), int(req.UserId), int(req.Raiting))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *CourseHandler) UploadFile(ctx context.Context, req *coursepb.UploadFileRequest) (*coursepb.UploadFileResponse, error) {
	file, fileHeader, err := ConvertToMultipart(req.FileData, req.FileName, req.ContentType)
	if err != nil {
		return nil, err
	}
	url, err := h.usecase.UploadFile(ctx, file, fileHeader)
	if err != nil {
		return nil, err
	}
	return &coursepb.UploadFileResponse{Url: url}, nil
}

func (h *CourseHandler) SaveCourseImage(ctx context.Context, req *coursepb.SaveCourseImageRequest) (*coursepb.SaveCourseImageResponse, error) {
	newURL, err := h.usecase.SaveCourseImage(ctx, req.Url, int(req.CourseId))
	if err != nil {
		return nil, err
	}
	return &coursepb.SaveCourseImageResponse{NewPhtotoUrl: newURL}, nil
}

func ConvertToMultipart(fileData []byte, fileName, contentType string) (multipart.File, *multipart.FileHeader, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Создаем заголовки для файла
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+fileName+`"`)
	h.Set("Content-Type", contentType)

	// Пишем файл в multipart.Writer
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, nil, err
	}
	if _, err := io.Copy(part, bytes.NewReader(fileData)); err != nil {
		return nil, nil, err
	}
	err = writer.Close()
	if err != nil {
		logs.PrintLog(context.Background(), "ConvertToMultipart", fmt.Sprintf("%+v", err))
	}

	// Парсим то, что получилось, как multipart/form-data
	r := multipart.NewReader(body, writer.Boundary())
	form, err := r.ReadForm(int64(len(fileData)) + 1024) // выделяем буфер
	if err != nil {
		return nil, nil, err
	}

	files := form.File["file"]
	if len(files) == 0 {
		return nil, nil, io.EOF
	}
	fh := files[0]
	f, err := fh.Open()
	if err != nil {
		return nil, nil, err
	}

	return f, fh, nil
}
