package app

import (
	"context"
	"net/http"
	"time"

	myhandlers "github.com/go-park-mail-ru/2024_1_IMAO/internal/handlers"
	"github.com/gorilla/handlers"
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

func (srv *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	if err := srv.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
