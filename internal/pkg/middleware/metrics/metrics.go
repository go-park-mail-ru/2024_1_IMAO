package metrics

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

var hits = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "hits",
}, []string{"status", "handler"})

func CreateMetricsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			code := new(int)
			*code = 200
			request = request.WithContext(context.WithValue(request.Context(), "code", code))
			next.ServeHTTP(writer, request)

			log.Println(*code)

			//route := mux.CurrentRoute(request)
			//path, _ := route.GetPathTemplate()
		})
	}
}
