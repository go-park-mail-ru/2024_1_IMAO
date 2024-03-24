package responses

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
)

type AdvertsOkResponse struct {
	Code  int                        `json:"code"`
	Items []*storage.ReturningAdvert `json:"items"`
}

type AdvertsErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewAdvertsOkResponse(adverts []*storage.ReturningAdvert) *AdvertsOkResponse {
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
