package delivery

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	profileproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery/protobuf"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"

	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
)

const (
	maxMemory = 2 << 20
)

type ProfileHandler struct {
	profileClient profileproto.ProfileClient
	authClient    authproto.AuthClient
}

func NewProfileHandler(profileClient profileproto.ProfileClient, authClient authproto.AuthClient) *ProfileHandler {
	return &ProfileHandler{
		profileClient: profileClient,
		authClient:    authClient,
	}
}

func (h *ProfileHandler) GetProfile(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	vars := mux.Vars(request)

	profileId, _ := strconv.Atoi(vars["id"])

	authClient := h.authClient

	session, cookieErr := request.Cookie("session_id")

	var sessionValue string

	if session != nil {
		sessionValue = session.Value
	}

	var userIDCookie uint

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: sessionValue})

	if cookieErr == nil && user.IsAuth {
		userIDCookie = uint(user.ID)
	}

	profile, err := h.profileClient.GetProfile(ctx, &profileproto.ProfileIDRequest{ProfileId: uint64(profileId), UserId: uint64(userIDCookie)})
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(CleanProfileData(profile)))
}

func (h *ProfileHandler) SetProfileCity(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	authClient := h.authClient
	profileClient := h.profileClient

	session, _ := request.Cookie("session_id")

	var data models.City

	reqData, _ := io.ReadAll(request.Body)

	err := data.UnmarshalJSON(reqData)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	profile, err := profileClient.SetProfileCity(ctx, &profileproto.SetCityRequest{
		ID:          user.ID,
		CityID:      uint64(data.ID),
		CityName:    data.CityName,
		Translation: data.Translation,
	})
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(CleanProfileData(profile)))
}

func (h *ProfileHandler) SetProfilePhone(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	authClient := h.authClient
	profileClient := h.profileClient

	session, _ := request.Cookie("session_id")

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	var data models.SetProfilePhoneNec

	reqData, _ := io.ReadAll(request.Body)

	err := data.UnmarshalJSON(reqData)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	profile, err := profileClient.SetProfilePhone(ctx, &profileproto.SetPhoneRequest{
		ID:    user.ID,
		Phone: data.Phone,
	})
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(CleanProfileData(profile)))
}

func (h *ProfileHandler) EditProfile(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	authClient := h.authClient
	profileClient := h.profileClient

	session, _ := request.Cookie("session_id")

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	err := request.ParseMultipartForm(maxMemory)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	avatar := request.MultipartForm.File["avatar"]

	data := &profileproto.EditProfileRequest{
		ID:      user.ID,
		Name:    request.PostFormValue("name"),
		Surname: request.PostFormValue("surname"),
	}

	var (
		pl       *profileproto.ProfileData
		fullPath string
	)

	if len(avatar) != 0 {
		fullPath, err = utils.WriteFile(avatar[0], "avatars")
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while writing file of the image, err=%w",
				err))
			log.Println(err, responses.StatusInternalServerError)
			responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
				responses.ErrInternalServer))

			return
		}

		data.Avatar = fullPath
	}

	pl, err = profileClient.EditProfile(ctx, data)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.ErrInternalServer)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(CleanProfileData(pl)))
}

func (h *ProfileHandler) ChangeSubscription(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	profileClient := h.profileClient
	authClient := h.authClient

	var data models.ReceivedMerchantItem

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
	isAppended, _ := profileClient.AppendSubByIDs(ctx, &profileproto.UserIdMerchantIdRequest{UserId: uint32(user.ID),
		MerchantId: uint32(data.MerchantID)})

	responses.SendOkResponse(writer, responses.NewOkResponse(models.ProfileAppended{IsAppended: isAppended.IsAppended}))

	if isAppended.IsAppended {
		log.Println("User", user.ID, "has been added to subscribers of merchant", data.MerchantID)
		logging.LogHandlerInfo(logger, fmt.Sprintf("User %s has been added to the subscribers of merchant %s",
			fmt.Sprint(user.ID), fmt.Sprint(data.MerchantID)), responses.StatusOk)
	} else {
		log.Println("User", user.ID, "has been added to subscribers of merchant", data.MerchantID)
		logging.LogHandlerInfo(logger,
			fmt.Sprintf("User %s has been removed from the subscribers of merchant %s",
				fmt.Sprint(user.ID), fmt.Sprint(data.MerchantID)), responses.StatusOk)
	}
}
