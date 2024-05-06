package server

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	myrouter "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery/routers"

	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	logger "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/usecases"

	advertrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/repository"
	cartrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/repository"
	cityrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/repository"
	favrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/favourites/repository"
	orderrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/repository"
	profilerepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/repository"
	surveyrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/survey/repository"
	authrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/repository"

	"github.com/gorilla/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	Timeout            = time.Second * 3
	Address            = "www.vol-4-ok.ru:8080" //"109.120.183.3:8080"
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

	defer logger.Sync()

	advertStorage := advertrepo.NewAdvertStorage(connPool, logger)
	cartStorage := cartrepo.NewCartStorage(connPool, logger)
	cityStorage := cityrepo.NewCityStorage(connPool, logger)
	orderStorage := orderrepo.NewOrderStorage(connPool)
	profileStorage := profilerepo.NewProfileStorage(connPool, logger)
	userStorage := authrepo.NewUserStorage(connPool, logger)
	surveyStorage := surveyrepo.NewSurveyStorage(connPool)
	favouritesStorage := favrepo.NewFavouritesStorage(connPool)

	router := myrouter.NewRouter(connPool, logger, advertStorage, cartStorage, cityStorage, orderStorage,
		profileStorage, userStorage, surveyStorage, favouritesStorage)

	credentials := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"http://www.vol-4-ok.ru:8008"}) // "http://109.120.183.3:8008"
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
