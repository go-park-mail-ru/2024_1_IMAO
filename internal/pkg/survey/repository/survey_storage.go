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

func (survey *SurveyStorage) insertAnswer(ctx context.Context, tx pgx.Tx, userID, surveyID, answerNum, answerValue uint) error {
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

func (survey *SurveyStorage) SaveSurveyResults(ctx context.Context, surveyAnswers []*models.SurveyAnswer) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	for i := 0; i < len(surveyAnswers); i++ {

		err := pgx.BeginFunc(ctx, survey.pool, func(tx pgx.Tx) error {
			err := survey.insertAnswer(ctx, tx, surveyAnswers[i].UserID, surveyAnswers[i].SurveyID, surveyAnswers[i].AnswerNum, surveyAnswers[i].AnswerValue)

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

func (survey *SurveyStorage) GetStatics() {

}
