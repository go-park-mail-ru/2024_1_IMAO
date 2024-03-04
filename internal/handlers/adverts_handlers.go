package myhandlers

import (
	"2024_1_IMAO/internal/responses"
	"log"
	"net/http"
)

func (advertsHandler *AdvertsHandler) Root(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	list := advertsHandler.List

	adsList, err := list.GetSeveralAdverts(50)

	if err != nil {
		log.Println(err)
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrTooManyAdverts), responses.StatusBadRequest)

		return
	}

	responses.SendOkResponse(writer, responses.NewAdvertsOkResponse(adsList))
}
