package routers

import (
	"log"

	createAuthCheckMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/auth_check"
	createCsrfMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/csrf"
	createLogMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/log"
	recoveryMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/recover"

	advdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/delivery"
	cartdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery"
	citydel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/delivery"
	favdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/favourites/delivery"
	orderdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/delivery"
	profdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery"
	surveydel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/survey/delivery"
	authdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"

	advusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	cartusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/usecases"
	cityusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/usecases"
	favusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/favourites/usecases"
	orderusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/usecases"
	profusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/usecases"
	surveyusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/survey/usecases"
	authusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const plug = ""

func NewRouter(pool *pgxpool.Pool, logger *zap.SugaredLogger,
	advertStorage advusecases.AdvertsStorageInterface,
	cartStorage cartusecases.CartStorageInterface,
	cityStorage cityusecases.CityStorageInterface,
	orderStorage orderusecases.OrderStorageInterface,
	profileStorage profusecases.ProfileStorageInterface,
	userStorage authusecases.UsersStorageInterface,
	surveyStorage surveyusecases.SurveyStorageInterface,
	favouritesStorage favusecases.FavouritesStorageInterface) *mux.Router {

	router := mux.NewRouter()
	router.Use(recoveryMiddleware.RecoveryMiddleware)

	logMiddleware := createLogMiddleware.CreateLogMiddleware(logger)
	router.Use(logMiddleware)

	csrfMiddleware := createCsrfMiddleware.CreateCsrfMiddleware()
	authCheckMiddleware := createAuthCheckMiddleware.CreateAuthCheckMiddleware(userStorage)

	advertsHandler := advdel.NewAdvertsHandler(advertStorage, userStorage, plug, plug, logger)
	cartHandler := cartdel.NewCartHandler(cartStorage, advertStorage, userStorage, plug, plug, logger)
	authHandler := authdel.NewAuthHandler(userStorage, profileStorage, plug, plug)
	profileHandler := profdel.NewProfileHandler(profileStorage, userStorage, plug, plug, logger)
	orderHandler := orderdel.NewOrderHandler(orderStorage, cartStorage, advertStorage, userStorage, plug, plug, logger)
	cityHandler := citydel.NewCityHandler(cityStorage, plug, plug, logger)
	surveyHandler := surveydel.NewSurveyHandler(userStorage, surveyStorage)
	favouritesHandler := favdel.NewFavouritesHandler(favouritesStorage, advertStorage, userStorage)

	log.Println("Server is running")

	rootRouter := router.PathPrefix("/api").Subrouter()
	ServeAuthRouter(rootRouter, authHandler, authCheckMiddleware)
	ServeAdvertsRouter(rootRouter, advertsHandler, authCheckMiddleware, csrfMiddleware)
	ServeProfileRouter(rootRouter, profileHandler, authCheckMiddleware, csrfMiddleware)
	ServeCartRouter(rootRouter, cartHandler, authCheckMiddleware)
	ServeOrderRouter(rootRouter, orderHandler, authCheckMiddleware)
	ServeSurveyRouter(rootRouter, surveyHandler, authCheckMiddleware)
	ServeFavouritesRouter(rootRouter, favouritesHandler, authCheckMiddleware)

	rootRouter.HandleFunc("/city", cityHandler.GetCityList)

	return router
}
