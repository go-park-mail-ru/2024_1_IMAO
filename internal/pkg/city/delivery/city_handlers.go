package delivery

import (
	"log"
	"net/http"

	advdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/delivery"
	cityrep "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/repository"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
)

type CityHandler struct {
	CityList *cityrep.CityListWrapper
}

func (h *CityHandler) GetCityList(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	ctx := request.Context()

	city, err := h.CityList.GetCityList(ctx)
	if err != nil {
		h.CityList.Logger.Error(err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, advdel.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewCityListOkResponse(city))
}
