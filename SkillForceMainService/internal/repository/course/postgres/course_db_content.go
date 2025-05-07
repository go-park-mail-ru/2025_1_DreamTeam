package postgres

import (
	"context"
	"fmt"

	"skillForce/pkg/logs"
)

func (d *Database) GetVideoUrl(ctx context.Context, lessonId int) (string, error) {
	var videoUrl string
	err := d.conn.QueryRow("SELECT video_src FROM video_lesson WHERE lesson_id = $1", lessonId).Scan(&videoUrl)
	if err != nil {
		logs.PrintLog(ctx, "GetVideo", fmt.Sprintf("%+v", err))
		return "", err
	}
	return videoUrl, nil
}
