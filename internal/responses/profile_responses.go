package responses

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
)

type ProfileOkResponse struct {
	Code    int              `json:"code"`
	Profile *storage.Profile `json:"profile"`
}

type ProfileErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewProfileOkResponse(profile *storage.Profile) *ProfileOkResponse {
	return &ProfileOkResponse{
		Code:    StatusOk,
		Profile: profile,
	}
}

func NewProfileErrResponse(code int, status string) *ProfileErrResponse {
	return &ProfileErrResponse{
		Code:   code,
		Status: status,
	}
}
