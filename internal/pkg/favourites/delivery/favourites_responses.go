package delivery

import (
	models "github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
)

type FavouritesOkResponse struct {
	Code  int                       `json:"code"`
	Items []*models.ReturningAdvert `json:"items"`
}

type FavouritesChangeResponse struct {
	Code       int  `json:"code"`
	IsAppended bool `json:"isAppended"`
}

type FavouritesErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewCartChangeResponse(isAppended bool) *FavouritesChangeResponse {
	return &FavouritesChangeResponse{
		Code:       responses.StatusOk,
		IsAppended: isAppended,
	}
}

func NewCartOkResponse(adverts []*models.ReturningAdvert) *FavouritesOkResponse {
	return &FavouritesOkResponse{
		Code:  responses.StatusOk,
		Items: adverts,
	}
}

func NewCartErrResponse(code int, status string) *FavouritesErrResponse {
	return &FavouritesErrResponse{
		Code:   code,
		Status: status,
	}
}
