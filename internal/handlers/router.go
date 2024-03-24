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

	log.Println("Server is running")

	router.HandleFunc("/api/adverts", advertsHandler.Root)
	router.HandleFunc("/api/adverts/{city:[a-zA-Z]+}", advertsHandler.Root)
	router.HandleFunc("/api/adverts/{city:[a-zA-Z]+}/{category:[a-zA-Z]+}", advertsHandler.GetCategoryAds)

	router.HandleFunc("/api/auth/login", authHandler.Login)
	router.HandleFunc("/api/auth/check_auth", authHandler.CheckAuth)
	router.HandleFunc("/api/auth/logout", authHandler.Logout)
	router.HandleFunc("/api/auth/signup", authHandler.Signup)

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
