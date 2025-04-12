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
	rows, err := d.conn.Query("SELECT id, creator_user_id, title, description, avatar_src, price, time_to_pass FROM course LIMIT 16")
	if err != nil {
		logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var course models.Course
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

func (d *Database) GetCoursesPurchases(ctx context.Context, bucketCoursesWithoutPurchases []*models.Course) (map[int]int, error) {
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

func (d *Database) GetCoursesRaitings(ctx context.Context, bucketCoursesWithoutRating []*models.Course) (map[int]float32, error) {
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
	err := d.conn.QueryRow("SELECT id, creator_user_id, title, description, avatar_src, price, time_to_pass FROM course WHERE id = $1", courseId).Scan(
		&course.Id, &course.CreatorId, &course.Title, &course.Description, &course.ScrImage, &course.Price, &course.TimeToPass)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseById", fmt.Sprintf("%+v", err))
		return nil, err
	}
	logs.PrintLog(ctx, "GetCourseById", fmt.Sprintf("get course %+v from db", course))
	return &course, nil
}

func (d *Database) MarkLessonCompleted(ctx context.Context, userId int, courseId int, lessonId int) error {
	exists := false
	err := d.conn.QueryRow("SELECT EXISTS (SELECT 1 FROM LESSON_CHECKPOINT WHERE user_id = $1 AND lesson_id = $2 AND course_id = $3)",
		userId, lessonId, courseId).Scan(&exists)
	if err != nil {
		logs.PrintLog(ctx, "MarkLessonComplete", fmt.Sprintf("%+v", err))
	}

	if exists {
		logs.PrintLog(ctx, "MarkLessonComplete", fmt.Sprintf("lesson id:%+v is already learned by the user id:%+v", lessonId, userId))
		return nil
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

func (d *Database) getLessonHeaderNewCourse(ctx context.Context, userId int, courseId int) (*dto.LessonDtoHeader, int, string, error) {
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

	var lessonHeader dto.LessonDtoHeader
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
		return nil, 0, "", err
	}
	defer rows.Close()

	var points []models.LessonPoint
	for rows.Next() {
		var point models.LessonPoint
		if err := rows.Scan(&point.LessonId, &point.Type); err != nil {
			logs.PrintLog(ctx, "fillLessonHeaderNewCourse", fmt.Sprintf("%+v", err))
			return nil, 0, "", err
		}
		points = append(points, point)
	}

	if len(points) == 0 {
		logs.PrintLog(ctx, "fillLessonHeaderNewCourse", fmt.Sprintf("no points found in bucket %+v", bucket.Id))
		return nil, 0, "", err
	}

	currentLessonId := points[0].LessonId
	currentLessonType := points[0].Type

	err = d.MarkLessonCompleted(ctx, userId, courseId, currentLessonId)
	if err != nil {
		logs.PrintLog(ctx, "fillLessonHeaderNewCourse", fmt.Sprintf("%+v", err))
		return nil, 0, "", err
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

	return &lessonHeader, currentLessonId, currentLessonType, nil
}

func (d *Database) GetLastLessonHeader(ctx context.Context, userId int, courseId int) (*dto.LessonDtoHeader, int, string, error) {
	var lessonHeader dto.LessonDtoHeader
	course, err := d.GetCourseById(ctx, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseLesson", fmt.Sprintf("%+v", err))
		return nil, 0, "", err
	}
	lessonHeader.CourseTitle = course.Title
	lessonHeader.CourseId = course.Id

	rows1, err := d.conn.Query(`
			SELECT cp.Lesson_ID, l.type
			FROM LESSON_CHECKPOINT cp
			JOIN LESSON l ON l.ID = cp.Lesson_ID
			WHERE cp.User_ID = $1 AND cp.Course_ID = $2
			ORDER BY cp.Updated_at DESC
		`, userId, courseId)

	if err != nil {
		logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("%+v", err))
		return nil, 0, "", err
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
				return nil, 0, "", err
			}
			visitedLessonPointsIds = append(visitedLessonPointsIds, lessonPoint.LessonId)
			continue
		}
		var lessonPointId int
		var lessonPointType string
		err := rows1.Scan(&lessonPointId, &lessonPointType)
		if err != nil {
			logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("%+v", err))
			return nil, 0, "", err
		}
		visitedLessonPointsIds = append(visitedLessonPointsIds, lessonPointId)
	}

	if len(visitedLessonPointsIds) == 0 {
		logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("no visited lesson points fo user %+v in course %+v", userId, courseId))
		return d.getLessonHeaderNewCourse(ctx, userId, courseId)
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
		return nil, 0, "", err
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
		return nil, 0, "", err
	}
	defer rows2.Close()

	var points []models.LessonPoint
	for rows2.Next() {
		var point models.LessonPoint
		if err := rows2.Scan(&point.LessonId, &point.Type); err != nil {
			logs.PrintLog(ctx, "FillLessonHeader", fmt.Sprintf("%+v", err))
			return nil, 0, "", err
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
		return nil, 0, "", err
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

	return &lessonHeader, lessonPoint.LessonId, lessonPoint.Type, nil
}

func (d *Database) GetLessonHeaderByLessonId(ctx context.Context, userId int, currentLessonId int) (*dto.LessonDtoHeader, error) {
	var course models.Course
	var part models.CoursePart
	var bucket models.LessonBucket
	err := d.conn.QueryRow(`
			SELECT p.Title, p.Part_Order, p.ID, lb.Title, lb.Lesson_Bucket_Order, lb.ID, c.ID, c.Title
			FROM lesson l
			JOIN LESSON_BUCKET lb ON l.LESSON_BUCKET_ID = lb.ID
			JOIN PART p ON lb.PART_ID = p.ID
			JOIN COURSE c ON p.COURSE_ID = c.ID
			WHERE l.ID = $1
		`, currentLessonId).Scan(
		&part.Title, &part.Order, &part.Id, &bucket.Title, &bucket.Order, &bucket.Id, &course.Id, &course.Title)

	if err != nil {
		logs.PrintLog(ctx, "FillLessonHeaderByLessonId", fmt.Sprintf("%+v", err))
		return nil, err
	}

	lessonHeader := &dto.LessonDtoHeader{}
	lessonHeader.CourseId = course.Id
	lessonHeader.CourseTitle = course.Title
	lessonHeader.Part.Order = part.Order
	lessonHeader.Part.Title = part.Title
	lessonHeader.Bucket.Order = bucket.Order
	lessonHeader.Bucket.Title = bucket.Title

	rows1, err := d.conn.Query(`
		SELECT lesson_id
		FROM LESSON_CHECKPOINT
		WHERE course_id = $1 and user_id = $2
	`, course.Id, userId)

	if err != nil {
		logs.PrintLog(ctx, "FillLessonHeaderByLessonId", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer rows1.Close()

	visitedLessonPointsIds := make(map[int]bool)
	for rows1.Next() {
		var visitedLessonPointId int
		if err := rows1.Scan(&visitedLessonPointId); err != nil {
			logs.PrintLog(ctx, "FillLessonHeaderByLessonId", fmt.Sprintf("%+v", err))
			return nil, err
		}
		visitedLessonPointsIds[visitedLessonPointId] = true
	}

	logs.PrintLog(ctx, "FillLessonHeaderByLessonId", fmt.Sprintf("visitedLessonPointsIds %+v", visitedLessonPointsIds))

	rows2, err := d.conn.Query(`
		SELECT id, type
		FROM LESSON
		WHERE lesson_bucket_id = $1
	`, bucket.Id)

	if err != nil {
		logs.PrintLog(ctx, "FillLessonHeaderByLessonId", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer rows1.Close()

	var points []*dto.LessonPointDTO
	for rows2.Next() {
		var point models.LessonPoint
		if err := rows2.Scan(&point.LessonId, &point.Type); err != nil {
			logs.PrintLog(ctx, "FillLessonHeaderByLessonId", fmt.Sprintf("%+v", err))
			return nil, err
		}

		_, ok := visitedLessonPointsIds[point.LessonId]
		if ok {
			point.IsDone = true
		}

		points = append(points, &dto.LessonPointDTO{
			LessonId: point.LessonId,
			Type:     point.Type,
			IsDone:   point.IsDone,
		})
	}

	if len(points) == 0 {
		logs.PrintLog(ctx, "FillLessonHeaderByLessonId", fmt.Sprintf("no points found in bucket %+v", bucket.Id))
		return nil, errors.New("no points found in bucket")
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

	return lessonHeader, nil
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

func (d *Database) GetBucketByLessonId(ctx context.Context, currentLessonId int) (*models.LessonBucket, error) {
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

	bucket := &models.LessonBucket{
		Id: bucketId,
	}
	return bucket, nil
}

func (d *Database) GetLessonFooters(ctx context.Context, currentLessonId int) ([]int, error) {
	footers := []int{-1, -1, -1}

	var currentLessonOrder int
	var currentBucket models.LessonBucket
	err := d.conn.QueryRow(`
			SELECT l.lesson_order, lb.id, lb.lesson_bucket_order
			FROM LESSON l
			JOIN LESSON_BUCKET lb ON l.lesson_bucket_id = lb.id
			WHERE l.id = $1
		`, currentLessonId).Scan(&currentLessonOrder, &currentBucket.Id, &currentBucket.Order)
	if err != nil {
		logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("%+v", err))
		return nil, err
	}

	footers[1] = currentLessonId
	logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("current order: %d", currentLessonOrder))

	rows, err := d.conn.Query(`
			SELECT id, lesson_order
			FROM LESSON
			WHERE lesson_bucket_id = $1
			ORDER BY Lesson_Order ASC
		`, currentBucket.Id)

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

		if footer.Order == currentLessonOrder-1 {
			footers[0] = footer.LessonId
		} else if footer.Order == currentLessonOrder+1 {
			footers[2] = footer.LessonId
		}
	}

	if footers[0] == -1 && currentBucket.Order > 1 {
		var prevLessonId int
		err := d.conn.QueryRow(`
				SELECT l.id
				FROM LESSON_BUCKET lb
				JOIN LESSON l ON l.lesson_bucket_id = lb.id
				WHERE lb.lesson_bucket_order = $1
				ORDER BY l.Lesson_Order DESC
				LIMIT 1
			`, currentBucket.Order-1).Scan(&prevLessonId)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("%+v", err))
				return nil, err
			}
		}
		footers[0] = prevLessonId
	}

	if footers[2] == -1 && currentBucket.Order < 2 {
		var nextLessonId int
		err := d.conn.QueryRow(`
				SELECT l.id
				FROM LESSON_BUCKET lb
				JOIN LESSON l ON l.lesson_bucket_id = lb.id
				WHERE lb.lesson_bucket_order = $1
				ORDER BY l.Lesson_Order ASC
				LIMIT 1
			`, currentBucket.Order+1).Scan(&nextLessonId)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("%+v", err))
				return nil, err
			}
		}
		footers[2] = nextLessonId
	}

	logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("footers: %+v", footers))

	return footers, nil
}

func (d *Database) GetCourseParts(ctx context.Context, courseId int) ([]*models.CoursePart, error) {
	var courseParts []*models.CoursePart
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
		var coursePart models.CoursePart
		if err := rows.Scan(&coursePart.Id, &coursePart.Title); err != nil {
			logs.PrintLog(ctx, "GetCourseParts", fmt.Sprintf("%+v", err))
			return nil, err
		}
		logs.PrintLog(ctx, "GetCourseParts", fmt.Sprintf("get course part %+v", coursePart))
		courseParts = append(courseParts, &coursePart)
	}
	return courseParts, nil
}

func (d *Database) GetPartBuckets(ctx context.Context, partId int) ([]*models.LessonBucket, error) {
	var buckets []*models.LessonBucket
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
		var bucket models.LessonBucket
		if err := rows.Scan(&bucket.Id, &bucket.Title); err != nil {
			logs.PrintLog(ctx, "GetPartBuckets", fmt.Sprintf("%+v", err))
			return nil, err
		}
		logs.PrintLog(ctx, "GetPartBuckets", fmt.Sprintf("get bucket %+v", bucket))
		buckets = append(buckets, &bucket)
	}
	return buckets, nil
}

func (d *Database) GetBucketLessons(ctx context.Context, userId int, courseId int, bucketId int) ([]*models.LessonPoint, error) {
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

	var lessons []*models.LessonPoint
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
		var lesson models.LessonPoint
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
