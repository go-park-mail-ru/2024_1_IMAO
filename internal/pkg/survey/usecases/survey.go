package usecases

import (
	"context"

	models "github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type SurveyStorageInterface interface {
	SaveSurveyResults(ctx context.Context, surveyAnswers models.SurveyAnswersList) error
	GetResults(ctx context.Context, userID, surveyID uint) (bool, error)
	GetStatics()
}
