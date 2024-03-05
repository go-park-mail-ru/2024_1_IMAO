package myhandlers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
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
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrTooManyAdverts), responses.StatusBadRequest)

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
