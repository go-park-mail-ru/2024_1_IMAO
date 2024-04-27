package models

type SurveyAnswer struct {
	UserID      uint `json:"userId"`
	SurveyID    uint `json:"surveyId"`
	AnswerNum   uint `json:"answerNum"`
	AnswerValue uint `json:"answerValue"`
}

type SurveyAnswersList struct {
	Survey []*SurveyAnswer `json:"survey"`
}
