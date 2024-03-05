package myhandlers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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

	router.HandleFunc("/", advertsHandler.Root)

	router.HandleFunc("/login", authHandler.Login)
	router.HandleFunc("/check_auth", authHandler.CheckAuth)
	router.HandleFunc("/logout", authHandler.Logout)
	router.HandleFunc("/signup", authHandler.Signup)

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
