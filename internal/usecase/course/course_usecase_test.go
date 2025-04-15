package usecase

import (
	"bytes"
	"context"
	"errors"
	"io"
	coursemodels "skillForce/internal/models/course"
	"skillForce/internal/models/dto"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/logs"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetBucketCourses_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	courses := []*coursemodels.Course{
		{
			Id:          1,
			CreatorId:   10,
			Title:       "Go Basics",
			Description: "Learn Go",
			ScrImage:    "image.png",
			Price:       100,
			TimeToPass:  5,
		},
	}

	ratings := map[int]float32{1: 4.2}
	tags := map[int][]string{1: {"go", "programming"}}
	purchases := map[int]int{1: 100}

	mockRepo.EXPECT().GetBucketCourses(ctx).Return(courses, nil)
	mockRepo.EXPECT().GetCoursesRaitings(ctx, courses).Return(ratings, nil)
	mockRepo.EXPECT().GetCoursesTags(ctx, courses).Return(tags, nil)
	mockRepo.EXPECT().GetCoursesPurchases(ctx, courses).Return(purchases, nil)

	result, err := uc.GetBucketCourses(ctx)
	require.NoError(t, err)
	require.Len(t, result, 1)

	c := result[0]
	require.Equal(t, 1, c.Id)
	require.Equal(t, "Go Basics", c.Title)
	require.Equal(t, float32(4.2), c.Rating)
	require.Equal(t, 100, c.PurchasesAmount)
	require.Equal(t, []string{"go", "programming"}, c.Tags)
}

func TestGetBucketCourses_GetBucketCourses_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	mockRepo.EXPECT().GetBucketCourses(ctx).Return(nil, errors.New("db error"))

	result, err := uc.GetBucketCourses(ctx)
	require.Error(t, err)
	require.Nil(t, result)
}

func TestGetBucketCourses_GetCoursesRaitings_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	courses := []*coursemodels.Course{
		{Id: 1, Title: "Go Basics"},
	}

	mockRepo.EXPECT().GetBucketCourses(ctx).Return(courses, nil)
	mockRepo.EXPECT().GetCoursesRaitings(ctx, courses).Return(nil, errors.New("ratings error"))

	result, err := uc.GetBucketCourses(ctx)
	require.Error(t, err)
	require.Nil(t, result)
}

func TestGetBucketCourses_GetCoursesTags_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	courses := []*coursemodels.Course{
		{Id: 1, Title: "Go Basics"},
	}
	ratings := map[int]float32{1: 4.5}

	mockRepo.EXPECT().GetBucketCourses(ctx).Return(courses, nil)
	mockRepo.EXPECT().GetCoursesRaitings(ctx, courses).Return(ratings, nil)
	mockRepo.EXPECT().GetCoursesTags(ctx, courses).Return(nil, errors.New("tags error"))

	result, err := uc.GetBucketCourses(ctx)
	require.Error(t, err)
	require.Nil(t, result)
}

func TestGetBucketCourses_GetCoursesPurchases_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	courses := []*coursemodels.Course{
		{Id: 1, Title: "Go Basics"},
	}
	ratings := map[int]float32{1: 4.5}
	tags := map[int][]string{1: {"go", "backend"}}

	mockRepo.EXPECT().GetBucketCourses(ctx).Return(courses, nil)
	mockRepo.EXPECT().GetCoursesRaitings(ctx, courses).Return(ratings, nil)
	mockRepo.EXPECT().GetCoursesTags(ctx, courses).Return(tags, nil)
	mockRepo.EXPECT().GetCoursesPurchases(ctx, courses).Return(nil, errors.New("purchases error"))

	result, err := uc.GetBucketCourses(ctx)
	require.Error(t, err)
	require.Nil(t, result)
}

func TestGetCourse_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	course := &coursemodels.Course{
		Id:          1,
		CreatorId:   10,
		Title:       "Go Mastery",
		Description: "<p>Learn Go</p>",
		ScrImage:    "image.png",
		Price:       1999,
		TimeToPass:  120,
	}

	ratings := map[int]float32{1: 4.7}
	tags := map[int][]string{1: {"go", "backend"}}
	purchases := map[int]int{1: 42}

	userProfile := &usermodels.UserProfile{
		Id: 100,
	}

	mockRepo.EXPECT().GetCourseById(ctx, 1).Return(course, nil)
	mockRepo.EXPECT().GetCoursesRaitings(ctx, []*coursemodels.Course{course}).Return(ratings, nil)
	mockRepo.EXPECT().GetCoursesTags(ctx, []*coursemodels.Course{course}).Return(tags, nil)
	mockRepo.EXPECT().GetCoursesPurchases(ctx, []*coursemodels.Course{course}).Return(purchases, nil)
	mockRepo.EXPECT().IsUserPurchasedCourse(ctx, 100, 1).Return(true, nil)

	result, err := uc.GetCourse(ctx, 1, userProfile)
	require.NoError(t, err)
	require.NotNil(t, result)

	require.Equal(t, course.Id, result.Id)
	require.Equal(t, "Go Mastery", result.Title)
	require.Equal(t, float32(4.7), result.Rating)
	require.Equal(t, []string{"go", "backend"}, result.Tags)
	require.Equal(t, 42, result.PurchasesAmount)
	require.Equal(t, true, result.IsPurchased)
}

func TestGetCourseLesson_TextLesson_FirstLesson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	userId := 1
	courseId := 2
	currentLessonId := 3

	header := &dto.LessonDtoHeader{
		CourseTitle: "Go Mastery",
		CourseId:    courseId,
	}
	lessonType := "text"
	first := true

	blocks := []string{"block 1", "block 2"}
	footers := []int{0, 3, 4}

	user := &usermodels.User{
		Id:        userId,
		HideEmail: false,
		Email:     "user@example.com",
	}

	mockRepo.EXPECT().AddUserToCourse(ctx, userId, courseId).Return(nil)
	mockRepo.EXPECT().GetLastLessonHeader(ctx, userId, courseId).Return(header, currentLessonId, lessonType, first, nil)
	mockRepo.EXPECT().GetLessonBlocks(ctx, currentLessonId).Return(blocks, nil)
	mockRepo.EXPECT().GetLessonFooters(ctx, currentLessonId).Return(footers, nil)
	mockRepo.EXPECT().GetUserById(ctx, userId).Return(user, nil)
	mockRepo.EXPECT().SendWelcomeCourseMail(ctx, user, courseId).Return(nil)

	result, err := uc.GetCourseLesson(ctx, userId, courseId)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.LessonBody.Blocks, 2)
	require.Equal(t, 3, result.LessonBody.Footer.CurrentLessonId)
	require.Equal(t, 4, result.LessonBody.Footer.NextLessonId)
	require.Equal(t, 0, result.LessonBody.Footer.PreviousLessonId)
}

func TestGetCourseLesson_VideoLesson_FirstLesson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	userId := 1
	courseId := 2
	currentLessonId := 3

	header := &dto.LessonDtoHeader{
		CourseTitle: "Go Mastery",
		CourseId:    courseId,
	}
	lessonType := "video"
	first := true

	videoBlocks := []string{"video_block_1", "video_block_2"}
	footers := []int{1, 3, 5}

	user := &usermodels.User{
		Id:        userId,
		HideEmail: false,
		Email:     "user@example.com",
	}

	mockRepo.EXPECT().AddUserToCourse(ctx, userId, courseId).Return(nil)
	mockRepo.EXPECT().GetLastLessonHeader(ctx, userId, courseId).Return(header, currentLessonId, lessonType, first, nil)
	mockRepo.EXPECT().GetLessonVideo(ctx, currentLessonId).Return(videoBlocks, nil)
	mockRepo.EXPECT().GetLessonFooters(ctx, currentLessonId).Return(footers, nil)
	mockRepo.EXPECT().GetUserById(ctx, userId).Return(user, nil)
	mockRepo.EXPECT().SendWelcomeCourseMail(ctx, user, courseId).Return(nil)

	result, err := uc.GetCourseLesson(ctx, userId, courseId)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, header.CourseId, result.LessonHeader.CourseId)
	require.Len(t, result.LessonBody.Blocks, 2)
	require.Equal(t, "video_block_1", result.LessonBody.Blocks[0].Body)
	require.Equal(t, 5, result.LessonBody.Footer.NextLessonId)
	require.Equal(t, 3, result.LessonBody.Footer.CurrentLessonId)
	require.Equal(t, 1, result.LessonBody.Footer.PreviousLessonId)
}

func TestGetNextLesson_VideoLesson_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	userId := 1
	courseId := 2
	lessonId := 3

	lesson := &coursemodels.LessonPoint{
		Type:  "video",
		Title: "Sample Video Lesson",
	}

	videoBlocks := []string{"block1", "block2"}
	footers := []int{2, 3, 4}

	header := &dto.LessonDtoHeader{
		CourseTitle: "Go Mastery",
		CourseId:    courseId,
	}

	mockRepo.EXPECT().GetLessonById(ctx, lessonId).Return(lesson, nil)
	mockRepo.EXPECT().GetLessonVideo(ctx, lessonId).Return(videoBlocks, nil)
	mockRepo.EXPECT().GetLessonFooters(ctx, lessonId).Return(footers, nil)
	mockRepo.EXPECT().MarkLessonCompleted(ctx, userId, courseId, lessonId).Return(nil)
	mockRepo.EXPECT().GetLessonHeaderByLessonId(ctx, userId, lessonId).Return(header, nil)

	result, err := uc.GetNextLesson(ctx, userId, courseId, lessonId)
	require.NoError(t, err)
	require.NotNil(t, result)

	require.Equal(t, header.CourseId, result.LessonHeader.CourseId)
	require.Len(t, result.LessonBody.Blocks, 2)
	require.Equal(t, "block1", result.LessonBody.Blocks[0].Body)
	require.Equal(t, footers[2], result.LessonBody.Footer.NextLessonId)
	require.Equal(t, footers[1], result.LessonBody.Footer.CurrentLessonId)
	require.Equal(t, footers[0], result.LessonBody.Footer.PreviousLessonId)
}

func TestGetNextLesson_TextLesson_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	userId := 1
	courseId := 2
	lessonId := 3

	lesson := &coursemodels.LessonPoint{
		Type:  "text",
		Title: "Sample Text Lesson",
	}

	textBlocks := []string{"block1", "block2"}
	footers := []int{10, 11, 12}

	header := &dto.LessonDtoHeader{
		CourseTitle: "Go Mastery",
		CourseId:    courseId,
	}

	mockRepo.EXPECT().GetLessonById(ctx, lessonId).Return(lesson, nil)
	mockRepo.EXPECT().GetLessonBlocks(ctx, lessonId).Return(textBlocks, nil)
	mockRepo.EXPECT().GetLessonFooters(ctx, lessonId).Return(footers, nil)
	mockRepo.EXPECT().MarkLessonCompleted(ctx, userId, courseId, lessonId).Return(nil)
	mockRepo.EXPECT().GetLessonHeaderByLessonId(ctx, userId, lessonId).Return(header, nil)

	result, err := uc.GetNextLesson(ctx, userId, courseId, lessonId)
	require.NoError(t, err)
	require.NotNil(t, result)

	require.Equal(t, header.CourseId, result.LessonHeader.CourseId)
	require.Len(t, result.LessonBody.Blocks, 2)
	require.Equal(t, "block1", result.LessonBody.Blocks[0].Body)
	require.Equal(t, footers[2], result.LessonBody.Footer.NextLessonId)
	require.Equal(t, footers[1], result.LessonBody.Footer.CurrentLessonId)
	require.Equal(t, footers[0], result.LessonBody.Footer.PreviousLessonId)
}

func TestGetCourseRoadmap_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	userId := 1
	courseId := 2

	part := &coursemodels.CoursePart{
		Id:    10,
		Title: "Part 1",
	}
	bucket := &coursemodels.LessonBucket{
		Id:    20,
		Title: "Bucket 1",
	}
	lesson := &coursemodels.LessonPoint{
		LessonId: 30,
		Title:    "Lesson 1",
		IsDone:   true,
		Type:     "text",
	}

	mockRepo.EXPECT().GetCourseParts(ctx, courseId).Return([]*coursemodels.CoursePart{part}, nil)
	mockRepo.EXPECT().GetPartBuckets(ctx, part.Id).Return([]*coursemodels.LessonBucket{bucket}, nil)
	mockRepo.EXPECT().GetBucketLessons(ctx, userId, courseId, bucket.Id).Return([]*coursemodels.LessonPoint{lesson}, nil)

	result, err := uc.GetCourseRoadmap(ctx, userId, courseId)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Parts, 1)
	require.Equal(t, part.Id, result.Parts[0].Id)
	require.Len(t, result.Parts[0].Buckets, 1)
	require.Equal(t, bucket.Id, result.Parts[0].Buckets[0].Id)
	require.Len(t, result.Parts[0].Buckets[0].Lessons, 1)
	require.Equal(t, lesson.LessonId, result.Parts[0].Buckets[0].Lessons[0].LessonId)
}

func TestGetVideoUrl_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.Background()
	lessonId := 42
	expectedUrl := "https://example.com/video.mp4"

	mockRepo.EXPECT().GetVideoUrl(ctx, lessonId).Return(expectedUrl, nil)

	url, err := uc.GetVideoUrl(ctx, lessonId)
	require.NoError(t, err)
	require.Equal(t, expectedUrl, url)
}

func TestGetMeta_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.Background()
	videoName := "lecture.mp4"
	expectedMeta := dto.VideoMeta{
		Name: videoName,
		Size: 123456,
	}

	mockRepo.EXPECT().Stat(ctx, videoName).Return(expectedMeta, nil)

	meta, err := uc.GetMeta(ctx, videoName)
	require.NoError(t, err)
	require.Equal(t, expectedMeta, meta)
}

func TestGetFragment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.Background()
	videoName := "video.mp4"
	start := int64(0)
	end := int64(1024)

	expectedReader := io.NopCloser(strings.NewReader("video fragment"))

	mockRepo.EXPECT().GetVideoRange(ctx, videoName, start, end).Return(expectedReader, nil)

	reader, err := uc.GetFragment(ctx, videoName, start, end)
	require.NoError(t, err)

	buf := new(bytes.Buffer)
	_, readErr := buf.ReadFrom(reader)
	require.NoError(t, readErr)
	require.Equal(t, "video fragment", buf.String())
}

func TestMarkLessonAsNotCompleted_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockCourseRepository(ctrl)
	uc := NewCourseUsecase(mockRepo)

	ctx := context.Background()
	userId := 1
	lessonId := 100

	mockRepo.EXPECT().MarkLessonAsNotCompleted(ctx, userId, lessonId).Return(nil)

	err := uc.MarkLessonAsNotCompleted(ctx, userId, lessonId)
	require.NoError(t, err)
}
