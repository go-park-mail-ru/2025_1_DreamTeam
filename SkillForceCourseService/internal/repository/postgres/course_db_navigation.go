package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	coursemodels "skillForce/internal/models/course"
	"skillForce/internal/models/dto"
	"skillForce/pkg/logs"
)

func (d *Database) getLessonHeaderNewCourse(ctx context.Context, userId int, courseId int) (*dto.LessonDtoHeader, int, string, bool, error) {
	var part coursemodels.CoursePart
	err := d.conn.QueryRow(`
					SELECT title, part_order, id
					FROM part
					WHERE course_id = $1
					ORDER BY part_order ASC
					LIMIT 1;
				`, courseId).Scan(
		&part.Title, &part.Order, &part.Id)

	if err != nil {
		logs.PrintLog(ctx, "getLessonHeaderNewCourse", fmt.Sprintf("course %+v has no parts:%+v", courseId, err))
	}

	var lessonHeader dto.LessonDtoHeader
	lessonHeader.Part.Order = part.Order
	lessonHeader.Part.Title = part.Title

	var bucket coursemodels.LessonBucket
	err = d.conn.QueryRow(`
					SELECT title, lesson_bucket_order, id
					FROM lesson_bucket
					WHERE part_id = $1
					ORDER BY lesson_bucket_order ASC
					LIMIT 1;
				`, part.Id).Scan(
		&bucket.Title, &bucket.Order, &bucket.Id)

	if err != nil {
		logs.PrintLog(ctx, "getLessonHeaderNewCourse", fmt.Sprintf("part %+v has no buckets:%+v", part.Id, err))
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
		logs.PrintLog(ctx, "getLessonHeaderNewCourse", fmt.Sprintf("%+v", err))
		return nil, 0, "", false, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logs.PrintLog(ctx, "SearchCoursesByTitle", fmt.Sprintf("%+v", err))
		}
	}()

	var points []coursemodels.LessonPoint
	for rows.Next() {
		var point coursemodels.LessonPoint
		if err := rows.Scan(&point.LessonId, &point.Type); err != nil {
			logs.PrintLog(ctx, "getLessonHeaderNewCourse", fmt.Sprintf("%+v", err))
			return nil, 0, "", false, err
		}
		points = append(points, point)
	}

	if len(points) == 0 {
		logs.PrintLog(ctx, "getlLessonHeaderNewCourse", fmt.Sprintf("no points found in bucket %+v", bucket.Id))
		return nil, 0, "", false, err
	}

	currentLessonId := points[0].LessonId
	currentLessonType := points[0].Type

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

	return &lessonHeader, currentLessonId, currentLessonType, true, nil
}

func (d *Database) GetLastLessonHeader(ctx context.Context, userId int, courseId int) (*dto.LessonDtoHeader, int, string, bool, error) {
	var lessonHeader dto.LessonDtoHeader
	course, err := d.GetCourseById(ctx, courseId)
	if err != nil {
		logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("%+v", err))
		return nil, 0, "", false, err
	}
	lessonHeader.CourseTitle = course.Title
	lessonHeader.CourseId = course.Id

	logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("course %+v", course))
	logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("user id: %+v", userId))

	rows1, err := d.conn.Query(`
			SELECT cp.Lesson_ID, l.type
			FROM LESSON_CHECKPOINT cp
			JOIN LESSON l ON l.ID = cp.Lesson_ID
			WHERE cp.User_ID = $1 AND cp.Course_ID = $2
			ORDER BY cp.Updated_at DESC
		`, userId, courseId)

	if err != nil {
		logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("%+v", err))
		return nil, 0, "", false, err
	}
	defer func() {
		if err := rows1.Close(); err != nil {
			logs.PrintLog(ctx, "SearchCoursesByTitle", fmt.Sprintf("%+v", err))
		}
	}()

	var lessonPoint coursemodels.LessonPoint
	var visitedLessonPointsIds []int
	firstRow := true
	for rows1.Next() {
		if firstRow {
			firstRow = false
			err := rows1.Scan(&lessonPoint.LessonId, &lessonPoint.Type)
			if err != nil {
				logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("%+v", err))
				return nil, 0, "", false, err
			}
			visitedLessonPointsIds = append(visitedLessonPointsIds, lessonPoint.LessonId)
			continue
		}
		var lessonPointId int
		var lessonPointType string
		err := rows1.Scan(&lessonPointId, &lessonPointType)
		if err != nil {
			logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("%+v", err))
			return nil, 0, "", false, err
		}
		logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("visit lesson point %+v", lessonPointId))
		visitedLessonPointsIds = append(visitedLessonPointsIds, lessonPointId)
	}

	if len(visitedLessonPointsIds) == 0 {
		logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("no visited lesson points fo user %+v in course %+v", userId, courseId))
		return d.getLessonHeaderNewCourse(ctx, userId, courseId)
	}

	logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("lesson point %+v", lessonPoint))
	logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("visited lesson points ids %+v", visitedLessonPointsIds))

	var part coursemodels.CoursePart
	var bucket coursemodels.LessonBucket
	err = d.conn.QueryRow(`
			SELECT p.Title, p.Part_Order, p.ID, lb.Title, lb.Lesson_Bucket_Order, lb.ID
			FROM PART p
			JOIN LESSON_BUCKET lb ON lb.Part_ID = p.ID
			JOIN LESSON l ON l.Lesson_Bucket_ID = lb.ID
			WHERE l.ID = $1
		`, lessonPoint.LessonId).Scan(
		&part.Title, &part.Order, &part.Id, &bucket.Title, &bucket.Order, &bucket.Id)

	if err != nil {
		logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("%+v", err))
		return nil, 0, "", false, err
	}

	lessonHeader.Part.Order = part.Order
	lessonHeader.Part.Title = part.Title
	lessonHeader.Bucket.Order = bucket.Order
	lessonHeader.Bucket.Title = bucket.Title

	logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("part %+v, bucket %+v", part, bucket))

	rows2, err := d.conn.Query(`
					SELECT id, type
					FROM LESSON
					WHERE lesson_bucket_id = $1
					ORDER BY Lesson_Order ASC
				`, bucket.Id)

	if err != nil {
		logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("%+v", err))
		return nil, 0, "", false, err
	}
	defer func() {
		if err := rows2.Close(); err != nil {
			logs.PrintLog(ctx, "SearchCoursesByTitle", fmt.Sprintf("%+v", err))
		}
	}()

	var points []coursemodels.LessonPoint
	for rows2.Next() {
		var point coursemodels.LessonPoint
		if err := rows2.Scan(&point.LessonId, &point.Type); err != nil {
			logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("%+v", err))
			return nil, 0, "", false, err
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
		logs.PrintLog(ctx, "GetLastLessonHeader", fmt.Sprintf("no points found in bucket %+v", bucket.Id))
		return nil, 0, "", false, err
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

	return &lessonHeader, lessonPoint.LessonId, lessonPoint.Type, false, nil
}

func (d *Database) GetLessonHeaderByLessonId(ctx context.Context, userId int, currentLessonId int) (*dto.LessonDtoHeader, error) {
	var course coursemodels.Course
	var part coursemodels.CoursePart
	var bucket coursemodels.LessonBucket
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
		logs.PrintLog(ctx, "GetLessonHeaderByLessonId", fmt.Sprintf("%+v", err))
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
		logs.PrintLog(ctx, "GetLessonHeaderByLessonId", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer func() {
		if err := rows1.Close(); err != nil {
			logs.PrintLog(ctx, "SearchCoursesByTitle", fmt.Sprintf("%+v", err))
		}
	}()

	visitedLessonPointsIds := make(map[int]bool)
	for rows1.Next() {
		var visitedLessonPointId int
		if err := rows1.Scan(&visitedLessonPointId); err != nil {
			logs.PrintLog(ctx, "GetLessonHeaderByLessonId", fmt.Sprintf("%+v", err))
			return nil, err
		}
		visitedLessonPointsIds[visitedLessonPointId] = true
	}

	logs.PrintLog(ctx, "GetLessonHeaderByLessonId", fmt.Sprintf("visitedLessonPointsIds %+v", visitedLessonPointsIds))

	rows1, err = d.conn.Query(`
		SELECT lesson_id
		FROM LESSON_CHECKPOINT
		WHERE course_id = $1 and user_id = $2
	`, course.Id, userId)

	if err != nil {
		logs.PrintLog(ctx, "GetLessonHeaderByLessonId", fmt.Sprintf("%+v", err))
		return nil, err
	}

	rows2, err := d.conn.Query(`
		SELECT id, type
		FROM LESSON
		WHERE lesson_bucket_id = $1
	`, bucket.Id)

	if err != nil {
		logs.PrintLog(ctx, "GetLessonHeaderByLessonId", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer func() {
		if err := rows1.Close(); err != nil {
			logs.PrintLog(ctx, "SearchCoursesByTitle", fmt.Sprintf("%+v", err))
		}
	}()

	var points []*dto.LessonPointDTO
	for rows2.Next() {
		var point coursemodels.LessonPoint
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

func (d *Database) GetLessonFooters(ctx context.Context, currentLessonId int) ([]int, error) {
	footers := []int{-1, -1, -1}

	var currentLessonOrder int
	var currentBucket coursemodels.LessonBucket
	var currentPart coursemodels.CoursePart
	var currentCourse coursemodels.Course
	err := d.conn.QueryRow(`
			SELECT l.lesson_order, lb.id, lb.lesson_bucket_order, lb.part_id, p.part_order, c.id
			FROM LESSON l
			JOIN LESSON_BUCKET lb ON l.lesson_bucket_id = lb.id
			JOIN PART p ON lb.part_id = p.id
			JOIN COURSE c ON p.course_id = c.id
			WHERE l.id = $1
		`, currentLessonId).Scan(&currentLessonOrder, &currentBucket.Id, &currentBucket.Order, &currentBucket.PartId, &currentPart.Order, &currentCourse.Id)
	if err != nil {
		logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("%+v", err))
		return nil, err
	}

	footers[1] = currentLessonId
	logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("current lesson order: %d", currentLessonOrder))
	logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("current bucket order: %d", currentBucket.Order))
	logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("current part order: %d", currentPart.Order))
	logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("current course id: %d", currentCourse.Id))

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
	defer func() {
		if err := rows.Close(); err != nil {
			logs.PrintLog(ctx, "SearchCoursesByTitle", fmt.Sprintf("%+v", err))
		}
	}()

	for rows.Next() {
		var footer struct {
			LessonId int
			Order    int
		}

		if err := rows.Scan(&footer.LessonId, &footer.Order); err != nil {
			logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("%+v", err))
			return nil, err
		}

		switch diff := footer.Order - currentLessonOrder; {
		case diff == -1:
			footers[0] = footer.LessonId
		case diff == 1:
			footers[2] = footer.LessonId
		}
	}

	if footers[0] == -1 {
		prevLessonId := -1
		err := d.conn.QueryRow(`
				SELECT l.id
				FROM LESSON_BUCKET lb
				JOIN LESSON l ON l.lesson_bucket_id = lb.id
				WHERE lb.lesson_bucket_order = $1 and lb.part_id = $2
				ORDER BY l.Lesson_Order DESC
				LIMIT 1
			`, currentBucket.Order-1, currentBucket.PartId).Scan(&prevLessonId)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("%+v", err))
				return nil, err
			}
		}

		if prevLessonId == -1 && currentPart.Order > 1 {
			err := d.conn.QueryRow(`
					SELECT l.id
					FROM LESSON_BUCKET lb
					JOIN LESSON l ON l.lesson_bucket_id = lb.id
					JOIN PART p ON p.id = lb.part_id
					JOIN COURSE c ON p.course_id = c.id
					WHERE p.part_order = $1 AND c.id = $2
					ORDER BY lb.lesson_bucket_order DESC, l.Lesson_Order DESC
					LIMIT 1
				`, currentPart.Order-1, currentCourse.Id).Scan(&prevLessonId)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("%+v", err))
					return nil, err
				}
			}
		}

		footers[0] = prevLessonId
	}

	if footers[2] == -1 {
		nextLessonId := -1
		err := d.conn.QueryRow(`
				SELECT l.id
				FROM LESSON_BUCKET lb
				JOIN LESSON l ON l.lesson_bucket_id = lb.id
				WHERE lb.lesson_bucket_order = $1 and lb.part_id = $2
				ORDER BY l.Lesson_Order ASC
				LIMIT 1
			`, currentBucket.Order+1, currentBucket.PartId).Scan(&nextLessonId)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("%+v", err))
				return nil, err
			}
		}

		if nextLessonId == -1 {
			err := d.conn.QueryRow(`
					SELECT l.id
					FROM LESSON_BUCKET lb
					JOIN LESSON l ON l.lesson_bucket_id = lb.id
					JOIN PART p ON p.id = lb.part_id
					JOIN COURSE c ON p.course_id = c.id
					WHERE p.part_order = $1 AND c.id = $2
					ORDER BY lb.lesson_bucket_order ASC, l.Lesson_Order ASC
					LIMIT 1
				`, currentPart.Order+1, currentCourse.Id).Scan(&nextLessonId)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("%+v", err))
					return nil, err
				}
			}
		}

		footers[2] = nextLessonId
	}

	logs.PrintLog(ctx, "GetLessonFooters", fmt.Sprintf("footers: %+v", footers))

	return footers, nil
}

func (d *Database) IsMiddle(ctx context.Context, userId int, courseId int) (bool, error) {
	var exists bool
	err := d.conn.QueryRow(`
			SELECT EXISTS (
			SELECT 1 FROM SENDED_MAILS 
			WHERE User_ID = $1 AND Course_ID = $2
		)
		`, userId, courseId).Scan(&exists)

	if err != nil {
		logs.PrintLog(ctx, "IsMiddle", fmt.Sprintf("%+v", err))
		return false, err
	}

	if exists {
		return false, nil
	}

	var countMarkedLessons int
	err = d.conn.QueryRow(`
			SELECT COUNT(Lesson_ID)
			FROM LESSON_CHECKPOINT
			WHERE User_ID = $1 AND Course_ID = $2
		`, userId, courseId).Scan(&countMarkedLessons)

	if err != nil {
		logs.PrintLog(ctx, "IsMiddle", fmt.Sprintf("%+v", err))
		return false, err
	}

	var countLessons int
	err = d.conn.QueryRow(`
			SELECT COUNT(l.ID)
			FROM LESSON l
			JOIN LESSON_BUCKET lb ON l.LESSON_BUCKET_ID = lb.ID
			JOIN PART p ON lb.PART_ID = p.ID
			WHERE p.Course_ID = $1
		`, courseId).Scan(&countLessons)

	if err != nil {
		logs.PrintLog(ctx, "IsMiddle", fmt.Sprintf("%+v", err))
		return false, err
	}

	if countLessons*2 > countMarkedLessons {
		_, err = d.conn.Exec(
			"INSERT INTO SENDED_MAILS (User_ID, Course_ID) VALUES ($1, $2)",
			userId, courseId)
		if err != nil {
			logs.PrintLog(ctx, "IsMiddle", fmt.Sprintf("%+v", err))
			return false, err
		}
		logs.PrintLog(ctx, "IsMiddle", "marked that sended mail")
		return true, nil
	}

	return false, nil
}

func (d *Database) IsWelcomeCourseMailSended(ctx context.Context, userId int, courseId int) (bool, error) {
	var exists bool
	err := d.conn.QueryRow(`
			SELECT EXISTS (
			SELECT 1 FROM WELCOME_COURSE_SENDED_MAILS 
			WHERE User_ID = $1 AND Course_ID = $2
		)
		`, userId, courseId).Scan(&exists)

	if err != nil {
		logs.PrintLog(ctx, "IsWelcomeCourseMailSended", fmt.Sprintf("%+v", err))
		return false, err
	}

	if exists {
		return false, nil
	}

	_, err = d.conn.Exec(
		"INSERT INTO WELCOME_COURSE_SENDED_MAILS (User_ID, Course_ID) VALUES ($1, $2)",
		userId, courseId)
	if err != nil {
		logs.PrintLog(ctx, "IsWelcomeCourseMailSended", fmt.Sprintf("%+v", err))
		return false, err
	}
	logs.PrintLog(ctx, "IsWelcomeCourseMailSended", "marked that sended mail")

	return true, nil
}
