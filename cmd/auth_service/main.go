package main

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	profilerepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/repository"
	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	authrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

func main() {
	cfg := config.ReadConfig()
	addr := cfg.Server.Host + cfg.Server.AuthServicePort

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("Error occurred while listening auth service", err)
		return
	}

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

	profileStorage := profilerepo.NewProfileStorage(connPool)
	userStorage := authrepo.NewUserStorage(connPool)
	authManager := delivery.NewAuthManager(userStorage, profileStorage)

	srv := grpc.NewServer()
	authproto.RegisterAuthServer(srv, authManager)
	log.Println("Auth service is running on port", cfg.Server.AuthServicePort)

	err = srv.Serve(listener)
	if err != nil {
		log.Println(err)
		return
	}
}
