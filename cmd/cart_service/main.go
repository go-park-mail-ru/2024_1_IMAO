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

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery"
	cartproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery/protobuf"
	cartrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/repository"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	timeout = 10 * time.Second
	port    = 7073
)

func main() {
	cfg := config.ReadConfig()
	addr := cfg.Server.CartIP + cfg.Server.CartServicePort // ДЛЯ ЗАПУСКА В КОНТЕЙНЕРЕ
	//addr := cfg.Server.Host + cfg.Server.CartServicePort // ДЛЯ ЛОКАЛЬНОГО ЗАПУСКА (НЕ В КОНТЕЙНЕРЕ)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("Error occurred while listening cart service", err)

		return
	}

	metrics, err := mymetrics.CreateGRPCMetrics("cart")
	if err != nil {
		log.Println("Error occurred while creating cart metrics", err)
	}

	interceptor := createMetricsMiddleware.CreateMetricsInterceptor(*metrics)

	grpcConn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println("Error occurred while starting grpc connection on cart service", err)

		return
	}
	defer grpcConn.Close()

	connPool, err := pgxpool.NewWithConfig(context.Background(), pgxpoolconfig.PGXPoolConfig())
	if err != nil {
		log.Println("Error while creating connection to the database!!")
	}

	postgresMetrics, err := mymetrics.CreateDatabaseMetrics("cart", "postgres")
	if err != nil {
		log.Println("Error while creating postgres metrics for cart service")
	}

	cartStorage := cartrepo.NewCartStorage(connPool, postgresMetrics)
	cartManager := delivery.NewCartManager(cartStorage)

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.ServeMetricsInterceptor))
	cartproto.RegisterCartServer(srv, cartManager)
	log.Println("Cart service is running on port", cfg.Server.CartServicePort)

	router := mux.NewRouter()
	router.PathPrefix("/metrics").Handler(promhttp.Handler())
	server := http.Server{Handler: router, Addr: fmt.Sprintf(":%d", port), ReadHeaderTimeout: timeout}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("fail cart.ListenAndServe")
		}
	}()

	err = srv.Serve(listener)
	if err != nil {
		log.Println(err)

		return
	}
}
