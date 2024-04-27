package routers

import (
	"github.com/gorilla/mux"

	delivery "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/survey/delivery"
)

func ServeSurveyRouter(router *mux.Router, surveyHandler *delivery.SurveyHandler,
	authCheckMiddleware mux.MiddlewareFunc) {
	subrouter := router.PathPrefix("/survey").Subrouter()
	subrouter.Use(authCheckMiddleware)

	subrouter.HandleFunc("/create", surveyHandler.CreateAnswer).Methods("POST")
	subrouter.HandleFunc("/statistics", surveyHandler.GetStatistics).Methods("GET")
	subrouter.HandleFunc("/check/{survey_id:[0-9]+}", surveyHandler.CheckIfAnswered).Methods("GET")
}
