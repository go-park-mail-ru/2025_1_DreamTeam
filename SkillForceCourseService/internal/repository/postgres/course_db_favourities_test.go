package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"

	coursemodels "skillForce/internal/models/course"
	"skillForce/pkg/logs"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) (*Database, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	return &Database{conn: db}, mock, func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close mock database: %v", err)
		}
	}
}

func TestAddCourseToFavourites_Success(t *testing.T) {
	db, mock, close := setupDB(t)
	defer close()

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO FAVOURITE_COURSES (user_id, course_id)
		VALUES ($1, $2)	
	`)).
		WithArgs(1, 10).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.AddCourseToFavourites(ctx, 10, 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddCourseToFavourites_Error(t *testing.T) {
	db, mock, close := setupDB(t)
	defer close()

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	mock.ExpectExec("INSERT INTO FAVOURITE_COURSES").
		WithArgs(1, 10).
		WillReturnError(errors.New("insert failed"))

	err := db.AddCourseToFavourites(ctx, 10, 1)
	assert.Error(t, err)
}

func TestDeleteCourseFromFavourites_Success(t *testing.T) {
	db, mock, close := setupDB(t)
	defer close()

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	mock.ExpectExec("DELETE FROM FAVOURITE_COURSES").
		WithArgs(10, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := db.DeleteCourseFromFavourites(ctx, 10, 1)
	assert.NoError(t, err)
}

func TestDeleteCourseFromFavourites_Error(t *testing.T) {
	db, mock, close := setupDB(t)
	defer close()

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	mock.ExpectExec("DELETE FROM FAVOURITE_COURSES").
		WithArgs(10, 1).
		WillReturnError(errors.New("delete failed"))

	err := db.DeleteCourseFromFavourites(ctx, 10, 1)
	assert.Error(t, err)
}

func TestGetFavouriteCourses_Success(t *testing.T) {
	db, mock, close := setupDB(t)
	defer close()

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	rows := sqlmock.NewRows([]string{"id", "creator_user_id", "title", "description", "avatar_src", "price", "time_to_pass"}).
		AddRow(1, 100, "Course 1", "Desc", "img.png", 500, 30).
		AddRow(2, 101, "Course 2", "Desc", "img2.png", 600, 45)

	mock.ExpectQuery("SELECT c.id, c.creator_user_id").
		WithArgs(1).
		WillReturnRows(rows)

	courses, err := db.GetFavouriteCourses(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, courses, 2)
	assert.Equal(t, "Course 1", courses[0].Title)
	assert.Equal(t, 101, courses[1].CreatorId)
}

func TestGetFavouriteCourses_QueryError(t *testing.T) {
	db, mock, close := setupDB(t)
	defer close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	mock.ExpectQuery("SELECT c.id, c.creator_user_id").
		WithArgs(1).
		WillReturnError(errors.New("query failed"))

	courses, err := db.GetFavouriteCourses(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, courses)
}

func TestGetFavouriteCourses_ScanError(t *testing.T) {
	db, mock, close := setupDB(t)
	defer close()

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	rows := sqlmock.NewRows([]string{"id", "creator_user_id", "title", "description", "avatar_src", "price", "time_to_pass"}).
		AddRow("bad", 100, "Course", "Desc", "img", 500, 30)

	mock.ExpectQuery("SELECT c.id, c.creator_user_id").
		WithArgs(1).
		WillReturnRows(rows)

	courses, err := db.GetFavouriteCourses(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, courses)
}

func TestGetCoursesFavouriteStatus_Success(t *testing.T) {
	db, mock, close := setupDB(t)
	defer close()

	input := []*coursemodels.Course{
		{Id: 1}, {Id: 2}, {Id: 3},
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	rows := sqlmock.NewRows([]string{"course_id"}).
		AddRow(1).
		AddRow(3)

	mock.ExpectQuery("SELECT course_id FROM FAVOURITE_COURSES").
		WithArgs(5).
		WillReturnRows(rows)

	result, err := db.GetCoursesFavouriteStatus(ctx, input, 5)
	assert.NoError(t, err)
	assert.True(t, result[1])
	assert.False(t, result[2])
	assert.True(t, result[3])
}

func TestGetCoursesFavouriteStatus_Error(t *testing.T) {
	db, mock, close := setupDB(t)
	defer close()

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	mock.ExpectQuery("SELECT course_id FROM FAVOURITE_COURSES").
		WithArgs(5).
		WillReturnError(errors.New("query error"))

	_, err := db.GetCoursesFavouriteStatus(ctx, []*coursemodels.Course{{Id: 1}}, 5)
	assert.Error(t, err)
}
