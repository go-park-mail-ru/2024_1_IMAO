package metrics

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateMetricsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			code := new(int)
			*code = 200
			request = request.WithContext(context.WithValue(request.Context(), "code", code))
			next.ServeHTTP(writer, request)

			//route := mux.CurrentRoute(request)
			//path, _ := route.GetPathTemplate()
		})
	}
}
