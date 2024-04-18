package server

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	myrouter "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery/router"
	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	logger "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/usecases"

	"github.com/gorilla/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	Timeout            = time.Second * 3
	Address            = "127.0.0.1:8080" //"109.120.183.3:8080"
	outputLogPath      = "stdout logs.json"
	errorOutputLogPath = "stderr err_logs.json"
)

type Server struct {
	server *http.Server
}

func (srv *Server) Run() error {
	connPool, err := pgxpool.NewWithConfig(context.Background(), pgxpoolconfig.PGXPoolConfig())
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	logger, err := logger.NewLogger(strings.Split(outputLogPath, " "),
		strings.Split(errorOutputLogPath, " "))
	if err != nil {
		return err //nolint:wrapcheck
	}

	router := myrouter.NewRouter(connPool, logger)

	credentials := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"http://127.0.0.1:8008"}) // "http://109.120.183.3:8008"
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
