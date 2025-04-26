package postgres

import (
	"context"
	"fmt"
	coursemodels "skillForce/internal/models/course"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/logs"
)

func (d *Database) CreateCourse(ctx context.Context, course *coursemodels.Course, userProfile *usermodels.UserProfile) (int, error) {
	logs.PrintLog(ctx, "CreateCourse", fmt.Sprintf("create course %+v", course))
	query := `
        INSERT INTO COURSE (Creator_User_ID, Title, Description, Price, Time_to_pass)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING ID
    `

	var courseID int

	err := d.conn.QueryRow(
		query,
		userProfile.Id,
		course.Title,
		course.Description,
		course.Price,
		course.TimeToPass,
	).Scan(&courseID)

	if err != nil {
		logs.PrintLog(ctx, "CreateCourse", fmt.Sprintf("%+v", err))
		return 0, err
	}

	return courseID, nil
}

func (d *Database) CreatePart(ctx context.Context, part *coursemodels.CoursePart, courseId int) (int, error) {
	logs.PrintLog(ctx, "CreatePart", fmt.Sprintf("create part %+v", part))
	query := `
		INSERT INTO PART (Course_ID, Part_order, Title)
		VALUES ($1, $2, $3)
		RETURNING ID
	`

	var partID int

	err := d.conn.QueryRow(
		query,
		courseId,
		part.Order,
		part.Title,
	).Scan(&partID)

	if err != nil {
		logs.PrintLog(ctx, "CreatePart", fmt.Sprintf("%+v", err))
		return 0, err
	}
	return partID, nil
}

func (d *Database) CreateBucket(ctx context.Context, bucket *coursemodels.LessonBucket, partId int) (int, error) {
	logs.PrintLog(ctx, "CreateBucket", fmt.Sprintf("create bucket %+v", bucket))
	query := `
		INSERT INTO LESSON_BUCKET (Part_ID, Lesson_Bucket_Order, Title)
		VALUES ($1, $2, $3)
		RETURNING ID
	`

	var bucketID int

	err := d.conn.QueryRow(
		query,
		partId,
		bucket.Order,
		bucket.Title,
	).Scan(&bucketID)

	if err != nil {
		logs.PrintLog(ctx, "CreateBucket", fmt.Sprintf("%+v", err))
		return 0, err
	}

	return bucketID, nil
}

func (d *Database) CreateTextLesson(ctx context.Context, lesson *coursemodels.LessonPoint, bucketId int) error {
	logs.PrintLog(ctx, "CreateTextLesson", fmt.Sprintf("create text lesson %+v", lesson))
	query1 := `
		INSERT INTO LESSON (Lesson_Bucket_ID, Lesson_Order, Title, Type)
		VALUES ($1, $2, $3, $4)
		RETURNING ID
	`

	var lessonID int

	err := d.conn.QueryRow(
		query1,
		bucketId,
		lesson.Order,
		lesson.Title,
		lesson.Type,
	).Scan(&lessonID)

	if err != nil {
		logs.PrintLog(ctx, "CreateTextLesson", fmt.Sprintf("%+v", err))
		return err
	}

	query2 := `
		INSERT INTO TEXT_LESSON (Lesson_ID)
		VALUES ($1)
		RETURNING ID
	`
	var textLessonID int

	err = d.conn.QueryRow(
		query2,
		lessonID,
	).Scan(&textLessonID)

	if err != nil {
		logs.PrintLog(ctx, "CreateTextLesson", fmt.Sprintf("%+v", err))
		return err
	}

	query3 := `
		INSERT INTO text_lesson_block (text_lesson_id, value, is_image, text_lesson_block_order)
		VALUES ($1, $2, $3, $4)	
	`

	_, err = d.conn.Exec(
		query3,
		textLessonID,
		lesson.Value,
		lesson.IsImage,
		1,
	)

	if err != nil {
		logs.PrintLog(ctx, "CreateTextLesson", fmt.Sprintf("%+v", err))
		return err
	}

	return nil
}

func (d *Database) CreateVideoLesson(ctx context.Context, lesson *coursemodels.LessonPoint, bucketId int) error {
	logs.PrintLog(ctx, "CreateVideoLesson", fmt.Sprintf("create video lesson %+v", lesson))
	query1 := `
		INSERT INTO LESSON (Lesson_Bucket_ID, Lesson_Order, Title, Type)
		VALUES ($1, $2, $3, $4)
		RETURNING ID
	`

	var lessonID int

	err := d.conn.QueryRow(
		query1,
		bucketId,
		lesson.Order,
		lesson.Title,
		lesson.Type,
	).Scan(&lessonID)

	if err != nil {
		logs.PrintLog(ctx, "CreateTextLesson", fmt.Sprintf("%+v", err))
		return err
	}

	query2 := `
		INSERT INTO VIDEO_LESSON (Lesson_ID, Video_src)
		VALUES ($1, $2)
	`
	_, err = d.conn.Exec(
		query2,
		lessonID,
		lesson.Value,
	)

	if err != nil {
		logs.PrintLog(ctx, "CreateVideoLesson", fmt.Sprintf("%+v", err))
		return err
	}

	return nil
}

func (d *Database) CreateSurvey(ctx context.Context, survey *coursemodels.Survey, userProfile *usermodels.UserProfile) error {
	err := d.conn.QueryRow("INSERT INTO survey DEFAULT VALUES RETURNING id").Scan(&survey.Id)
	if err != nil {
		logs.PrintLog(ctx, "CreateSurvey", fmt.Sprintf("%+v", err))
		return err
	}

	for _, question := range survey.Questions {
		query := `
		INSERT INTO survey_question (survey_id, question, left_desc, right_desc, metric_type)
		VALUES ($1, $2, $3, $4, $5)	
	`

		_, err = d.conn.Exec(
			query,
			survey.Id,
			question.Question,
			question.LeftLebal,
			question.RightLebal,
			question.Metric,
		)

		if err != nil {
			logs.PrintLog(ctx, "CreateSurvey", fmt.Sprintf("%+v", err))
			return err
		}
	}

	return nil
}
