package server

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	cartproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery/protobuf"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	profileproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery/protobuf"
	myrouter "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery/routers"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	logger "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/usecases"

	advertrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/repository"
	cartrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/repository"
	cityrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/repository"
	favrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/favourites/repository"
	orderrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/repository"
	surveyrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/survey/repository"
	"github.com/gorilla/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	Timeout            = time.Second * 15
	Address            = ":8080" //"109.120.183.3:8080"
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

	advertStorage := advertrepo.NewAdvertStorage(connPool)
	cartStorage := cartrepo.NewCartStorage(connPool)
	cityStorage := cityrepo.NewCityStorage(connPool)
	orderStorage := orderrepo.NewOrderStorage(connPool)
	surveyStorage := surveyrepo.NewSurveyStorage(connPool)
	favouritesStorage := favrepo.NewFavouritesStorage(connPool)

	cfg := config.ReadConfig()

	authAddr := cfg.Server.AuthIP + cfg.Server.AuthServicePort
	grpcConnAuth, err := grpc.Dial(
		authAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println("Error occurred while starting grpc connection on auth service", err)
		return err
	}
	defer grpcConnAuth.Close()
	authClient := authproto.NewAuthClient(grpcConnAuth)

	profileAddr := cfg.Server.ProfileIP + cfg.Server.ProfileServicePort
	grpcConnProfile, err := grpc.Dial(
		profileAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println("Error occurred while starting grpc connection on profile service", err)
		return err
	}
	defer grpcConnProfile.Close()
	profileClient := profileproto.NewProfileClient(grpcConnProfile)

	cartAddr := cfg.Server.CartIP + cfg.Server.CartServicePort
	grpcConnCart, err := grpc.Dial(
		cartAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println("Error occurred while starting grpc connection on cart service", err)
		return err
	}
	defer grpcConnCart.Close()
	cartClient := cartproto.NewCartClient(grpcConnCart)

	router := myrouter.NewRouter(logger, advertStorage, cartClient, cartStorage, cityStorage, orderStorage,
		surveyStorage, authClient, profileClient, favouritesStorage)

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

	log.Println("Server is running on port", Address)
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
