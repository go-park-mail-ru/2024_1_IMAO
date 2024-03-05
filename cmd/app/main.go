package main

import (
	"log"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/app"
)

//	@title      YULA project API
//	@version    1.0
//	@description  This is a server of YULA server.
//
// @Schemes http
// @host  109.120.183.3:8080
// @BasePath  /api/v1
func main() {
	srv := new(app.Server)

	if err := srv.Run(); err != nil {
		log.Fatal("Error occurred while starting server:", err.Error())
	}
}
