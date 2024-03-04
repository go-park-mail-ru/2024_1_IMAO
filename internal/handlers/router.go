package myhandlers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
	"github.com/gorilla/mux"
	"log"
)

type AuthHandler struct {
	List *storage.UsersList
}

type AdvertsHandler struct {
	List *storage.AdvertsList
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	log.Println("Server is running")

	usersList := storage.NewActiveUser()
	authHandler := &AuthHandler{
		List: usersList,
	}

	advertsList := storage.NewAdvertsList()
	storage.FillAdvertsList(advertsList)
	advertsHandler := &AdvertsHandler{
		List: advertsList,
	}

	router.HandleFunc("/", advertsHandler.Root)

	router.HandleFunc("/login", authHandler.Login)
	router.HandleFunc("/check_auth", authHandler.CheckAuth)
	router.HandleFunc("/logout", authHandler.Logout)
	router.HandleFunc("/signup", authHandler.Signup)

	return router
}
