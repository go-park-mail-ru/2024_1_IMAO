package routers

import (
	createAuthCheckMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/auth_check"
	createCsrfMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/csrf"
	createLogMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/log"
	recoveryMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/recover"
	"log"

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

func NewRouter(pool *pgxpool.Pool, logger *zap.SugaredLogger,
	advertStorage advusecases.AdvertsStorageInterface,
	cartStorage cartusecases.CartStorageInterface,
	cityStorage cityusecases.CityStorageInterface,
	orderStorage orderusecases.OrderStorageInterface,
	profileStorage profusecases.ProfileStorageInterface,
	userStorage authusecases.UsersStorageInterface) *mux.Router {

	router := mux.NewRouter()
	router.Use(recoveryMiddleware.RecoveryMiddleware)

	logMiddleware := createLogMiddleware.CreateLogMiddleware(logger)
	router.Use(logMiddleware)

	csrfMiddleware := createCsrfMiddleware.CreateCsrfMiddleware()
	authCheckMiddleware := createAuthCheckMiddleware.CreateAuthCheckMiddleware(userStorage)

	advertsHandler := advdel.NewAdvertsHandler(advertStorage, userStorage, plug, plug, logger)
	cartHandler := cartdel.NewCartHandler(cartStorage, advertStorage, userStorage, plug, plug, logger)
	authHandler := authdel.NewAuthHandler(userStorage, profileStorage, plug, plug, logger)
	profileHandler := profdel.NewProfileHandler(profileStorage, userStorage, plug, plug, logger)
	orderHandler := orderdel.NewOrderHandler(orderStorage, cartStorage, advertStorage, userStorage, plug, plug, logger)
	cityHandler := citydel.NewCityHandler(cityStorage, plug, plug, logger)

	log.Println("Server is running")

	rootRouter := router.PathPrefix("/api").Subrouter()
	ServeAuthRouter(rootRouter, authHandler, authCheckMiddleware)
	ServeAdvertsRouter(rootRouter, advertsHandler, authCheckMiddleware, csrfMiddleware)
	ServeProfileRouter(rootRouter, profileHandler, authCheckMiddleware, csrfMiddleware)
	ServeCartRouter(rootRouter, cartHandler)
	ServeOrderRouter(rootRouter, orderHandler)

	rootRouter.HandleFunc("/city", cityHandler.GetCityList)

	return router
}
