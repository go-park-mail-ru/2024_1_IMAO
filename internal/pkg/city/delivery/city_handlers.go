package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"log"
	"net/http"
	"os"

	cityusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/usecases"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"go.uber.org/zap"
)

const (
	geocoderURL = "http://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address"
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

func (h *CityHandler) GetCityName(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	apiKey := os.Getenv("GEOAPI_KEY")

	req, err := http.NewRequest("POST", geocoderURL, request.Body)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))

		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apiKey))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))

		return
	}
	defer resp.Body.Close()

	var data models.GeocoderResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))

		return
	}

	var city string

	if len(data.Suggestions) != 0 {
		city = data.Suggestions[0].Data.City
	}

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(models.CityResponse{CityName: city}))
}
