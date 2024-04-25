package delivery

import (
	"log"
	"net/http"

	advdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/delivery"
	cityusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/usecases"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"go.uber.org/zap"
)

type CityHandler struct {
	storage    cityusecases.CityStorageInterface
	addrOrigin string
	schema     string
	logger     *zap.SugaredLogger
}

func NewCityHandler(storage cityusecases.CityStorageInterface, addrOrigin string, schema string, logger *zap.SugaredLogger) *CityHandler {
	return &CityHandler{
		storage:    storage,
		addrOrigin: addrOrigin,
		schema:     schema,
		logger:     logger,
	}
}

func (h *CityHandler) GetCityList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	city, err := h.storage.GetCityList(ctx)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, advdel.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewCityListOkResponse(city))
}
