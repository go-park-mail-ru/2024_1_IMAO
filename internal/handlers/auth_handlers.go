package myhandlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
	"github.com/go-park-mail-ru/2024_1_IMAO/pkg"
)

const (
	sessionTime = 24 * time.Hour
)

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
// @Router /login [post]
func (authHandler *AuthHandler) Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := authHandler.List

	session, cookieErr := request.Cookie("session_id")

	if cookieErr == nil && usersList.SessionExists(session.Value) {
		log.Println("User already authorized", responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrAuthorized))

		return
	}

	user := storage.UnauthorizedUser{
		Email:          request.PostFormValue("email"),
		Password:       request.PostFormValue("password"),
		PasswordRepeat: "",
	}

	email := user.Email
	password := user.Password

	expectedUser, err := usersList.GetUserByEmail(email)
	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrWrongCredentials))

		return
	}

	if !pkg.CheckPassword(password, expectedUser.PasswordHash) {
		log.Println("Passwords do not match", responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrWrongCredentials))

		return
	}

	sessionID := usersList.AddSession(email)

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(sessionTime),
		HttpOnly: true,
	}
	http.SetCookie(writer, cookie)

	userData := responses.NewAuthOkResponse(*expectedUser, sessionID, true)
	responses.SendOkResponse(writer, userData)
	log.Println("You have been authorized with session ID:", sessionID)
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
// @Router /logout [post]
func (authHandler *AuthHandler) Logout(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := authHandler.List

	session, err := request.Cookie("session_id")
	if err != nil {
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	err = usersList.RemoveSession(session.Value)
	if err != nil {
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(writer, session)

	userData := responses.NewAuthOkResponse(storage.User{}, "", false)
	responses.SendOkResponse(writer, userData)

	log.Println("You have been logged out")
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
// @Router /signup [post]
func (authHandler *AuthHandler) Signup(writer http.ResponseWriter, request *http.Request) { //nolint:funlen
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := authHandler.List

	session, cookieErr := request.Cookie("session_id")

	if cookieErr == nil && usersList.SessionExists(session.Value) {
		log.Println("User already authorized", responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrAuthorized))

		return
	}

	newUser := storage.UnauthorizedUser{
		Email:          request.PostFormValue("email"),
		Password:       request.PostFormValue("password"),
		PasswordRepeat: request.PostFormValue("passwordRepeat"),
	}

	email := newUser.Email
	password := newUser.Password
	passwordRepeat := newUser.PasswordRepeat

	if !pkg.ValidateEmail(email) {
		log.Println("Bad email format", responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrWrongEmailFormat))

		return
	}

	if usersList.UserExists(email) {
		log.Println("User already exists", responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrUserAlreadyExists))

		return
	}

	if password != passwordRepeat {
		log.Println("Passwords do not match")
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrDifferentPasswords))

		return
	}

	if pkg.ValidatePassword(password) != "" {
		missed := pkg.ValidatePassword(password)
		log.Println("Bad password format:", missed, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			missed))

		return
	}

	user, _ := usersList.CreateUser(email, pkg.HashPassword(password))

	sessionID := usersList.AddSession(email)

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(sessionTime),
		HttpOnly: true,
	}
	http.SetCookie(writer, cookie)

	responses.SendOkResponse(writer, responses.NewAuthOkResponse(*user, sessionID, true))

	log.Println("You have been authorized with session ID:", sessionID)
}

// CheckAuth godoc
// @Summary Check user authentication
// @Description Verify if the user is authenticated by checking the session
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} responses.AuthOkResponse
// @Router /check_auth [get]
func (authHandler *AuthHandler) CheckAuth(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := authHandler.List

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		log.Println("User not authorized")
		responses.SendOkResponse(writer, responses.NewAuthOkResponse(storage.User{}, "", false))

		return
	}

	user, _ := usersList.GetUserBySession(session.Value)

	log.Println("User authorized")
	responses.SendOkResponse(writer, responses.NewAuthOkResponse(*user, session.Value, true))
}
