package routers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/favourites/delivery"
	"github.com/gorilla/mux"
)

func ServeFavouritesRouter(router *mux.Router, favouritesHandler *delivery.FavouritesHandler,
	authCheckMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/favourites").Subrouter()

	subrouter.Use(authCheckMiddleware)

	subrouter.HandleFunc("/list", favouritesHandler.GetFavouritesList)
	subrouter.HandleFunc("/change", favouritesHandler.ChangeFavourites)
	subrouter.HandleFunc("/delete", favouritesHandler.DeleteFromFavourites)
	subrouter.HandleFunc("/subscribed", favouritesHandler.GetSubscribedAdverts).Methods("GET")
}
