package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	advrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/repository"
	profrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/repository"
	authrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/repository"

	advdel "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/delivery"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
)

type ProfileHandler struct {
	AdvertsList *advrepo.AdvertsListWrapper
	ProfileList *profrepo.ProfileListWrapper
	UsersList   *authrepo.UsersListWrapper
}

func (h *ProfileHandler) GetProfile(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	ctx := request.Context()

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	p, err := h.ProfileList.GetProfileByUserID(ctx, uint(id))
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, advdel.NewAdvertsErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfileCity(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	ctx := request.Context()

	fmt.Println("aboba")

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(ctx, session.Value)

	var data models.City

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.SetProfileCity(ctx, user.ID, data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewProfileOkResponse(p))
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
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	var data models.SetProfileRatingNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.SetProfileRating(uint(userID), data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfilePhone(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	ctx := request.Context()

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(ctx, session.Value)

	var data models.SetProfilePhoneNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.SetProfilePhone(ctx, user.ID, data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewProfileOkResponse(p))
}

func (h *ProfileHandler) EditProfile(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	ctx := request.Context()

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(ctx, session.Value)

	var data models.EditProfileNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.EditProfile(user.ID, data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfileApproved(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	ctx := request.Context()

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(ctx, session.Value)

	p, err := h.ProfileList.SetProfileApproved(user.ID)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfile(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	ctx := request.Context()

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(ctx, session.Value)

	var data models.SetProfileNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.SetProfile(user.ID, data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfilePassword(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	ctx := request.Context()

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(ctx, session.Value)

	var data models.SetProfileNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.SetProfile(user.ID, data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewProfileOkResponse(p))
}

func (h *ProfileHandler) SetProfileEmail(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	ctx := request.Context()

	usersList := h.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(ctx, session.Value)

	var data models.SetProfileNec

	err = json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	p, err := h.ProfileList.SetProfile(user.ID, data)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, NewProfileOkResponse(p))
}

func (h *ProfileHandler) ProfileAdverts(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	filter, _ := strconv.Atoi(request.URL.Query().Get("filter"))

	var ads []*models.ReturningAdvert
	var err error

	switch models.AdvertsFilter(filter) {
	case models.FilterAll:
		ads, err = h.AdvertsList.GetAdvertsByUserIDFiltered(uint(id),
			func(ad *models.Advert) bool {
				return true
			})
	case models.FilterActive:
		ads, err = h.AdvertsList.GetAdvertsByUserIDFiltered(uint(id),
			func(ad *models.Advert) bool {
				return ad.Active
			})
	default:
		ads, err = h.AdvertsList.GetAdvertsByUserIDFiltered(uint(id),
			func(ad *models.Advert) bool {
				return !ad.Active
			})
	}

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewProfileErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest))

		return
	}

	responses.SendOkResponse(writer, advdel.NewAdvertsOkResponse(ads))
}
