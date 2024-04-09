package myhandlers

import (
	"log"
	"net/http"

	advrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/repository"
	cartrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/repository"
	cityrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/repository"
	orderrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/repository"
	profrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/repository"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	authrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/repository"

	advdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/delivery"
	cartdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery"
	citydel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/delivery"
	orderdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/delivery"
	profdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery"
	authdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewRouter(pool *pgxpool.Pool, logger *zap.SugaredLogger) *mux.Router {
	router := mux.NewRouter()
	router.Use(recoveryMiddleware)

	advertsList := advrepo.NewAdvertsList(pool, logger)
	advrepo.FillAdvertsList(advertsList)
	profileList := profrepo.NewProfileList(pool, logger)
	usersList := authrepo.NewActiveUser(pool, logger)
	cityList := cityrepo.NewCityList(pool, logger)

	advertsHandler := &advdel.AdvertsHandler{
		List: advertsList,
	}

	profileHandler := &profdel.ProfileHandler{
		ProfileList: profileList,
		AdvertsList: advertsList,
		UsersList:   usersList,
	}

	authHandler := &authdel.AuthHandler{
		UsersList:   usersList,
		ProfileList: profileList,
	}

	cartList := cartrepo.NewCartList(pool, logger)
	cartHandler := &cartdel.CartHandler{
		ListCart:    cartList,
		ListAdverts: advertsList,
		ListUsers:   usersList,
	}

	orderList := orderrepo.NewOrderList(pool, logger)
	orderHandler := &orderdel.OrderHandler{
		ListOrder:   orderList,
		ListCart:    cartList,
		ListAdverts: advertsList,
		ListUsers:   usersList,
	}

	cityHandler := &citydel.CityHandler{
		CityList: cityList,
	}

	log.Println("Server is running")

	router.HandleFunc("/api/adverts/create", advertsHandler.CreateAdvert)
	router.HandleFunc("/api/adverts/edit", advertsHandler.EditAdvert)
	router.HandleFunc("/api/adverts/", advertsHandler.GetAdsList)
	router.HandleFunc("/api/adverts/{city:[a-zA-Z]+}", advertsHandler.GetAdsList)
	router.HandleFunc("/api/adverts/{city:[a-zA-Z]+}/{category:[a-zA-Z]+}", advertsHandler.GetAdsList)
	router.HandleFunc("/api/adverts/{city:[a-zA-Z]+}/{category:[a-zA-Z]+}/{id:[0-9]+}", advertsHandler.GetAdvert)
	router.HandleFunc("/api/adverts/{id:[0-9]+}", advertsHandler.GetAdvertByID)
	router.HandleFunc("/api/adverts/delete/{id:[0-9]+}", advertsHandler.DeleteAdvert)
	router.HandleFunc("/api/adverts/close/{id:[0-9]+}", advertsHandler.CloseAdvert)

	router.HandleFunc("/api/auth/login", authHandler.Login)
	router.HandleFunc("/api/auth/check_auth", authHandler.CheckAuth)
	router.HandleFunc("/api/auth/logout", authHandler.Logout)
	router.HandleFunc("/api/auth/signup", authHandler.Signup)
	router.HandleFunc("/api/auth/edit/email", authHandler.EditUserEmail)

	router.HandleFunc("/api/profile/{id:[0-9]+}", profileHandler.GetProfile)
	router.HandleFunc("/api/profile/{id:[0-9]+}/rating", profileHandler.SetProfileRating)
	router.HandleFunc("/api/profile/approved", profileHandler.SetProfileApproved)
	router.HandleFunc("/api/profile/edit", profileHandler.EditProfile)
	router.HandleFunc("/api/profile/set", profileHandler.SetProfile)
	router.HandleFunc("/api/profile/phone", profileHandler.SetProfilePhone)
	router.HandleFunc("/api/profile/city", profileHandler.SetProfileCity)

	router.HandleFunc("/api/city", cityHandler.GetCityList)

	router.HandleFunc("/api/profile/{id:[0-9]+}/adverts", profileHandler.SetProfileCity)

	router.HandleFunc("/api/cart/list", cartHandler.GetCartList)
	router.HandleFunc("/api/cart/change", cartHandler.ChangeCart)
	router.HandleFunc("/api/cart/delete", cartHandler.DeleteFromCart)

	router.HandleFunc("/api/order/list", orderHandler.GetOrderList)
	router.HandleFunc("/api/order/create", orderHandler.CreateOrder)

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
