package profile_service

import (
	"context"
	"fmt"
	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"
	createMetricsMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/metrics"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery"
	profileproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery/protobuf"
	profilerepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/repository"
	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunProfile() {
	cfg := config.ReadConfig()
	addr := cfg.Server.Host + cfg.Server.ProfileServicePort

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("Error occurred while listening profile service", err)
		return
	}

	metrics, err := mymetrics.CreateGRPCMetrics("profile")
	if err != nil {
		log.Println("Error occurred while creating auth metrics", err)
	}

	interceptor := createMetricsMiddleware.CreateMetricsInterceptor(*metrics)

	grpcConn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println("Error occurred while starting grpc connection on profile service", err)
		return
	}
	defer grpcConn.Close()

	connPool, err := pgxpool.NewWithConfig(context.Background(), pgxpoolconfig.PGXPoolConfig())
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	profileStorage := profilerepo.NewProfileStorage(connPool)
	profileManager := delivery.NewProfileManager(profileStorage)

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.ServeMetricsInterceptor))
	profileproto.RegisterProfileServer(srv, profileManager)
	log.Println("Profile service is running on port", cfg.Server.ProfileServicePort)

	go func() {
		err = srv.Serve(listener)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	router := mux.NewRouter()
	router.PathPrefix("/metrics").Handler(promhttp.Handler())

	server := http.Server{Handler: router, Addr: fmt.Sprintf(":%d", 7072)}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("fail profile.ListenAndServe")
		}
	}()
}
