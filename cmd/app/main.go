package main

import (
	"github.com/joho/godotenv"
	"log"

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

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading env file", err)
	}

	if err := srv.Run(); err != nil {
		log.Fatal("Error occurred while starting server:", err.Error())
	}
}
