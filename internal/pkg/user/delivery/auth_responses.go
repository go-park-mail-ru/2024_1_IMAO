package delivery

import (
	models "github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
)

type AuthOkResponseLogged struct {
	Code    int         `json:"code"`
	User    models.User `json:"user"`
	Avatar  string      `json:"avatarImg"`
	IsAuth  bool        `json:"isAuth"`
	CartNum uint        `json:"cartNum"`
	FavNum  uint        `json:"favNum"`
}

type AuthOkResponse struct {
	Code   int         `json:"code"`
	User   models.User `json:"user"`
	Avatar string      `json:"avatarImg"`
	IsAuth bool        `json:"isAuth"`
}

type AuthErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

type SessionOKResponse struct {
	Code      int    `json:"code"`
	TokenBody string `json:"tokenBody"`
}

type ValidationErrResponse struct {
	Code   int      `json:"code"`
	Status []string `json:"status"`
}

func NewAuthOkResponseLogged(user models.User, avatar string, isAuth bool, cartNum, favNum uint) *AuthOkResponseLogged {
	return &AuthOkResponseLogged{
		Code:    responses.StatusOk,
		User:    user,
		Avatar:  avatar,
		IsAuth:  isAuth,
		CartNum: cartNum,
		FavNum:  favNum,
	}
}

func NewAuthOkResponse(user models.User, avatar string, isAuth bool) *AuthOkResponse {
	return &AuthOkResponse{
		Code:   responses.StatusOk,
		User:   user,
		Avatar: avatar,
		IsAuth: isAuth,
	}
}

func NewAuthErrResponse(code int, status string) *AuthErrResponse {
	return &AuthErrResponse{
		Code:   code,
		Status: status,
	}
}

func NewValidationErrResponse(code int, status []string) *ValidationErrResponse {
	return &ValidationErrResponse{
		Code:   code,
		Status: status,
	}
}

func NewSessionOkResponse(tokenBody string) *SessionOKResponse {
	return &SessionOKResponse{
		Code:      responses.StatusOk,
		TokenBody: tokenBody,
	}
}
