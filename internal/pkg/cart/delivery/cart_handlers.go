package delivery

import (
	"encoding/json"
	"fmt"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"go.uber.org/zap"

	cartusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/usecases"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
)

type CartHandler struct {
	storage    cartusecases.CartStorageInterface
	authClient authproto.AuthClient
}

func NewCartHandler(storage cartusecases.CartStorageInterface, authClient authproto.AuthClient) *CartHandler {
	return &CartHandler{
		storage:    storage,
		authClient: authClient,
	}
}

// const (
// 	advertsPerPage = 30
// 	defaultCity    = "Moskva"
// )

// GetAdsList godoc
// @Summary Retrieve a list of adverts
// @Description Get a paginated list of adverts
// @Tags adverts
// @Accept json
// @Produce json
// @Success 200 {object} responses.AdvertsOkResponse
// @Failure 400 {object} responses.AdvertsErrResponse "Too many adverts specified"
// @Failure 405 {object} responses.AdvertsErrResponse "Method not allowed"
// @Failure 500 {object} responses.AdvertsErrResponse "Internal server error"
// @Router /api/adverts/list [get]
func (cartHandler *CartHandler) GetCartList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	storage := cartHandler.storage
	authClient := cartHandler.authClient

	session, _ := request.Cookie("session_id")

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})
	adsList, err := storage.GetCartByUserID(ctx, uint(user.ID))

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewCartErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}
	log.Println("Get cart for user", user.ID)
	responses.SendOkResponse(writer, NewCartOkResponse(adsList))
	logging.LogHandlerInfo(logger, fmt.Sprintf("Get cart for user %s", fmt.Sprint(user.ID)), responses.StatusOk)
}

func (cartHandler *CartHandler) ChangeCart(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	storage := cartHandler.storage
	authClient := cartHandler.authClient
	var data models.ReceivedCartItem

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewCartErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, _ := request.Cookie("session_id")

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})
	isAppended := storage.AppendAdvByIDs(ctx, uint(user.ID), data.AdvertID)

	responses.SendOkResponse(writer, NewCartChangeResponse(isAppended))

	if isAppended {
		log.Println("Advert", data.AdvertID, "has been added to cart of user", user.ID)
		logging.LogHandlerInfo(logger, fmt.Sprintf("Advert %s has been added to the cart of user %s", fmt.Sprint(data.AdvertID), fmt.Sprint(user.ID)), responses.StatusOk)
	} else {
		log.Println("Advert", data.AdvertID, "has been removed from cart of user", user.ID)
		logging.LogHandlerInfo(logger, fmt.Sprintf("Advert %s has been removed from thecart of user %s", fmt.Sprint(data.AdvertID), fmt.Sprint(user.ID)), responses.StatusOk)
	}
}

func (cartHandler *CartHandler) DeleteFromCart(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	storage := cartHandler.storage
	authClient := cartHandler.authClient
	var data models.ReceivedCartItems

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewCartErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, _ := request.Cookie("session_id")

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	for _, item := range data.AdvertIDs {
		err = storage.DeleteAdvByIDs(ctx, uint(user.ID), item)

		if err != nil {
			log.Println(err, responses.StatusBadRequest)
			logging.LogHandlerError(logger, err, responses.StatusBadRequest)
			responses.SendErrResponse(writer, NewCartErrResponse(responses.StatusBadRequest,
				responses.ErrBadRequest))

			return
		}
	}

	log.Println("Adverts", data.AdvertIDs, "has been removed from cart of user", user.ID)

	responses.SendOkResponse(writer, NewCartChangeResponse(false))

	logging.LogHandlerInfo(logger, fmt.Sprintf("Adverts %s has been removed from cart of user %s", fmt.Sprint(data.AdvertIDs), fmt.Sprint(user.ID)), responses.StatusOk)
}
