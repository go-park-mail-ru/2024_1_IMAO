package main

import (
	"log"

	"github.com/go-park-mail-ru/2024_1_IMAO/cmd/auth_service"
	"github.com/go-park-mail-ru/2024_1_IMAO/cmd/cart_service"
	"github.com/go-park-mail-ru/2024_1_IMAO/cmd/profile_service"

	app "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server"
)

//	@title      YULA project API
//	@version    1.0
//	@description  This is a server of YULA server.
//
// @Schemes http
// @host  109.120.183.3:8008
func main() {
	srv := new(app.Server)

	go auth_service.RunAuth()
	go profile_service.RunProfile()
	go cart_service.RunCart()

	if err := srv.Run(); err != nil {
		log.Fatal("Error occurred while starting server:", err.Error())
	}
}
