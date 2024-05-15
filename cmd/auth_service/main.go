package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"
	createMetricsMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/metrics"
	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	authrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/repository"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.ReadConfig()
	//addr := cfg.Server.AuthIP + cfg.Server.AuthServicePort // ДЛЯ ЗАПУСКА В КОНТЕЙНЕРЕ
	addr := cfg.Server.Host + cfg.Server.AuthServicePort // ДЛЯ ЛОКАЛЬНОГО ЗАПУСКА (НЕ В КОНТЕЙНЕРЕ)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("Error occurred while listening auth service", err)
		return
	}

	metrics, err := mymetrics.CreateGRPCMetrics("auth")
	if err != nil {
		log.Println("Error occurred while creating auth metrics", err)
	}

	interceptor := createMetricsMiddleware.CreateMetricsInterceptor(*metrics)

	grpcConn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println("Error occurred while starting grpc connection on auth service", err)
		return
	}
	defer grpcConn.Close()

	connPool, err := pgxpool.NewWithConfig(context.Background(), pgxpoolconfig.PGXPoolConfig())
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	postgresMetrics, err := mymetrics.CreateDatabaseMetrics("auth", "postgres")
	if err != nil {
		log.Fatal("Error while creating postgres metrics for auth service")
	}

	userStorage := authrepo.NewUserStorage(connPool, postgresMetrics)
	authManager := delivery.NewAuthManager(userStorage)

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.ServeMetricsInterceptor))
	authproto.RegisterAuthServer(srv, authManager)
	log.Println("Auth service is running on port", cfg.Server.AuthServicePort)

	router := mux.NewRouter()
	router.PathPrefix("/metrics").Handler(promhttp.Handler())
	server := http.Server{Handler: router, Addr: fmt.Sprintf(":%d", 7071)}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("fail auth.ListenAndServe")
		}
	}()

	err = srv.Serve(listener)
	if err != nil {
		log.Println(err)
		return
	}
}
