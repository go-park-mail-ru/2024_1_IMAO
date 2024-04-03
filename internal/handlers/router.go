package myhandlers

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	List *storage.UsersList
}

type AdvertsHandler struct {
	List *storage.AdvertsList
}

type CartHandler struct {
	ListCart    *storage.CartList
	ListAdverts *storage.AdvertsList
	ListUsers   *storage.UsersList
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(recoveryMiddleware)

	usersList := storage.NewActiveUser()
	authHandler := &AuthHandler{
		List: usersList,
	}

	advertsList := storage.NewAdvertsList()
	storage.FillAdvertsList(advertsList)
	advertsHandler := &AdvertsHandler{
		List: advertsList,
	}

	cartList := storage.NewCartList()
	cartHandler := &CartHandler{
		ListCart:    cartList,
		ListAdverts: advertsList,
		ListUsers:   usersList,
	}

	log.Println("Server is running")

	router.HandleFunc("/api/adverts/create", advertsHandler.CreateAdvert)
	router.HandleFunc("/api/adverts/list", advertsHandler.GetAdsList)
	router.HandleFunc("/api/adverts/{city:[a-zA-Z]+}", advertsHandler.GetAdsList)
	router.HandleFunc("/api/adverts/{city:[a-zA-Z]+}/{category:[a-zA-Z]+}", advertsHandler.GetAdsList)

	router.HandleFunc("/api/auth/login", authHandler.Login)
	router.HandleFunc("/api/auth/check_auth", authHandler.CheckAuth)
	router.HandleFunc("/api/auth/logout", authHandler.Logout)
	router.HandleFunc("/api/auth/signup", authHandler.Signup)

	router.HandleFunc("/api/cart/list", cartHandler.GetCartList)
	router.HandleFunc("/api/cart/change", cartHandler.ChangeCart)
	router.HandleFunc("/api/cart/delete", cartHandler.DeleteFromCart)

	return router
}

// Обработка паник
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Panic occurred:", err)
				http.Error(writer, responses.ErrInternalServer, responses.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
