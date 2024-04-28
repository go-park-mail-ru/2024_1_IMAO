package routers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/delivery"
	"github.com/gorilla/mux"
)

func ServeOrderRouter(router *mux.Router, orderHandler *delivery.OrderHandler) {
	subrouter := router.PathPrefix("/order").Subrouter()

	subrouter.HandleFunc("/list", orderHandler.GetOrderList).Methods("GET")
	subrouter.HandleFunc("/create", orderHandler.CreateOrder).Methods("POST")
}
