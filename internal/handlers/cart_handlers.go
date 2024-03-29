package myhandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
)

// const (
// 	advertsPerPage = 30
// 	defaultCity    = "Moskva"
// )

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
func (cartHandler *CartHandler) GetCartList(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	// vars := mux.Vars(request)
	// city := vars["city"]
	// category := vars["category"]

	list := cartHandler.ListCart

	userId, err := strconv.Atoi(request.URL.Query().Get("userId"))
	// startID, _ := strconv.Atoi(request.URL.Query().Get("startId"))

	// if city == "" && request.URL.Query().Get("city") != "" {
	// 	city = request.URL.Query().Get("city")
	// } else {
	// 	city = defaultCity
	// }

	var adsList []*storage.Advert
	//var err error

	if err == nil {
		adsList, err = list.GetCartByUserID(uint(userId), cartHandler.ListUsers, cartHandler.ListAdverts)
	}

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewCartErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewCartOkResponse(adsList))
}

func (cartHandler *CartHandler) AppendCart(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	list := cartHandler.ListCart
	var data storage.ReceivedCartItem

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewCartErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	_, err = list.AppendAdvByIDs(data, cartHandler.ListUsers, cartHandler.ListAdverts)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewCartErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	adsList, _ := cartHandler.ListCart.GetCartByUserID(data.UserID, cartHandler.ListUsers, cartHandler.ListAdverts)

	responses.SendOkResponse(writer, responses.NewCartOkResponse(adsList))
}

// func (cartHandler *AdvertsHandler) EditAdvert(writer http.ResponseWriter, request *http.Request) {

// }

// func (cartHandler *AdvertsHandler) DeleteAdvert(writer http.ResponseWriter, request *http.Request) {

// }

// func (cartHandler *AdvertsHandler) CloseAdvert(writer http.ResponseWriter, request *http.Request) {

// }
