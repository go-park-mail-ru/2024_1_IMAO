package main

import (
	"context"
	"log"
	"net"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery"
	cartproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery/protobuf"
	cartrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/repository"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.ReadConfig()
	addr := cfg.Server.Host + cfg.Server.CartServicePort

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("Error occurred while listening cart service", err)
		return
	}

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

	srv := grpc.NewServer()
	cartproto.RegisterCartServer(srv, cartManager)
	log.Println("Cart service is running on port", cfg.Server.CartServicePort)

	err = srv.Serve(listener)
	if err != nil {
		log.Println(err)
		return
	}
}
