package routers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"
	"github.com/gorilla/mux"
)

func ServeAuthRouter(router *mux.Router, authHandler *delivery.AuthHandler,
	authCheckMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/auth").Subrouter()

	subrouter.HandleFunc("/login", authHandler.Login)
	subrouter.HandleFunc("/check_auth", authHandler.CheckAuth)
	subrouter.HandleFunc("/logout", authHandler.Logout)
	subrouter.HandleFunc("/signup", authHandler.Signup)
	subrouter.HandleFunc("/edit/email", authHandler.EditUserEmail)

	subrouterCSRF := subrouter.PathPrefix("/csrf").Subrouter()
	subrouterCSRF.Use(authCheckMiddleware)
	subrouterCSRF.HandleFunc("", authHandler.GetCSRFToken).Methods("GET")
}
