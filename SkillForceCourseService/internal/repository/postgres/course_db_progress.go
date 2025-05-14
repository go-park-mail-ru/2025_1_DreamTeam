package postgres

import (
	"context"
	"fmt"
	"skillForce/pkg/logs"
)

func (d *Database) MarkLessonCompleted(ctx context.Context, userId int, lessonId int) error {
	exists := false
	err := d.conn.QueryRow("SELECT EXISTS (SELECT 1 FROM LESSON_CHECKPOINT WHERE user_id = $1 AND lesson_id = $2)",
		userId, lessonId).Scan(&exists)
	if err != nil {
		logs.PrintLog(ctx, "MarkLessonComplete", fmt.Sprintf("%+v", err))
	}

	if exists {
		logs.PrintLog(ctx, "MarkLessonComplete", fmt.Sprintf("lesson id:%+v is already learned by the user id:%+v", lessonId, userId))
		return nil
	}

	var courseId int
	err = d.conn.QueryRow(`
			SELECT p.Course_ID
			FROM LESSON l
			JOIN LESSON_BUCKET lb ON l.Lesson_Bucket_ID = lb.id
			JOIN PART p ON p.id = lb.Part_ID
			WHERE l.id = $1
		`, lessonId).Scan(&courseId)
	if err != nil {
		logs.PrintLog(ctx, "MarkLessonComplete", fmt.Sprintf("%+v", err))
		return err
	}

	_, err = d.conn.Exec(
		"INSERT INTO LESSON_CHECKPOINT (user_id, lesson_id, course_id) VALUES ($1, $2, $3)",
		userId, lessonId, courseId)

	if err != nil {
		logs.PrintLog(ctx, "MarkLessonComplete", fmt.Sprintf("%+v", err))
		return err
	}

	logs.PrintLog(ctx, "MarkLessonComplete", fmt.Sprintf("mark that lesson id:%+v is learned by the user id:%+v", lessonId, userId))
	return nil
}

func (d *Database) MarkLessonAsNotCompleted(ctx context.Context, userId int, lessonId int) error {
	_, err := d.conn.Exec(
		"DELETE FROM LESSON_CHECKPOINT WHERE user_id = $1 AND lesson_id = $2",
		userId, lessonId)

	if err != nil {
		logs.PrintLog(ctx, "markLessonAsNotComplete", fmt.Sprintf("%+v", err))
		return err
	}

	logs.PrintLog(ctx, "markLessonAsNotComplete", fmt.Sprintf("mark that lesson id:%+v is not learned by the user id:%+v", lessonId, userId))
	return nil
}

func (d *Database) AddUserToCourse(ctx context.Context, userId int, courseId int) error {
	exists := false
	err := d.conn.QueryRow("SELECT EXISTS (SELECT 1 FROM SIGNUPS WHERE user_id = $1 AND course_id = $2)",
		userId, courseId).Scan(&exists)
	if err != nil {
		logs.PrintLog(ctx, "AddUserToCourse", fmt.Sprintf("%+v", err))
	}

	if exists {
		logs.PrintLog(ctx, "AddUserToCourse", fmt.Sprintf("user with id %+v is already in course with id %+v", userId, courseId))
		return nil
	}

	_, err = d.conn.Exec(
		"INSERT INTO SIGNUPS (user_id, course_id) VALUES ($1, $2)",
		userId, courseId)
	if err != nil {
		logs.PrintLog(ctx, "AddUserToCourse", fmt.Sprintf("%+v", err))
		return err
	}
	logs.PrintLog(ctx, "AddUserToCourse", fmt.Sprintf("add user with id %+v to course with id %+v", userId, courseId))
	return nil
}
