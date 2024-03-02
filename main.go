package main

import (
	"2024_1_IMAO/auth"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	active := auth.NewActiveUser()

	credentials := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	//r.HandleFunc("/", active.Root)
	r.HandleFunc("/login", active.Login)
	r.HandleFunc("/logout", active.Logout)
	r.HandleFunc("/signup", active.Signup)

	http.ListenAndServe(":8080", handlers.CORS(credentials, originsOk, headersOk, methodsOk)(r))
}
