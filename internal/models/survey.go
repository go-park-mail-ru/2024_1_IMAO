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

type SurveyResults struct {
	SurveyTitle       uint               `json:"surveyTitle"`
	SurveyDescription uint               `json:"surveyDescription"`
	Results           []*QuestionResults `json:"results"`
}

type QuestionResults struct {
	QuestionResults []uint
}
