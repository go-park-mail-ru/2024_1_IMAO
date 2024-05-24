package delivery

import (
	models "github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
)

type PaymentFormOkResponse struct {
	Code           int    `json:"code"`
	PaymentFormUrl string `json:"paymentFormUrl"`
}

type CityErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewPaymentFormOkResponse(paymentFormUrl string) *PaymentFormOkResponse {
	return &PaymentFormOkResponse{
		Code:           responses.StatusOk,
		PaymentFormUrl: paymentFormUrl,
	}
}

func NewCityErrResponse(code int, status string) *CityErrResponse {
	return &CityErrResponse{
		Code:   code,
		Status: status,
	}
}

type CityListOkResponse struct {
	Code     int              `json:"code"`
	CityList *models.CityList `json:"city_list"`
}

type CityListErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewCityListOkResponse(cityList *models.CityList) *CityListOkResponse {
	return &CityListOkResponse{
		Code:     responses.StatusOk,
		CityList: cityList,
	}
}

func NewCityListErrResponse(code int, status string) *CityListErrResponse {
	return &CityListErrResponse{
		Code:   code,
		Status: status,
	}
}
