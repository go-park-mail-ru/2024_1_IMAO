package delivery

import (
	models "github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
)

type FavouritesOkResponse struct {
	Code  int                         `json:"code"`
	Items []*models.ReturningAdInList `json:"items"`
}

type FavouritesChangeResponse struct {
	Code       int  `json:"code"`
	IsAppended bool `json:"isAppended"`
}

type FavouritesErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewFavouritesChangeResponse(isAppended bool) *FavouritesChangeResponse {
	return &FavouritesChangeResponse{
		Code:       responses.StatusOk,
		IsAppended: isAppended,
	}
}

func NewFavouritesOkResponse(adverts []*models.ReturningAdInList) *FavouritesOkResponse {
	return &FavouritesOkResponse{
		Code:  responses.StatusOk,
		Items: adverts,
	}
}

func NewFavouritesErrResponse(code int, status string) *FavouritesErrResponse {
	return &FavouritesErrResponse{
		Code:   code,
		Status: status,
	}
}
