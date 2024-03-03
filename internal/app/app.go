package app

import (
	myhandlers "2024_1_IMAO/internal/handlers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	router *mux.Router
}

func (srv *Server) Run() error {
	srv.router = myhandlers.NewRouter()

	credentials := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	return http.ListenAndServe(":8080", handlers.CORS(credentials, originsOk, headersOk, methodsOk)(srv.router))
}
