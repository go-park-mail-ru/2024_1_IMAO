package csrf

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	"go.uber.org/zap"

	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/gorilla/mux"
)

const (
	tokenLen = 2
)

type HashToken struct {
	Secret []byte
}

func NewHMACHashToken(secret string) (*HashToken, error) {
	return &HashToken{Secret: []byte(secret)}, nil
}

func (tk *HashToken) Create(s *models.Session, tokenExpTime int64) (string, error) {
	h := hmac.New(sha256.New, tk.Secret)
	data := fmt.Sprintf("%s:%d:%d", s.Value, s.UserID, tokenExpTime)
	h.Write([]byte(data))
	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(tokenExpTime, 10)

	return token, nil
}

func (tk *HashToken) Check(s *models.Session, inputToken string) (bool, error) {
	tokenData := strings.Split(inputToken, ":")
	if len(tokenData) != tokenLen {
		return false, fmt.Errorf("bad token data")
	}

	tokenExp, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return false, fmt.Errorf("bad token time")
	}

	if tokenExp < time.Now().Unix() {
		return false, fmt.Errorf("token expired")
	}

	h := hmac.New(sha256.New, tk.Secret)
	data := fmt.Sprintf("%s:%d:%d", s.Value, s.UserID, tokenExp)

	h.Write([]byte(data))

	expectedMAC := h.Sum(nil)

	messageMAC, err := hex.DecodeString(tokenData[0])
	if err != nil {
		return false, fmt.Errorf("cand hex decode token")
	}

	return hmac.Equal(messageMAC, expectedMAC), nil
}

func CreateCsrfMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

			secret := "Vol4okSecretKey"
			hashToken, _ := NewHMACHashToken(secret)

			sessionInstance, ok := ctx.Value(config.SessionContextKey).(models.Session)
			if !ok {
				err := errors.New("error while getting sessionInstance from context")
				logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
				responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
					responses.ErrInternalServer))

				return
			}

			inputToken := request.PostFormValue("CSRFToken")

			isValid, err := hashToken.Check(&sessionInstance, inputToken)
			if err != nil || !isValid {
				logging.LogInfo(logger, "csrf is not valid")
				logging.LogHandlerError(logger, err, responses.StatusForbidden)
				responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusForbidden,
					responses.ErrForbidden))

				return
			}

			logging.LogInfo(logger, "success")
			next.ServeHTTP(writer, request)
		})
	}
}
