package myhandlers

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
	"github.com/go-park-mail-ru/2024_1_IMAO/pkg"
	"log"
	"net/http"
	"time"
)

const (
	sessionTime = 24 * time.Hour
)

func (authHandler *AuthHandler) Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := authHandler.List

	session, err := request.Cookie("session_id")

	if err == nil && usersList.SessionExists(session.Value) {
		log.Println("User already authorized", responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrAuthorized), responses.StatusBadRequest)

		return
	}

	user := storage.UnauthorizedUser{
		Email:    request.PostFormValue("email"),
		Password: request.PostFormValue("password"),
	}

	email := user.Email
	password := user.Password
	expectedUser, err := usersList.GetUserByEmail(email)

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrWrongCredentials), responses.StatusBadRequest)

		return
	}

	if !pkg.CheckPassword(password, expectedUser.PasswordHash) {
		log.Println("Passwords do not match", responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrWrongCredentials), responses.StatusBadRequest)

		return
	}

	sessionID := usersList.AddSession(email)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(sessionTime),
	}

	userData := responses.NewAuthOkResponse(*expectedUser, sessionID, true)
	responses.SendOkResponse(writer, userData)
	http.SetCookie(writer, cookie)

	log.Println("You have been authorized with session ID:", sessionID)
}

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
			responses.ErrUnauthorized), responses.StatusUnauthorized)

		return
	}

	err = usersList.RemoveSession(session.Value)

	if err != nil {
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized), responses.StatusUnauthorized)

		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(writer, session)

	log.Println("You have been logged out")
}

func (authHandler *AuthHandler) Signup(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	usersList := authHandler.List

	session, err := request.Cookie("session_id")

	if err == nil || usersList.SessionExists(session.Value) {
		log.Println("User already authorized", responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrAuthorized), responses.StatusBadRequest)

		return
	}

	newUser := storage.UnauthorizedUser{
		Email:          request.PostFormValue("email"),
		Password:       request.PostFormValue("password"),
		PasswordRepeat: request.PostFormValue("passwordRepeat"),
	}

	if err != nil {
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrBadRequest), responses.StatusBadRequest)

		return
	}

	email := newUser.Email
	password := newUser.Password
	passwordRepeat := newUser.PasswordRepeat

	if password != passwordRepeat {
		log.Println("Passwords do not match")
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrDifferentPasswords), responses.StatusBadRequest)

		return
	}

	if pkg.ValidatePassword(password) != "" {
		missed := pkg.ValidatePassword(password)
		log.Println("Bad password format:", missed, responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			missed), responses.StatusBadRequest)

		return
	}

	if !pkg.ValidateEmail(email) {
		log.Println("Bad email format", responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrWrongEmailFormat), responses.StatusBadRequest)

		return
	}

	user, err := usersList.CreateUser(email, pkg.HashPassword(password))

	if err != nil {
		log.Println("User already exists", responses.StatusBadRequest)
		responses.SendErrResponse(writer, responses.NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrUserAlreadyExists), responses.StatusBadRequest)

		return
	}

	sessionID := usersList.AddSession(email)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(sessionTime),
	}
	http.SetCookie(writer, cookie)

	responses.SendOkResponse(writer, responses.NewAuthOkResponse(*user, sessionID, true))

	log.Println("You have been authorized with session ID:", sessionID)
}

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
