package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"go.uber.org/zap"

	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	authresp "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"

	advertusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	cartusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/usecases"
	orderusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/usecases"
	userusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
)

type OrderHandler struct {
	storage       orderusecases.OrderStorageInterface
	cartStorage   cartusecases.CartStorageInterface
	advertStorage advertusecases.AdvertsStorageInterface
	userStorage   userusecases.UsersStorageInterface
	addrOrigin    string
	schema        string
	logger        *zap.SugaredLogger
}

func NewOrderHandler(storage orderusecases.OrderStorageInterface, cartStorage cartusecases.CartStorageInterface, advertStorage advertusecases.AdvertsStorageInterface, userStorage userusecases.UsersStorageInterface,
	addrOrigin string, schema string, logger *zap.SugaredLogger) *OrderHandler {
	return &OrderHandler{
		storage:       storage,
		cartStorage:   cartStorage,
		advertStorage: advertStorage,
		userStorage:   userStorage,
		addrOrigin:    addrOrigin,
		schema:        schema,
		logger:        logger,
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
func (orderHandler *OrderHandler) GetOrderList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	storage := orderHandler.storage
	userStorage := orderHandler.userStorage

	session, err := request.Cookie("session_id")

	if err != nil || !userStorage.SessionExists(session.Value) {
		if err == nil {
			err = errors.New("no such cookie in userStorage")
		}
		logging.LogHandlerError(logger, err, responses.StatusUnauthorized)
		log.Println("User not authorized")
		responses.SendOkResponse(writer, authresp.NewAuthOkResponse(models.User{}, "", false))

		return
	}

	user, _ := userStorage.GetUserBySession(ctx, session.Value)

	var ordersList []*models.ReturningOrder

	ordersList, err = storage.GetReturningOrderByUserID(ctx, uint(user.ID), orderHandler.advertStorage)

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewOrderErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	log.Println("Get orders for user", user.ID)
	logging.LogHandlerInfo(logger, fmt.Sprintf("Get orders for user %s", fmt.Sprint(user.ID)), responses.StatusOk)
	responses.SendOkResponse(writer, NewOrderOkResponse(ordersList))
}

func (orderHandler *OrderHandler) CreateOrder(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	storage := orderHandler.storage
	cartStorage := orderHandler.cartStorage
	userStorage := orderHandler.userStorage

	var data models.ReceivedOrderItems

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewOrderErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, err := request.Cookie("session_id")

	if err != nil || !userStorage.SessionExists(session.Value) {
		if err == nil {
			err = errors.New("no such cookie in userStorage")
		}
		logging.LogHandlerError(logger, err, responses.StatusUnauthorized)
		log.Println("User not authorized")
		responses.SendOkResponse(writer, authresp.NewAuthOkResponse(models.User{}, "", false))

		return
	}

	user, _ := userStorage.GetUserBySession(ctx, session.Value)

	// НИЖЕ БИЗНЕС ЛОГИКА, ЕЁ НУЖНО ВЫТЕСТИ В REPOSITORY
	for _, receivedOrderItem := range data.Adverts {
		isDeleted := cartStorage.DeleteAdvByIDs(ctx, uint(user.ID), receivedOrderItem.AdvertID, userStorage, orderHandler.advertStorage)

		if isDeleted != nil {
			log.Println("Can not create an order", receivedOrderItem.AdvertID, "for user", user.ID)
			logging.LogHandlerInfo(logger, fmt.Sprintf("Can not create an order %s for user %s", fmt.Sprint(receivedOrderItem.AdvertID), fmt.Sprint(user.ID)),
				responses.StatusInternalServerError)
			responses.SendErrResponse(writer, NewOrderErrResponse(responses.StatusInternalServerError,
				responses.ErrInternalServer))
			return
		}

		storage.CreateOrderByID(uint(user.ID), receivedOrderItem, orderHandler.advertStorage)
		log.Println("An order", receivedOrderItem.AdvertID, "for user", user.ID, "successfully created")
	}

	// isAppended := list.AppendAdvByIDs(user.ID, data.AdvertID, cartHandler.ListUsers, cartHandler.ListAdverts)

	// if isAppended {
	// 	log.Println("Advert", data.AdvertID, "has been added to cart of user", user.ID)
	// } else {
	// 	log.Println("Advert", data.AdvertID, "has been removed from cart of user", user.ID)
	// }

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewOrderCreateResponse(true))

}
