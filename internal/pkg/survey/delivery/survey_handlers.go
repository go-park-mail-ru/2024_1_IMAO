package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	advdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/delivery"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	surveyusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/survey/usecases"
	userusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
)

type SurveyHandler struct {
	userStorage   userusecases.UsersStorageInterface
	surveyStorage surveyusecases.SurveyStorageInterface
}

func NewSurveyHandler(userStorage userusecases.UsersStorageInterface,
	surveyStorage surveyusecases.SurveyStorageInterface) *SurveyHandler {
	return &SurveyHandler{
		userStorage:   userStorage,
		surveyStorage: surveyStorage,
	}
}

func (surveyHandler *SurveyHandler) CreateAnswer(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var survey models.SurveyAnswersList

	err := json.NewDecoder(request.Body).Decode(&survey)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewErrResponse(responses.StatusInternalServerError,
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
	userStorage := surveyHandler.userStorage

	session, _ := request.Cookie("session_id")

	user, _ := userStorage.GetUserBySession(ctx, session.Value)

	isChecked, err := storage.GetResults(ctx, user.ID, uint(surveyID))
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, advdel.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewSurveyCheckResponse(isChecked))
}

func (surveyHandler *SurveyHandler) GetStatistics(writer http.ResponseWriter, request *http.Request) {
	responses.SendOkResponse(writer, NewSurveyOkResponse(nil))
}
