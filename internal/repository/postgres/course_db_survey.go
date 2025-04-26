package postgres

import (
	"context"
	"fmt"
	coursemodels "skillForce/internal/models/course"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/logs"
)

func (d *Database) CreateSurvey(ctx context.Context, survey *coursemodels.Survey, userProfile *usermodels.UserProfile) error {
	err := d.conn.QueryRow("INSERT INTO survey DEFAULT VALUES RETURNING id").Scan(&survey.Id)
	if err != nil {
		logs.PrintLog(ctx, "CreateSurvey", fmt.Sprintf("%+v", err))
		return err
	}

	for _, question := range survey.Questions {
		query := `
		INSERT INTO survey_question (survey_id, question, left_desc, right_desc, metric_type)
		VALUES ($1, $2, $3, $4, $5)	
	`

		_, err = d.conn.Exec(
			query,
			survey.Id,
			question.Question,
			question.LeftLebal,
			question.RightLebal,
			question.Metric,
		)

		if err != nil {
			logs.PrintLog(ctx, "CreateSurvey", fmt.Sprintf("%+v", err))
			return err
		}
	}

	return nil
}

func (d *Database) GetSurvey(ctx context.Context) (*coursemodels.Survey, error) {
	survey := coursemodels.Survey{}
	err := d.conn.QueryRow("SELECT id FROM survey ORDER BY id DESC LIMIT 1").Scan(&survey.Id)
	if err != nil {
		logs.PrintLog(ctx, "GetSurvey", fmt.Sprintf("%+v", err))
		return nil, err
	}

	logs.PrintLog(ctx, "GetSurvey", fmt.Sprintf("get survey id: %+v", survey))

	rows, err := d.conn.Query(`
			SELECT id, metric_type, question, left_desc, right_desc
			FROM survey_question
			WHERE survey_id = $1
		`, survey.Id)
	if err != nil {
		logs.PrintLog(ctx, "GetSurvey", fmt.Sprintf("%+v", err))
		return nil, err
	}

	for rows.Next() {
		var question coursemodels.Question
		if err := rows.Scan(&question.QuestionId, &question.Metric, &question.Question, &question.LeftLebal, &question.RightLebal); err != nil {
			logs.PrintLog(ctx, "GetSurvey", fmt.Sprintf("%+v", err))
			return nil, err
		}
		logs.PrintLog(ctx, "GetSurvey", fmt.Sprintf("get question: %+v", question))
		survey.Questions = append(survey.Questions, question)
	}

	return &survey, nil
}
