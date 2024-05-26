package routers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery"
	"github.com/gorilla/mux"
)

func ServeProfileRouter(router *mux.Router, profileHandler *delivery.ProfileHandler,
	authCheckMiddleware, csrfMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/profile").Subrouter()

	subrouter.HandleFunc("/{id:[0-9]+}", profileHandler.GetProfile).Methods("GET")

	subrouterChange := subrouter.PathPrefix("/change").Subrouter()
	subrouterChange.Use(authCheckMiddleware)
	subrouterChange.HandleFunc("", profileHandler.ChangeSubscription).Methods("POST")

	subrouterEdit := subrouter.PathPrefix("/edit").Subrouter()
	subrouterEdit.Use(authCheckMiddleware, csrfMiddleware)
	subrouterEdit.HandleFunc("", profileHandler.EditProfile).Methods("POST")

	subrouterPhone := subrouter.PathPrefix("/phone").Subrouter()
	subrouterPhone.Use(authCheckMiddleware)
	subrouterPhone.HandleFunc("", profileHandler.SetProfilePhone).Methods("POST")

	subrouterCity := subrouter.PathPrefix("/city").Subrouter()
	subrouterCity.Use(authCheckMiddleware)
	subrouterCity.HandleFunc("", profileHandler.SetProfileCity).Methods("POST")
}
