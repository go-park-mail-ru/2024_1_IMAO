package app

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/handlers"
	"github.com/gorilla/handlers"
	"net/http"
	"time"
)

const (
	Timeout = time.Second * 3
	Address = ":8080"
)

type Server struct {
	server *http.Server
}

func (srv *Server) Run() error {
	router := myhandlers.NewRouter()

	credentials := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	muxWithCORS := handlers.CORS(credentials, originsOk, headersOk, methodsOk)(router)

	srv.server = &http.Server{
		Addr:         Address,
		Handler:      muxWithCORS,
		ReadTimeout:  Timeout,
		WriteTimeout: Timeout,
	}

	return srv.server.ListenAndServe()
}
