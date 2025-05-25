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

func (d *Database) GetRating(ctx context.Context, userId int, courseId int) (*dto.Raiting, error) {
	query := `
		SELECT 
			u.name, u.avatar_src, u.id,
			COUNT(lc.id) AS lessons_completed
		FROM 
			lesson_checkpoint lc
		JOIN 
			usertable u ON lc.user_id = u.id
		WHERE 
			lc.course_id = $1
		GROUP BY 
			u.id, u.name
		ORDER BY 
			lessons_completed DESC
	`

	rows, err := d.conn.Query(query, courseId)

	if err != nil {
		return nil, fmt.Errorf("failed to query database: %v", err)
	}
	defer rows.Close()

	var rating []dto.RaitingItem

	for rows.Next() {
		var item dto.RaitingItem
		var amountCompletedLessons int
		var newUserId int
		err := rows.Scan(&item.User.Name, &item.User.AvatarSrc, &newUserId, &amountCompletedLessons)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		item.Rating = amountCompletedLessons
		rating = append(rating, item)
	}

	if err := rows.Err(); err != nil {
		logs.PrintLog(ctx, "GetRating", fmt.Sprintf("%+v", err))
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	resultRatingList := dto.Raiting{Rating: rating}

	return &resultRatingList, nil
}

func (d *Database) GetStatistic(ctx context.Context, userId int, courseId int) (*dto.UserStats, error) {
	stats := &dto.UserStats{}

	err := d.conn.QueryRowContext(ctx, `
		SELECT 
            COUNT(CASE WHEN l.type = 'text' THEN 1 END) as text_lessons,
            COUNT(CASE WHEN l.type = 'video' THEN 1 END) as video_lessons
        FROM lesson l
        JOIN lesson_bucket lb ON l.lesson_bucket_id = lb.id
        JOIN part p ON lb.part_id = p.id
        WHERE p.course_id = $1`, courseId).Scan(&stats.AmountTextLessons, &stats.AmountVideoLessons)
	if err != nil {
		return nil, fmt.Errorf("failed to get total lessons count: %w", err)
	}

	err = d.conn.QueryRowContext(ctx, `
		SELECT 
			COUNT(CASE WHEN l.type = 'text' THEN 1 END) as completed_text,
			COUNT(CASE WHEN l.type = 'video' THEN 1 END) as completed_video
		FROM lesson_checkpoint lc
		JOIN lesson l ON lc.lesson_id = l.id
		WHERE lc.user_id = $1 AND lc.course_id = $2`, userId, courseId).Scan(&stats.CompletedTextLessons, &stats.CompletedVideoLessons)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed lessons count: %w", err)
	}

	totalLessons := stats.AmountTextLessons + stats.AmountVideoLessons
	completedLessons := stats.CompletedTextLessons + stats.CompletedVideoLessons

	if totalLessons > 0 {
		stats.Percentage = (completedLessons * 100) / totalLessons
	}

	logs.PrintLog(ctx, "GetStatistic", fmt.Sprintf("stats: %+v", stats))

	return stats, nil
}
