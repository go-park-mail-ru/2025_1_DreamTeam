package courseCourseInfrastructure

import (
	"context"
	"io"
	coursemodels "skillForce/internal/models/course"
	"skillForce/internal/models/dto"
	usermodels "skillForce/internal/models/user"
	"skillForce/internal/repository/mail"
	"skillForce/internal/repository/minio"
	"skillForce/internal/repository/postgres"
)

type CourseInfrastructure struct {
	Database *postgres.Database
	Minio    *minio.Minio
	Mail     *mail.Mail
}

func NewCourseInfrastructure(db *postgres.Database, mail *mail.Mail, minio *minio.Minio) *CourseInfrastructure {
	return &CourseInfrastructure{
		Database: db,
		Minio:    minio,
		Mail:     mail,
	}
}

func (i *CourseInfrastructure) GetBucketCourses(ctx context.Context) ([]*coursemodels.Course, error) {
	return i.Database.GetBucketCourses(ctx)
}

func (i *CourseInfrastructure) GetCoursesRaitings(ctx context.Context, bucketCoursesWithoutRating []*coursemodels.Course) (map[int]float32, error) {
	return i.Database.GetCoursesRaitings(ctx, bucketCoursesWithoutRating)
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

func (i *CourseInfrastructure) MarkLessonCompleted(ctx context.Context, userId int, courseId int, lessonId int) error {
	return i.Database.MarkLessonCompleted(ctx, userId, courseId, lessonId)
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

func (i *CourseInfrastructure) GetVideoUrl(ctx context.Context, lesson_id int) (string, error) {
	return i.Database.GetVideoUrl(ctx, lesson_id)
}

func (i *CourseInfrastructure) GetVideoRange(ctx context.Context, name string, start, end int64) (io.ReadCloser, error) {
	return i.Minio.GetVideoRange(ctx, name, start, end)
}

func (i *CourseInfrastructure) Stat(ctx context.Context, name string) (dto.VideoMeta, error) {
	return i.Minio.Stat(ctx, name)
}

func (i *CourseInfrastructure) GetUserById(ctx context.Context, userId int) (*usermodels.User, error) {
	return i.Database.GetUserById(ctx, userId)
}

func (i *CourseInfrastructure) SendWelcomeCourseMail(ctx context.Context, user *usermodels.User, courseId int) error {
	return i.Mail.SendWelcomeCourseMail(ctx, user, courseId)
}

func (i *CourseInfrastructure) IsUserPurchasedCourse(ctx context.Context, userId int, courseId int) (bool, error) {
	return i.Database.IsUserPurchasedCourse(ctx, userId, courseId)
}

func (i *CourseInfrastructure) GetLessonVideo(ctx context.Context, lessonId int) ([]string, error) {
	return i.Database.GetLessonVideo(ctx, lessonId)
}

func (i *CourseInfrastructure) GetLessonById(ctx context.Context, lessonId int) (*coursemodels.LessonPoint, error) {
	return i.Database.GetLessonById(ctx, lessonId)
}
