package postgres

import (
	"context"
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

func (d *Database) markLessonCompleted(ctx context.Context, userId int, courseId int, lessonId int) error {
	_, err := d.conn.Exec(
		"INSERT INTO LESSON_CHECKPOINT (user_id, lesson_id, course_id) VALUES ($1, $2, $3)",
		userId, lessonId, courseId)

	if err != nil {
		logs.PrintLog(ctx, "markLessonComplete", fmt.Sprintf("%+v", err))
		return err
	}

	logs.PrintLog(ctx, "markLessonComplete", fmt.Sprintf("mark that lesson id:%+v is learned by the user id:%+v", lessonId, userId))
	return nil
}

func (d *Database) fillLessonHeaderNewCourse(ctx context.Context, userId int, courseId int, lessonHeader *dto.LessonDtoHeader) (int, int, string, error) {
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
		logs.PrintLog(ctx, "fillLessonHeaderNewCourse", fmt.Sprintf("course %+v has no parts:%+v", courseId, err))
	}

	lessonHeader.Part.Order = part.Order
	lessonHeader.Part.Title = part.Title

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
		logs.PrintLog(ctx, "fillLessonHeaderNewCourse", fmt.Sprintf("part %+v has no buckets:%+v", part.Id, err))
	}

	lessonHeader.Bucket.Order = bucket.Order
	lessonHeader.Bucket.Title = bucket.Title

	rows, err := d.conn.Query(`
					SELECT id, type
					FROM LESSON
					WHERE lesson_bucket_id = $1
					ORDER BY Lesson_Order ASC
				`, bucket.Id)

	if err != nil {
		logs.PrintLog(ctx, "fillLessonHeaderNewCourse", fmt.Sprintf("%+v", err))
		return 0, 0, "", err
	}
	defer rows.Close()

	var points []models.LessonPoint
	for rows.Next() {
		var point models.LessonPoint
		if err := rows.Scan(&point.LessonId, &point.Type); err != nil {
			logs.PrintLog(ctx, "fillLessonHeaderNewCourse", fmt.Sprintf("%+v", err))
			return 0, 0, "", err
		}
		points = append(points, point)
	}

	if len(points) == 0 {
		logs.PrintLog(ctx, "fillLessonHeaderNewCourse", fmt.Sprintf("no points found in bucket %+v", bucket.Id))
		return 0, 0, "", errors.New("no points found in bucket")
	}

	currentLessonId := points[0].LessonId
	currentLessonType := points[0].Type

	err = d.markLessonCompleted(ctx, userId, courseId, currentLessonId)
	if err != nil {
		logs.PrintLog(ctx, "fillLessonHeaderNewCourse", fmt.Sprintf("%+v", err))
		return 0, 0, "", err
	}
	points[0].IsDone = true

	for _, point := range points {
		lessonHeader.Points = append(lessonHeader.Points, struct {
			LessonId int    `json:"lesson_id"`
			Type     string `json:"type"`
			IsDone   bool   `json:"is_done"`
		}{
			LessonId: point.LessonId,
			Type:     point.Type,
			IsDone:   point.IsDone,
		})
	}

	return currentLessonId, bucket.Id, currentLessonType, nil
}

func (d *Database) FillLessonHeader(ctx context.Context, userId int, courseId int, lessonHeader *dto.LessonDtoHeader) (int, int, string, error) {
	rows1, err := d.conn.Query(`
			SELECT cp.Lesson_ID, l.type
			FROM LESSON_CHECKPOINT cp
			JOIN LESSON l ON l.ID = cp.Lesson_ID
			WHERE cp.User_ID = $1 AND cp.Course_ID = $2
			ORDER BY cp.Updated_at DESC
		`, userId, courseId)

	if err != nil {
		logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("%+v", err))
		return 0, 0, "", err
	}
	defer rows1.Close()

	var lessonPoint models.LessonPoint
	var visitedLessonPointsIds []int
	firstRow := true
	for rows1.Next() {
		if firstRow {
			firstRow = false
			err := rows1.Scan(&lessonPoint.LessonId, &lessonPoint.Type)
			if err != nil {
				logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("%+v", err))
				return 0, 0, "", err
			}
			visitedLessonPointsIds = append(visitedLessonPointsIds, lessonPoint.LessonId)
			continue
		}
		var lessonPointId int
		var lessonPointType string
		err := rows1.Scan(&lessonPointId, &lessonPointType)
		if err != nil {
			logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("%+v", err))
			return 0, 0, "", err
		}
		visitedLessonPointsIds = append(visitedLessonPointsIds, lessonPointId)
	}

	if len(visitedLessonPointsIds) == 0 {
		logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("no visited lesson points fo user %+v in course %+v", userId, courseId))
		return d.fillLessonHeaderNewCourse(ctx, userId, courseId, lessonHeader)
	}

	logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("lesson point %+v", lessonPoint))
	logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("visited lesson points ids %+v", visitedLessonPointsIds))

	var part models.CoursePart
	var bucket models.LessonBucket
	err = d.conn.QueryRow(`
			SELECT p.Title, p.Part_Order, p.ID, lb.Title, lb.Lesson_Bucket_Order, lb.ID
			FROM PART p
			JOIN LESSON_BUCKET lb ON lb.Part_ID = p.ID
			JOIN LESSON l ON l.Lesson_Bucket_ID = lb.ID
			WHERE l.ID = $1
		`, lessonPoint.LessonId).Scan(
		&part.Title, &part.Order, &part.Id, &bucket.Title, &bucket.Order, &bucket.Id)

	if err != nil {
		logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("%+v", err))
		return 0, 0, "", err
	}

	lessonHeader.Part.Order = part.Order
	lessonHeader.Part.Title = part.Title
	lessonHeader.Bucket.Order = bucket.Order
	lessonHeader.Bucket.Title = bucket.Title

	logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("part %+v, bucket %+v", part, bucket))

	rows2, err := d.conn.Query(`
					SELECT id, type
					FROM LESSON
					WHERE lesson_bucket_id = $1
					ORDER BY Lesson_Order ASC
				`, bucket.Id)

	if err != nil {
		logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("%+v", err))
		return 0, 0, "", err
	}
	defer rows2.Close()

	var points []models.LessonPoint
	for rows2.Next() {
		var point models.LessonPoint
		if err := rows2.Scan(&point.LessonId, &point.Type); err != nil {
			logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("%+v", err))
			return 0, 0, "", err
		}
		for _, visitedLessonPointId := range visitedLessonPointsIds {
			if point.LessonId == visitedLessonPointId {
				point.IsDone = true
				continue
			}
		}
		points = append(points, point)
	}

	if len(points) == 0 {
		logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("no points found in bucket %+v", bucket.Id))
		return 0, 0, "", errors.New("no points found in bucket")
	}

	for _, point := range points {
		lessonHeader.Points = append(lessonHeader.Points, struct {
			LessonId int    `json:"lesson_id"`
			Type     string `json:"type"`
			IsDone   bool   `json:"is_done"`
		}{
			LessonId: point.LessonId,
			Type:     point.Type,
			IsDone:   point.IsDone,
		})
	}

	return lessonPoint.LessonId, bucket.Id, lessonPoint.Type, nil
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

func (d *Database) GetLessonFooters(ctx context.Context, currentLessonId int, currentBucketId int) ([]int, error) {
	footers := []int{-1, -1}

	var currentOrder int
	err := d.conn.QueryRow(`
			SELECT lesson_order
			FROM LESSON
			WHERE id = $1
		`, currentLessonId).Scan(&currentOrder)
	if err != nil {
		logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("%+v", err))
		return nil, err
	}

	rows, err := d.conn.Query(`
			SELECT id, lesson_order
			FROM LESSON
			WHERE lesson_bucket_id = $1
			ORDER BY Lesson_Order ASC
		`, currentBucketId)

	if err != nil {
		logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var footer struct {
			LessonId int
			Order    int
		}

		if err := rows.Scan(&footer.LessonId, &footer.Order); err != nil {
			logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("%+v", err))
			return nil, err
		}

		if footer.Order == currentOrder-1 {
			footers[0] = footer.LessonId
			continue
		} else if footer.Order == currentOrder+1 {
			footers[1] = footer.LessonId
		}
	}

	return footers, nil
}
