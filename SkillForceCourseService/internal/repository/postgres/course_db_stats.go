package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	coursemodels "skillForce/internal/models/course"
	"skillForce/pkg/logs"
)

func (d *Database) GetCoursesPurchases(ctx context.Context, bucketCoursesWithoutPurchases []*coursemodels.Course) (map[int]int, error) {
	coursesPurchases := make(map[int]int, len(bucketCoursesWithoutPurchases))

	for _, course := range bucketCoursesWithoutPurchases {
		var purchases int
		err := d.conn.QueryRow("SELECT COUNT(*) FROM SIGNUPS WHERE course_id = $1", course.Id).Scan(&purchases)
		if err != nil {
			logs.PrintLog(ctx, "GetBucketCoursesPurchases", fmt.Sprintf("%+v", err))
			return nil, err
		}

		coursesPurchases[course.Id] = purchases
	}
	return coursesPurchases, nil
}

func (d *Database) GetCoursesRaitings(ctx context.Context, bucketCoursesWithoutRating []*coursemodels.Course) (map[int]float32, error) {
	coursesRatings := make(map[int]float32, len(bucketCoursesWithoutRating))

	for _, course := range bucketCoursesWithoutRating {
		rows, err := d.conn.Query("SELECT rating FROM course_metrik WHERE course_id = $1", course.Id)
		if err != nil {
			logs.PrintLog(ctx, "GetCoursesRaitings", fmt.Sprintf("%+v", err))
			return nil, err
		}
		defer rows.Close()

		var sumMetrics float32
		var countMetrics float32

		for rows.Next() {
			var metric float32
			if err := rows.Scan(&metric); err != nil {
				logs.PrintLog(ctx, "GetCoursesRaitings", fmt.Sprintf("%+v", err))
				return nil, err
			}
			sumMetrics += metric
			countMetrics++
		}

		if countMetrics == 0 {
			continue
		}

		coursesRatings[course.Id] = sumMetrics / countMetrics
	}

	logs.PrintLog(ctx, "GetCoursesRaitings", "get courses ratings from db")
	return coursesRatings, nil
}

func (d *Database) GetCoursesTags(ctx context.Context, bucketCoursesWithoutTags []*coursemodels.Course) (map[int][]string, error) {
	coursesTags := make(map[int][]string, len(bucketCoursesWithoutTags))

	for _, course := range bucketCoursesWithoutTags {
		rows, err := d.conn.Query(`
			SELECT vt.Title
			FROM TAGS t
			JOIN VALID_TAGS vt ON t.Tag_ID = vt.ID
			WHERE t.Course_ID = $1
		`, course.Id)

		if err != nil {
			logs.PrintLog(ctx, "GetCoursesTags", fmt.Sprintf("%+v", err))
			return nil, err
		}
		defer rows.Close()

		var tags []string

		for rows.Next() {
			var tag string
			if err := rows.Scan(&tag); err != nil {
				logs.PrintLog(ctx, "GetCoursesTags", fmt.Sprintf("%+v", err))
				return nil, err
			}
			tags = append(tags, tag)
		}

		if len(tags) == 0 {
			continue
		}

		coursesTags[course.Id] = tags

	}
	logs.PrintLog(ctx, "GetCoursesTags", "get courses tags from db")
	return coursesTags, nil
}

func (d *Database) IsUserPurchasedCourse(ctx context.Context, userId int, courseId int) (bool, error) {
	var exists bool
	err := d.conn.QueryRow("SELECT EXISTS (SELECT 1 FROM SIGNUPS WHERE user_id = $1 AND course_id = $2)",
		userId, courseId).Scan(&exists)

	if errors.Is(err, sql.ErrNoRows) {
		logs.PrintLog(ctx, "IsUserPurchasedCourse", fmt.Sprintf("%+v", err))
		return false, nil
	}

	if err != nil {
		logs.PrintLog(ctx, "IsUserPurchasedCourse", fmt.Sprintf("%+v", err))
		return false, err
	}
	return exists, nil
}

func (d *Database) IsUserCompletedCourse(ctx context.Context, userId int, courseId int) (bool, error) {
	var exists bool
	err := d.conn.QueryRow("SELECT EXISTS (SELECT 1 FROM COMPLETED_COURSES WHERE user_id = $1 AND course_id = $2)",
		userId, courseId).Scan(&exists)

	if errors.Is(err, sql.ErrNoRows) {
		logs.PrintLog(ctx, "IsUserCompletedCourse", fmt.Sprintf("%+v", err))
		return false, nil
	}

	if err != nil {
		logs.PrintLog(ctx, "IsUserCompletedCourse", fmt.Sprintf("%+v", err))
		return false, err
	}
	return exists, nil
}
