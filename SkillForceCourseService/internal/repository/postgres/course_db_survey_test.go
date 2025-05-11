package postgres

import (
	"context"
	"testing"

	coursemodels "skillForce/internal/models/course"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/logs"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateSurvey(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	d := &Database{conn: db}

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	// Mock survey data
	survey := &coursemodels.Survey{
		Questions: []coursemodels.Question{
			{Question: "What is your favorite color?", LeftLebal: "Red", RightLebal: "Blue", Metric: "color"},
		},
	}

	// Mock SQL queries
	mock.ExpectQuery("INSERT INTO survey DEFAULT VALUES RETURNING id").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec("INSERT INTO survey_question").
		WithArgs(1, survey.Questions[0].Question, survey.Questions[0].LeftLebal, survey.Questions[0].RightLebal, survey.Questions[0].Metric).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the function
	err = d.CreateSurvey(ctx, survey, &usermodels.UserProfile{})
	assert.NoError(t, err)

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSurvey(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	d := &Database{conn: db}

	// Mock SQL queries
	mock.ExpectQuery("SELECT id FROM survey ORDER BY id DESC LIMIT 1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	mock.ExpectQuery("SELECT id, metric_type, question, left_desc, right_desc").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "metric_type", "question", "left_desc", "right_desc"}).
			AddRow(1, "color", "What is your favorite color?", "Red", "Blue"))

	// Call the function
	survey, err := d.GetSurvey(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, survey)
	assert.Equal(t, 1, survey.Id)
	assert.Equal(t, 1, len(survey.Questions))
	assert.Equal(t, "What is your favorite color?", survey.Questions[0].Question)

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMetricCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	d := &Database{conn: db}

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	// Mock SQL query
	mock.ExpectQuery("SELECT COUNT\\(sq.id\\) FROM survey_question sq").
		WithArgs(1, "color").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	// Call the function
	count, err := d.GetMetricCount(ctx, 1, "color")
	assert.NoError(t, err)
	assert.Equal(t, 5, count)

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMetricAvg(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	d := &Database{conn: db}

	// Mock SQL query
	mock.ExpectQuery("SELECT COALESCE\\(AVG\\(sa.answer\\), 0.0\\) AS avg").
		WithArgs(1, "color").
		WillReturnRows(sqlmock.NewRows([]string{"avg"}).AddRow(4.5))

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	// Call the function
	avg, err := d.GetMetricAvg(ctx, 1, "color")
	assert.NoError(t, err)
	assert.Equal(t, 4.5, avg)

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMetricDistribution(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	d := &Database{conn: db}

	// Mock SQL queries
	mock.ExpectQuery("SELECT COUNT\\(sa.answer\\) FROM survey_question sq").
		WithArgs(1, "color").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(100))

	for i := 0; i < 11; i++ {
		mock.ExpectQuery("SELECT COUNT\\(sa.answer\\) FROM survey_question sq").
			WithArgs(1, "color", i).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	// Call the function
	distribution, err := d.GetMetricDistribution(ctx, 1, "color")
	assert.NoError(t, err)
	assert.Equal(t, 11, len(distribution))
	assert.Equal(t, 10, distribution[0]) // Example check

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMetricAnswers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	d := &Database{conn: db}

	ctx := context.Background()
	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})

	// Mock SQL query
	mock.ExpectQuery("SELECT sa.answer, u.name FROM survey_question sq").
		WithArgs(1, "color").
		WillReturnRows(sqlmock.NewRows([]string{"answer", "name"}).
			AddRow(5, "John").
			AddRow(3, "Jane"))

	// Call the function
	answers, err := d.GetMetricAnswers(ctx, 1, "color")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(answers))
	assert.Equal(t, "John", answers[0].Username)
	assert.Equal(t, 5, answers[0].Answer)

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

// func TestGetMetrics(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	d := &Database{conn: db}

// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, logs.LogsKey, &logs.CtxLog{
// 		Data: make([]*logs.LogString, 0),
// 	})

// 	// Mock SQL queries
// 	mock.ExpectQuery("SELECT id FROM survey ORDER BY id DESC LIMIT 1").
// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

// 	mock.ExpectQuery("SELECT COUNT\\(sq.id\\) FROM survey_question sq").
// 		WithArgs(1, "color").
// 		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

// 	mock.ExpectQuery("SELECT COALESCE\\(AVG\\(sa.answer\\), 0.0\\) AS avg").
// 		WithArgs(1, "color").
// 		WillReturnRows(sqlmock.NewRows([]string{"avg"}).AddRow(4.5))

// 	for i := 0; i < 11; i++ {
// 		mock.ExpectQuery("SELECT COUNT\\(sa.answer\\) FROM survey_question sq").
// 			WithArgs(1, "color", i).
// 			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
// 	}

// 	mock.ExpectQuery("SELECT sa.answer, u.name FROM survey_question sq").
// 		WithArgs(1, "color").
// 		WillReturnRows(sqlmock.NewRows([]string{"answer", "name"}).
// 			AddRow(5, "John").
// 			AddRow(3, "Jane"))

// 	// Call the function
// 	metrics, err := d.GetMetrics(ctx, "color")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, metrics)
// 	assert.Equal(t, "color", metrics.Type)
// 	assert.Equal(t, 5, metrics.Count)
// 	assert.Equal(t, 4.5, metrics.Avg)
// 	assert.Equal(t, 11, len(metrics.Distribution))
// 	assert.Equal(t, 2, len(metrics.Answers))

// 	// Verify expectations
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }
