package myhandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
	"github.com/gorilla/mux"
)

func (h *ProfileHandler) GetProfile(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	p, err := h.ProfileList.GetProfileByUserID(uint(id))
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfileCity(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(session.Value)

	var data storage.SetProfileCityNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.SetProfileCity(user.ID, data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfileRating(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	userID, _ := strconv.Atoi(vars["id"])

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	var data storage.SetProfileRatingNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.SetProfileRating(uint(userID), data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfilePhone(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(session.Value)

	var data storage.SetProfilePhoneNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.SetProfilePhone(user.ID, data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewProfileOkResponse(p))
}

func (h *ProfileHandler) EditProfile(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(session.Value)

	var data storage.EditProfileNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.EditProfile(user.ID, data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfileApproved(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(session.Value)

	p, err := h.ProfileList.SetProfileApproved(user.ID)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfile(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(session.Value)

	var data storage.SetProfileNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.SetProfile(user.ID, data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfilePassword(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(session.Value)

	var data storage.SetProfileNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.SetProfile(user.ID, data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfileEmail(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(session.Value)

	var data storage.SetProfileNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.SetProfile(user.ID, data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewProfileOkResponse(p))
}

func (h *ProfileHandler) ProfileAdverts(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	filter, _ := strconv.Atoi(request.URL.Query().Get("filter"))
	
	var ads []*storage.ReturningAdvert
	var err error
	
	switch storage.AdvertsFilter(filter) {
	case storage.FilterAll:
		ads, err = h.AdvertsList.GetAdvertsByUserID(uint(id))
	case storage.FilterActive:
		ads, err = h.AdvertsList.GetToggledAdvertsByUserID(uint(id), true)
	default:
		ads, err = h.AdvertsList.GetToggledAdvertsByUserID(uint(id), false)
	}
	
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, responses.NewAdvertsOkResponse(ads))
}
