package postgres

import (
	"context"
	"fmt"

	coursemodels "skillForce/internal/models/course"
	"skillForce/pkg/logs"
)

func (d *Database) GetBucketCourses(ctx context.Context) ([]*coursemodels.Course, error) {
	//TODO: можно заморочиться и сделать самописную пагинацию через LIMIT OFFSET
	var bucketCourses []*coursemodels.Course
	rows, err := d.conn.Query("SELECT id, creator_user_id, title, description, avatar_src, price, time_to_pass FROM course LIMIT 16")
	if err != nil {
		logs.PrintLog(ctx, "GetBucketCourses", fmt.Sprintf("%+v", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var course coursemodels.Course
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

func (d *Database) GetBucketByLessonId(ctx context.Context, currentLessonId int) (*coursemodels.LessonBucket, error) {
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

	bucket := &coursemodels.LessonBucket{
		Id: bucketId,
	}
	return bucket, nil
}

func (d *Database) GetBucketLessons(ctx context.Context, userId int, courseId int, bucketId int) ([]*coursemodels.LessonPoint, error) {
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

	var lessons []*coursemodels.LessonPoint
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
		var lesson coursemodels.LessonPoint
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

func (d *Database) GetCourseById(ctx context.Context, courseId int) (*coursemodels.Course, error) {
	var course coursemodels.Course
	err := d.conn.QueryRow("SELECT id, creator_user_id, title, description, avatar_src, price, time_to_pass FROM course WHERE id = $1", courseId).Scan(
		&course.Id, &course.CreatorId, &course.Title, &course.Description, &course.ScrImage, &course.Price, &course.TimeToPass)
	if err != nil {
		logs.PrintLog(ctx, "GetCourseById", fmt.Sprintf("%+v", err))
		return nil, err
	}
	logs.PrintLog(ctx, "GetCourseById", fmt.Sprintf("get course %+v from db", course))
	return &course, nil
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

func (d *Database) GetLessonVideo(ctx context.Context, currentLessonId int) ([]string, error) {
	var videoSrc string
	err := d.conn.QueryRow(`
			SELECT video_src
			FROM VIDEO_LESSON
			WHERE lesson_ID = $1
		`, currentLessonId).Scan(&videoSrc)
	if err != nil {
		logs.PrintLog(ctx, "GetLessonVideos", fmt.Sprintf("%+v", err))
		return nil, err
	}
	return []string{videoSrc}, nil
}

func (d *Database) GetLessonById(ctx context.Context, lessonId int) (*coursemodels.LessonPoint, error) {
	var lesson coursemodels.LessonPoint
	err := d.conn.QueryRow(`
			SELECT id, title, type
			FROM LESSON
			WHERE id = $1
		`, lessonId).Scan(&lesson.LessonId, &lesson.Title, &lesson.Type)
	if err != nil {
		logs.PrintLog(ctx, "GetLessonById", fmt.Sprintf("%+v", err))
		return nil, err
	}
	return &lesson, nil
}

func (d *Database) GetCourseParts(ctx context.Context, courseId int) ([]*coursemodels.CoursePart, error) {
	var courseParts []*coursemodels.CoursePart
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
		var coursePart coursemodels.CoursePart
		if err := rows.Scan(&coursePart.Id, &coursePart.Title); err != nil {
			logs.PrintLog(ctx, "GetCourseParts", fmt.Sprintf("%+v", err))
			return nil, err
		}
		logs.PrintLog(ctx, "GetCourseParts", fmt.Sprintf("get course part %+v", coursePart))
		courseParts = append(courseParts, &coursePart)
	}
	return courseParts, nil
}

func (d *Database) GetPartBuckets(ctx context.Context, partId int) ([]*coursemodels.LessonBucket, error) {
	var buckets []*coursemodels.LessonBucket
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
		var bucket coursemodels.LessonBucket
		if err := rows.Scan(&bucket.Id, &bucket.Title); err != nil {
			logs.PrintLog(ctx, "GetPartBuckets", fmt.Sprintf("%+v", err))
			return nil, err
		}
		logs.PrintLog(ctx, "GetPartBuckets", fmt.Sprintf("get bucket %+v", bucket))
		buckets = append(buckets, &bucket)
	}
	return buckets, nil
}

func (d *Database) GetVideoUrl(ctx context.Context, lessonId int) (string, error) {
	var videoUrl string
	err := d.conn.QueryRow("SELECT video_src FROM video_lesson WHERE lesson_id = $1", lessonId).Scan(&videoUrl)
	if err != nil {
		logs.PrintLog(ctx, "GetVideo", fmt.Sprintf("%+v", err))
		return "", err
	}
	return videoUrl, nil
}
