package routers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/payments/delivery"
	"github.com/gorilla/mux"
)

func ServePaymentsRouter(router *mux.Router, favouritesHandler *delivery.PaymentsHandler,
	authCheckMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/payments").Subrouter()

	subrouter.Use(authCheckMiddleware)

	subrouter.HandleFunc("/form", favouritesHandler.GetPaymentForm)

}
