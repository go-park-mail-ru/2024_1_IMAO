package routers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/delivery"
	"github.com/gorilla/mux"
)

func ServeAdvertsRouter(router *mux.Router, advertsHandler *delivery.AdvertsHandler,
	authCheckMiddleware, csrfMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/adverts").Subrouter()

	subrouterCreate := subrouter.PathPrefix("/create").Subrouter()
	subrouterCreate.Use(authCheckMiddleware, csrfMiddleware)
	subrouterCreate.HandleFunc("", advertsHandler.CreateAdvert).Methods("POST")

	subrouterEdit := subrouter.PathPrefix("/edit").Subrouter()
	subrouterEdit.Use(authCheckMiddleware, csrfMiddleware)
	subrouterEdit.HandleFunc("", advertsHandler.EditAdvert).Methods("POST")

	subrouter.HandleFunc("/search", advertsHandler.GetAdsListWithSearch).Methods("GET")
	subrouter.HandleFunc("/suggestions", advertsHandler.GetSuggestions).Methods("GET")
	subrouter.HandleFunc("/price_history/{id:[0-9]+}", advertsHandler.GetAdvertPriceHistoryByID).
		Methods("GET")
	subrouter.HandleFunc("/", advertsHandler.GetAdsList).Methods("GET")
	subrouter.HandleFunc("/{city:[a-zA-Z_]+}", advertsHandler.GetAdsList).Methods("GET")
	subrouter.HandleFunc("/{city:[a-zA-Z_]+}/{category:[a-zA-Z_]+}", advertsHandler.GetAdsList).
		Methods("GET")
	subrouter.HandleFunc("/{city:[a-zA-Z_]+}/{category:[a-zA-Z_]+}/{id:[0-9]+}", advertsHandler.GetAdvert).
		Methods("GET")
	subrouter.HandleFunc("/{id:[0-9]+}", advertsHandler.GetAdvertByID).Methods("GET")
	subrouter.HandleFunc("/close/{id:[0-9]+}", advertsHandler.CloseAdvert).Methods("POST")
	subrouter.HandleFunc("/promotion/{id:[0-9]+}", advertsHandler.GetPromotionData).Methods("GET")
}
