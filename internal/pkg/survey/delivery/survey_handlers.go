package delivery

import (
	surveyusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/survey/usecases"
	userusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
)

type SurveyHandler struct {
	userStorage   userusecases.UsersStorageInterface
	surveyStorage surveyusecases.SurveyStorageInterface
}

func NewSurveyHandler(userStorage userusecases.UsersStorageInterface,
	surveyStorage surveyusecases.SurveyStorageInterface) *SurveyHandler {
	return &SurveyHandler{
		userStorage:   userStorage,
		surveyStorage: surveyStorage,
	}
}
