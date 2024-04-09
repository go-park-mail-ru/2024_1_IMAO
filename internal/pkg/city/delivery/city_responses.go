package delivery

import (
	models "github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
)

type CityOkResponse struct {
	Code int          `json:"code"`
	City *models.City `json:"city"`
}

type CityErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewCityOkResponse(city *models.City) *CityOkResponse {
	return &CityOkResponse{
		Code: responses.StatusOk,
		City: city,
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
