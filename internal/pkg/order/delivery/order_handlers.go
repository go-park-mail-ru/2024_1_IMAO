package delivery

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"go.uber.org/zap"

	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	authresp "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"

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
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	ctx := request.Context()

	storage := orderHandler.storage
	userStorage := orderHandler.userStorage

	session, err := request.Cookie("session_id")

	if err != nil || !userStorage.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendOkResponse(writer, authresp.NewAuthOkResponse(models.User{}, "", false))

		return
	}

	user, _ := userStorage.GetUserBySession(ctx, session.Value)

	var ordersList []*models.ReturningOrder

	ordersList, err = storage.GetReturningOrderByUserID(ctx, uint(user.ID), orderHandler.advertStorage)

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewOrderErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}
	log.Println("Get orders for user", user.ID)
	responses.SendOkResponse(writer, NewOrderOkResponse(ordersList))
}

func (orderHandler *OrderHandler) CreateOrder(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	ctx := request.Context()

	storage := orderHandler.storage
	cartStorage := orderHandler.cartStorage
	userStorage := orderHandler.userStorage

	var data models.ReceivedOrderItems

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewOrderErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, err := request.Cookie("session_id")

	if err != nil || !userStorage.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendOkResponse(writer, authresp.NewAuthOkResponse(models.User{}, "", false))

		return
	}

	user, _ := userStorage.GetUserBySession(ctx, session.Value)

	for _, receivedOrderItem := range data.Adverts {
		isDeleted := cartStorage.DeleteAdvByIDs(ctx, uint(user.ID), receivedOrderItem.AdvertID, userStorage, orderHandler.advertStorage)

		if isDeleted != nil {
			log.Println("Can not create an order", receivedOrderItem.AdvertID, "for user", user.ID)
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

	responses.SendOkResponse(writer, NewOrderCreateResponse(true))
}
