package postgres

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	coursemodels "skillForce/internal/models/course"
	"skillForce/pkg/logs"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetBucketCourses_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}

	rows := sqlmock.NewRows([]string{"id", "creator_user_id", "title", "description", "avatar_src", "price", "time_to_pass"}).
		AddRow(1, 101, "Course 1", "Desc 1", "img1.jpg", 100, 10).
		AddRow(2, 102, "Course 2", "Desc 2", "img2.jpg", 200, 15)

	mock.ExpectQuery("SELECT id, creator_user_id, title, description, avatar_src, price, time_to_pass FROM course LIMIT 16").
		WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	courses, err := database.GetBucketCourses(ctx)

	require.NoError(t, err)
	require.Len(t, courses, 2)

	require.Equal(t, 1, courses[0].Id)
	require.Equal(t, 101, courses[0].CreatorId)
	require.Equal(t, "Course 1", courses[0].Title)
	require.Equal(t, "Desc 1", courses[0].Description)
	require.Equal(t, "img1.jpg", courses[0].ScrImage)
	require.Equal(t, 100, courses[0].Price)
	require.Equal(t, 10, courses[0].TimeToPass)

	require.Equal(t, 2, courses[1].Id)
	require.Equal(t, "Course 2", courses[1].Title)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCoursesPurchases_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{
		Data: make([]*logs.LogString, 0),
	})
	course := &coursemodels.Course{Id: 1}
	courses := []*coursemodels.Course{course}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM SIGNUPS WHERE course_id = $1")).
		WithArgs(course.Id).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	database := &Database{conn: db}
	purchases, err := database.GetCoursesPurchases(ctx, courses)

	require.NoError(t, err)
	require.Equal(t, 5, purchases[1])
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCoursesRaitings_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	course := &coursemodels.Course{Id: 1}
	courses := []*coursemodels.Course{course}

	query := regexp.QuoteMeta("SELECT rating FROM course_metrik WHERE course_id = $1")
	mock.ExpectQuery(query).
		WithArgs(course.Id).
		WillReturnRows(sqlmock.NewRows([]string{"rating"}).
			AddRow(4).
			AddRow(5),
		)

	database := &Database{conn: db}

	result, err := database.GetCoursesRaitings(ctx, courses)
	require.NoError(t, err)

	expectedAvg := float32(4.5)
	require.Contains(t, result, course.Id)
	require.InDelta(t, expectedAvg, result[course.Id], 0.01)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCoursesTags_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	course := &coursemodels.Course{Id: 1}
	courses := []*coursemodels.Course{course}

	query := regexp.QuoteMeta(`
			SELECT vt.Title
			FROM TAGS t
			JOIN VALID_TAGS vt ON t.Tag_ID = vt.ID
			WHERE t.Course_ID = $1
		`)

	mock.ExpectQuery(query).
		WithArgs(course.Id).
		WillReturnRows(sqlmock.NewRows([]string{"Title"}).
			AddRow("go").
			AddRow("backend"),
		)

	database := &Database{conn: db}

	result, err := database.GetCoursesTags(ctx, courses)
	require.NoError(t, err)

	expectedTags := []string{"go", "backend"}
	require.Contains(t, result, course.Id)
	require.ElementsMatch(t, expectedTags, result[course.Id])

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCourseById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	courseId := 1

	expectedCourse := coursemodels.Course{
		Id:          courseId,
		CreatorId:   42,
		Title:       "Go Backend",
		Description: "Learn Go backend development",
		ScrImage:    "avatar.png",
		Price:       9900,
		TimeToPass:  120,
	}

	query := regexp.QuoteMeta(`SELECT id, creator_user_id, title, description, avatar_src, price, time_to_pass FROM course WHERE id = $1`)

	mock.ExpectQuery(query).
		WithArgs(courseId).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "creator_user_id", "title", "description", "avatar_src", "price", "time_to_pass",
		}).AddRow(
			expectedCourse.Id,
			expectedCourse.CreatorId,
			expectedCourse.Title,
			expectedCourse.Description,
			expectedCourse.ScrImage,
			expectedCourse.Price,
			expectedCourse.TimeToPass,
		))

	database := &Database{conn: db}

	result, err := database.GetCourseById(ctx, courseId)
	require.NoError(t, err)
	require.Equal(t, &expectedCourse, result)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestMarkLessonAsNotCompleted(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 1
	lessonId := 100

	database := &Database{conn: db}

	t.Run("successfully deleted", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(
			"DELETE FROM LESSON_CHECKPOINT WHERE user_id = $1 AND lesson_id = $2",
		)).
			WithArgs(userId, lessonId).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := database.MarkLessonAsNotCompleted(ctx, userId, lessonId)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("delete returns error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(
			"DELETE FROM LESSON_CHECKPOINT WHERE user_id = $1 AND lesson_id = $2",
		)).
			WithArgs(userId, lessonId).
			WillReturnError(sql.ErrConnDone)

		err := database.MarkLessonAsNotCompleted(ctx, userId, lessonId)
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetLessonHeaderNewCourse_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 1
	courseId := 10
	partId := 100
	bucketId := 200
	lessonId := 300
	lessonType := "text"

	// Настройка ожиданий для запросов
	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT title, part_order, id
        FROM part
        WHERE course_id = $1
        ORDER BY part_order ASC
        LIMIT 1;
    `)).
		WithArgs(courseId).
		WillReturnRows(sqlmock.NewRows([]string{"title", "part_order", "id"}).
			AddRow("Part Title", 1, partId))

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT title, lesson_bucket_order, id
        FROM lesson_bucket
        WHERE part_id = $1
        ORDER BY lesson_bucket_order ASC
        LIMIT 1;
    `)).
		WithArgs(partId).
		WillReturnRows(sqlmock.NewRows([]string{"title", "lesson_bucket_order", "id"}).
			AddRow("Bucket Title", 1, bucketId))

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT id, type
        FROM LESSON
        WHERE lesson_bucket_id = $1
        ORDER BY Lesson_Order ASC
    `)).
		WithArgs(bucketId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type"}).
			AddRow(lessonId, lessonType))

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT EXISTS (SELECT 1 FROM LESSON_CHECKPOINT WHERE user_id = $1 AND lesson_id = $2 AND course_id = $3)")).
		WithArgs(userId, lessonId, courseId). // Используем реальные значения
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO LESSON_CHECKPOINT (user_id, lesson_id, course_id) VALUES ($1, $2, $3)")).
		WithArgs(userId, lessonId, courseId). // Используем реальные значения
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Вызов тестируемого метода
	header, gotLessonId, gotType, gotIsNew, err := database.getLessonHeaderNewCourse(ctx, userId, courseId)

	// Проверки
	require.NoError(t, err)
	require.NotNil(t, header)
	require.Equal(t, lessonId, gotLessonId)
	require.Equal(t, lessonType, gotType)
	require.True(t, gotIsNew)
	require.Len(t, header.Points, 1)
	require.Equal(t, "Part Title", header.Part.Title)
	require.Equal(t, "Bucket Title", header.Bucket.Title)
}

func TestGetLastLessonHeader_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	database := &Database{conn: db}

	userId := 1
	courseId := 2

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, creator_user_id, title, description, avatar_src, price, time_to_pass FROM course WHERE id = $1")).
		WithArgs(courseId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "creator_user_id", "title", "description", "avatar_src", "price", "time_to_pass"}).
			AddRow(courseId, 1, "Course Title", "Desc", "img.png", 100, 120))

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT cp.Lesson_ID, l.type
		FROM LESSON_CHECKPOINT cp
		JOIN LESSON l ON l.ID = cp.Lesson_ID
		WHERE cp.User_ID = $1 AND cp.Course_ID = $2
		ORDER BY cp.Updated_at DESC
	`)).
		WithArgs(userId, courseId).
		WillReturnRows(sqlmock.NewRows([]string{"Lesson_ID", "type"}).
			AddRow(10, "video").
			AddRow(11, "quiz"))

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT p.Title, p.Part_Order, p.ID, lb.Title, lb.Lesson_Bucket_Order, lb.ID
		FROM PART p
		JOIN LESSON_BUCKET lb ON lb.Part_ID = p.ID
		JOIN LESSON l ON l.Lesson_Bucket_ID = lb.ID
		WHERE l.ID = $1
	`)).
		WithArgs(10).
		WillReturnRows(sqlmock.NewRows([]string{
			"Title", "Part_Order", "ID", "Title", "Lesson_Bucket_Order", "ID"}).
			AddRow("Part 1", 1, 100, "Bucket A", 1, 200))

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, type
		FROM LESSON
		WHERE lesson_bucket_id = $1
		ORDER BY Lesson_Order ASC
	`)).
		WithArgs(200).
		WillReturnRows(sqlmock.NewRows([]string{"id", "type"}).
			AddRow(10, "video").
			AddRow(11, "quiz").
			AddRow(12, "text"))

	header, lessonId, lessonType, isNew, err := database.GetLastLessonHeader(ctx, userId, courseId)
	require.NoError(t, err)
	require.NotNil(t, header)
	require.Equal(t, 10, lessonId)
	require.Equal(t, "video", lessonType)
	require.False(t, isNew)

	require.Equal(t, "Course Title", header.CourseTitle)
	require.Equal(t, 3, len(header.Points))
	require.True(t, header.Points[0].IsDone)
	require.True(t, header.Points[1].IsDone)
	require.False(t, header.Points[2].IsDone)

	require.NoError(t, mock.ExpectationsWereMet())
}

// func TestGetLessonHeaderByLessonId_Success(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	require.NoError(t, err)
// 	defer  db.Close()

// 	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
// 	database := &Database{conn: db}

// 	userId := 1
// 	currentLessonId := 10

// 	mock.ExpectQuery(regexp.QuoteMeta(`
// 		SELECT p.Title, p.Part_Order, p.ID, lb.Title, lb.Lesson_Bucket_Order, lb.ID, c.ID, c.Title
// 		FROM lesson l
// 		JOIN LESSON_BUCKET lb ON l.LESSON_BUCKET_ID = lb.ID
// 		JOIN PART p ON lb.PART_ID = p.ID
// 		JOIN COURSE c ON p.COURSE_ID = c.ID
// 		WHERE l.ID = $1
// 	`)).
// 		WithArgs(currentLessonId).
// 		WillReturnRows(sqlmock.NewRows([]string{
// 			"p.Title", "p.Part_Order", "p.ID",
// 			"lb.Title", "lb.Lesson_Bucket_Order", "lb.ID",
// 			"c.ID", "c.Title"}).
// 			AddRow("Part A", 1, 101, "Bucket A", 1, 201, 301, "Course Title"))

// 	mock.ExpectQuery(regexp.QuoteMeta(`
// 		SELECT lesson_id
// 		FROM LESSON_CHECKPOINT
// 		WHERE course_id = $1 and user_id = $2
// 	`)).
// 		WithArgs(301, userId).
// 		WillReturnRows(sqlmock.NewRows([]string{"lesson_id"}).
// 			AddRow(10).
// 			AddRow(11))

// 	mock.ExpectQuery(regexp.QuoteMeta(`
// 		SELECT id, type
// 		FROM LESSON
// 		WHERE lesson_bucket_id = $1
// 	`)).
// 		WithArgs(201).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "type"}).
// 			AddRow(10, "video").
// 			AddRow(11, "quiz").
// 			AddRow(12, "text"))

// 	header, err := database.GetLessonHeaderByLessonId(ctx, userId, currentLessonId)
// 	require.NoError(t, err)
// 	require.NotNil(t, header)

// 	require.Equal(t, 301, header.CourseId)
// 	require.Equal(t, "Course Title", header.CourseTitle)
// 	require.Equal(t, 3, len(header.Points))

// 	require.Equal(t, 10, header.Points[0].LessonId)
// 	require.Equal(t, true, header.Points[0].IsDone)

// 	require.Equal(t, 11, header.Points[1].LessonId)
// 	require.Equal(t, true, header.Points[1].IsDone)

// 	require.Equal(t, 12, header.Points[2].LessonId)
// 	require.Equal(t, false, header.Points[2].IsDone)

// 	require.NoError(t, mock.ExpectationsWereMet())
// }

func TestGetLessonBlocks_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	lessonID := 123

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT tlb.value
		FROM TEXT_LESSON_BLOCK tlb
		JOIN TEXT_LESSON tl ON tlb.Text_Lesson_ID = tl.ID
		WHERE tl.Lesson_ID = $1
		ORDER BY tlb.Text_Lesson_Block_Order ASC
	`)).
		WithArgs(lessonID).
		WillReturnRows(sqlmock.NewRows([]string{"value"}).
			AddRow("Block 1").
			AddRow("Block 2").
			AddRow("Block 3"))

	blocks, err := database.GetLessonBlocks(ctx, lessonID)
	require.NoError(t, err)
	require.Len(t, blocks, 3)
	require.Equal(t, []string{"Block 1", "Block 2", "Block 3"}, blocks)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLessonVideo_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	lessonID := 123
	expectedVideoSrc := "https://example.com/video.mp4"

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT video_src
		FROM VIDEO_LESSON
		WHERE lesson_ID = $1
	`)).
		WithArgs(lessonID).
		WillReturnRows(sqlmock.NewRows([]string{"video_src"}).
			AddRow(expectedVideoSrc))

	video, err := database.GetLessonVideo(ctx, lessonID)

	require.NoError(t, err)
	require.Equal(t, []string{expectedVideoSrc}, video)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLessonVideo_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	lessonID := 123

	mock.ExpectQuery("SELECT video_src").
		WithArgs(lessonID).
		WillReturnError(errors.New("query failed"))

	video, err := database.GetLessonVideo(ctx, lessonID)

	require.Error(t, err)
	require.Nil(t, video)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLessonById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	lessonID := 123
	expectedLesson := &coursemodels.LessonPoint{
		LessonId: lessonID,
		Title:    "Test Lesson",
		Type:     "video",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, title, type
		FROM LESSON
		WHERE id = $1
	`)).
		WithArgs(lessonID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "type"}).
			AddRow(expectedLesson.LessonId, expectedLesson.Title, expectedLesson.Type))

	lesson, err := database.GetLessonById(ctx, lessonID)

	require.NoError(t, err)
	require.Equal(t, expectedLesson, lesson)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLessonById_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	lessonID := 123

	mock.ExpectQuery("SELECT id, title, type").
		WithArgs(lessonID).
		WillReturnError(errors.New("query failed"))

	lesson, err := database.GetLessonById(ctx, lessonID)

	require.Error(t, err)
	require.Nil(t, lesson)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBucketByLessonId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	lessonID := 123
	expectedBucket := &coursemodels.LessonBucket{
		Id: 456,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT lesson_bucket_id
		FROM LESSON
		WHERE id = $1
	`)).
		WithArgs(lessonID).
		WillReturnRows(sqlmock.NewRows([]string{"lesson_bucket_id"}).
			AddRow(expectedBucket.Id))

	bucket, err := database.GetBucketByLessonId(ctx, lessonID)

	require.NoError(t, err)
	require.Equal(t, expectedBucket, bucket)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBucketByLessonId_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	lessonID := 123

	mock.ExpectQuery("SELECT lesson_bucket_id").
		WithArgs(lessonID).
		WillReturnError(errors.New("query failed"))

	bucket, err := database.GetBucketByLessonId(ctx, lessonID)

	require.Error(t, err)
	require.Nil(t, bucket)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBucketByLessonId_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	lessonID := 123

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT lesson_bucket_id
		FROM LESSON
		WHERE id = $1
	`)).
		WithArgs(lessonID).
		WillReturnRows(sqlmock.NewRows([]string{"lesson_bucket_id"}))

	bucket, err := database.GetBucketByLessonId(ctx, lessonID)

	require.Error(t, err)
	require.Nil(t, bucket)
	require.NoError(t, mock.ExpectationsWereMet())
}

// func TestGetLessonFooters_Success(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	require.NoError(t, err)
// 	defer db.Close()

// 	database := &Database{conn: db}
// 	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
// 	lessonID := 123
// 	expectedFooters := []int{100, 123, 124}

// 	mock.ExpectQuery(regexp.QuoteMeta(`
// 		SELECT l.lesson_order, lb.id, lb.lesson_bucket_order
// 		FROM LESSON l
// 		JOIN LESSON_BUCKET lb ON l.lesson_bucket_id = lb.id
// 		WHERE l.id = $1
// 	`)).
// 		WithArgs(lessonID).
// 		WillReturnRows(sqlmock.NewRows([]string{"lesson_order", "id", "lesson_bucket_order"}).
// 			AddRow(2, 5, 1))

// 	mock.ExpectQuery(regexp.QuoteMeta(`
// 		SELECT id, lesson_order
// 		FROM LESSON
// 		WHERE lesson_bucket_id = $1
// 		ORDER BY Lesson_Order ASC
// 	`)).
// 		WithArgs(5).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "lesson_order"}).
// 			AddRow(100, 1).
// 			AddRow(123, 2).
// 			AddRow(124, 3))

// 	footers, err := database.GetLessonFooters(ctx, lessonID)

// 	require.NoError(t, err)
// 	require.Equal(t, expectedFooters, footers)
// 	require.NoError(t, mock.ExpectationsWereMet())
// }

// func TestGetLessonFooters_QueryErrorFirst(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	require.NoError(t, err)
// 	defer db.Close()

// 	database := &Database{conn: db}
// 	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
// 	lessonID := 123

// 	mock.ExpectQuery(regexp.QuoteMeta(`
// 		SELECT l.lesson_order, lb.id, lb.lesson_bucket_order
// 		FROM LESSON l
// 		JOIN LESSON_BUCKET lb ON l.lesson_bucket_id = lb.id
// 		WHERE l.id = $1
// 	`)).
// 		WithArgs(lessonID).
// 		WillReturnError(errors.New("query failed"))

// 	footers, err := database.GetLessonFooters(ctx, lessonID)

// 	require.Error(t, err)
// 	require.Nil(t, footers)
// 	require.NoError(t, mock.ExpectationsWereMet())
// }

// func TestGetLessonFooters_QueryErrorSecond(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	require.NoError(t, err)
// 	defer db.Close()

// 	database := &Database{conn: db}
// 	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
// 	lessonID := 123

// 	mock.ExpectQuery(regexp.QuoteMeta(`
// 		SELECT l.lesson_order, lb.id, lb.lesson_bucket_order
// 		FROM LESSON l
// 		JOIN LESSON_BUCKET lb ON l.lesson_bucket_id = lb.id
// 		WHERE l.id = $1
// 	`)).
// 		WithArgs(lessonID).
// 		WillReturnRows(sqlmock.NewRows([]string{"lesson_order", "id", "lesson_bucket_order"}).
// 			AddRow(2, 5, 1))

// 	mock.ExpectQuery(regexp.QuoteMeta(`
// 		SELECT id, lesson_order
// 		FROM LESSON
// 		WHERE lesson_bucket_id = $1
// 		ORDER BY Lesson_Order ASC
// 	`)).
// 		WithArgs(5).
// 		WillReturnError(errors.New("query failed"))

// 	footers, err := database.GetLessonFooters(ctx, lessonID)

// 	require.Error(t, err)
// 	require.Nil(t, footers)
// 	require.NoError(t, mock.ExpectationsWereMet())
// }

// func TestGetLessonFooters_NoLessonsInBucket(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	require.NoError(t, err)
// 	defer db.Close()

// 	database := &Database{conn: db}
// 	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
// 	lessonID := 123

// 	mock.ExpectQuery(regexp.QuoteMeta(`
// 		SELECT l.lesson_order, lb.id, lb.lesson_bucket_order
// 		FROM LESSON l
// 		JOIN LESSON_BUCKET lb ON l.lesson_bucket_id = lb.id
// 		WHERE l.id = $1
// 	`)).
// 		WithArgs(lessonID).
// 		WillReturnRows(sqlmock.NewRows([]string{"lesson_order", "id", "lesson_bucket_order"}).
// 			AddRow(2, 5, 1))

// 	mock.ExpectQuery(regexp.QuoteMeta(`
// 		SELECT id, lesson_order
// 		FROM LESSON
// 		WHERE lesson_bucket_id = $1
// 		ORDER BY Lesson_Order ASC
// 	`)).
// 		WithArgs(5).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "lesson_order"}))

// 	footers, err := database.GetLessonFooters(ctx, lessonID)

// 	require.Error(t, err)
// 	require.Nil(t, footers)
// 	require.NoError(t, mock.ExpectationsWereMet())
// }

func TestGetCourseParts_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	courseID := 1
	expectedParts := []*coursemodels.CoursePart{
		{Id: 1, Title: "Part 1"},
		{Id: 2, Title: "Part 2"},
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, title
		FROM PART
		WHERE course_id = $1
		ORDER BY part_order ASC
	`)).
		WithArgs(courseID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).
			AddRow(1, "Part 1").
			AddRow(2, "Part 2"))

	parts, err := database.GetCourseParts(ctx, courseID)

	require.NoError(t, err)
	require.Equal(t, expectedParts, parts)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCourseParts_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	courseID := 1

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, title
		FROM PART
		WHERE course_id = $1
		ORDER BY part_order ASC
	`)).
		WithArgs(courseID).
		WillReturnError(errors.New("query failed"))

	parts, err := database.GetCourseParts(ctx, courseID)

	require.Error(t, err)
	require.Nil(t, parts)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCourseParts_NoParts(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	courseID := 1

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, title
		FROM PART
		WHERE course_id = $1
		ORDER BY part_order ASC
	`)).
		WithArgs(courseID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))

	parts, err := database.GetCourseParts(ctx, courseID)

	require.NoError(t, err)
	require.Empty(t, parts)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPartBuckets_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	partID := 1
	expectedBuckets := []*coursemodels.LessonBucket{
		{Id: 1, Title: "Bucket 1"},
		{Id: 2, Title: "Bucket 2"},
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, title
		FROM LESSON_BUCKET
		WHERE part_id = $1
		ORDER BY lesson_bucket_order ASC
	`)).
		WithArgs(partID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).
			AddRow(1, "Bucket 1").
			AddRow(2, "Bucket 2"))

	buckets, err := database.GetPartBuckets(ctx, partID)

	require.NoError(t, err)
	require.Equal(t, expectedBuckets, buckets)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPartBuckets_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	partID := 1

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, title
		FROM LESSON_BUCKET
		WHERE part_id = $1
		ORDER BY lesson_bucket_order ASC
	`)).
		WithArgs(partID).
		WillReturnError(errors.New("query failed"))

	buckets, err := database.GetPartBuckets(ctx, partID)

	require.Error(t, err)
	require.Nil(t, buckets)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBucketLessons_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 1
	courseId := 1
	bucketId := 10

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT lesson_id
		FROM LESSON_CHECKPOINT
		WHERE user_id = $1 AND course_id = $2
	`)).
		WithArgs(userId, courseId).
		WillReturnRows(sqlmock.NewRows([]string{"lesson_id"}).
			AddRow(101))

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, title, type
		FROM LESSON
		WHERE lesson_bucket_id = $1
		ORDER BY lesson_order ASC
	`)).
		WithArgs(bucketId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "type"}).
			AddRow(101, "Lesson 1", "text").
			AddRow(102, "Lesson 2", "video"))

	lessons, err := database.GetBucketLessons(ctx, userId, courseId, bucketId)
	require.NoError(t, err)
	require.Len(t, lessons, 2)

	require.True(t, lessons[0].IsDone)
	require.Equal(t, "Lesson 1", lessons[0].Title)

	require.False(t, lessons[1].IsDone)
	require.Equal(t, "Lesson 2", lessons[1].Title)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBucketLessons_ErrorOnCompletedLessonsQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	mock.ExpectQuery("SELECT lesson_id FROM LESSON_CHECKPOINT").
		WithArgs(1, 1).
		WillReturnError(errors.New("db error"))

	lessons, err := database.GetBucketLessons(ctx, 1, 1, 1)
	require.Error(t, err)
	require.Nil(t, lessons)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBucketLessons_ErrorOnLessonQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	mock.ExpectQuery("SELECT lesson_id FROM LESSON_CHECKPOINT").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"lesson_id"}))

	mock.ExpectQuery("SELECT id, title, type FROM LESSON").
		WithArgs(1).
		WillReturnError(errors.New("query failed"))

	lessons, err := database.GetBucketLessons(ctx, 1, 1, 1)
	require.Error(t, err)
	require.Nil(t, lessons)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddUserToCourse_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	userId := 1
	courseId := 10

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT EXISTS (SELECT 1 FROM SIGNUPS WHERE user_id = $1 AND course_id = $2)
	`)).
		WithArgs(userId, courseId).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO SIGNUPS (user_id, course_id) VALUES ($1, $2)
	`)).
		WithArgs(userId, courseId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = database.AddUserToCourse(ctx, userId, courseId)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddUserToCourse_AlreadyExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(1, 10).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	err = database.AddUserToCourse(ctx, 1, 10)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddUserToCourse_QueryRowError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(1, 10).
		WillReturnError(errors.New("db error"))

	mock.ExpectExec("INSERT INTO SIGNUPS").
		WithArgs(1, 10).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = database.AddUserToCourse(ctx, 1, 10)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetVideoUrl_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})
	lessonId := 1
	expectedURL := "https://example.com/video.mp4"

	mock.ExpectQuery("SELECT video_src FROM video_lesson WHERE lesson_id = \\$1").
		WithArgs(lessonId).
		WillReturnRows(sqlmock.NewRows([]string{"video_src"}).AddRow(expectedURL))

	url, err := database.GetVideoUrl(ctx, lessonId)
	require.NoError(t, err)
	require.Equal(t, expectedURL, url)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetVideoUrl_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	mock.ExpectQuery("SELECT video_src FROM video_lesson WHERE lesson_id = \\$1").
		WithArgs(1).
		WillReturnError(errors.New(("query failed")))

	url, err := database.GetVideoUrl(ctx, 1)
	require.Error(t, err)
	require.Empty(t, url)
	require.EqualError(t, err, "query failed")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestIsUserPurchasedCourse_ExistsTrue(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.Background()

	mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM SIGNUPS WHERE user_id = \\$1 AND course_id = \\$2\\)").
		WithArgs(1, 10).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	ok, err := database.IsUserPurchasedCourse(ctx, 1, 10)
	require.NoError(t, err)
	require.True(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestIsUserPurchasedCourse_ExistsFalse(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.Background()

	mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM SIGNUPS WHERE user_id = \\$1 AND course_id = \\$2\\)").
		WithArgs(1, 10).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	ok, err := database.IsUserPurchasedCourse(ctx, 1, 10)
	require.NoError(t, err)
	require.False(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestIsUserPurchasedCourse_ErrNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM SIGNUPS WHERE user_id = \\$1 AND course_id = \\$2\\)").
		WithArgs(1, 10).
		WillReturnError(sql.ErrNoRows)

	ok, err := database.IsUserPurchasedCourse(ctx, 1, 10)
	require.NoError(t, err)
	require.False(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestIsUserPurchasedCourse_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	database := &Database{conn: db}
	ctx := context.WithValue(context.Background(), logs.LogsKey, &logs.CtxLog{})

	mock.ExpectQuery("SELECT EXISTS \\(SELECT 1 FROM SIGNUPS WHERE user_id = \\$1 AND course_id = \\$2\\)").
		WithArgs(1, 10).
		WillReturnError(errors.New("query error"))

	ok, err := database.IsUserPurchasedCourse(ctx, 1, 10)
	require.Error(t, err)
	require.False(t, ok)
	require.EqualError(t, err, "query error")
	require.NoError(t, mock.ExpectationsWereMet())
}
