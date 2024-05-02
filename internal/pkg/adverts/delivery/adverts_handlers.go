package delivery

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"

	advertusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	userusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
)

const (
	defaultCity = "Moscow"
)

type AdvertsHandler struct {
	storage     advertusecases.AdvertsStorageInterface
	userStorage userusecases.UsersStorageInterface
	addrOrigin  string
	schema      string
	logger      *zap.SugaredLogger
}

func NewAdvertsHandler(storage advertusecases.AdvertsStorageInterface, userStorage userusecases.UsersStorageInterface,
	addrOrigin string, schema string, logger *zap.SugaredLogger) *AdvertsHandler {
	return &AdvertsHandler{
		storage:     storage,
		userStorage: userStorage,
		addrOrigin:  addrOrigin,
		schema:      schema,
		logger:      logger,
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

	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	city := vars["city"]
	category := vars["category"]

	storage := advertsHandler.storage
	userStorage := advertsHandler.userStorage

	count, errCount := strconv.Atoi(request.URL.Query().Get("count"))
	startID, errstartID := strconv.Atoi(request.URL.Query().Get("startId"))
	userID, errUser := strconv.Atoi(request.URL.Query().Get("userId"))
	deleted, errdeleted := strconv.Atoi(request.URL.Query().Get("deleted"))

	if city == "" && request.URL.Query().Get("city") != "" {
		city = request.URL.Query().Get("city")
	} else if city == "" {
		city = defaultCity
	}

	var adsList []*models.ReturningAdInList
	var err error

	session, cookieErr := request.Cookie("session_id")

	var userIdCookie uint = 0

	if cookieErr == nil && userStorage.SessionExists(session.Value) {
		userIdCookie = userStorage.MAP_GetUserIDBySession(session.Value)

	}

	if category != "" {
		adsList, err = storage.GetAdvertsByCategory(ctx, category, city, userIdCookie, uint(startID), uint(count))
	} else if errCount == nil && errstartID == nil {
		adsList, err = storage.GetAdvertsByCity(ctx, city, userIdCookie, uint(startID), uint(count))
	} else if errUser == nil && errdeleted == nil {
		adsList, err = storage.GetAdvertsForUserWhereStatusIs(ctx, userIdCookie, uint(userID), uint(deleted))
	}
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewAdvertsOkResponse(adsList))
}

func (advertsHandler *AdvertsHandler) GetAdsListWithSearch(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	storage := advertsHandler.storage
	userStorage := advertsHandler.userStorage

	count, errCount := strconv.Atoi(request.URL.Query().Get("count"))
	startID, errStartID := strconv.Atoi(request.URL.Query().Get("startId"))
	title := request.URL.Query().Get("title")
	city := request.URL.Query().Get("city")

	if city == "" {
		city = defaultCity
	}

	var adsList []*models.ReturningAdInList
	var err error

	session, cookieErr := request.Cookie("session_id")

	var userIdCookie uint = 0

	if cookieErr == nil && userStorage.SessionExists(session.Value) {
		userIdCookie = userStorage.MAP_GetUserIDBySession(session.Value)

	}

	if errCount == nil && errStartID == nil && title != "" {
		adsList, err = storage.SearchAdvertByTitle(ctx, title, userIdCookie, uint(startID), uint(count))
	} else {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewAdvertsOkResponse(adsList))
}

func (advertsHandler *AdvertsHandler) GetSuggestions(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	storage := advertsHandler.storage
	num, errCount := strconv.Atoi(request.URL.Query().Get("num"))
	title := request.URL.Query().Get("title")
	city := request.URL.Query().Get("city")

	if city == "" {
		city = defaultCity
	}

	var suggestions []string
	var err error

	if errCount == nil && title != "" {
		suggestions, err = storage.GetSuggestions(ctx, title, uint(num))
	} else {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	if err != nil {
		fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewAdvertsOkResponse(suggestions))
}

func (advertsHandler *AdvertsHandler) GetAdvert(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	city := vars["city"]
	category := vars["category"]
	id, _ := strconv.Atoi(vars["id"])

	storage := advertsHandler.storage
	userStorage := advertsHandler.userStorage

	session, cookieErr := request.Cookie("session_id")

	var userIdCookie uint = 0

	if cookieErr == nil && userStorage.SessionExists(session.Value) {
		userIdCookie = userStorage.MAP_GetUserIDBySession(session.Value)

	}

	ad, err := storage.GetAdvert(ctx, userIdCookie, uint(id), city, category)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	if cookieErr == nil && userStorage.SessionExists(session.Value) {
		userID := userStorage.MAP_GetUserIDBySession(session.Value)
		err = storage.InsertView(ctx, userID, uint(id))
		if err != nil {
			logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
			log.Println(err, responses.StatusInternalServerError)
			responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusInternalServerError,
				responses.ErrInternalServer))

			return
		}
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewAdvertsOkResponse(ad))
}

func (advertsHandler *AdvertsHandler) GetAdvertByID(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	storage := advertsHandler.storage

	ad, err := storage.GetAdvertByOnlyByID(ctx, uint(id))
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewAdvertsOkResponse(ad))
}

func (advertsHandler *AdvertsHandler) CreateAdvert(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	err := request.ParseMultipartForm(2 << 28)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusInternalServerError,
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
	}

	var advert *models.ReturningAdvert
	advert, err = storage.CreateAdvert(ctx, photos, data)

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewAdvertsOkResponse(advert))
}

func (advertsHandler *AdvertsHandler) EditAdvert(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	err := request.ParseMultipartForm(2 << 28)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusInternalServerError,
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
	}

	var advert *models.ReturningAdvert
	advert, err = storage.EditAdvert(ctx, photos, data)

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewAdvertsOkResponse(advert))
}

// func (advertsHandler *AdvertsHandler) DeleteAdvert(writer http.ResponseWriter, request *http.Request) {
// 	if request.Method != http.MethodPost {
// 		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

// 		return
// 	}

// 	ctx := request.Context()
// 	requestUUID := uuid.New().String()

// 	ctx = context.WithValue(ctx, "requestUUID", requestUUID)

// 	childLogger := advertsHandler.logger.With(
// 		zap.String("requestUUID", requestUUID),
// 	)

// 	vars := mux.Vars(request)
// 	id, _ := strconv.Atoi(vars["id"])

// 	storage := advertsHandler.storage

// 	err := storage.DeleteAdvert(uint(id))
// 	if err != nil {
// 		childLogger.Error(err, responses.StatusBadRequest)
// 		log.Println(err, responses.StatusBadRequest)
// 		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
// 			responses.ErrBadRequest))

// 		return
// 	}

// 	adResponse := NewAdvertsOkResponse(nil)

// 	responses.SendOkResponse(writer, adResponse)
// }

func (advertsHandler *AdvertsHandler) CloseAdvert(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	storage := advertsHandler.storage

	err := storage.CloseAdvert(ctx, uint(id))
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	adResponse := NewAdvertsOkResponse(nil)

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, adResponse)
}
