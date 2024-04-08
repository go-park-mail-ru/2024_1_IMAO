package delivery

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	advrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/repository"
	cartrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/repository"
	orderrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/order/repository"
	authrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/repository"

	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	authresp "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"
)

type OrderHandler struct {
	ListOrder   *orderrepo.OrderListWrapper
	ListCart    *cartrepo.CartListWrapper
	ListAdverts *advrepo.AdvertsListWrapper
	ListUsers   *authrepo.UsersListWrapper
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

	list := orderHandler.ListOrder
	usersList := orderHandler.ListUsers

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendOkResponse(writer, authresp.NewAuthOkResponse(models.User{}, "", false))

		return
	}

	user, _ := usersList.GetUserBySession(ctx, session.Value)

	var ordersList []*models.ReturningOrder

	ordersList, err = list.GetReturningOrderByUserID(uint(user.ID), orderHandler.ListAdverts)

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

	list := orderHandler.ListOrder
	cartlist := orderHandler.ListCart
	usersList := orderHandler.ListUsers

	var data models.ReceivedOrderItems

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewOrderErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendOkResponse(writer, authresp.NewAuthOkResponse(models.User{}, "", false))

		return
	}

	user, _ := usersList.GetUserBySession(ctx, session.Value)

	for _, receivedOrderItem := range data.Adverts {
		isDeleted := cartlist.DeleteAdvByIDs(uint(user.ID), receivedOrderItem.AdvertID, usersList, orderHandler.ListAdverts)

		if isDeleted != nil {
			log.Println("Can not create an order", receivedOrderItem.AdvertID, "for user", user.ID)
			responses.SendErrResponse(writer, NewOrderErrResponse(responses.StatusInternalServerError,
				responses.ErrInternalServer))
			return
		}

		list.CreateOrderByID(uint(user.ID), receivedOrderItem, orderHandler.ListAdverts)
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
