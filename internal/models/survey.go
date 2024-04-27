package models

type SurveyAnswer struct {
	AnswerNum   uint `json:"answerNum"`
	AnswerValue uint `json:"answerValue"`
}

type SurveyAnswersList struct {
	UserID   uint            `json:"userId"`
	SurveyID uint            `json:"surveyId"`
	Survey   []*SurveyAnswer `json:"survey"`
}
