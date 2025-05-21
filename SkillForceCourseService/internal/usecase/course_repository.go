package usecase

import (
	"context"
	"mime/multipart"
	coursemodels "skillForce/internal/models/course"
	"skillForce/internal/models/dto"
	usermodels "skillForce/internal/models/user"
)

type CourseRepository interface {
	GetBucketCourses(ctx context.Context) ([]*coursemodels.Course, error)
	GetPurchasedBucketCourses(ctx context.Context, userId int) ([]*coursemodels.Course, error)
	GetCompletedBucketCourses(ctx context.Context, userId int) ([]*coursemodels.Course, error)
	SearchCoursesByTitle(ctx context.Context, keywords string) ([]*coursemodels.Course, error)
	GetCourseById(ctx context.Context, courseId int) (*coursemodels.Course, error)
	GetCourseParts(ctx context.Context, courseId int) ([]*coursemodels.CoursePart, error)
	GetPartBuckets(ctx context.Context, partId int) ([]*coursemodels.LessonBucket, error)
	GetBucketLessons(ctx context.Context, userId int, courseId int, bucketId int) ([]*coursemodels.LessonPoint, error)
	GetLessonById(ctx context.Context, lessonId int) (*coursemodels.LessonPoint, error)
	GetBucketByLessonId(ctx context.Context, lessonId int) (*coursemodels.LessonBucket, error)
	GetLessonVideo(ctx context.Context, currentLessonId int) ([]string, error)
	GetLessonBlocks(ctx context.Context, currentLessonId int) ([]string, error)
	GetLessonTest(ctx context.Context, currentLessonId int, user_id int) (*dto.Test, error)
	AnswerQuiz(ctx context.Context, question_id int, answer_id int, user_id int, course_id int) (*dto.QuizResult, error)
	GetQuestionTestLesson(ctx context.Context, currentLessonId int, user_id int) (*dto.QuestionTest, error)
	AnswerQuestion(ctx context.Context, question_id int, user_id int, answer string) error
	GetRating(ctx context.Context, userId int, courseId int) (*dto.Raiting, error)
	GetGeneratedSertificate(ctx context.Context, userProfile *usermodels.UserProfile, courseId int) (string, error)

	GetLastLessonHeader(ctx context.Context, userId int, courseId int) (*dto.LessonDtoHeader, int, string, bool, error)
	GetLessonHeaderByLessonId(ctx context.Context, userId int, currentLessonId int) (*dto.LessonDtoHeader, error)
	GetLessonFooters(ctx context.Context, currentLessonId int) ([]int, error)
	IsMiddle(ctx context.Context, userId int, courseId int) (bool, error)

	MarkLessonCompleted(ctx context.Context, userId int, lessonId int) error
	MarkLessonAsNotCompleted(ctx context.Context, userId int, lessonId int) error
	IsUserPurchasedCourse(ctx context.Context, userId int, courseId int) (bool, error)
	IsUserCompletedCourse(ctx context.Context, userId int, courseId int) (bool, error)
	AddUserToCourse(ctx context.Context, userId int, courseId int) error
	MarkCourseAsCompleted(ctx context.Context, userId int, courseId int) error

	GetCoursesRaitings(ctx context.Context, bucketCoursesWithoutRating []*coursemodels.Course) (map[int]float32, error)
	GetCoursesTags(ctx context.Context, bucketCoursesWithoutTags []*coursemodels.Course) (map[int][]string, error)
	GetCoursesPurchases(ctx context.Context, bucketCoursesWithoutPurchases []*coursemodels.Course) (map[int]int, error)

	GetUserById(ctx context.Context, userId int) (*usermodels.User, error)

	CreateCourse(ctx context.Context, course *coursemodels.Course, userProfile *usermodels.UserProfile) (int, error)
	CreatePart(ctx context.Context, part *coursemodels.CoursePart, courseId int) (int, error)
	CreateBucket(ctx context.Context, bucket *coursemodels.LessonBucket, partId int) (int, error)
	CreateTextLesson(ctx context.Context, lesson *coursemodels.LessonPoint, bucketId int) error
	CreateVideoLesson(ctx context.Context, lesson *coursemodels.LessonPoint, bucketId int) error

	AddCourseToFavourites(ctx context.Context, courseId int, userId int) error
	DeleteCourseFromFavourites(ctx context.Context, courseId int, userId int) error
	GetFavouriteCourses(ctx context.Context, userId int) ([]*coursemodels.Course, error)
	GetCoursesFavouriteStatus(ctx context.Context, bucketCourses []*coursemodels.Course, userId int) (map[int]bool, error)

	UploadFileToMinIO(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}
