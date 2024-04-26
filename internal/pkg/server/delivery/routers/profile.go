package routers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery"
	"github.com/gorilla/mux"
)

func ServeProfileRouter(router *mux.Router, profileHandler *delivery.ProfileHandler,
	authCheckMiddleware, csrfMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/profile").Subrouter()

	subrouter.HandleFunc("/{id:[0-9]+}", profileHandler.GetProfile)
	//subrouter.HandleFunc("/{id:[0-9]+}/adverts", profileHandler.SetProfileCity)
	//subrouter.HandleFunc("/api/profile/{id:[0-9]+}/rating", profileHandler.SetProfileRating)
	//subrouter.HandleFunc("/api/profile/approved", profileHandler.SetProfileApproved)
	//subrouter.HandleFunc("/api/profile/edit", profileHandler.EditProfile)
	//subrouter.HandleFunc("/api/profile/set", profileHandler.SetProfile)
	//subrouter.HandleFunc("/api/profile/phone", profileHandler.SetProfilePhone)
	//subrouter.HandleFunc("/api/profile/city", profileHandler.SetProfileCity)

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
