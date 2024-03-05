package main

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/app"
	"log"
)

func main() {
	srv := new(app.Server)

	if err := srv.Run(); err != nil {
		log.Fatal("Error occurred while starting server:", err.Error())
	}
}
