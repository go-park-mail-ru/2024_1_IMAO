package myrecover

import (
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	"log"
	"net/http"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
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
