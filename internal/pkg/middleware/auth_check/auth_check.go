package authcheck

import (
	"context"
	"errors"

	"net/http"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	authusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func CreateAuthCheckMiddleware(storage authusecases.UsersStorageInterface) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))
			session, err := request.Cookie("session_id")

			if err != nil || !storage.SessionExists(session.Value) {
				if err == nil {
					err = errors.New("no such cookie in userStorage")
				}
				logging.LogHandlerError(logger, err, responses.StatusUnauthorized)
				responses.SendErrResponse(writer, responses.NewErrResponse(responses.StatusUnauthorized,
					responses.ErrUnauthorized))

				return
			}

			userID := storage.MAP_GetUserIDBySession(session.Value)

			sessionInstance := models.Session{
				UserID: uint32(userID),
				Value:  session.Value,
			}

			ctx = context.WithValue(ctx, config.SessionContextKey, sessionInstance)
			request = request.WithContext(ctx)
			next.ServeHTTP(writer, request)
		})
	}
}
