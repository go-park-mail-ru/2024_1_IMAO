package routers

import (
	"github.com/gorilla/mux"

	delivery "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/survey/delivery"
)

func ServeSurveyRouter(router *mux.Router, surveyHandler *delivery.SurveyHandler,
	authCheckMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/survey").Subrouter()

	subrouter.
}
