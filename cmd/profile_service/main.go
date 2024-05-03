package main

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery"
	profileproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery/protobuf"
	profilerepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/repository"
	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

func main() {
	cfg := config.ReadConfig()
	addr := cfg.Server.Host + cfg.Server.ProfileServicePort

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("Error occurred while listening profile service", err)
		return
	}

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

	srv := grpc.NewServer()
	profileproto.RegisterProfileServer(srv, profileManager)
	log.Println("Profile service is running on port", cfg.Server.ProfileServicePort)

	err = srv.Serve(listener)
	if err != nil {
		log.Println(err)
		return
	}
}
