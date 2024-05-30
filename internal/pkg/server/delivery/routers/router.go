package routers

import (
	"log"

	advdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/delivery"
	cartdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery"
	cartproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery/protobuf"
	citydel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/delivery"
	favdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/favourites/delivery"
	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"
	createAuthCheckMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/auth_check"
	createCsrfMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/csrf"
	createLogMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/log"
	createMetricsMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/metrics"
	recoveryMiddleware "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/recover"
	orderdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/delivery"
	paydel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/payments/delivery"
	profdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery"
	profileproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery/protobuf"
	surveydel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/survey/delivery"
	authdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	advusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	cartusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/usecases"
	cityusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/usecases"
	favusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/favourites/usecases"
	orderusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/usecases"
	paymentsusescases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/payments/usecases"
	surveyusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/survey/usecases"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.SugaredLogger,
	advertStorage advusecases.AdvertsStorageInterface,
	cartClient cartproto.CartClient,
	cartStorage cartusecases.CartStorageInterface,
	cityStorage cityusecases.CityStorageInterface,
	orderStorage orderusecases.OrderStorageInterface,
	surveyStorage surveyusecases.SurveyStorageInterface,
	authClient authproto.AuthClient,
	profileClient profileproto.ProfileClient,
	favouritesStorage favusecases.FavouritesStorageInterface,
	paymentsStorage paymentsusescases.PaymentsStorageInterface) *mux.Router {
	router := mux.NewRouter()
	router.Use(recoveryMiddleware.RecoveryMiddleware)

	metrics, err := mymetrics.CreateHTTPMetrics("main")
	if err != nil {
		log.Println("error occurred while creating metrics", err)

		return nil
	}

	logMiddleware := createLogMiddleware.CreateLogMiddleware(logger)
	metricsMiddleware := createMetricsMiddleware.CreateMetricsMiddleware(metrics)

	router.Use(metricsMiddleware)
	router.Use(logMiddleware)

	csrfMiddleware := createCsrfMiddleware.CreateCsrfMiddleware()
	authCheckMiddleware := createAuthCheckMiddleware.CreateAuthCheckMiddleware(authClient)

	advertsHandler := advdel.NewAdvertsHandler(advertStorage, authClient)
	cartHandler := cartdel.NewCartHandler(cartClient, authClient)
	authHandler := authdel.NewAuthHandler(authClient, profileClient)
	profileHandler := profdel.NewProfileHandler(profileClient, authClient)
	orderHandler := orderdel.NewOrderHandler(orderStorage, cartStorage, authClient, advertStorage)
	cityHandler := citydel.NewCityHandler(cityStorage)
	surveyHandler := surveydel.NewSurveyHandler(authClient, surveyStorage)
	favouritesHandler := favdel.NewFavouritesHandler(favouritesStorage, advertStorage, authClient)
	paymentsHandler := paydel.NewPaymentsHandler(paymentsStorage, authClient)

	rootRouter := router.PathPrefix("/api").Subrouter()
	ServeAuthRouter(rootRouter, authHandler, authCheckMiddleware)
	ServeAdvertsRouter(rootRouter, advertsHandler, authCheckMiddleware, csrfMiddleware)
	ServeProfileRouter(rootRouter, profileHandler, authCheckMiddleware, csrfMiddleware)
	ServeCartRouter(rootRouter, cartHandler, authCheckMiddleware)
	ServeOrderRouter(rootRouter, orderHandler, authCheckMiddleware)
	ServeSurveyRouter(rootRouter, surveyHandler, authCheckMiddleware)
	ServeFavouritesRouter(rootRouter, favouritesHandler, authCheckMiddleware)
	ServePaymentsRouter(rootRouter, paymentsHandler, authCheckMiddleware)
	ServeCityRouter(rootRouter, cityHandler)

	router.PathPrefix("/metrics").Handler(promhttp.Handler())

	return router
}
