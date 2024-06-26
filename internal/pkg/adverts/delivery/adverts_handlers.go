//nolint:ineffassign,gocritic,cyclop
package delivery

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"

	advertusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
)

const (
	defaultCity       = "Moscow"
	defaultAdverCount = 20
	maxMemory         = 2 << 28
)

type AdvertsHandler struct {
	storage    advertusecases.AdvertsStorageInterface
	authClient authproto.AuthClient
}

func NewAdvertsHandler(storage advertusecases.AdvertsStorageInterface,
	authClient authproto.AuthClient) *AdvertsHandler {
	return &AdvertsHandler{
		storage:    storage,
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
func (advertsHandler *AdvertsHandler) GetAdsList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	vars := mux.Vars(request)
	city := vars["city"]
	category := vars["category"]

	storage := advertsHandler.storage
	authClient := advertsHandler.authClient

	count, errCount := strconv.Atoi(request.URL.Query().Get("count"))
	startID, errstartID := strconv.Atoi(request.URL.Query().Get("startId"))
	userID, errUser := strconv.Atoi(request.URL.Query().Get("userId"))
	deleted, errdeleted := strconv.Atoi(request.URL.Query().Get("deleted"))

	if count == 0 {
		count = defaultAdverCount
	}

	if city == "" && request.URL.Query().Get("city") != "" {
		city = request.URL.Query().Get("city")
	} else if city == "" {
		city = defaultCity
	}

	var (
		adsList      []*models.ReturningAdInList
		err          error
		sessionValue string
	)

	session, cookieErr := request.Cookie("session_id")

	if session != nil {
		sessionValue = session.Value
	}

	var userIDCookie uint

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: sessionValue})

	if cookieErr == nil && user.IsAuth {
		userIDCookie = uint(user.ID)
	}

	if category != "" {
		adsList, err = storage.GetAdvertsByCategory(ctx, category, city, userIDCookie, uint(startID), uint(count))
	} else if errCount == nil && errstartID == nil {
		adsList, err = storage.GetAdvertsByCity(ctx, city, userIDCookie, uint(startID), uint(count))
	} else if errUser == nil && errdeleted == nil {
		adsList, err = storage.GetAdvertsForUserWhereStatusIs(ctx, userIDCookie, uint(userID),
			uint(deleted), uint(count))
	}

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(adsList))
}

func (advertsHandler *AdvertsHandler) GetAdsListWithSearch(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	storage := advertsHandler.storage
	authClient := advertsHandler.authClient
	count, errCount := strconv.Atoi(request.URL.Query().Get("count"))
	startID, errStartID := strconv.Atoi(request.URL.Query().Get("startId"))
	title := request.URL.Query().Get("title")

	var (
		adsList      []*models.ReturningAdInList
		err          error
		sessionValue string
	)

	session, cookieErr := request.Cookie("session_id")

	if session != nil {
		sessionValue = session.Value
	}

	var userIDCookie uint

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: sessionValue})

	if cookieErr == nil && user.IsAuth {
		userIDCookie = uint(user.ID)
	}

	if errCount == nil && errStartID == nil && title != "" {
		adsList, err = storage.SearchAdvertByTitle(ctx, title, userIDCookie, uint(startID), uint(count))
	} else {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(adsList))
}

func (advertsHandler *AdvertsHandler) GetSuggestions(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	storage := advertsHandler.storage
	num, errCount := strconv.Atoi(request.URL.Query().Get("num"))
	title := request.URL.Query().Get("title")

	var (
		suggestions []string
		err         error
	)

	if errCount == nil && title != "" {
		suggestions, err = storage.GetSuggestions(ctx, title, uint(num))
	} else {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(suggestions))
}

func (advertsHandler *AdvertsHandler) GetAdvert(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	storage := advertsHandler.storage
	authClient := advertsHandler.authClient

	sessionValue := ""

	session, cookieErr := request.Cookie("session_id")

	if session != nil {
		sessionValue = session.Value
	}

	var userIDCookie uint

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: sessionValue})

	if cookieErr == nil && user.IsAuth {
		userIDCookie = uint(user.ID)

		ownership := storage.CheckAdvertOwnership(ctx, uint(id), userIDCookie)

		if ownership {
			paymentList, err := utils.YuKassaUpdates()

			if err == nil {
				_ = storage.YuKassaUpdateDB(ctx, paymentList, uint(id)) // ПЕРЕПИСАТЬ ЧЕРЕЗ ПЕРЕСЕЧЕНИЕ МНОЖЕСТВ И BULK UPDATE
			}
		}
	}

	ad, err := storage.GetAdvert(ctx, userIDCookie, uint(id))
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	if cookieErr == nil && user.IsAuth {
		err = storage.InsertView(ctx, uint(user.ID), uint(id))
		if err != nil {
			logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
			log.Println(err, responses.StatusInternalServerError)
			responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
				responses.ErrInternalServer))

			return
		}
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(ad))
}

func (advertsHandler *AdvertsHandler) GetAdvertByID(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	storage := advertsHandler.storage

	ad, err := storage.GetAdvertOnlyByID(ctx, uint(id))
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(ad))
}

func (advertsHandler *AdvertsHandler) GetAdvertPriceHistoryByID(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	storage := advertsHandler.storage

	priceHistory, err := storage.GetPriceHistory(ctx, uint(id))
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(priceHistory))
}

func (advertsHandler *AdvertsHandler) CreateAdvert(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := request.ParseMultipartForm(maxMemory)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	storage := advertsHandler.storage

	isUsed := true
	if request.PostFormValue("condition") == "1" {
		isUsed = false
	}

	price, _ := strconv.Atoi(request.PostFormValue("price"))
	userID, _ := strconv.Atoi(request.PostFormValue("userId"))

	photos := request.MultipartForm.File["photos"]
	data := models.ReceivedAdData{
		UserID:      uint(userID),
		City:        request.PostFormValue("city"),
		Category:    request.PostFormValue("category"),
		Title:       request.PostFormValue("title"),
		Description: request.PostFormValue("description"),
		Price:       uint(price),
		IsUsed:      isUsed,
		Phone:       request.PostFormValue("phone"),
	}

	var advert *models.ReturningAdvert
	advert, err = storage.CreateAdvert(ctx, photos, data)

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(advert))
}

func (advertsHandler *AdvertsHandler) EditAdvert(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := request.ParseMultipartForm(maxMemory)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	storage := advertsHandler.storage
	isUsed := true

	if request.PostFormValue("condition") == "1" {
		isUsed = false
	}

	price, _ := strconv.Atoi(request.PostFormValue("price"))
	id, _ := strconv.Atoi(request.PostFormValue("id"))
	userID, _ := strconv.Atoi(request.PostFormValue("userId"))

	photos := request.MultipartForm.File["photos"]
	data := models.ReceivedAdData{
		ID:          uint(id),
		UserID:      uint(userID),
		City:        request.PostFormValue("city"),
		Category:    request.PostFormValue("category"),
		Title:       request.PostFormValue("title"),
		Description: request.PostFormValue("description"),
		Price:       uint(price),
		IsUsed:      isUsed,
		Phone:       request.PostFormValue("phone"),
	}

	var advert *models.ReturningAdvert
	advert, err = storage.EditAdvert(ctx, photos, data)

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(advert))
}

func (advertsHandler *AdvertsHandler) CloseAdvert(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	storage := advertsHandler.storage

	err := storage.CloseAdvert(ctx, uint(id))
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(nil))
}

func (advertsHandler *AdvertsHandler) GetPromotionData(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	session, cookieErr := request.Cookie("session_id")

	var userIDCookie uint

	authClient := advertsHandler.authClient
	storage := advertsHandler.storage

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	if cookieErr == nil && user.IsAuth {
		userIDCookie = uint(user.ID)

		ownership := storage.CheckAdvertOwnership(ctx, uint(id), userIDCookie)

		if ownership {
			paymentList, err := utils.YuKassaUpdates()

			if err == nil {
				_ = storage.YuKassaUpdateDB(ctx, paymentList, uint(id))
			}
		}
	}

	promotionData, err := storage.GetPromotionData(ctx, uint(id))
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(promotionData))
}
