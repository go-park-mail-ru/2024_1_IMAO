package routers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/delivery"
	"github.com/gorilla/mux"
)

func ServeCityRouter(router *mux.Router, cityHandler *delivery.CityHandler) {
	subrouter := router.PathPrefix("/city").Subrouter()

	subrouter.HandleFunc("/list", cityHandler.GetCityList).Methods("GET")
	subrouter.HandleFunc("/name", cityHandler.GetCityName).Methods("GET")
}
