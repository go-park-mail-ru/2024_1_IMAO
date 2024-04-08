package myhandlers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
)

const (
	defaultCity = "Moskva"
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
	} else if city == "" {
		city = defaultCity
	}

	var adsList []*storage.ReturningAdInList
	var err error

	if category != "" {
		adsList, err = list.GetAdvertsByCategory(category, city, uint(startID), uint(count))
	} else {
		adsList, err = list.GetAdvertsByCity(city, uint(startID), uint(count))
	}

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewAdvertsOkResponse(adsList))
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
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewAdvertsOkResponse(ad))
}

func (advertsHandler *AdvertsHandler) GetAdvertByID(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	list := advertsHandler.List

	ad, err := list.GetAdvertByID(uint(id))
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewAdvertsOkResponse(ad))
}

func (advertsHandler *AdvertsHandler) CreateAdvert(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	err := request.ParseMultipartForm(0)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	list := advertsHandler.List

	isUsed := true
	if request.PostFormValue("condition") == "1" {
		isUsed = false
	}
	price, _ := strconv.Atoi(request.PostFormValue("price"))
	userID, _ := strconv.Atoi(request.PostFormValue("userId"))

	data := storage.ReceivedAdData{
		UserID:      uint(userID),
		City:        request.PostFormValue("city"),
		Category:    request.PostFormValue("category"),
		Title:       request.PostFormValue("title"),
		Description: request.PostFormValue("description"),
		Price:       uint(price),
		IsUsed:      isUsed,
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
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	err := request.ParseMultipartForm(0)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	list := advertsHandler.List

	isUsed := true
	if request.PostFormValue("condition") == "1" {
		isUsed = false
	}
	price, _ := strconv.Atoi(request.PostFormValue("price"))
	id, _ := strconv.Atoi(request.PostFormValue("id"))
	userID, _ := strconv.Atoi(request.PostFormValue("userId"))

	data := storage.ReceivedAdData{
		ID:          uint(id),
		UserID:      uint(userID),
		City:        request.PostFormValue("city"),
		Category:    request.PostFormValue("category"),
		Title:       request.PostFormValue("title"),
		Description: request.PostFormValue("description"),
		Price:       uint(price),
		IsUsed:      isUsed,
	}

	ad, err := list.EditAdvert(data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewAdvertsOkResponse(ad))
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
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	adResponse := responses.NewAdvertsOkResponse(nil)

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
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	adResponse := responses.NewAdvertsOkResponse(nil)

	responses.SendOkResponse(writer, adResponse)
}
