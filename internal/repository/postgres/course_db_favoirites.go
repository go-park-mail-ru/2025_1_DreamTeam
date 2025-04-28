package postgres

import (
	"context"
	"fmt"
	"skillForce/pkg/logs"
)

func (d *Database) AddCourseToFavourites(ctx context.Context, courseId int, userId int) error {
	query := `
		INSERT INTO FAVOURITE_COURSES (user_id, course_id)
		VALUES ($1, $2)	
	`

	_, err := d.conn.Exec(
		query,
		userId,
		courseId,
	)

	if err != nil {
		logs.PrintLog(ctx, "CreateTextLesson", fmt.Sprintf("%+v", err))
		return err
	}

	return nil
}
