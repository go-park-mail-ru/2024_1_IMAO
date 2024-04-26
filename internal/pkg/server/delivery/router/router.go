package myhandlers

import (
	"log"
	"net/http"

	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"

	createAuthCheckMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/auth_check"
	createCsrfMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/csrf"
	createLogMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/log"

	advdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/delivery"
	cartdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery"
	citydel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/delivery"
	orderdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/delivery"
	profdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery"
	authdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"

	advusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	cartusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/usecases"
	cityusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/usecases"
	orderusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/usecases"
	profusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/usecases"
	authusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const plug = ""

func NewRouter(pool *pgxpool.Pool, logger *zap.SugaredLogger, advertStorage advusecases.AdvertsStorageInterface,
	cartStorage cartusecases.CartStorageInterface, cityStorage cityusecases.CityStorageInterface, orderStorage orderusecases.OrderStorageInterface,
	profileStorage profusecases.ProfileStorageInterface, userStorage authusecases.UsersStorageInterface) *mux.Router {
	router := mux.NewRouter()
	router.Use(recoveryMiddleware)
	logMiddleware := createLogMiddleware.CreateLogMiddleware(logger)
	router.Use(logMiddleware)
	csrfMiddleware := createCsrfMiddleware.CreateCsrfMiddleware()
	//router.Use(csrfMiddleware)
	AuthCheckMiddleware := createAuthCheckMiddleware.CreateAuthCheckMiddleware(userStorage)

	advertsHandler := advdel.NewAdvertsHandler(advertStorage, userStorage, plug, plug, logger)

	cartHandler := cartdel.NewCartHandler(cartStorage, advertStorage, userStorage, plug, plug, logger)

	authHandler := authdel.NewAuthHandler(userStorage, profileStorage, plug, plug, logger)

	profileHandler := profdel.NewProfileHandler(profileStorage, userStorage, plug, plug, logger)

	orderHandler := orderdel.NewOrderHandler(orderStorage, cartStorage, advertStorage, userStorage, plug, plug, logger)

	cityHandler := citydel.NewCityHandler(cityStorage, plug, plug, logger)

	log.Println("Server is running")

	subrouterAdvertsCreate := router.PathPrefix("/api/adverts").Subrouter()
	subrouterAdvertsCreate.Use(AuthCheckMiddleware)
	subrouterAdvertsCreate.Use(csrfMiddleware)
	subrouterAdvertsCreate.HandleFunc("/create", advertsHandler.CreateAdvert).Methods("POST")

	subrouterAuthCSRF := router.PathPrefix("/api/auth").Subrouter()
	subrouterAuthCSRF.Use(AuthCheckMiddleware)
	subrouterAuthCSRF.HandleFunc("/csrf", authHandler.GetCSRFToken).Methods("GET")

	subrouterAdvertsEdit := router.PathPrefix("/api/adverts").Subrouter()
	subrouterAdvertsEdit.Use(AuthCheckMiddleware)
	subrouterAdvertsEdit.Use(csrfMiddleware)
	subrouterAdvertsEdit.HandleFunc("/edit", advertsHandler.EditAdvert).Methods("POST")

	//router.HandleFunc("/api/adverts/create", advertsHandler.CreateAdvert)
	//router.HandleFunc("/api/adverts/edit", advertsHandler.EditAdvert)
	router.HandleFunc("/api/adverts/", advertsHandler.GetAdsList)
	router.HandleFunc("/api/adverts/{city:[a-zA-Z_]+}", advertsHandler.GetAdsList)
	router.HandleFunc("/api/adverts/{city:[a-zA-Z_]+}/{category:[a-zA-Z_]+}", advertsHandler.GetAdsList)
	router.HandleFunc("/api/adverts/{city:[a-zA-Z_]+}/{category:[a-zA-Z_]+}/{id:[0-9]+}", advertsHandler.GetAdvert)
	router.HandleFunc("/api/adverts/{id:[0-9]+}", advertsHandler.GetAdvertByID)
	//router.HandleFunc("/api/adverts/delete/{id:[0-9]+}", advertsHandler.DeleteAdvert)
	router.HandleFunc("/api/adverts/close/{id:[0-9]+}", advertsHandler.CloseAdvert)

	router.HandleFunc("/api/auth/login", authHandler.Login)
	router.HandleFunc("/api/auth/check_auth", authHandler.CheckAuth)
	router.HandleFunc("/api/auth/logout", authHandler.Logout)
	router.HandleFunc("/api/auth/signup", authHandler.Signup)
	router.HandleFunc("/api/auth/edit/email", authHandler.EditUserEmail)
	//router.HandleFunc("/api/auth/csrf", authHandler.GetCSRFToken)

	router.HandleFunc("/api/profile/{id:[0-9]+}", profileHandler.GetProfile)
	//router.HandleFunc("/api/profile/{id:[0-9]+}/rating", profileHandler.SetProfileRating)
	//router.HandleFunc("/api/profile/approved", profileHandler.SetProfileApproved)
	router.HandleFunc("/api/profile/edit", profileHandler.EditProfile)
	//router.HandleFunc("/api/profile/set", profileHandler.SetProfile)
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
