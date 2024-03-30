package responses

type AdvertsOkResponse struct {
	Code  int `json:"code"`
	Items any `json:"items"`
}

type AdvertsErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewAdvertsOkResponse(adverts any) *AdvertsOkResponse {
	return &AdvertsOkResponse{
		Code:  StatusOk,
		Items: adverts,
	}
}

func NewAdvertsErrResponse(code int, status string) *AdvertsErrResponse {
	return &AdvertsErrResponse{
		Code:   code,
		Status: status,
	}
}
