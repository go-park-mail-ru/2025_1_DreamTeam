package postgres

import (
	"context"
	"fmt"
	coursemodels "skillForce/internal/models/course"
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

func (d *Database) DeleteCourseFromFavourites(ctx context.Context, courseId int, userId int) error {
	_, err := d.conn.Exec("DELETE FROM FAVOURITE_COURSES WHERE course_id = $1 AND user_id = $2", courseId, userId)
	if err != nil {
		logs.PrintLog(ctx, "DeleteCourseFromFavourites", fmt.Sprintf("%+v", err))
		return err
	}
	return nil
}

func (d *Database) GetFavouriteCourses(ctx context.Context, userId int) ([]*coursemodels.Course, error) {
	var bucketCourses []*coursemodels.Course
	rows, err := d.conn.Query(`
			SELECT c.id, c.creator_user_id, c.title, c.description, c.avatar_src, c.price, c.time_to_pass 
			FROM course c
			JOIN FAVOURITE_COURSES fc ON c.id = fc.course_id
			WHERE fc.user_id = $1
		`, userId)

	if err != nil {
		logs.PrintLog(ctx, "GetFavouriteCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logs.PrintLog(ctx, "SearchCoursesByTitle", fmt.Sprintf("%+v", err))
		}
	}()

	for rows.Next() {
		var course coursemodels.Course
		if err := rows.Scan(&course.Id, &course.CreatorId, &course.Title, &course.Description, &course.ScrImage, &course.Price, &course.TimeToPass); err != nil {
			logs.PrintLog(ctx, "GetFavouriteCourses", fmt.Sprintf("%+v", err))
			return nil, err
		}
		logs.PrintLog(ctx, "GetFavouriteCourses", fmt.Sprintf("get course %+v from db", course))
		bucketCourses = append(bucketCourses, &course)
	}

	logs.PrintLog(ctx, "GetFavouriteCourses", "get bucket ourses from db")

	return bucketCourses, nil
}

func (d *Database) GetCoursesFavouriteStatus(ctx context.Context, bucketCourses []*coursemodels.Course, userId int) (map[int]bool, error) {
	result := make(map[int]bool, len(bucketCourses))
	for _, course := range bucketCourses {
		result[course.Id] = false
	}

	rows, err := d.conn.Query(`
		SELECT course_id FROM FAVOURITE_COURSES WHERE user_id = $1
	`, userId)
	if err != nil {
		logs.PrintLog(ctx, "GetCoursesFavouriteStatus", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logs.PrintLog(ctx, "SearchCoursesByTitle", fmt.Sprintf("%+v", err))
		}
	}()

	for rows.Next() {
		var courseId int
		if err := rows.Scan(&courseId); err != nil {
			logs.PrintLog(ctx, "GetCoursesFavouriteStatus", fmt.Sprintf("%+v", err))
			return nil, err
		}
		if _, ok := result[courseId]; ok {
			result[courseId] = true
		}
	}
	return result, nil
}
