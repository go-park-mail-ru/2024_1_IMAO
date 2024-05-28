package delivery

import (
	"log"
	"net/http"

	cityusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/usecases"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"go.uber.org/zap"
)

type CityHandler struct {
	storage cityusecases.CityStorageInterface
}

func NewCityHandler(storage cityusecases.CityStorageInterface) *CityHandler {
	return &CityHandler{
		storage: storage,
	}
}

func (h *CityHandler) GetCityList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	city, err := h.storage.GetCityList(ctx)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(city))
}
