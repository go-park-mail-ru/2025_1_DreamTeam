package postgres

import (
	"context"
	"database/sql"
	"regexp"
	"skillForce/pkg/logs"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetLessonFooters(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	pg := &Database{conn: db}
	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	currentLessonId := 1
	// Mock for initial query to get lesson order and context
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT l.lesson_order, lb.id, lb.lesson_bucket_order, lb.part_id, p.part_order, c.id
		FROM LESSON l
		JOIN LESSON_BUCKET lb ON l.lesson_bucket_id = lb.id
		JOIN PART p ON lb.part_id = p.id
		JOIN COURSE c ON p.course_id = c.id
		WHERE l.id = $1
	`)).WithArgs(currentLessonId).WillReturnRows(sqlmock.NewRows([]string{
		"lesson_order", "bucket_id", "bucket_order", "part_id", "part_order", "course_id",
	}).AddRow(2, 10, 1, 20, 1, 100))

	// Mock for lessons in the current bucket
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, lesson_order
		FROM LESSON
		WHERE lesson_bucket_id = $1
		ORDER BY Lesson_Order ASC
	`)).WithArgs(10).WillReturnRows(sqlmock.NewRows([]string{
		"id", "lesson_order",
	}).AddRow(5, 1).AddRow(1, 2).AddRow(6, 3)) // current: 2, prev: 1, next: 3

	// Mock for logging
	mock.ExpectQuery("SELECT l.id FROM LESSON_BUCKET lb JOIN LESSON l ON l.lesson_bucket_id = lb.id WHERE lb.lesson_bucket_order = .*").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery("SELECT l.id FROM LESSON_BUCKET lb JOIN LESSON l ON l.lesson_bucket_id = lb.id JOIN PART p ON p.id = lb.part_id JOIN COURSE c ON p.course_id = c.id WHERE p.part_order = .*").
		WillReturnError(sql.ErrNoRows)

	footers, err := pg.GetLessonFooters(ctx, currentLessonId)
	assert.NoError(t, err)
	assert.Equal(t, []int{5, 1, 6}, footers)
}

func TestGetLessonHeaderNewCourse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a mock database connection", err)
	}

	d := &Database{conn: db}

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	userId := 1
	courseId := 1

	// Mocking the first query for getting course part
	mock.ExpectQuery(`SELECT title, part_order, id FROM part WHERE course_id = \$1 ORDER BY part_order ASC LIMIT 1`).
		WithArgs(courseId).
		WillReturnRows(sqlmock.NewRows([]string{"title", "part_order", "id"}).
			AddRow("Part 1", 1, 1))

	// Mocking the second query for getting lesson bucket
	mock.ExpectQuery(`SELECT title, lesson_bucket_order, id FROM lesson_bucket WHERE part_id = \$1 ORDER BY lesson_bucket_order ASC LIMIT 1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"title", "lesson_bucket_order", "id"}).
			AddRow("Bucket 1", 1, 1))

	// Mocking the third query for getting lessons
	mock.ExpectQuery(`SELECT id, type FROM LESSON WHERE lesson_bucket_id = \$1 ORDER BY Lesson_Order ASC`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type"}).
			AddRow(1, "video"))

	// Mocking the INSERT INTO LESSON_CHECKPOINT query
	mock.ExpectExec(`INSERT INTO LESSON_CHECKPOINT \(user_id, lesson_id, course_id\) VALUES \(\$1, \$2, \$3\)`).
		WithArgs(userId, 1, courseId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	lessonHeader, currentLessonId, currentLessonType, success, err := d.getLessonHeaderNewCourse(ctx, userId, courseId)

	assert.NoError(t, err)
	assert.True(t, success)
	assert.NotNil(t, lessonHeader)
	assert.Equal(t, "Part 1", lessonHeader.Part.Title)
	assert.Equal(t, "Bucket 1", lessonHeader.Bucket.Title)
	assert.Equal(t, 1, currentLessonId)
	assert.Equal(t, "video", currentLessonType)
}
