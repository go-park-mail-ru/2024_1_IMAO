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

	subrouter.HandleFunc("/search", advertsHandler.GetAdsListWithSearch)
	subrouter.HandleFunc("/", advertsHandler.GetAdsList)
	subrouter.HandleFunc("/{city:[a-zA-Z_]+}", advertsHandler.GetAdsList)
	subrouter.HandleFunc("/{city:[a-zA-Z_]+}/{category:[a-zA-Z_]+}", advertsHandler.GetAdsList)
	subrouter.HandleFunc("/{city:[a-zA-Z_]+}/{category:[a-zA-Z_]+}/{id:[0-9]+}", advertsHandler.GetAdvert)
	subrouter.HandleFunc("/{id:[0-9]+}", advertsHandler.GetAdvertByID)

	//subrouter.HandleFunc("/api/adverts/delete/{id:[0-9]+}", advertsHandler.DeleteAdvert)
	subrouter.HandleFunc("/close/{id:[0-9]+}", advertsHandler.CloseAdvert)
}
