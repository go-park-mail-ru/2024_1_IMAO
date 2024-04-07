package delivery

import (
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
)

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
		Code:  responses.StatusOk,
		Items: adverts,
	}
}

func NewAdvertsErrResponse(code int, status string) *AdvertsErrResponse {
	return &AdvertsErrResponse{
		Code:   code,
		Status: status,
	}
}
