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

func (d *Database) GetMetricCount(ctx context.Context, surveyId int, metric string) (int, error) {
	var count int
	err := d.conn.QueryRow(`
		SELECT COUNT(sq.id) 
		FROM survey_question sq
		JOIN survey s ON sq.survey_id = s.id
		JOIN survey_answer sa ON sq.id = sa.question_id
		WHERE s.id = $1 AND sq.metric_type = $2
		`, surveyId, metric).Scan(&count)
	if err != nil {
		logs.PrintLog(ctx, "GetMetricCount", fmt.Sprintf("%+v", err))
		return 0, err
	}
	return count, nil
}

func (d *Database) GetMetricAvg(ctx context.Context, surveyId int, metric string) (float64, error) {
	var avg float64
	err := d.conn.QueryRow(`
		SELECT COALESCE(AVG(sa.answer), 0.0) AS avg
		FROM survey_question sq
		JOIN survey s ON sq.survey_id = s.id
		JOIN survey_answer sa ON sq.id = sa.question_id
		WHERE s.id = $1 AND sq.metric_type = $2
		`, surveyId, metric).Scan(&avg)
	if err != nil {
		logs.PrintLog(ctx, "GetMetricAvg", fmt.Sprintf("%+v", err))
		return 0, err
	}
	return avg, nil
}

func (d *Database) GetMetricDistribution(ctx context.Context, surveyId int, metric string) ([]int, error) {
	var distribution []int

	logs.PrintLog(ctx, "GetMetricDistribution", fmt.Sprintf("metric: %+v", metric))
	logs.PrintLog(ctx, "GetMetricDistribution", fmt.Sprintf("surveyId: %+v", surveyId))

	var countAnswers int
	err := d.conn.QueryRow(`
		SELECT COUNT(sa.answer) 
		FROM survey_question sq
		JOIN survey s ON sq.survey_id = s.id
		JOIN survey_answer sa ON sq.id = sa.question_id
		WHERE s.id = $1 AND sq.metric_type = $2
		`, surveyId, metric).Scan(&countAnswers)
	if err != nil {
		logs.PrintLog(ctx, "GetMetricDistribution", fmt.Sprintf("%+v", err))
		return nil, err
	}

	for i := 0; i < 11; i++ {
		var countAnswer int
		err := d.conn.QueryRow(`
			SELECT COUNT(sa.answer) 
			FROM survey_question sq
			JOIN survey s ON sq.survey_id = s.id
			JOIN survey_answer sa ON sq.id = sa.question_id
			WHERE s.id = $1 AND sq.metric_type = $2 AND sa.answer = $3
			`, surveyId, metric, i).Scan(&countAnswer)
		if err != nil {
			logs.PrintLog(ctx, "GetMetricDistribution", fmt.Sprintf("%+v", err))
			return nil, err
		}
		distribution = append(distribution, (countAnswer*100)/countAnswers)
	}
	return distribution, nil
}

func (d *Database) GetMetrics(ctx context.Context, metric string) (*coursemodels.SurveyMetric, error) {
	survey := coursemodels.Survey{}
	err := d.conn.QueryRow("SELECT id FROM survey ORDER BY id DESC LIMIT 1").Scan(&survey.Id)
	if err != nil {
		logs.PrintLog(ctx, "GetSurvey", fmt.Sprintf("%+v", err))
		return nil, err
	}

	surveyMetric := coursemodels.SurveyMetric{}

	count, err := d.GetMetricCount(ctx, survey.Id, metric)
	if err != nil {
		logs.PrintLog(ctx, "GetCSATMetrics", fmt.Sprintf("%+v", err))
		return nil, err
	}
	surveyMetric.Type = metric
	surveyMetric.Count = count

	if count == 0 {
		return &surveyMetric, nil
	}

	avg, err := d.GetMetricAvg(ctx, survey.Id, metric)
	if err != nil {
		logs.PrintLog(ctx, "GetCSATMetrics", fmt.Sprintf("%+v", err))
		return nil, err
	}
	surveyMetric.Avg = avg

	distribution, err := d.GetMetricDistribution(ctx, survey.Id, metric)
	if err != nil {
		logs.PrintLog(ctx, "GetCSATMetrics", fmt.Sprintf("%+v", err))
		return nil, err
	}
	surveyMetric.Distribution = distribution

	return &surveyMetric, nil
}
