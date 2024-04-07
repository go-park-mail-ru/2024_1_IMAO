package delivery

import (
	models "github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
)

type OrderOkResponse struct {
	Code  int                      `json:"code"`
	Items []*models.ReturningOrder `json:"items"`
}

type OrderErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

type OrderCreateResponse struct {
	Code      int  `json:"code"`
	IsCreated bool `json:"isCreated"`
}

func NewOrderOkResponse(orderItems []*models.ReturningOrder) *OrderOkResponse {
	return &OrderOkResponse{
		Code:  responses.StatusOk,
		Items: orderItems,
	}
}

func NewOrderErrResponse(code int, status string) *OrderErrResponse {
	return &OrderErrResponse{
		Code:   code,
		Status: status,
	}
}

func NewOrderCreateResponse(isCreated bool) *OrderCreateResponse {
	return &OrderCreateResponse{
		Code:      responses.StatusCreated,
		IsCreated: isCreated,
	}
}
