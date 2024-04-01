package myhandlers

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	UsersList   *storage.UsersList
	ProfileList *storage.ProfileList
}

type AdvertsHandler struct {
	List *storage.AdvertsList
}

type ProfileHandler struct {
	AdvertsList *storage.AdvertsList
	ProfileList *storage.ProfileList
	UsersList   *storage.UsersList
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(recoveryMiddleware)

	advertsList := storage.NewAdvertsList()
	storage.FillAdvertsList(advertsList)
	profileList := storage.NewProfileList()
	usersList := storage.NewActiveUser()

	advertsHandler := &AdvertsHandler{
		List: advertsList,
	}

	profileHandler := &ProfileHandler{
		ProfileList: profileList,
		AdvertsList: advertsList,
		UsersList:   usersList,
	}

	authHandler := &AuthHandler{
		UsersList:   usersList,
		ProfileList: profileList,
	}

	log.Println("Server is running")

	router.HandleFunc("/api/adverts/create", advertsHandler.CreateAdvert)
	router.HandleFunc("/api/adverts/edit", advertsHandler.EditAdvert)
	router.HandleFunc("/api/adverts/", advertsHandler.GetAdsList)
	router.HandleFunc("/api/adverts/{city:[a-zA-Z]+}", advertsHandler.GetAdsList)
	router.HandleFunc("/api/adverts/{city:[a-zA-Z]+}/{category:[a-zA-Z]+}", advertsHandler.GetAdsList)
	router.HandleFunc("/api/adverts/{city:[a-zA-Z]+}/{category:[a-zA-Z]+}/{id:[0-9]+}", advertsHandler.GetAdvert)
	router.HandleFunc("/api/adverts/delete/{id:[0-9]+}", advertsHandler.DeleteAdvert)
	router.HandleFunc("/api/adverts/close/{id:[0-9]+}", advertsHandler.CloseAdvert)

	router.HandleFunc("/api/auth/login", authHandler.Login)
	router.HandleFunc("/api/auth/check_auth", authHandler.CheckAuth)
	router.HandleFunc("/api/auth/logout", authHandler.Logout)
	router.HandleFunc("/api/auth/signup", authHandler.Signup)
	router.HandleFunc("/api/auth/edit", authHandler.EditUser)

	router.HandleFunc("/api/profile/{id:[0-9]+}", profileHandler.GetProfile)
	router.HandleFunc("/api/profile/{id:[0-9]+}/rating", profileHandler.SetProfileRating)
	router.HandleFunc("/api/profile/approved", profileHandler.SetProfileApproved)
	router.HandleFunc("/api/profile/edit", profileHandler.EditProfile)
	router.HandleFunc("/api/profile/set", profileHandler.SetProfile)
	router.HandleFunc("/api/profile/phone", profileHandler.SetProfilePhone)
	router.HandleFunc("/api/profile/city", profileHandler.SetProfileCity)

	router.HandleFunc("/api/profile/{id:[0-9]+}/adverts", profileHandler.SetProfileCity)

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
