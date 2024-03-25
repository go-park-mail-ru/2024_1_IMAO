package myhandlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
)

const (
	advertsPerPage = 30
	defaultCity    = "Moskva"
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

	vars := mux.Vars(request)
	city := vars["city"]

	list := advertsHandler.List
	var data storage.GettingAdsData

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	if city == "" && data.City != "" {
		city = data.City
	} else if city == "" {
		city = defaultCity
	}

	count := data.Count
	startID := data.StartID

	adsList, err := list.GetAdvertsByCity(city, count, startID)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewAdvertsOkResponse(adsList))
}

func (advertsHandler *AdvertsHandler) GetCategoryAds(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	city := vars["city"]
	category := vars["category"]

	list := advertsHandler.List
	var data storage.GettingAdsData

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	count := data.Count
	startID := data.StartID

	adsList, err := list.GetAdvertsByCategory(city, category, count, startID)
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
