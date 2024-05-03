package routers

import (
	createAuthCheckMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/auth_check"
	createCsrfMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/csrf"
	createLogMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/log"
	recoveryMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/recover"
	profileproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery/protobuf"
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
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.SugaredLogger,
	advertStorage advusecases.AdvertsStorageInterface,
	cartStorage cartusecases.CartStorageInterface,
	cityStorage cityusecases.CityStorageInterface,
	orderStorage orderusecases.OrderStorageInterface,
	profileStorage profusecases.ProfileStorageInterface,
	surveyStorage surveyusecases.SurveyStorageInterface,
	authClient authproto.AuthClient,
	profileClient profileproto.ProfileClient) *mux.Router {

	router := mux.NewRouter()
	router.Use(recoveryMiddleware.RecoveryMiddleware)

	logMiddleware := createLogMiddleware.CreateLogMiddleware(logger)
	router.Use(logMiddleware)

	csrfMiddleware := createCsrfMiddleware.CreateCsrfMiddleware()
	authCheckMiddleware := createAuthCheckMiddleware.CreateAuthCheckMiddleware(authClient)

	advertsHandler := advdel.NewAdvertsHandler(advertStorage)
	cartHandler := cartdel.NewCartHandler(cartStorage, authClient)
	authHandler := authdel.NewAuthHandler(authClient, profileStorage)
	profileHandler := profdel.NewProfileHandler(profileClient, authClient)
	orderHandler := orderdel.NewOrderHandler(orderStorage, cartStorage, authClient, advertStorage)
	cityHandler := citydel.NewCityHandler(cityStorage)
	surveyHandler := surveydel.NewSurveyHandler(authClient, surveyStorage)

	rootRouter := router.PathPrefix("/api").Subrouter()
	ServeAuthRouter(rootRouter, authHandler, authCheckMiddleware)
	ServeAdvertsRouter(rootRouter, advertsHandler, authCheckMiddleware, csrfMiddleware)
	ServeProfileRouter(rootRouter, profileHandler, authCheckMiddleware, csrfMiddleware)
	ServeCartRouter(rootRouter, cartHandler, authCheckMiddleware)
	ServeOrderRouter(rootRouter, orderHandler, authCheckMiddleware)
	ServeSurveyRouter(rootRouter, surveyHandler, authCheckMiddleware)

	rootRouter.HandleFunc("/city", cityHandler.GetCityList)

	return router
}
