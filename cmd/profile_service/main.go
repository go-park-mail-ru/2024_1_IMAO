package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"
	createMetricsMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/metrics"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery"
	profileproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery/protobuf"
	profilerepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/repository"
	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	timeout = 10 * time.Second
	port    = 7072
)

func main() {
	cfg := config.ReadConfig()
	addr := cfg.Server.ProfileIP + cfg.Server.ProfileServicePort // ДЛЯ ЗАПУСКА В КОНТЕЙНЕРЕ
	//addr := cfg.Server.Host + cfg.Server.ProfileServicePort // ДЛЯ ЛОКАЛЬНОГО ЗАПУСКА (НЕ В КОНТЕЙНЕРЕ)

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
		log.Println("Error while creating connection to the database!!")
	}

	postgresMetrics, err := mymetrics.CreateDatabaseMetrics("profile", "postgres")
	if err != nil {
		log.Println("Error while creating postgres metrics")
	}

	profileStorage := profilerepo.NewProfileStorage(connPool, postgresMetrics)
	profileManager := delivery.NewProfileManager(profileStorage)

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.ServeMetricsInterceptor))
	profileproto.RegisterProfileServer(srv, profileManager)
	log.Println("Profile service is running on port", cfg.Server.ProfileServicePort)

	router := mux.NewRouter()
	router.PathPrefix("/metrics").Handler(promhttp.Handler())
	server := http.Server{Handler: router, Addr: fmt.Sprintf(":%d", port), ReadHeaderTimeout: timeout}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("fail profile.ListenAndServe")
		}
	}()

	err = srv.Serve(listener)
	if err != nil {
		log.Println(err)

		return
	}
}
