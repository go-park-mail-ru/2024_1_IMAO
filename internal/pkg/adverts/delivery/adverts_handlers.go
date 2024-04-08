package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	advrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/repository"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	authresp "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"
)

const (
	defaultCity = "Moskva"
)

type AdvertsHandler struct {
	List *advrepo.AdvertsListWrapper
}

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
	} else if city == "" {
		city = defaultCity
	}

	var adsList []*models.ReturningAdInList
	var err error

	if category != "" {
		adsList, err = list.GetAdvertsByCategory(category, city, uint(startID), uint(count))
	} else {
		adsList, err = list.GetAdvertsByCity(city, uint(startID), uint(count))
	}

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewAdvertsOkResponse(adsList))
}

func (advertsHandler *AdvertsHandler) GetAdvert(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	city := vars["city"]
	category := vars["category"]
	id, _ := strconv.Atoi(vars["id"])

	list := advertsHandler.List

	ad, err := list.GetAdvert(uint(id), city, category)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewAdvertsOkResponse(ad))
}

func (advertsHandler *AdvertsHandler) CreateAdvert(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	list := advertsHandler.List
	var data models.ReceivedAdData

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, authresp.NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	adsList, err := list.CreateAdvert(data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewAdvertsOkResponse(adsList))
}

func (advertsHandler *AdvertsHandler) EditAdvert(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	list := advertsHandler.List
	var data models.ReceivedAdData

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, authresp.NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	ad, err := list.EditAdvert(data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewAdvertsOkResponse(ad))
}

func (advertsHandler *AdvertsHandler) DeleteAdvert(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	list := advertsHandler.List

	err := list.DeleteAdvert(uint(id))
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	adResponse := NewAdvertsOkResponse(nil)

	responses.SendOkResponse(writer, adResponse)
}

func (advertsHandler *AdvertsHandler) CloseAdvert(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	list := advertsHandler.List

	err := list.CloseAdvert(uint(id))
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	adResponse := NewAdvertsOkResponse(nil)

	responses.SendOkResponse(writer, adResponse)
}