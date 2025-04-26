package usecase

import (
	"context"
	"io"
	coursemodels "skillForce/internal/models/course"
	"skillForce/internal/models/dto"
	usermodels "skillForce/internal/models/user"
)

type CourseRepository interface {
	GetBucketCourses(ctx context.Context) ([]*coursemodels.Course, error)
	GetCourseById(ctx context.Context, courseId int) (*coursemodels.Course, error)
	GetCourseParts(ctx context.Context, courseId int) ([]*coursemodels.CoursePart, error)
	GetPartBuckets(ctx context.Context, partId int) ([]*coursemodels.LessonBucket, error)
	GetBucketLessons(ctx context.Context, userId int, courseId int, bucketId int) ([]*coursemodels.LessonPoint, error)
	GetLessonById(ctx context.Context, lessonId int) (*coursemodels.LessonPoint, error)
	GetBucketByLessonId(ctx context.Context, lessonId int) (*coursemodels.LessonBucket, error)
	GetLessonVideo(ctx context.Context, currentLessonId int) ([]string, error)
	GetLessonBlocks(ctx context.Context, currentLessonId int) ([]string, error)

	GetLastLessonHeader(ctx context.Context, userId int, courseId int) (*dto.LessonDtoHeader, int, string, bool, error)
	GetLessonHeaderByLessonId(ctx context.Context, userId int, currentLessonId int) (*dto.LessonDtoHeader, error)
	GetLessonFooters(ctx context.Context, currentLessonId int) ([]int, error)
	IsMiddle(ctx context.Context, userId int, courseId int) (bool, error)

	MarkLessonCompleted(ctx context.Context, userId int, courseId int, lessonId int) error
	MarkLessonAsNotCompleted(ctx context.Context, userId int, lessonId int) error
	IsUserPurchasedCourse(ctx context.Context, userId int, courseId int) (bool, error)
	AddUserToCourse(ctx context.Context, userId int, courseId int) error

	GetCoursesRaitings(ctx context.Context, bucketCoursesWithoutRating []*coursemodels.Course) (map[int]float32, error)
	GetCoursesTags(ctx context.Context, bucketCoursesWithoutTags []*coursemodels.Course) (map[int][]string, error)
	GetCoursesPurchases(ctx context.Context, bucketCoursesWithoutPurchases []*coursemodels.Course) (map[int]int, error)

	GetVideoUrl(ctx context.Context, lessonId int) (string, error)
	GetVideoRange(ctx context.Context, name string, start, end int64) (io.ReadCloser, error)
	Stat(ctx context.Context, name string) (dto.VideoMeta, error)

	SendWelcomeCourseMail(ctx context.Context, user *usermodels.User, courseId int) error
	SendMiddleCourseMail(ctx context.Context, user *usermodels.User, courseId int) error

	GetUserById(ctx context.Context, userId int) (*usermodels.User, error)

	CreateCourse(ctx context.Context, course *coursemodels.Course, userProfile *usermodels.UserProfile) (int, error)
	CreatePart(ctx context.Context, part *coursemodels.CoursePart, courseId int) (int, error)
	CreateBucket(ctx context.Context, bucket *coursemodels.LessonBucket, partId int) (int, error)
	CreateTextLesson(ctx context.Context, lesson *coursemodels.LessonPoint, bucketId int) error
	CreateVideoLesson(ctx context.Context, lesson *coursemodels.LessonPoint, bucketId int) error
	CreateSurvey(ctx context.Context, survey *coursemodels.Survey, userProfile *usermodels.UserProfile) error
	SendSurveyQuestionAnswer(ctx context.Context, surveyAnswerDto *coursemodels.SurveyAnswer, userProfile *usermodels.UserProfile) error
}
