package myhandlers

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
)

const (
	advertsPerPage = 30
)

// Root godoc
// @Summary Retrieve a list of adverts
// @Description Get a paginated list of adverts
// @Tags adverts
// @Accept json
// @Produce json
// @Success 200 {object} responses.AdvertsOkResponse
// @Failure 400 {object} responses.AdvertsErrResponse "Too many adverts specified"
// @Failure 405 {object} responses.AdvertsErrResponse "Method not allowed"
// @Failure 500 {object} responses.AdvertsErrResponse "Internal server error"
// @Router /api/adverts [get]
func (advertsHandler *AdvertsHandler) Root(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	list := advertsHandler.List

	adsList, err := list.GetSeveralAdverts(advertsPerPage)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrTooManyAdverts))

		return
	}

	responses.SendOkResponse(writer, responses.NewAdvertsOkResponse(adsList))
}

func PanicHandler(writer http.ResponseWriter) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Panic happened:", err)
			http.Error(writer, responses.ErrInternalServer, responses.StatusInternalServerError)
		}
	}()
}
