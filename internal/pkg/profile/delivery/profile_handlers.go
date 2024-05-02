package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"

	advdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/delivery"
	profileusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/usecases"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	userusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
)

type ProfileHandler struct {
	storage     profileusecases.ProfileStorageInterface
	userStorage userusecases.UsersStorageInterface
}

func NewProfileHandler(storage profileusecases.ProfileStorageInterface,
	userStorage userusecases.UsersStorageInterface) *ProfileHandler {
	return &ProfileHandler{
		storage:     storage,
		userStorage: userStorage,
	}
}

func (h *ProfileHandler) GetProfile(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	vars := mux.Vars(request)

	id, _ := strconv.Atoi(vars["id"])

	p, err := h.storage.GetProfileByUserID(ctx, uint(id))
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, advdel.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfileCity(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	userStorage := h.userStorage
	storage := h.storage

	session, err := request.Cookie("session_id")

	user, _ := userStorage.GetUserBySession(ctx, session.Value)

	var data models.City

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := storage.SetProfileCity(ctx, user.ID, data)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewProfileOkResponse(p))
}

// func (h *ProfileHandler) SetProfileRating(writer http.ResponseWriter, request *http.Request) {
// 	if request.Method != http.MethodPost {
// 		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

// 		return
// 	}

// 	vars := mux.Vars(request)
// 	userID, _ := strconv.Atoi(vars["id"])

// 	usersList := h.UsersList

// 	session, err := request.Cookie("session_id")

// 	if err != nil || !usersList.SessionExists(session.Value) {
// 		h.UsersList.Logger.Info("User not authorized")
// 		log.Println("User not authorized")
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusUnauthorized,
// 			responses.ErrUnauthorized))

// 		return
// 	}

// 	var data models.SetProfileRatingNec

// 	err = json.NewDecoder(request.Body).Decode(&data)
// 	if err != nil {
// 		h.UsersList.Logger.Error(err, responses.StatusInternalServerError)
// 		log.Println(err, responses.StatusInternalServerError)
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
// 			responses.ErrInternalServer))
// 	}

// 	p, err := h.ProfileList.SetProfileRating(uint(userID), data)
// 	if err != nil {
// 		h.UsersList.Logger.Error(err, responses.StatusBadRequest)
// 		log.Println(err, responses.StatusBadRequest)
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
// 			responses.ErrBadRequest))

// 		return
// 	}

// 	responses.SendOkResponse(writer, NewProfileOkResponse(p))
// }

func (h *ProfileHandler) SetProfilePhone(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	userStorage := h.userStorage
	storage := h.storage

	session, err := request.Cookie("session_id")

	user, _ := userStorage.GetUserBySession(ctx, session.Value)

	var data models.SetProfilePhoneNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := storage.SetProfilePhone(ctx, user.ID, data)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewProfileOkResponse(p))
}

func (h *ProfileHandler) EditProfile(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	userStorage := h.userStorage
	storage := h.storage

	session, err := request.Cookie("session_id")

	user, _ := userStorage.GetUserBySession(ctx, session.Value)

	err = request.ParseMultipartForm(2 << 20)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	avatar := request.MultipartForm.File["avatar"]
	data := models.EditProfileNec{
		Name:    request.PostFormValue("name"),
		Surname: request.PostFormValue("surname"),
	}

	var pl *models.Profile
	if len(avatar) != 0 {
		pl, err = storage.SetProfileInfo(ctx, user.ID, avatar[0], data)
	} else {
		pl, err = storage.SetProfileInfo(ctx, user.ID, nil, data)
	}

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.ErrInternalServer)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewProfileOkResponse(pl))
}

// func (h *ProfileHandler) SetProfileApproved(writer http.ResponseWriter, request *http.Request) {
// 	if request.Method != http.MethodPost {
// 		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

// 		return
// 	}

// 	ctx := request.Context()
// 	requestUUID := uuid.New().String()

// 	ctx = context.WithValue(ctx, "requestUUID", requestUUID)

// 	childLogger := h.logger.With(
// 		zap.String("requestUUID", requestUUID),
// 	)

// 	userStorage := h.userStorage
// 	storage := h.storage

// 	session, err := request.Cookie("session_id")

// 	if err != nil || !userStorage.SessionExists(session.Value) {
// 		childLogger.Info("User not authorized")
// 		log.Println("User not authorized")
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusUnauthorized,
// 			responses.ErrUnauthorized))

// 		return
// 	}

// 	user, _ := userStorage.GetUserBySession(ctx, session.Value)

// 	p, err := storage.SetProfileApproved(user.ID)
// 	if err != nil {
// 		childLogger.Error(err, responses.StatusBadRequest)
// 		log.Println(err, responses.StatusBadRequest)
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
// 			responses.ErrBadRequest))

// 		return
// 	}

// 	responses.SendOkResponse(writer, NewProfileOkResponse(p))
// }

// func (h *ProfileHandler) SetProfile(writer http.ResponseWriter, request *http.Request) {
// 	if request.Method != http.MethodPost {
// 		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

// 		return
// 	}

// 	ctx := request.Context()
// 	requestUUID := uuid.New().String()

// 	ctx = context.WithValue(ctx, "requestUUID", requestUUID)

// 	childLogger := h.UsersList.Logger.With(
// 		zap.String("requestUUID", requestUUID),
// 	)

// 	usersList := h.UsersList

// 	session, err := request.Cookie("session_id")

// 	if err != nil || !usersList.SessionExists(session.Value) {
// 		childLogger.Info("User not authorized")
// 		log.Println("User not authorized")
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusUnauthorized,
// 			responses.ErrUnauthorized))

// 		return
// 	}

// 	user, _ := usersList.GetUserBySession(ctx, session.Value)

// 	var data models.SetProfileNec

// 	err = json.NewDecoder(request.Body).Decode(&data)
// 	if err != nil {
// 		childLogger.Error(err, responses.StatusInternalServerError)
// 		log.Println(err, responses.StatusInternalServerError)
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
// 			responses.ErrInternalServer))
// 	}

// 	p, err := h.ProfileList.SetProfile(user.ID, data)
// 	if err != nil {
// 		childLogger.Error(err, responses.StatusBadRequest)
// 		log.Println(err, responses.StatusBadRequest)
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
// 			responses.ErrBadRequest))

// 		return
// 	}

// 	responses.SendOkResponse(writer, NewProfileOkResponse(p))
// }

// func (h *ProfileHandler) SetProfilePassword(writer http.ResponseWriter, request *http.Request) {
// 	if request.Method != http.MethodPost {
// 		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

// 		return
// 	}

// 	ctx := request.Context()
// 	requestUUID := uuid.New().String()

// 	ctx = context.WithValue(ctx, "requestUUID", requestUUID)

// 	childLogger := h.UsersList.Logger.With(
// 		zap.String("requestUUID", requestUUID),
// 	)

// 	usersList := h.UsersList

// 	session, err := request.Cookie("session_id")

// 	if err != nil || !usersList.SessionExists(session.Value) {
// 		childLogger.Info("User not authorized")
// 		log.Println("User not authorized")
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusUnauthorized,
// 			responses.ErrUnauthorized))

// 		return
// 	}

// 	user, _ := usersList.GetUserBySession(ctx, session.Value)

// 	var data models.SetProfileNec

// 	err = json.NewDecoder(request.Body).Decode(&data)
// 	if err != nil {
// 		childLogger.Error(err, responses.StatusInternalServerError)
// 		log.Println(err, responses.StatusInternalServerError)
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
// 			responses.ErrInternalServer))
// 	}

// 	p, err := h.ProfileList.SetProfile(user.ID, data)
// 	if err != nil {
// 		childLogger.Error(err, responses.StatusBadRequest)
// 		log.Println(err, responses.StatusBadRequest)
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
// 			responses.ErrBadRequest))

// 		return
// 	}

// 	responses.SendOkResponse(writer, NewProfileOkResponse(p))
// }

// func (h *ProfileHandler) SetProfileEmail(writer http.ResponseWriter, request *http.Request) {
// 	if request.Method != http.MethodPost {
// 		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

// 		return
// 	}

// 	ctx := request.Context()
// 	requestUUID := uuid.New().String()

// 	ctx = context.WithValue(ctx, "requestUUID", requestUUID)

// 	childLogger := h.UsersList.Logger.With(
// 		zap.String("requestUUID", requestUUID),
// 	)

// 	usersList := h.UsersList

// 	session, err := request.Cookie("session_id")

// 	if err != nil || !usersList.SessionExists(session.Value) {
// 		childLogger.Info("User not authorized")
// 		log.Println("User not authorized")
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusUnauthorized,
// 			responses.ErrUnauthorized))

// 		return
// 	}

// 	user, _ := usersList.GetUserBySession(ctx, session.Value)

// 	var data models.SetProfileNec

// 	err = json.NewDecoder(request.Body).Decode(&data)
// 	if err != nil {
// 		childLogger.Error(err, responses.StatusInternalServerError)
// 		log.Println(err, responses.StatusInternalServerError)
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
// 			responses.ErrInternalServer))
// 	}

// 	p, err := h.ProfileList.SetProfile(user.ID, data)
// 	if err != nil {
// 		childLogger.Error(err, responses.StatusBadRequest)
// 		log.Println(err, responses.StatusBadRequest)
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
// 			responses.ErrBadRequest))

// 		return
// 	}

// 	responses.SendOkResponse(writer, NewProfileOkResponse(p))
// }

// func (h *ProfileHandler) ProfileAdverts(writer http.ResponseWriter, request *http.Request) {
// 	if request.Method != http.MethodGet {
// 		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

// 		return
// 	}

// 	vars := mux.Vars(request)
// 	id, _ := strconv.Atoi(vars["id"])

// 	filter, _ := strconv.Atoi(request.URL.Query().Get("filter"))

// 	var ads []*models.ReturningAdvert
// 	var err error

// 	switch models.AdvertsFilter(filter) {
// 	case models.FilterAll:
// 		ads, err = h.AdvertsList.GetAdvertsByUserIDFiltered(uint(id),
// 			func(ad *models.Advert) bool {
// 				return true
// 			})
// 	case models.FilterActive:
// 		ads, err = h.AdvertsList.GetAdvertsByUserIDFiltered(uint(id),
// 			func(ad *models.Advert) bool {
// 				return ad.Active
// 			})
// 	default:
// 		ads, err = h.AdvertsList.GetAdvertsByUserIDFiltered(uint(id),
// 			func(ad *models.Advert) bool {
// 				return !ad.Active
// 			})
// 	}

// 	if err != nil {
// 		h.UsersList.Logger.Error(err, responses.StatusBadRequest)
// 		log.Println(err, responses.StatusBadRequest)
// 		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
// 			responses.ErrBadRequest))

// 		return
// 	}

// 	responses.SendOkResponse(writer, advdel.NewAdvertsOkResponse(ads))
// }
