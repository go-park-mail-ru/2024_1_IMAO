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
	favouritesusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/favourites/usecases"
	userusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"

	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
)

type FavouritesHandler struct {
	storage       favouritesusecases.FavouritesStorageInterface
	advertStorage advertusecases.AdvertsStorageInterface
	userStorage   userusecases.UsersStorageInterface
}

func NewFavouritesHandler(storage favouritesusecases.FavouritesStorageInterface, advertStorage advertusecases.AdvertsStorageInterface, userStorage userusecases.UsersStorageInterface,
) *FavouritesHandler {
	return &FavouritesHandler{
		storage:       storage,
		advertStorage: advertStorage,
		userStorage:   userStorage,
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
func (favouritesHandler *FavouritesHandler) GetFavouritesList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	storage := favouritesHandler.storage
	userStorage := favouritesHandler.userStorage

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

	var adsList []*models.ReturningAdInList

	adsList, err = storage.GetFavouritesByUserID(ctx, uint(user.ID))

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewFavouritesErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}
	log.Println("Get favourites for user", user.ID)
	responses.SendOkResponse(writer, NewFavouritesOkResponse(adsList))
	logging.LogHandlerInfo(logger, fmt.Sprintf("Get favourites for user %s", fmt.Sprint(user.ID)), responses.StatusOk)
}

func (favouritesHandler *FavouritesHandler) ChangeFavourites(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	storage := favouritesHandler.storage
	userStorage := favouritesHandler.userStorage
	var data models.ReceivedCartItem

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewFavouritesErrResponse(responses.StatusInternalServerError,
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

	isAppended := storage.AppendAdvByIDs(ctx, user.ID, data.AdvertID, favouritesHandler.userStorage, favouritesHandler.advertStorage)

	responses.SendOkResponse(writer, NewFavouritesChangeResponse(isAppended))

	if isAppended {
		log.Println("Advert", data.AdvertID, "has been added to favourites of user", user.ID)
		logging.LogHandlerInfo(logger, fmt.Sprintf("Advert %s has been added to the favourites of user %s", fmt.Sprint(data.AdvertID), fmt.Sprint(user.ID)), responses.StatusOk)
	} else {
		log.Println("Advert", data.AdvertID, "has been removed from favourites of user", user.ID)
		logging.LogHandlerInfo(logger, fmt.Sprintf("Advert %s has been removed from the favourites of user %s", fmt.Sprint(data.AdvertID), fmt.Sprint(user.ID)), responses.StatusOk)
	}
}

func (favouritesHandler *FavouritesHandler) DeleteFromFavourites(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	storage := favouritesHandler.storage
	userStorage := favouritesHandler.userStorage
	var data models.ReceivedCartItems

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewFavouritesErrResponse(responses.StatusInternalServerError,
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
		err = storage.DeleteAdvByIDs(ctx, user.ID, item, favouritesHandler.userStorage, favouritesHandler.advertStorage)

		if err != nil {
			log.Println(err, responses.StatusBadRequest)
			logging.LogHandlerError(logger, err, responses.StatusBadRequest)
			responses.SendErrResponse(writer, NewFavouritesErrResponse(responses.StatusBadRequest,
				responses.ErrBadRequest))

			return
		}
	}

	log.Println("Adverts", data.AdvertIDs, "has been removed from favourites of user", user.ID)

	responses.SendOkResponse(writer, NewFavouritesChangeResponse(false))

	logging.LogHandlerInfo(logger, fmt.Sprintf("Adverts %s has been removed from favourites of user %s", fmt.Sprint(data.AdvertIDs), fmt.Sprint(user.ID)), responses.StatusOk)
}
