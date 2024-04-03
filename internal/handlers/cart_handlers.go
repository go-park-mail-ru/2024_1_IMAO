package myhandlers

import (
	"encoding/json"
	"log"
	"net/http"

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

	list := cartHandler.ListCart
	usersList := cartHandler.ListUsers

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendOkResponse(writer, responses.NewAuthOkResponse(storage.User{}, "", false))

		return
	}

	user, _ := usersList.GetUserBySession(session.Value)

	var adsList []*storage.Advert

	adsList, err = list.GetCartByUserID(uint(user.ID), cartHandler.ListUsers, cartHandler.ListAdverts)

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewCartErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewCartOkResponse(adsList))
}

func (cartHandler *CartHandler) ChangeCart(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	list := cartHandler.ListCart
	usersList := cartHandler.ListUsers
	var data storage.ReceivedCartItem

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewCartErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendOkResponse(writer, responses.NewAuthOkResponse(storage.User{}, "", false))

		return
	}

	user, _ := usersList.GetUserBySession(session.Value)

	isAppended := list.AppendAdvByIDs(user.ID, data.AdvertID, cartHandler.ListUsers, cartHandler.ListAdverts)

	responses.SendOkResponse(writer, responses.NewCartChangeResponse(isAppended))
}

func (cartHandler *CartHandler) DeleteFromCart(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	list := cartHandler.ListCart
	usersList := cartHandler.ListUsers
	var data storage.ReceivedCartItems

	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewCartErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendOkResponse(writer, responses.NewAuthOkResponse(storage.User{}, "", false))

		return
	}

	user, _ := usersList.GetUserBySession(session.Value)

	for _, item := range data.AdvertIDs {
		err = list.DeleteAdvByIDs(user.ID, item, cartHandler.ListUsers, cartHandler.ListAdverts)

		if err != nil {
			log.Println(err, responses.StatusBadRequest)
			responses.SendErrResponse(writer, responses.NewCartErrResponse(responses.StatusBadRequest,
				responses.ErrBadRequest))

			return
		}
	}

	responses.SendOkResponse(writer, responses.NewCartChangeResponse(false))
}
