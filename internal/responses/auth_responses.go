package responses

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
)

type AuthOkResponse struct {
	Code      int          `json:"code"`
	User      storage.User `json:"user"`
	SessionID string       `json:"session_id"`
	IsAuth    bool         `json:"is_auth"`
}

type AuthErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewAuthOkResponse(user storage.User, sessionID string, isAuth bool) *AuthOkResponse {
	return &AuthOkResponse{
		Code:      StatusOk,
		User:      user,
		SessionID: sessionID,
		IsAuth:    isAuth,
	}
}

func NewAuthErrResponse(code int, status string) *AuthErrResponse {
	return &AuthErrResponse{
		Code:   code,
		Status: status,
	}
}
