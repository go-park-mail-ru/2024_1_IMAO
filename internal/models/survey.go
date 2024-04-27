package models

type SurveyAnswer struct {
	UserID      uint `json:"userID"`
	SurveyID    uint `json:"surveyID"`
	AnswerNum   uint `json:"answerNum"`
	AnswerValue uint `json:"answerValue"`
}
