package repository

import (
	"context"
	"fmt"
	"time"

	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	questionsNum = 5
)

type SurveyStorage struct {
	pool    *pgxpool.Pool
	metrics *mymetrics.DatabaseMetrics
}

func NewSurveyStorage(pool *pgxpool.Pool, metrics *mymetrics.DatabaseMetrics) *SurveyStorage {
	return &SurveyStorage{
		pool:    pool,
		metrics: metrics,
	}
}

func (survey *SurveyStorage) insertUserSurvey(ctx context.Context, tx pgx.Tx, userID, surveyID uint) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLInsertAnswer := `
	INSERT INTO public.user_survey(
		user_id, survey_id)
		VALUES ($1, $2);`

	logging.LogInfo(logger, "INSERT INTO user_survey")

	var err error

	start := time.Now()

	_, err = tx.Exec(ctx, SQLInsertAnswer, userID, surveyID)

	survey.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while executing insertUserSurvey query, err=%w", err))
		survey.metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

func (survey *SurveyStorage) insertAnswer(ctx context.Context, tx pgx.Tx, userID, surveyID, answerNum,
	answerValue uint) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLInsertAnswer := `
	INSERT INTO public.answer(
		user_id, survey_id, answer_num, answer_value)
		VALUES ($1, $2, $3, $4);`

	logging.LogInfo(logger, "INSERT INTO answer")

	var err error

	start := time.Now()

	_, err = tx.Exec(ctx, SQLInsertAnswer, userID, surveyID, answerNum, answerValue)

	survey.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing insertAnswer query, err=%w",
			err))
		survey.metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

func (survey *SurveyStorage) InsertUserSurvey(ctx context.Context, surveyAnswersList models.SurveyAnswersList) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, survey.pool, func(tx pgx.Tx) error {
		err := survey.insertUserSurvey(ctx, tx, surveyAnswersList.UserID, surveyAnswersList.SurveyID)

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while inserting user_survey, err=%w", err))

		return err
	}

	return nil
}

func (survey *SurveyStorage) SaveSurveyResults(ctx context.Context, surveyAnswersList models.SurveyAnswersList) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	fail := survey.InsertUserSurvey(ctx, surveyAnswersList)

	if fail != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while inserting UserSurvey, err=%w", fail))

		return fail
	}

	for i := 0; i < len(surveyAnswersList.Survey); i++ {
		err := pgx.BeginFunc(ctx, survey.pool, func(tx pgx.Tx) error {
			err := survey.insertAnswer(ctx, tx, surveyAnswersList.UserID, surveyAnswersList.SurveyID,
				surveyAnswersList.Survey[i].AnswerNum, surveyAnswersList.Survey[i].AnswerValue)

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while inserting answers, err=%w", err))

			return err
		}
	}

	return nil
}

func (survey *SurveyStorage) selectFromUserSurvey(ctx context.Context, tx pgx.Tx, userID, surveyID uint) (bool, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLSelectFromUserSurvey := `SELECT EXISTS(SELECT 1 FROM public.user_survey WHERE user_id=$1 AND survey_id = $2);`

	logging.LogInfo(logger, "SELECT FROM user_survey")

	start := time.Now()

	userLine := tx.QueryRow(ctx, SQLSelectFromUserSurvey, userID, surveyID)

	survey.metrics.AddDuration(funcName, time.Since(start))

	var exists bool

	if err := userLine.Scan(&exists); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning UserSurvey exists, err=%w", err))

		return false, err
	}

	return exists, nil
}

func (survey *SurveyStorage) GetResults(ctx context.Context, userID, surveyID uint) (bool, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var exists bool

	err := pgx.BeginFunc(ctx, survey.pool, func(tx pgx.Tx) error {
		userExists, err := survey.selectFromUserSurvey(ctx, tx, userID, surveyID)
		exists = userExists

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing user exists query, err=%w", err))

		return true, err
	}

	return exists, nil
}

func (survey *SurveyStorage) getStatics(ctx context.Context, tx pgx.Tx,
	surveyInstance *models.Survey) (*models.SurveyResults, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLGetStatics := `
	SELECT answer_num, answer_value, COUNT(*) AS answer_count
	FROM answer
	GROUP BY answer_num, answer_value
	ORDER BY answer_num, answer_value;`

	logging.LogInfo(logger, "SELECT FROM answer")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLGetStatics)

	survey.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select statistics, err=%w", err))
		survey.metrics.IncreaseErrors(funcName)

		return nil, err
	}

	defer rows.Close()

	surveyResults := &models.SurveyResults{
		SurveyTitle:       surveyInstance.SurveyTitle,
		SurveyDescription: surveyInstance.SurveyTitle,
		Results:           make([]*models.QuestionResults, surveyInstance.QuestionNumber),
	}

	for i := 0; i < len(surveyResults.Results); i++ {
		surveyResults.Results[i] = &models.QuestionResults{
			QuestionResults: make([]uint, questionsNum),
		}
	}

	for rows.Next() {
		var (
			answerNumber uint
			answerValue  uint
			answerCount  uint
		)

		if err := rows.Scan(&answerNumber, &answerValue, &answerCount); err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while scanning rows of answers, err=%w", err))

			return nil, err
		}

		surveyResults.Results[answerNumber-1].QuestionResults[answerValue-1] = answerCount
	}

	return surveyResults, nil
}

func (survey *SurveyStorage) getSurvey(ctx context.Context, tx pgx.Tx, surveyID uint) (*models.Survey, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLSelectFromSurvey := `SELECT title, description, question_number
	FROM public.survey WHERE id = $1;`

	logging.LogInfo(logger, "SELECT FROM survey")

	start := time.Now()

	userLine := tx.QueryRow(ctx, SQLSelectFromSurvey, surveyID)

	survey.metrics.AddDuration(funcName, time.Since(start))

	surveyInstance := &models.Survey{}

	if err := userLine.Scan(&surveyInstance.SurveyTitle, &surveyInstance.SurveyDescription,
		&surveyInstance.QuestionNumber); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning survey, err=%w", err))

		return nil, err
	}

	return surveyInstance, nil
}

func (survey *SurveyStorage) GetStatics(ctx context.Context) (*models.SurveyResults, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var (
		surveyResults  *models.SurveyResults
		surveyInstance *models.Survey
	)

	err := pgx.BeginFunc(ctx, survey.pool, func(tx pgx.Tx) error {
		surveyInstanceInternal, err := survey.getSurvey(ctx, tx, 1)
		surveyInstance = surveyInstanceInternal

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing user exists query, err=%w", err))

		return nil, err
	}

	err = pgx.BeginFunc(ctx, survey.pool, func(tx pgx.Tx) error {
		surveyResultsInternal, err := survey.getStatics(ctx, tx, surveyInstance)
		surveyResults = surveyResultsInternal

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing user exists query, err=%w", err))

		return nil, err
	}

	return surveyResults, nil
}
