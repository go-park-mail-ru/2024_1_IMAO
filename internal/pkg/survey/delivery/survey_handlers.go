package delivery

import (
	"encoding/json"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	surveyusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/survey/usecases"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type SurveyHandler struct {
	authClient    authproto.AuthClient
	surveyStorage surveyusecases.SurveyStorageInterface
}

func NewSurveyHandler(authClient authproto.AuthClient,
	surveyStorage surveyusecases.SurveyStorageInterface) *SurveyHandler {
	return &SurveyHandler{
		authClient:    authClient,
		surveyStorage: surveyStorage,
	}
}

func (surveyHandler *SurveyHandler) CreateAnswer(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))
	storage := surveyHandler.surveyStorage

	var survey models.SurveyAnswersList

	err := json.NewDecoder(request.Body).Decode(&survey)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	err = storage.SaveSurveyResults(ctx, survey)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	responses.SendOkResponse(writer, NewSurveyOkResponse(survey))
}

func (surveyHandler *SurveyHandler) CheckIfAnswered(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	vars := mux.Vars(request)
	surveyID, _ := strconv.Atoi(vars["survey_id"])

	storage := surveyHandler.surveyStorage
	authClient := surveyHandler.authClient

	session, _ := request.Cookie("session_id")

	user, _ := authClient.GetCurrentUser(ctx, &authproto.SessionData{SessionID: session.Value})

	isChecked, err := storage.GetResults(ctx, uint(user.ID), uint(surveyID))
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewSurveyCheckResponse(isChecked))
}

func (surveyHandler *SurveyHandler) GetStatistics(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	storage := surveyHandler.surveyStorage

	surveyResults, err := storage.GetStatics(ctx)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewSurveyOkResponse(surveyResults))
}
