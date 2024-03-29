package myhandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
)

const (
	advertsPerPage = 30
	defaultCity    = "Moskva"
)

// GetAdsList godoc
// @Summary Retrieve a list of adverts
// @Description Get a paginated list of adverts
// @Tags adverts
// @Accept json
// @Produce json
// @Success 200 {object} responses.AdvertsOkResponse
// @Failure 400 {object} responses.AdvertsErrResponse "Too many adverts specified"
// @Failure 405 {object} responses.AdvertsErrResponse "Method not allowed"
// @Failure 500 {object} responses.AdvertsErrResponse "Internal server error"
// @Router /api/adverts/list [get]
func (advertsHandler *AdvertsHandler) GetAdsList(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	city := vars["city"]
	category := vars["category"]

	list := advertsHandler.List

	count, _ := strconv.Atoi(request.URL.Query().Get("count"))
	startID, _ := strconv.Atoi(request.URL.Query().Get("startId"))

	if city == "" && request.URL.Query().Get("city") != "" {
		city = request.URL.Query().Get("city")
	} else {
		city = defaultCity
	}

	var adsList []*storage.ReturningAdvert
	var err error

	if category != "" {
		adsList, err = list.GetAdvertsByCategory(city, category, uint(count), uint(startID))
	} else {
		adsList, err = list.GetAdvertsByCity(city, uint(count), uint(startID))
	}

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewAdvertsOkResponse(adsList))
}

func (advertsHandler *AdvertsHandler) CreateAdvert(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	list := advertsHandler.List
	var data storage.ReceivedAdData

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	adsList, err := list.CreateAdvert(data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewAdvertsOkResponse(adsList))
}

func (advertsHandler *AdvertsHandler) EditAdvert(writer http.ResponseWriter, request *http.Request) {

}

func (advertsHandler *AdvertsHandler) DeleteAdvert(writer http.ResponseWriter, request *http.Request) {

}

func (advertsHandler *AdvertsHandler) CloseAdvert(writer http.ResponseWriter, request *http.Request) {

}
