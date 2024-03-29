package responses

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
)

type CartOkResponse struct {
	Code  int               `json:"code"`
	Items []*storage.Advert `json:"items"`
}

type CartErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewCartOkResponse(adverts []*storage.Advert) *CartOkResponse {
	return &CartOkResponse{
		Code:  StatusOk,
		Items: adverts,
	}
}

func NewCartErrResponse(code int, status string) *CartErrResponse {
	return &CartErrResponse{
		Code:   code,
		Status: status,
	}
}
