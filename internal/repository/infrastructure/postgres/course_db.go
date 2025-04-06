package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"skillForce/internal/models"
	"skillForce/internal/models/dto"
	"skillForce/pkg/logs"
)

// GetBucketCourses - извлекает список курсов из базы данных
func (d *Database) GetBucketCourses(ctx context.Context) ([]*models.Course, error) {
	//TODO: можно заморочиться и сделать самописную пагинацию через LIMIT OFFSET
	var bucketCourses []*models.Course
	rows, err := d.conn.Query("SELECT id, creator_user_id, title, description, avatar_src, price, time_to_pass, purchases_amount FROM course LIMIT 16")
	if err != nil {
		logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.Id, &course.CreatorId, &course.Title, &course.Description, &course.ScrImage, &course.Price, &course.TimeToPass, &course.PurchasesAmount); err != nil {
			logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
			return nil, err
		}
		logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("get course %+v from db", course))
		bucketCourses = append(bucketCourses, &course)
	}

	logs.PrintLog(ctx, "GetBucketCourses", "get bucket ourses from db")

	return bucketCourses, nil
}

func (d *Database) GetCoursesRaitings(ctx context.Context, bucketCoursesWithoutRating []*models.Course) (map[int]models.CourseRating, error) {
	coursesRatings := make(map[int]models.CourseRating, len(bucketCoursesWithoutRating))

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

		coursesRatings[course.Id] = models.CourseRating{
			CourseId: course.Id,
			Rating:   sumMetrics / countMetrics,
		}
	}
	logs.PrintLog(ctx, "GetCoursesRaitings", "get courses ratings from db")
	return coursesRatings, nil
}

func (d *Database) GetCoursesTags(ctx context.Context, bucketCoursesWithoutTags []*models.Course) (map[int][]string, error) {
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

func (d *Database) GetCourseById(ctx context.Context, courseId int) (*models.Course, error) {
	var course models.Course
	err := d.conn.QueryRow("SELECT id, creator_user_id, title, description, avatar_src, price, time_to_pass, purchases_amount FROM course WHERE id = $1", courseId).Scan(
		&course.Id, &course.CreatorId, &course.Title, &course.Description, &course.ScrImage, &course.Price, &course.TimeToPass, &course.PurchasesAmount)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseById", fmt.Sprintf("%+v", err))
		return nil, err
	}
	logs.PrintLog(ctx, "GetCourseById", fmt.Sprintf("get course %+v from db", course))
	return &course, nil
}

func (d *Database) FillLessonHeader(ctx context.Context, userId int, courseId int, LessonHeader *dto.LessonDtoHeader) error {
	var LessonId int
	err := d.conn.QueryRow(`
			SELECT Lesson_ID
			FROM LESSON_CHECKPOINT
			WHERE User_ID = $1 AND Course_ID = $2
			ORDER BY Updated_at DESC
			LIMIT 1;
		`, userId, courseId).Scan(
		&LessonId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("no lesson found for user id:%+v and course id:%+v", userId, courseId))

			var part models.CoursePart
			err := d.conn.QueryRow(`
					SELECT title, part_order, id
					FROM part
					WHERE course_id = $1
					ORDER BY part_order ASC
					LIMIT 1;
				`, courseId).Scan(
				&part.Title, &part.Order, &part.Id)

			if err != nil {
				logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("course %+v has no parts:%+v", courseId, err))
			}

			LessonHeader.Part.Order = part.Order
			LessonHeader.Part.Title = part.Title

			var bucket models.LessonBucket
			err = d.conn.QueryRow(`
					SELECT title, lesson_bucket_order, id
					FROM lesson_bucket
					WHERE part_id = $1
					ORDER BY lesson_bucket_order ASC
					LIMIT 1;
				`, part.Id).Scan(
				&bucket.Title, &bucket.Order, &bucket.Id)

			if err != nil {
				logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("part %+v has no buckets:%+v", part.Id, err))
			}

			LessonHeader.Bucket.Order = bucket.Order
			LessonHeader.Bucket.Title = bucket.Title

			rows, err := d.conn.Query(`
					SELECT id, type
					FROM LESSON
					WHERE lesson_bucket_id = $1
					ORDER BY Lesson_Order ASC
				`, bucket.Id)

			if err != nil {
				logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("%+v", err))
				return err
			}
			defer rows.Close()

			var points []models.LessonPoint
			for rows.Next() {
				var point models.LessonPoint
				if err := rows.Scan(&point.LessonId, &point.Type); err != nil {
					logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("%+v", err))
					return err
				}
				points = append(points, point)
			}

			if len(points) == 0 {
				logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("no points found in bucket %+v", bucket.Id))
				return errors.New("no points found in bucket")
			}

			for _, point := range points {
				err := d.conn.QueryRow(`
					SELECT EXISTS(SELECT 1
					FROM LESSON_CHECKPOINT
					WHERE Lesson_ID = $1)
				`, point.LessonId).Scan(&point.IsDone)
				if err != nil {
					logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("%+v", err))
					return err
				}
			}

			for _, point := range points {
				LessonHeader.Points = append(LessonHeader.Points, struct {
					LessonId int    `json:"lesson_id"`
					Type     string `json:"type"`
					IsDone   bool   `json:"is_done"`
				}{
					LessonId: point.LessonId,
					Type:     point.Type,
					IsDone:   point.IsDone,
				})
			}

			return nil
		}

		logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("%+v", err))
		return err
	}

	return nil
}
