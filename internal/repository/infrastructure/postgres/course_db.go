package postgres

import (
	"context"
	"skillForce/internal/models"
	"skillForce/pkg/logs"
)

// GetBucketCourses - извлекает список курсов из базы данных
func (d *Database) GetBucketCourses(ctx context.Context) ([]*models.Course, error) {
	//TODO: можно заморочиться и сделать самописную пагинацию через LIMIT OFFSET
	var bucketCourses []*models.Course
	rows, err := d.conn.Query("SELECT id, creator_user_id, title, description, avatar_src, price, time_to_pass FROM course LIMIT 16")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.Id, &course.CreatorId, &course.Title, &course.Description, &course.ScrImage, &course.Price, &course.TimeToPass); err != nil {
			return nil, err
		}
		bucketCourses = append(bucketCourses, &course)
	}

	logs.PrintLog(ctx, "GetBucketCourses", "get bucket courses from db")

	return bucketCourses, nil
}
