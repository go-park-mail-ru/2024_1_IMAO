package delivery

import (
	"encoding/json"
	"fmt"
	advertusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"go.uber.org/zap"

	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"

	cartusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/usecases"
	orderusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/usecases"
)

type OrderHandler struct {
	storage       orderusecases.OrderStorageInterface
	cartStorage   cartusecases.CartStorageInterface
	authClient    authproto.AuthClient
	advertStorage advertusecases.AdvertsStorageInterface
}

func NewOrderHandler(storage orderusecases.OrderStorageInterface, cartStorage cartusecases.CartStorageInterface,
	authClient authproto.AuthClient,
	advertStorage advertusecases.AdvertsStorageInterface) *OrderHandler {
	return &OrderHandler{
		storage:       storage,
		cartStorage:   cartStorage,
		advertStorage: advertStorage,
		authClient:    authClient,
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

	storage := orderHandler.storage
	authClient := orderHandler.authClient

	session, err := request.Cookie("session_id")

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	var ordersList []*models.ReturningOrder

	ordersList, err = storage.GetReturningOrderByUserID(ctx, uint(user.ID), orderHandler.advertStorage)

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
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

	storage := orderHandler.storage
	cartStorage := orderHandler.cartStorage
	authClient := orderHandler.authClient

	var data models.ReceivedOrderItems

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, err := request.Cookie("session_id")

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	// НИЖЕ БИЗНЕС ЛОГИКА, ЕЁ НУЖНО ВЫТЕСТИ В REPOSITORY
	for _, receivedOrderItem := range data.Adverts {
		isDeleted := cartStorage.DeleteAdvByIDs(ctx, uint(user.ID), receivedOrderItem.AdvertID)

		if isDeleted != nil {
			log.Println("Can not create an order", receivedOrderItem.AdvertID, "for user", user.ID)
			logging.LogHandlerInfo(logger, fmt.Sprintf("Can not create an order %s for user %s", fmt.Sprint(receivedOrderItem.AdvertID), fmt.Sprint(user.ID)),
				responses.StatusInternalServerError)
			responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
				responses.ErrInternalServer))
			return
		}

		storage.CreateOrderByID(uint(user.ID), receivedOrderItem)
		log.Println("An order", receivedOrderItem.AdvertID, "for user", user.ID, "successfully created")
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewOrderCreateResponse(true))

}
