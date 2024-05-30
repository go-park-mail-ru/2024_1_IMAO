package delivery

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"go.uber.org/zap"

	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"

	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"

	advertusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	favouritesusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/favourites/usecases"

	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
)

type FavouritesHandler struct {
	storage       favouritesusecases.FavouritesStorageInterface
	advertStorage advertusecases.AdvertsStorageInterface
	authClient    authproto.AuthClient
}

func NewFavouritesHandler(storage favouritesusecases.FavouritesStorageInterface,
	advertStorage advertusecases.AdvertsStorageInterface, authClient authproto.AuthClient) *FavouritesHandler {
	return &FavouritesHandler{
		storage:       storage,
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
func (favouritesHandler *FavouritesHandler) GetFavouritesList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	storage := favouritesHandler.storage
	authClient := favouritesHandler.authClient

	session, err := request.Cookie("session_id")

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	var adsList []*models.ReturningAdInList

	adsList, err = storage.GetFavouritesByUserID(ctx, uint(user.ID))

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	log.Println("Get favourites for user", user.ID)
	responses.SendOkResponse(writer, responses.NewOkResponse(adsList))
	logging.LogHandlerInfo(logger, fmt.Sprintf("Get favourites for user %s", fmt.Sprint(user.ID)),
		responses.StatusOk)
}

func (favouritesHandler *FavouritesHandler) ChangeFavourites(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	storage := favouritesHandler.storage
	authClient := favouritesHandler.authClient

	var data models.ReceivedCartItem

	reqData, _ := io.ReadAll(request.Body)

	err := data.UnmarshalJSON(reqData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, err := request.Cookie("session_id")

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	isAppended := storage.AppendAdvByIDs(ctx, uint(user.ID), data.AdvertID)

	responses.SendOkResponse(writer, responses.NewOkResponse(models.Appended{IsAppended: isAppended}))

	if isAppended {
		log.Println("Advert", data.AdvertID, "has been added to favourites of user", user.ID)
		logging.LogHandlerInfo(logger, fmt.Sprintf("Advert %s has been added to the favourites of user %s",
			fmt.Sprint(data.AdvertID), fmt.Sprint(user.ID)), responses.StatusOk)
	} else {
		log.Println("Advert", data.AdvertID, "has been removed from favourites of user", user.ID)
		logging.LogHandlerInfo(logger, fmt.Sprintf("Advert %s has been removed from the favourites of user %s",
			fmt.Sprint(data.AdvertID), fmt.Sprint(user.ID)), responses.StatusOk)
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
	authClient := favouritesHandler.authClient

	var data models.ReceivedCartItems

	reqData, _ := io.ReadAll(request.Body)

	err := data.UnmarshalJSON(reqData)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, err := request.Cookie("session_id")

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	for _, item := range data.AdvertIDs {
		err = storage.DeleteAdvByIDs(ctx, uint(user.ID), item)

		if err != nil {
			log.Println(err, responses.StatusBadRequest)
			logging.LogHandlerError(logger, err, responses.StatusBadRequest)
			responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
				responses.ErrBadRequest))

			return
		}
	}

	log.Println("Adverts", data.AdvertIDs, "has been removed from favourites of user", user.ID)

	responses.SendOkResponse(writer, responses.NewOkResponse(models.Appended{IsAppended: false}))

	logging.LogHandlerInfo(logger, fmt.Sprintf("Adverts %s has been removed from favourites of user %s",
		fmt.Sprint(data.AdvertIDs), fmt.Sprint(user.ID)), responses.StatusOk)
}

func (favouritesHandler *FavouritesHandler) GetSubscribedAdverts(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	storage := favouritesHandler.storage
	authClient := favouritesHandler.authClient

	session, err := request.Cookie("session_id")

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	var adsList []*models.ReturningAdInList

	adsList, err = storage.GetSubscribedAdvertsByUserID(ctx, uint(user.ID))

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	log.Println("Get subscribed adverts for user", user.ID)
	responses.SendOkResponse(writer, responses.NewOkResponse(adsList))
	logging.LogHandlerInfo(logger, fmt.Sprintf("Get subscribed adverts for user %s", fmt.Sprint(user.ID)),
		responses.StatusOk)
}
