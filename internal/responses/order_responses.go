package responses

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
)

type OrderOkResponse struct {
	Code  int                       `json:"code"`
	Items []*storage.ReturningOrder `json:"items"`
}

type OrderErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

type OrderCreateResponse struct {
	Code      int  `json:"code"`
	IsCreated bool `json:"isCreated"`
}

func NewOrderOkResponse(orderItems []*storage.ReturningOrder) *OrderOkResponse {
	return &OrderOkResponse{
		Code:  StatusOk,
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
		Code:      StatusCreated,
		IsCreated: isCreated,
	}
}
