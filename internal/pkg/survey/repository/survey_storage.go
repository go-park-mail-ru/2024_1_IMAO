package repository

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type SurveyStorage struct {
	pool *pgxpool.Pool
}

func NewSurveyStorage(pool *pgxpool.Pool) *SurveyStorage {
	return &SurveyStorage{
		pool: pool,
	}
}

func (survey *SurveyStorage) insertAnswer(ctx context.Context, tx pgx.Tx, userID, surveyID, answerNum,
	answerValue uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLInsertAnswer := `
	INSERT INTO public.answer(
		user_id, survey_id, answer_num, answer_value)
		VALUES ($1, $2, $3, $4);`

	logging.LogInfo(logger, "INSERT INTO answer")

	var err error

	_, err = tx.Exec(ctx, SQLInsertAnswer, userID, surveyID, answerNum, answerValue)

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing insertAnswer query, err=%v", err))

		return err
	}

	return nil
}

func (survey *SurveyStorage) SaveSurveyResults(ctx context.Context, surveyAnswersList models.SurveyAnswersList) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	for i := 0; i < len(surveyAnswersList.Survey); i++ {

		err := pgx.BeginFunc(ctx, survey.pool, func(tx pgx.Tx) error {
			err := survey.insertAnswer(ctx, tx, surveyAnswersList.UserID, surveyAnswersList.SurveyID,
				surveyAnswersList.Survey[i].AnswerNum, surveyAnswersList.Survey[i].AnswerValue)

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while inserting answers, err=%v", err))

			return err
		}
	}

	return nil
}

func (survey *SurveyStorage) selectFromUserSurvey(ctx context.Context, tx pgx.Tx, userID, surveyID uint) (bool, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLSelectFromUserSurvey := `SELECT EXISTS(SELECT 1 FROM public.user_survey WHERE user_id=$1 AND survey_id = $2);`

	logging.LogInfo(logger, "SELECT FROM user_survey")

	userLine := tx.QueryRow(ctx, SQLSelectFromUserSurvey, userID, surveyID)

	var exists bool

	if err := userLine.Scan(&exists); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning UserSurvey exists, err=%v", err))

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
		logging.LogError(logger, fmt.Errorf("error while executing user exists query, err=%v", err))

		return true, err
	}

	return exists, nil
}

func (survey *SurveyStorage) getStatics(ctx context.Context, tx pgx.Tx) (*models.SurveyResults, error) {

	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLGetStatics := `
	SELECT answer_num, answer_value, COUNT(*) AS answer_count
	FROM answer
	GROUP BY answer_num, answer_value
	ORDER BY answer_num, answer_value;`

	logging.LogInfo(logger, "SELECT FROM advert_image")

	rows, err := tx.Query(ctx, SQLGetStatics)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select statistics, err=%v", err))

		return nil, err
	}
	defer rows.Close()

	var surveyResults *models.SurveyResults // ПОТЕНЦИАЛЬНАЯ ПРОБЛЕМА

	surveyResults.SurveyTitle = "Название опроса"
	surveyResults.SurveyDescription = "Описание опроса"

	var questionNumber uint = 1

	var questionResults *models.QuestionResults

	for rows.Next() {
		var answerNumber uint
		var answerValue uint
		var answerCount uint

		if err := rows.Scan(&answerNumber, &answerValue, &answerCount); err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while scanning rows of answers, err=%v", err))

			return nil, err
		}

		if answerNumber == questionNumber {
			questionResults.QuestionResults = append(questionResults.QuestionResults, answerCount)
		} else {
			surveyResults.Results = append(surveyResults.Results, questionResults)

			questionResults = new(models.QuestionResults)
		}

	}

	return surveyResults, nil
}

func (survey *SurveyStorage) GetStatics(ctx context.Context) (*models.SurveyResults, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var surveyResults *models.SurveyResults

	err := pgx.BeginFunc(ctx, survey.pool, func(tx pgx.Tx) error {
		surveyResultsInternal, err := survey.getStatics(ctx, tx)
		surveyResults = surveyResultsInternal

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing user exists query, err=%v", err))

		return nil, err
	}

	return surveyResults, nil

}
