package usecases

type SurveyStorageInterface interface {
	SaveSurveyResults()
	GetResults()
	GetStatics()
}
