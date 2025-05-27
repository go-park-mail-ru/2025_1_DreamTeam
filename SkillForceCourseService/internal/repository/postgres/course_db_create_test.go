package postgres

import (
	"context"
	"regexp"
	"testing"

	coursemodels "skillForce/internal/models/course"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/logs"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*Database, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	return &Database{conn: db}, mock

}

func TestCreateCourse(t *testing.T) {
	db, mock := setupMockDB(t)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	course := &coursemodels.Course{
		Title:       "Go Basics",
		Description: "Learn Go",
		Price:       100,
		TimeToPass:  5,
	}
	user := &usermodels.UserProfile{Id: 1}

	mock.ExpectQuery(`INSERT INTO COURSE`).
		WithArgs(user.Id, course.Title, course.Description, course.Price, course.TimeToPass).
		WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(42))

	id, err := db.CreateCourse(ctx, course, user)

	assert.NoError(t, err)
	assert.Equal(t, 42, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePart(t *testing.T) {
	db, mock := setupMockDB(t)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	part := &coursemodels.CoursePart{Order: 1, Title: "Introduction"}
	courseId := 10

	mock.ExpectQuery(`INSERT INTO PART`).
		WithArgs(courseId, part.Order, part.Title).
		WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(11))

	id, err := db.CreatePart(ctx, part, courseId)

	assert.NoError(t, err)
	assert.Equal(t, 11, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateBucket(t *testing.T) {
	db, mock := setupMockDB(t)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	bucket := &coursemodels.LessonBucket{Order: 2, Title: "Basics"}
	partId := 11

	mock.ExpectQuery(`INSERT INTO LESSON_BUCKET`).
		WithArgs(partId, bucket.Order, bucket.Title).
		WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(22))

	id, err := db.CreateBucket(ctx, bucket, partId)

	assert.NoError(t, err)
	assert.Equal(t, 22, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTextLesson(t *testing.T) {
	db, mock := setupMockDB(t)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	lesson := &coursemodels.LessonPoint{
		Order:   1,
		Title:   "Intro",
		Type:    "text",
		Value:   "Welcome",
		IsImage: false,
	}
	bucketId := 5

	mock.ExpectQuery(`INSERT INTO LESSON`).
		WithArgs(bucketId, lesson.Order, lesson.Title, lesson.Type).
		WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(101))

	mock.ExpectQuery(`INSERT INTO TEXT_LESSON`).
		WithArgs(101).
		WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(201))

	mock.ExpectExec(`INSERT INTO text_lesson_block`).
		WithArgs(201, lesson.Value, lesson.IsImage, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.CreateTextLesson(ctx, lesson, bucketId)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateVideoLesson(t *testing.T) {
	db, mock := setupMockDB(t)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	lesson := &coursemodels.LessonPoint{
		Order: 1,
		Title: "Video Intro",
		Type:  "video",
		Value: "https://video.url",
	}
	bucketId := 6

	mock.ExpectQuery(`INSERT INTO LESSON`).
		WithArgs(bucketId, lesson.Order, lesson.Title, lesson.Type).
		WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(101))

	mock.ExpectExec(`INSERT INTO VIDEO_LESSON`).
		WithArgs(101, lesson.Value).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.CreateVideoLesson(ctx, lesson, bucketId)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSendSurveyQuestionAnswer(t *testing.T) {
	db, mock := setupMockDB(t)

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	answer := &coursemodels.SurveyAnswer{
		QuestionId: 1,
		Answer:     3,
	}
	user := &usermodels.UserProfile{Id: 7}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO survey_answer`)).
		WithArgs(answer.QuestionId, answer.Answer, user.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.SendSurveyQuestionAnswer(ctx, answer, user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
