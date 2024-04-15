package delivery

import (
	models "github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
)

type ProfileOkResponse struct {
	Code    int             `json:"code"`
	Profile *models.Profile `json:"profile"`
}

type ProfileErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewProfileOkResponse(profile *models.Profile) *ProfileOkResponse {
	return &ProfileOkResponse{
		Code:    responses.StatusOk,
		Profile: profile,
	}
}

func NewProfileErrResponse(code int, status string) *ProfileErrResponse {
	return &ProfileErrResponse{
		Code:   code,
		Status: status,
	}
}