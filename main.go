package main

import (
	"2024_1_IMAO/adverts"
	"2024_1_IMAO/auth"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	active := auth.NewActiveUser()
	ads := adverts.NewAdvertsStorage()
	adverts.FillAdvertsStorage(ads)

	log.Println("Server is running")

	credentials := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	r.HandleFunc("/", ads.Root)
	r.HandleFunc("/login", active.Login)
	r.HandleFunc("/check_auth", active.CheckAuth)
	r.HandleFunc("/logout", active.Logout)
	r.HandleFunc("/signup", active.Signup)

	http.ListenAndServe(":8080", handlers.CORS(credentials, originsOk, headersOk, methodsOk)(r))
}
