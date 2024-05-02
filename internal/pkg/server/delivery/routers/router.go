package routers

import (
	createAuthCheckMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/auth_check"
	createCsrfMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/csrf"
	createLogMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/log"
	recoveryMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/recover"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"

	advdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/delivery"
	cartdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery"
	citydel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/delivery"
	orderdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/delivery"
	profdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery"
	surveydel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/survey/delivery"
	authdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"

	advusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	cartusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/usecases"
	cityusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/usecases"
	orderusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/usecases"
	profusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/usecases"
	surveyusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/survey/usecases"
	authusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.SugaredLogger,
	advertStorage advusecases.AdvertsStorageInterface,
	cartStorage cartusecases.CartStorageInterface,
	cityStorage cityusecases.CityStorageInterface,
	orderStorage orderusecases.OrderStorageInterface,
	profileStorage profusecases.ProfileStorageInterface,
	userStorage authusecases.UsersStorageInterface,
	surveyStorage surveyusecases.SurveyStorageInterface,
	authClient authproto.AuthClient) *mux.Router {

	router := mux.NewRouter()
	router.Use(recoveryMiddleware.RecoveryMiddleware)

	logMiddleware := createLogMiddleware.CreateLogMiddleware(logger)
	router.Use(logMiddleware)

	csrfMiddleware := createCsrfMiddleware.CreateCsrfMiddleware()
	authCheckMiddleware := createAuthCheckMiddleware.CreateAuthCheckMiddleware(authClient)

	advertsHandler := advdel.NewAdvertsHandler(advertStorage)
	cartHandler := cartdel.NewCartHandler(cartStorage, advertStorage, userStorage)
	authHandler := authdel.NewAuthHandler(authClient, profileStorage)
	profileHandler := profdel.NewProfileHandler(profileStorage, userStorage)
	orderHandler := orderdel.NewOrderHandler(orderStorage, cartStorage, advertStorage, userStorage)
	cityHandler := citydel.NewCityHandler(cityStorage)
	surveyHandler := surveydel.NewSurveyHandler(userStorage, surveyStorage)

	rootRouter := router.PathPrefix("/api").Subrouter()
	ServeAuthRouter(rootRouter, authHandler, authCheckMiddleware)
	ServeAdvertsRouter(rootRouter, advertsHandler, authCheckMiddleware, csrfMiddleware)
	ServeProfileRouter(rootRouter, profileHandler, authCheckMiddleware, csrfMiddleware)
	ServeCartRouter(rootRouter, cartHandler)
	ServeOrderRouter(rootRouter, orderHandler)
	ServeSurveyRouter(rootRouter, surveyHandler, authCheckMiddleware)

	rootRouter.HandleFunc("/city", cityHandler.GetCityList)

	return router
}
