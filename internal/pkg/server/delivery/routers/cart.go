package routers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery"
	"github.com/gorilla/mux"
)

func ServeCartRouter(router *mux.Router, cartHandler *delivery.CartHandler) {
	subrouter := router.PathPrefix("/cart").Subrouter()

	subrouter.HandleFunc("/list", cartHandler.GetCartList)
	subrouter.HandleFunc("/change", cartHandler.ChangeCart)
	subrouter.HandleFunc("/delete", cartHandler.DeleteFromCart)
}
