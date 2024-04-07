package delivery

import (
	models "github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
)

type CartOkResponse struct {
	Code  int              `json:"code"`
	Items []*models.Advert `json:"items"`
}

type CartChangeResponse struct {
	Code       int  `json:"code"`
	IsAppended bool `json:"isAppended"`
}

type CartErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewCartChangeResponse(isAppended bool) *CartChangeResponse {
	return &CartChangeResponse{
		Code:       responses.StatusOk,
		IsAppended: isAppended,
	}
}

func NewCartOkResponse(adverts []*models.Advert) *CartOkResponse {
	return &CartOkResponse{
		Code:  responses.StatusOk,
		Items: adverts,
	}
}

func NewCartErrResponse(code int, status string) *CartErrResponse {
	return &CartErrResponse{
		Code:   code,
		Status: status,
	}
}
