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

	advertusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	cartusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/cart/usecases"
	userusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"

	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
)

type CartHandler struct {
	storage       cartusecases.CartStorageInterface
	advertStorage advertusecases.AdvertsStorageInterface
	userStorage   userusecases.UsersStorageInterface
	addrOrigin    string
	schema        string
	logger        *zap.SugaredLogger
}

func NewCartHandler(storage cartusecases.CartStorageInterface, advertStorage advertusecases.AdvertsStorageInterface, userStorage userusecases.UsersStorageInterface,
	addrOrigin string, schema string, logger *zap.SugaredLogger) *CartHandler {
	return &CartHandler{
		storage:       storage,
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
func (cartHandler *CartHandler) GetCartList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	storage := cartHandler.storage
	userStorage := cartHandler.userStorage

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

	var adsList []*models.ReturningAdvert

	adsList, err = storage.GetCartByUserID(ctx, uint(user.ID), cartHandler.userStorage, cartHandler.advertStorage)

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
	userStorage := cartHandler.userStorage
	var data models.ReceivedCartItem

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewCartErrResponse(responses.StatusInternalServerError,
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

	isAppended := storage.AppendAdvByIDs(ctx, user.ID, data.AdvertID, cartHandler.userStorage, cartHandler.advertStorage)

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
	userStorage := cartHandler.userStorage
	var data models.ReceivedCartItems

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewCartErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, err := request.Cookie("session_id")

	if err != nil || !userStorage.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendOkResponse(writer, authresp.NewAuthOkResponse(models.User{}, "", false))

		return
	}

	user, _ := userStorage.GetUserBySession(ctx, session.Value)

	for _, item := range data.AdvertIDs {
		err = storage.DeleteAdvByIDs(ctx, user.ID, item, cartHandler.userStorage, cartHandler.advertStorage)

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
