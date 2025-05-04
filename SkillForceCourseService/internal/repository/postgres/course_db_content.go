package postgres

import (
	"context"
	"database/sql"
	"fmt"

	coursemodels "skillForce/internal/models/course"
	"skillForce/internal/models/dto"
	"skillForce/pkg/logs"
)

func (d *Database) GetBucketCourses(ctx context.Context) ([]*coursemodels.Course, error) {
	//TODO: можно заморочиться и сделать самописную пагинацию через LIMIT OFFSET
	var bucketCourses []*coursemodels.Course
	rows, err := d.conn.Query("SELECT id, creator_user_id, title, description, avatar_src, price, time_to_pass FROM course LIMIT 16")
	if err != nil {
		logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var course coursemodels.Course
		if err := rows.Scan(&course.Id, &course.CreatorId, &course.Title, &course.Description, &course.ScrImage, &course.Price, &course.TimeToPass); err != nil {
			logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
			return nil, err
		}
		logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("get course %+v from db", course))
		bucketCourses = append(bucketCourses, &course)
	}

	logs.PrintLog(ctx, "GetBucketCourses", "get bucket ourses from db")

	return bucketCourses, nil
}

func (d *Database) GetBucketByLessonId(ctx context.Context, currentLessonId int) (*coursemodels.LessonBucket, error) {
	var bucketId int
	err := d.conn.QueryRow(`
			SELECT lesson_bucket_id
			FROM LESSON
			WHERE id = $1
		`, currentLessonId).Scan(&bucketId)
	if err != nil {
		logs.PrintLog(ctx, "GetBucketByLessonId", fmt.Sprintf("%+v", err))
		return nil, err
	}

	bucket := &coursemodels.LessonBucket{
		Id: bucketId,
	}
	return bucket, nil
}

func (d *Database) GetBucketLessons(ctx context.Context, userId int, courseId int, bucketId int) ([]*coursemodels.LessonPoint, error) {
	completedLessons := make(map[int]bool)
	rows1, err := d.conn.Query(`
			SELECT lesson_id
			FROM LESSON_CHECKPOINT
			WHERE user_id = $1 AND course_id = $2
		`, userId, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetBucketLessons", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer rows1.Close()

	for rows1.Next() {
		var completedLessonId int
		if err := rows1.Scan(&completedLessonId); err != nil {
			logs.PrintLog(ctx, "GetBucketLessons", fmt.Sprintf("%+v", err))
			return nil, err
		}
		completedLessons[completedLessonId] = true
	}

	var lessons []*coursemodels.LessonPoint
	rows2, err := d.conn.Query(`
			SELECT id, title, type
			FROM LESSON
			WHERE lesson_bucket_id = $1
			ORDER BY lesson_order ASC
		`, bucketId)
	if err != nil {
		logs.PrintLog(ctx, "GetBucketLessons", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer rows2.Close()

	for rows2.Next() {
		var lesson coursemodels.LessonPoint
		if err := rows2.Scan(&lesson.LessonId, &lesson.Title, &lesson.Type); err != nil {
			logs.PrintLog(ctx, "GetBucketLessons", fmt.Sprintf("%+v", err))
			return nil, err
		}

		if _, ok := completedLessons[lesson.LessonId]; ok {
			lesson.IsDone = true
		}

		logs.PrintLog(ctx, "GetBucketLessons", fmt.Sprintf("get lesson %+v", lesson))
		lessons = append(lessons, &lesson)
	}
	return lessons, nil
}

func (d *Database) GetCourseById(ctx context.Context, courseId int) (*coursemodels.Course, error) {
	var course coursemodels.Course
	err := d.conn.QueryRow("SELECT id, creator_user_id, title, description, avatar_src, price, time_to_pass FROM course WHERE id = $1", courseId).Scan(
		&course.Id, &course.CreatorId, &course.Title, &course.Description, &course.ScrImage, &course.Price, &course.TimeToPass)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseById", fmt.Sprintf("%+v", err))
		return nil, err
	}
	logs.PrintLog(ctx, "GetCourseById", fmt.Sprintf("get course %+v from db", course))
	return &course, nil
}

func (d *Database) GetLessonBlocks(ctx context.Context, currentLessonId int) ([]string, error) {
	var blocks []string
	rows, err := d.conn.Query(`
			SELECT tlb.value
			FROM TEXT_LESSON_BLOCK tlb
			JOIN TEXT_LESSON tl ON tlb.Text_Lesson_ID = tl.ID
			WHERE tl.Lesson_ID = $1
			ORDER BY tlb.Text_Lesson_Block_Order ASC
		`, currentLessonId)
	if err != nil {
		logs.PrintLog(ctx, "GetLessonBlocks", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var block string
		if err := rows.Scan(&block); err != nil {
			logs.PrintLog(ctx, "GetLessonBlocks", fmt.Sprintf("%+v", err))
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func (d *Database) GetLessonVideo(ctx context.Context, currentLessonId int) ([]string, error) {
	var videoSrc string
	err := d.conn.QueryRow(`
			SELECT video_src
			FROM VIDEO_LESSON
			WHERE lesson_ID = $1
			`, currentLessonId).Scan(&videoSrc)
	if err != nil {
		logs.PrintLog(ctx, "GetLessonVideos", fmt.Sprintf("%+v", err))
		return nil, err
	}
	return []string{videoSrc}, nil
}

func (d *Database) GetLessonTest(ctx context.Context, currentLessonId int, user_id int) (*dto.Test, error) {
	query := `
        SELECT 
            qt.id AS question_id,
            qt.question,
            av.id AS answer_id,
            av.answer,
			av.is_true
        FROM test_lesson tl
        JOIN quiz_task qt ON qt.lesson_test_id = tl.id
        JOIN answer_variant av ON av.quiz_task_id = qt.id
        WHERE tl.lesson_id = $1
        AND qt.id = (
            SELECT id FROM quiz_task 
            WHERE lesson_test_id = tl.id 
            ORDER BY id LIMIT 1
        )
        ORDER BY av.id;
    `

	rows, err := d.conn.Query(query, currentLessonId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var q *dto.Test

	for rows.Next() {
		var qID int64
		var questionText string
		var aID int64
		var answerText string
		var isRightAnswer bool

		err := rows.Scan(&qID, &questionText, &aID, &answerText, &isRightAnswer)
		if err != nil {
			return nil, err
		}

		if q == nil {
			q = &dto.Test{
				QuestionID: qID,
				Question:   questionText,
				Answers:    []*dto.QuizAnswer{},
			}
		}

		q.Answers = append(q.Answers, &dto.QuizAnswer{
			AnswerID: aID,
			Answer:   answerText,
			IsRight:  isRightAnswer,
		})
	}

	if q == nil {
		return nil, sql.ErrNoRows
	}

	query = `
	SELECT question_lesson_id, answer_id, is_right
	FROM USER_ANSWERS
	WHERE user_id = $1 AND question_lesson_id = $2
	LIMIT 1
	`

	var Answer dto.UserAnswer
	err = d.conn.QueryRow(query, user_id, currentLessonId).Scan(&Answer.QuestionID, &Answer.AnswerID, &Answer.IsRight)
	if err == sql.ErrNoRows {
		return q, nil
	}
	if err != nil {
		return nil, err
	}
	q.UserAnswer = Answer

	return q, nil
}

func (d *Database) GetQuestionTestLesson(ctx context.Context, currentLessonId int, user_id int) (*dto.QuestionTest, error) {
	query := `
        SELECT qt.ID, qt.Question
        FROM Question_task qt
        WHERE qt.Lesson_test_id = $1
        LIMIT 1
    `

	var qt *dto.QuestionTest
	var QuestionID int64
	var Question string
	err := d.conn.QueryRow(query, currentLessonId).Scan(&QuestionID, &Question)
	qt = &dto.QuestionTest{
		QuestionID: QuestionID,
		Question:   Question,
	}
	if err != nil {
		return nil, err
	}

	query = `
	SELECT qta.answer, qta.status
	FROM question_task_answers qta
	JOIN question_task qt ON qt.id = qta.question_test_id
	WHERE qta.user_id = $1 AND qt.lesson_test_id = $2
	LIMIT 1;
`

	var Answer dto.UserQuestionAnswer
	err = d.conn.QueryRow(query, user_id, currentLessonId).Scan(&Answer.Answer, &Answer.Status)
	if err == sql.ErrNoRows {
		Answer.Status = "not passed"
		Answer.Answer = ""
		qt.UserAnswer = Answer
		logs.PrintLog(ctx, "GetQuestionTestLesson", fmt.Sprintf("%+v", qt))
		return qt, nil
	}
	if err != nil {
		return nil, err
	}
	qt.UserAnswer = Answer

	return qt, nil
}

func (d *Database) AnswerQuiz(ctx context.Context, question_id int, answer_id int, user_id int, course_id int) (*dto.QuizResult, error) {
	var isTrue bool
	err := d.conn.QueryRow(`
		SELECT av.Is_True 
		FROM ANSWER_VARIANT av
		JOIN QUIZ_TASK qt ON av.Quiz_Task_ID = qt.ID
		JOIN TEST_LESSON tl ON qt.Lesson_Test_ID = tl.ID
		WHERE av.ID = $1 AND tl.Lesson_ID = $2
	`, answer_id, question_id).Scan(&isTrue)

	if err != nil {
		return nil, err
	}

	_, err = d.conn.Exec(`
		INSERT INTO USER_ANSWERS (User_ID, Question_lesson_ID, Answer_ID, is_right)
		VALUES ($1, $2, $3, $4)
	`, user_id, question_id, answer_id, isTrue)

	if err != nil {
		return nil, err
	}

	res := &dto.QuizResult{
		Result: isTrue,
	}

	d.MarkLessonCompleted(ctx, user_id, course_id, question_id)

	return res, nil

}

func (d *Database) GetLessonById(ctx context.Context, lessonId int) (*coursemodels.LessonPoint, error) {
	var lesson coursemodels.LessonPoint
	err := d.conn.QueryRow(`
			SELECT id, title, type
			FROM LESSON
			WHERE id = $1
		`, lessonId).Scan(&lesson.LessonId, &lesson.Title, &lesson.Type)
	if err != nil {
		logs.PrintLog(ctx, "GetLessonById", fmt.Sprintf("%+v", err))
		return nil, err
	}
	return &lesson, nil
}

func (d *Database) AnswerQuestion(ctx context.Context, question_id int, user_id int, answer string) error {
	query := `
	INSERT INTO Question_task_answers (User_id, Question_test_id, Answer)
	VALUES ($1, $2, $3)
`
	_, err := d.conn.Exec(query, user_id, question_id, answer)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetCourseParts(ctx context.Context, courseId int) ([]*coursemodels.CoursePart, error) {
	var courseParts []*coursemodels.CoursePart
	rows, err := d.conn.Query(`
			SELECT id, title
			FROM PART
			WHERE course_id = $1
			ORDER BY part_order ASC
		`, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseParts", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var coursePart coursemodels.CoursePart
		if err := rows.Scan(&coursePart.Id, &coursePart.Title); err != nil {
			logs.PrintLog(ctx, "GetCourseParts", fmt.Sprintf("%+v", err))
			return nil, err
		}
		logs.PrintLog(ctx, "GetCourseParts", fmt.Sprintf("get course part %+v", coursePart))
		courseParts = append(courseParts, &coursePart)
	}
	return courseParts, nil
}

func (d *Database) GetPartBuckets(ctx context.Context, partId int) ([]*coursemodels.LessonBucket, error) {
	var buckets []*coursemodels.LessonBucket
	rows, err := d.conn.Query(`
			SELECT id, title
			FROM LESSON_BUCKET
			WHERE part_id = $1
			ORDER BY lesson_bucket_order ASC
		`, partId)
	if err != nil {
		logs.PrintLog(ctx, "GetPartBuckets", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bucket coursemodels.LessonBucket
		if err := rows.Scan(&bucket.Id, &bucket.Title); err != nil {
			logs.PrintLog(ctx, "GetPartBuckets", fmt.Sprintf("%+v", err))
			return nil, err
		}
		logs.PrintLog(ctx, "GetPartBuckets", fmt.Sprintf("get bucket %+v", bucket))
		buckets = append(buckets, &bucket)
	}
	return buckets, nil
}

func (d *Database) GetVideoUrl(ctx context.Context, lessonId int) (string, error) {
	var videoUrl string
	err := d.conn.QueryRow("SELECT video_src FROM video_lesson WHERE lesson_id = $1", lessonId).Scan(&videoUrl)
	if err != nil {
		logs.PrintLog(ctx, "GetVideo", fmt.Sprintf("%+v", err))
		return "", err
	}
	return videoUrl, nil
}
