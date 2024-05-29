package delivery

import (
	"fmt"
	"io"
	"log"
	"net/http"

	cartproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/delivery/protobuf"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"go.uber.org/zap"

	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
)

type CartHandler struct {
	cartClient cartproto.CartClient
	authClient authproto.AuthClient
}

func NewCartHandler(cartClient cartproto.CartClient, authClient authproto.AuthClient) *CartHandler {
	return &CartHandler{
		cartClient: cartClient,
		authClient: authClient,
	}
}

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

	cartClient := cartHandler.cartClient
	authClient := cartHandler.authClient

	session, _ := request.Cookie("session_id")

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})
	adsList, err := cartClient.GetCartByUserID(ctx, &cartproto.UserIdRequest{UserId: uint32(user.ID)})

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	log.Println("Get cart for user", user.ID)
	responses.SendOkResponse(writer, responses.NewOkResponse(ReturningAdvertItem(adsList)))
	logging.LogHandlerInfo(logger, fmt.Sprintf("Get cart for user %s", fmt.Sprint(user.ID)), responses.StatusOk)
}

func (cartHandler *CartHandler) ChangeCart(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	cartClient := cartHandler.cartClient
	authClient := cartHandler.authClient

	var data models.ReceivedCartItem

	reqData, _ := io.ReadAll(request.Body)

	err := data.UnmarshalJSON(reqData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, _ := request.Cookie("session_id")

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})
	isAppended, _ := cartClient.AppendAdvByIDs(ctx,
		&cartproto.UserIdAdvertIdRequest{UserId: uint32(user.ID), AdvertId: uint32(data.AdvertID)})

	responses.SendOkResponse(writer, responses.NewOkResponse(isAppended.IsAppended))

	if isAppended.IsAppended {
		log.Println("Advert", data.AdvertID, "has been added to cart of user", user.ID)
		logging.LogHandlerInfo(logger, fmt.Sprintf("Advert %s has been added to the cart of user %s",
			fmt.Sprint(data.AdvertID), fmt.Sprint(user.ID)), responses.StatusOk)
	} else {
		log.Println("Advert", data.AdvertID, "has been removed from cart of user", user.ID)
		logging.LogHandlerInfo(logger, fmt.Sprintf("Advert %s has been removed from thecart of user %s",
			fmt.Sprint(data.AdvertID), fmt.Sprint(user.ID)), responses.StatusOk)
	}
}

func (cartHandler *CartHandler) DeleteFromCart(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	cartClient := cartHandler.cartClient
	authClient := cartHandler.authClient

	var data models.ReceivedCartItems

	reqData, _ := io.ReadAll(request.Body)

	err := data.UnmarshalJSON(reqData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, _ := request.Cookie("session_id")

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	for _, item := range data.AdvertIDs {
		_, err := cartClient.DeleteAdvByIDs(ctx,
			&cartproto.UserIdAdvertIdRequest{UserId: uint32(user.ID), AdvertId: uint32(item)})

		if err != nil {
			log.Println(err, responses.StatusBadRequest)
			logging.LogHandlerError(logger, err, responses.StatusBadRequest)
			responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
				responses.ErrBadRequest))

			return
		}
	}

	log.Println("Adverts", data.AdvertIDs, "has been removed from cart of user", user.ID)

	responses.SendOkResponse(writer, responses.NewOkResponse(models.Appended{IsAppended: false}))

	logging.LogHandlerInfo(logger, fmt.Sprintf("Adverts %s has been removed from cart of user %s",
		fmt.Sprint(data.AdvertIDs), fmt.Sprint(user.ID)), responses.StatusOk)
}
