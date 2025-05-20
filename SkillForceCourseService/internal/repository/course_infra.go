package repository

import (
	"context"
	"fmt"
	"log"
	"skillForce/config"
	coursemodels "skillForce/internal/models/course"
	"skillForce/internal/models/dto"
	usermodels "skillForce/internal/models/user"
	"skillForce/internal/repository/postgres"
)

type CourseInfrastructure struct {
	Database *postgres.Database
}

func NewCourseInfrastructure(conf *config.Config) *CourseInfrastructure {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Name)
	database, err := postgres.NewDatabase(dsn, conf.Secrets.JwtSessionSecret)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return &CourseInfrastructure{
		Database: database,
	}
}

func (i *CourseInfrastructure) Close() {
	i.Database.Close()
}

func (i *CourseInfrastructure) GetBucketCourses(ctx context.Context) ([]*coursemodels.Course, error) {
	return i.Database.GetBucketCourses(ctx)
}

func (i *CourseInfrastructure) GetPurchasedBucketCourses(ctx context.Context, userId int) ([]*coursemodels.Course, error) {
	return i.Database.GetPurchasedBucketCourses(ctx, userId)
}

func (i *CourseInfrastructure) GetCompletedBucketCourses(ctx context.Context, userId int) ([]*coursemodels.Course, error) {
	return i.Database.GetCompletedBucketCourses(ctx, userId)
}

func (i *CourseInfrastructure) GetCoursesRaitings(ctx context.Context, bucketCoursesWithoutRating []*coursemodels.Course) (map[int]float32, error) {
	return i.Database.GetCoursesRaitings(ctx, bucketCoursesWithoutRating)
}

func (i *CourseInfrastructure) GetRating(ctx context.Context, userId int, courseId int) (*dto.Raiting, error) {
	return i.Database.GetRating(ctx, userId, courseId)
}

func (i *CourseInfrastructure) GetCoursesTags(ctx context.Context, bucketCoursesWithoutTags []*coursemodels.Course) (map[int][]string, error) {
	return i.Database.GetCoursesTags(ctx, bucketCoursesWithoutTags)
}

func (i *CourseInfrastructure) GetCourseById(ctx context.Context, courseId int) (*coursemodels.Course, error) {
	return i.Database.GetCourseById(ctx, courseId)
}

func (i *CourseInfrastructure) GetLastLessonHeader(ctx context.Context, userId int, courseId int) (*dto.LessonDtoHeader, int, string, bool, error) {
	return i.Database.GetLastLessonHeader(ctx, userId, courseId)
}

func (i *CourseInfrastructure) GetLessonBlocks(ctx context.Context, currentLessonId int) ([]string, error) {
	return i.Database.GetLessonBlocks(ctx, currentLessonId)
}

func (i *CourseInfrastructure) GetLessonFooters(ctx context.Context, currentLessonId int) ([]int, error) {
	return i.Database.GetLessonFooters(ctx, currentLessonId)
}

func (i *CourseInfrastructure) MarkLessonCompleted(ctx context.Context, userId int, lessonId int) error {
	return i.Database.MarkLessonCompleted(ctx, userId, lessonId)
}

func (i *CourseInfrastructure) MarkCourseAsCompleted(ctx context.Context, userId int, courseId int) error {
	return i.Database.MarkCourseAsCompleted(ctx, userId, courseId)
}

func (i *CourseInfrastructure) MarkLessonAsNotCompleted(ctx context.Context, userId int, lessonId int) error {
	return i.Database.MarkLessonAsNotCompleted(ctx, userId, lessonId)
}

func (i *CourseInfrastructure) GetCourseParts(ctx context.Context, courseId int) ([]*coursemodels.CoursePart, error) {
	return i.Database.GetCourseParts(ctx, courseId)
}

func (i *CourseInfrastructure) GetPartBuckets(ctx context.Context, partId int) ([]*coursemodels.LessonBucket, error) {
	return i.Database.GetPartBuckets(ctx, partId)
}

func (i *CourseInfrastructure) GetBucketLessons(ctx context.Context, userId int, courseId int, bucketId int) ([]*coursemodels.LessonPoint, error) {
	return i.Database.GetBucketLessons(ctx, userId, courseId, bucketId)
}

func (i *CourseInfrastructure) AddUserToCourse(ctx context.Context, userId int, courseId int) error {
	return i.Database.AddUserToCourse(ctx, userId, courseId)
}

func (i *CourseInfrastructure) GetCoursesPurchases(ctx context.Context, bucketCoursesWithoutPurchases []*coursemodels.Course) (map[int]int, error) {
	return i.Database.GetCoursesPurchases(ctx, bucketCoursesWithoutPurchases)
}

func (i *CourseInfrastructure) GetBucketByLessonId(ctx context.Context, lessonId int) (*coursemodels.LessonBucket, error) {
	return i.Database.GetBucketByLessonId(ctx, lessonId)
}

func (i *CourseInfrastructure) GetLessonHeaderByLessonId(ctx context.Context, userId int, currentLessonId int) (*dto.LessonDtoHeader, error) {
	return i.Database.GetLessonHeaderByLessonId(ctx, userId, currentLessonId)
}

func (i *CourseInfrastructure) GetUserById(ctx context.Context, userId int) (*usermodels.User, error) {
	return i.Database.GetUserById(ctx, userId)
}

func (i *CourseInfrastructure) IsUserPurchasedCourse(ctx context.Context, userId int, courseId int) (bool, error) {
	return i.Database.IsUserPurchasedCourse(ctx, userId, courseId)
}

func (i *CourseInfrastructure) IsUserCompletedCourse(ctx context.Context, userId int, courseId int) (bool, error) {
	return i.Database.IsUserCompletedCourse(ctx, userId, courseId)
}

func (i *CourseInfrastructure) GetLessonVideo(ctx context.Context, lessonId int) ([]string, error) {
	return i.Database.GetLessonVideo(ctx, lessonId)
}

func (i *CourseInfrastructure) GetLessonById(ctx context.Context, lessonId int) (*coursemodels.LessonPoint, error) {
	return i.Database.GetLessonById(ctx, lessonId)
}

func (i *CourseInfrastructure) IsMiddle(ctx context.Context, userId int, courseId int) (bool, error) {
	return i.Database.IsMiddle(ctx, userId, courseId)
}

func (i *CourseInfrastructure) CreateCourse(ctx context.Context, course *coursemodels.Course, userProfile *usermodels.UserProfile) (int, error) {
	return i.Database.CreateCourse(ctx, course, userProfile)
}

func (i *CourseInfrastructure) CreatePart(ctx context.Context, coursePart *coursemodels.CoursePart, courseId int) (int, error) {
	return i.Database.CreatePart(ctx, coursePart, courseId)
}

func (i *CourseInfrastructure) CreateBucket(ctx context.Context, bucket *coursemodels.LessonBucket, partId int) (int, error) {
	return i.Database.CreateBucket(ctx, bucket, partId)
}

func (i *CourseInfrastructure) CreateTextLesson(ctx context.Context, lesson *coursemodels.LessonPoint, bucketId int) error {
	return i.Database.CreateTextLesson(ctx, lesson, bucketId)
}

func (i *CourseInfrastructure) CreateVideoLesson(ctx context.Context, lesson *coursemodels.LessonPoint, bucketId int) error {
	return i.Database.CreateVideoLesson(ctx, lesson, bucketId)
}

func (i *CourseInfrastructure) SendSurveyQuestionAnswer(ctx context.Context, surveyAnswerDto *coursemodels.SurveyAnswer, userProfile *usermodels.UserProfile) error {
	return i.Database.SendSurveyQuestionAnswer(ctx, surveyAnswerDto, userProfile)
}

func (i *CourseInfrastructure) AddCourseToFavourites(ctx context.Context, courseId int, userId int) error {
	return i.Database.AddCourseToFavourites(ctx, courseId, userId)
}

func (i *CourseInfrastructure) DeleteCourseFromFavourites(ctx context.Context, courseId int, userId int) error {
	return i.Database.DeleteCourseFromFavourites(ctx, courseId, userId)
}

func (i *CourseInfrastructure) GetFavouriteCourses(ctx context.Context, userId int) ([]*coursemodels.Course, error) {
	return i.Database.GetFavouriteCourses(ctx, userId)
}

func (i *CourseInfrastructure) GetCoursesFavouriteStatus(ctx context.Context, bucketCourses []*coursemodels.Course, userId int) (map[int]bool, error) {
	return i.Database.GetCoursesFavouriteStatus(ctx, bucketCourses, userId)
}

func (i *CourseInfrastructure) GetLessonTest(ctx context.Context, currentLessonId int, user_id int) (*dto.Test, error) {
	return i.Database.GetLessonTest(ctx, currentLessonId, user_id)
}

func (i *CourseInfrastructure) AnswerQuiz(ctx context.Context, question_id int, answer_id int, user_id int, course_id int) (*dto.QuizResult, error) {
	return i.Database.AnswerQuiz(ctx, question_id, answer_id, user_id, course_id)
}

func (i *CourseInfrastructure) GetQuestionTestLesson(ctx context.Context, currentLessonId int, user_id int) (*dto.QuestionTest, error) {
	return i.Database.GetQuestionTestLesson(ctx, currentLessonId, user_id)
}

func (i *CourseInfrastructure) AnswerQuestion(ctx context.Context, question_id int, user_id int, answer string) error {
	return i.Database.AnswerQuestion(ctx, question_id, user_id, answer)
}

func (i *CourseInfrastructure) SearchCoursesByTitle(ctx context.Context, keyword string) ([]*coursemodels.Course, error) {
	return i.Database.SearchCoursesByTitle(ctx, keyword)
}
