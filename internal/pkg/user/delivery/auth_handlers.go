package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	profrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/repository"
	authrepo "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/repository"

	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	utils "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
)

const (
	sessionTime = 24 * time.Hour
)

type AuthHandler struct {
	UsersList   *authrepo.UsersListWrapper
	ProfileList *profrepo.ProfileListWrapper
}

// Login godoc
// @Summary User login
// @Description Authenticate user and create a new session
// @Tags auth
// @Accept json
// @Produce json
// @Param email formData string true "User email"
// @Param password formData string true "User password"
// @Success 200 {object} responses.AuthOkResponse
// @Failure 400 {object} responses.AuthErrResponse "Bad request"
// @Failure 401 {object} responses.AuthErrResponse "Unauthorized"
// @Failure 405 {object} responses.AuthErrResponse "Method not allowed"
// @Failure 500 {object} responses.AuthErrResponse "Internal server error"
// @Router /api/auth/login [post]
func (authHandler *AuthHandler) Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := authHandler.UsersList

	session, cookieErr := request.Cookie("session_id")

	if cookieErr == nil && usersList.SessionExists(session.Value) {
		log.Println(responses.ErrAuthorized, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrAuthorized))

		return
	}

	var user models.UnauthorizedUser

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	email := user.Email
	password := user.Password

	expectedUser, err := usersList.GetUserByEmail(email)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrWrongCredentials))

		return
	}

	if !utils.CheckPassword(password, expectedUser.PasswordHash) {
		log.Println("Passwords do not match", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrWrongCredentials))

		return
	}

	sessionID := usersList.AddSession(expectedUser.ID)

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(sessionTime),
		HttpOnly: true,
		SameSite: 4,
		Secure:   true,
	}
	http.SetCookie(writer, cookie)

	userData := NewAuthOkResponse(*expectedUser, sessionID, true)
	responses.SendOkResponse(writer, userData)
	log.Println("User", user.Email, "have been authorized with session ID:", sessionID)

}

// Logout godoc
// @Summary User logout
// @Description Invalidate the user session and log the user out
// @Tags auth
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {object} responses.AuthErrResponse "Bad request"
// @Failure 401 {object} responses.AuthErrResponse "User not authorized"
// @Failure 405 {object} responses.AuthErrResponse "Method not allowed"
// @Failure 500 {object} responses.AuthErrResponse "Internal server error"
// @Router /api/auth/logout [post]
func (authHandler *AuthHandler) Logout(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := authHandler.UsersList

	session, err := request.Cookie("session_id")
	if err != nil {
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	err = usersList.RemoveSession(session.Value)
	if err != nil {
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(writer, session)

	userData := NewAuthOkResponse(models.User{}, "", false)
	responses.SendOkResponse(writer, userData)

	log.Println("User have been logged out")
}

// Signup godoc
// @Summary User signup
// @Description Register a new user and create a new session
// @Tags auth
// @Accept json
// @Produce json
// @Param email formData string true "User email"
// @Param password formData string true "User password"
// @Param passwordRepeat formData string true "Password confirmation"
// @Success 201 {object} responses.AuthOkResponse
// @Failure 400 {object} responses.AuthErrResponse
// @Router /api/auth/signup [post]
func (authHandler *AuthHandler) Signup(writer http.ResponseWriter, request *http.Request) { //nolint:funlen
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	ctx := request.Context()
	usersList := authHandler.UsersList
	profileList := authHandler.ProfileList

	session, cookieErr := request.Cookie("session_id")

	if cookieErr == nil && usersList.SessionExists(session.Value) {
		log.Println("User already authorized", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrAuthorized))

		return
	}

	var newUser models.UnauthorizedUser

	err := json.NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	email := newUser.Email
	password := newUser.Password
	passwordRepeat := newUser.PasswordRepeat

	if usersList.UserExists(ctx, email) {
		log.Println("User already exists", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrUserAlreadyExists))

		return
	}

	errors := utils.Validate(email, password)
	if errors != nil {
		log.Println("Bad user data", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewValidationErrResponse(responses.StatusBadRequest,
			errors))

		return
	}

	if password != passwordRepeat {
		log.Println("Passwords do not match")
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrDifferentPasswords))

		return
	}

	user, _ := usersList.CreateUser(ctx, email, utils.HashPassword(password))
	profileList.CreateProfile(user.ID)

	sessionID := usersList.AddSession(user.ID)

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(sessionTime),
		HttpOnly: true,
		SameSite: 4,
		Secure:   true,
	}
	http.SetCookie(writer, cookie)

	responses.SendOkResponse(writer, NewAuthOkResponse(*user, sessionID, true))

	log.Println("User", user.Email, "have been authorized with session ID:", sessionID)
}

// CheckAuth godoc
// @Summary Check user authentication
// @Description Verify if the user is authenticated by checking the session
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} responses.AuthOkResponse
// @Router /api/auth/check_auth [get]
func (authHandler *AuthHandler) CheckAuth(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := authHandler.UsersList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendOkResponse(writer, NewAuthOkResponse(models.User{}, "", false))

		return
	}

	user, _ := usersList.GetUserBySession(session.Value)

	log.Println("User authorized")
	responses.SendOkResponse(writer, NewAuthOkResponse(*user, session.Value, true))
}

func (authHandler *AuthHandler) EditUser(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := authHandler.UsersList

	session, err := request.Cookie("session_id")

	if err != nil {
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(session.Value)

	var newUser models.UnauthorizedUser

	err = json.NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	email := newUser.Email
	password := newUser.Password
	passwordRepeat := newUser.PasswordRepeat

	errors := utils.Validate(email, password)
	if errors != nil {
		log.Println("Bad user data", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewValidationErrResponse(responses.StatusBadRequest,
			errors))

		return
	}

	if password != passwordRepeat {
		log.Println("Passwords do not match")
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrDifferentPasswords))

		return
	}

	usersList.EditUser(user.ID, email, password)

	responses.SendOkResponse(writer, NewAuthOkResponse(*user, session.Value, true))

	log.Println("User", user.Email, "successfully changed his authorization data.")
}