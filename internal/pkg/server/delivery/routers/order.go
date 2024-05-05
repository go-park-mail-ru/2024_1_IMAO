package routers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/delivery"
	"github.com/gorilla/mux"
)

func ServeOrderRouter(router *mux.Router, orderHandler *delivery.OrderHandler,
	authCheckMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/order").Subrouter()

	subrouter.Use(authCheckMiddleware)

	subrouter.HandleFunc("/list", orderHandler.GetSoldOrderList)
	subrouter.HandleFunc("/create", orderHandler.CreateOrder)
}
