package delivery

type SurveyOkResponse struct {
	Code int `json:"code"`
	Body any `json:"body"`
}

type SurveyCheckResponse struct {
	Code      int  `json:"code"`
	IsChecked bool `json:"isChecked"`
}

func NewSurveyOkResponse(body any) *SurveyOkResponse {
	return &SurveyOkResponse{
		Code: 200,
		Body: body,
	}
}

func NewSurveyCheckResponse(isChecked bool) *SurveyCheckResponse {
	return &SurveyCheckResponse{
		Code:      200,
		IsChecked: isChecked,
	}
}
