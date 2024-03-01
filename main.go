package main

import (
	"2024_1_IMAO/auth"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	active := auth.NewActiveUser()

	r.HandleFunc("/", active.Root)
	r.HandleFunc("/login", active.Login)
	r.HandleFunc("/logout", active.Logout)
	r.HandleFunc("/signup", active.Signup)

	http.ListenAndServe(":8080", r)
}
