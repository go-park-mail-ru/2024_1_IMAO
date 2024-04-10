package delivery

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	advrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/repository"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
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

	ctx := request.Context()

	vars := mux.Vars(request)
	city := vars["city"]
	category := vars["category"]

	list := advertsHandler.List

	count, errCount := strconv.Atoi(request.URL.Query().Get("count"))
	startID, errstartID := strconv.Atoi(request.URL.Query().Get("startId"))
	userID, errUser := strconv.Atoi(request.URL.Query().Get("userId"))
	deleted, errdeleted := strconv.Atoi(request.URL.Query().Get("deleted"))

	fmt.Println("errCount", errCount)
	fmt.Println("errstartID", errstartID)
	fmt.Println("errUser", errUser)
	fmt.Println("errdeleted", errdeleted)

	if city == "" && request.URL.Query().Get("city") != "" {
		city = request.URL.Query().Get("city")
	} else if city == "" {
		city = defaultCity
	}

	var adsList []*models.ReturningAdInList
	var err error

	if category != "" {
		adsList, err = list.GetAdvertsByCategory(category, city, uint(startID), uint(count))
	} else if errCount == nil && errstartID == nil {
		adsList, err = list.GetAdvertsByCity(city, uint(startID), uint(count))
	} else if errUser == nil && errdeleted == nil {
		adsList, err = list.GetAdvertsForUserWhereStatusIs(ctx, uint(userID), uint(deleted))
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

func (advertsHandler *AdvertsHandler) GetAdvertByID(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	list := advertsHandler.List

	ad, err := list.GetAdvertByOnlyByID(uint(id))
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

	err := request.ParseMultipartForm(0)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	list := advertsHandler.List

	isUsed := true
	if request.PostFormValue("condition") == "1" {
		isUsed = false
	}
	price, _ := strconv.Atoi(request.PostFormValue("price"))

	data := models.ReceivedAdData{
		UserID:      1,
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

	err := request.ParseMultipartForm(0)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAdvertsErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	list := advertsHandler.List

	isUsed := true
	if request.PostFormValue("condition") == "1" {
		isUsed = false
	}
	price, _ := strconv.Atoi(request.PostFormValue("price"))

	data := models.ReceivedAdData{
		UserID:      1,
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
