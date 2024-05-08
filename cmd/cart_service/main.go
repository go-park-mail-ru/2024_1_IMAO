package cart_service

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

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery"
	cartproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery/protobuf"
	cartrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/repository"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunCart() {
	cfg := config.ReadConfig()
	addr := cfg.Server.Host + cfg.Server.CartServicePort

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
		log.Fatal("Error while creating connection to the database!!")
	}

	cartStorage := cartrepo.NewCartStorage(connPool)
	cartManager := delivery.NewCartManager(cartStorage)

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.ServeMetricsInterceptor))
	cartproto.RegisterCartServer(srv, cartManager)
	log.Println("Cart service is running on port", cfg.Server.CartServicePort)

	go func() {
		err = srv.Serve(listener)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	router := mux.NewRouter()
	router.PathPrefix("/metrics").Handler(promhttp.Handler())

	server := http.Server{Handler: router, Addr: fmt.Sprintf(":%d", 7073)}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("fail cart.ListenAndServe")
		}
	}()
}
