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

	ctx := request.Context()

	usersList := authHandler.UsersList

	session, cookieErr := request.Cookie("session_id")

	if cookieErr == nil && usersList.SessionExists(session.Value) {
		authHandler.UsersList.Logger.Info(responses.ErrAuthorized, responses.StatusBadRequest)
		log.Println(responses.ErrAuthorized, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrAuthorized))

		return
	}

	var user models.UnauthorizedUser

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		authHandler.UsersList.Logger.Error(err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	email := user.Email
	password := user.Password

	expectedUser, err := usersList.GetUserByEmail(ctx, email)
	if err != nil {
		authHandler.UsersList.Logger.Error(err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrWrongCredentials))

		return
	}

	if !utils.CheckPassword(password, expectedUser.PasswordHash) {
		authHandler.UsersList.Logger.Info("Passwords do not match", responses.StatusBadRequest)
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
	}
	http.SetCookie(writer, cookie)

	userData := NewAuthOkResponse(*expectedUser, sessionID, true)
	responses.SendOkResponse(writer, userData)
	authHandler.UsersList.Logger.Info("User", user.Email, "have been authorized with session ID:", sessionID)
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
		authHandler.UsersList.Logger.Error(err, responses.StatusUnauthorized)
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	err = usersList.RemoveSession(session.Value)
	if err != nil {
		authHandler.UsersList.Logger.Error(err, responses.StatusUnauthorized)
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(writer, session)

	userData := NewAuthOkResponse(models.User{}, "", false)
	responses.SendOkResponse(writer, userData)

	authHandler.UsersList.Logger.Info("User have been logged out")
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
		authHandler.UsersList.Logger.Info("User already authorized", responses.StatusBadRequest)
		log.Println("User already authorized", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrAuthorized))

		return
	}

	var newUser models.UnauthorizedUser

	err := json.NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		authHandler.UsersList.Logger.Error(err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	email := newUser.Email
	password := newUser.Password
	passwordRepeat := newUser.PasswordRepeat

	if usersList.UserExists(ctx, email) {
		authHandler.UsersList.Logger.Info("User already exists", responses.StatusBadRequest)
		log.Println("User already exists", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrUserAlreadyExists))

		return
	}

	errors := utils.Validate(email, password)
	if errors != nil {
		authHandler.UsersList.Logger.Error("Bad user data", responses.StatusBadRequest)
		log.Println("Bad user data", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewValidationErrResponse(responses.StatusBadRequest,
			errors))

		return
	}

	if password != passwordRepeat {
		authHandler.UsersList.Logger.Info("Passwords do not match")
		log.Println("Passwords do not match")
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrDifferentPasswords))

		return
	}

	user, _ := usersList.CreateUser(ctx, email, utils.HashPassword(password))
	profileList.CreateProfile(ctx, user.ID)

	sessionID := usersList.AddSession(user.ID)

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(sessionTime),
		HttpOnly: true,
	}
	http.SetCookie(writer, cookie)

	responses.SendOkResponse(writer, NewAuthOkResponse(*user, sessionID, true))

	authHandler.UsersList.Logger.Info("User", user.Email, "have been authorized with session ID:", sessionID)
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

	ctx := request.Context()

	usersList := authHandler.UsersList
	profileList := authHandler.ProfileList

	session, err := request.Cookie("session_id")

	if err != nil || !usersList.SessionExists(session.Value) {
		authHandler.UsersList.Logger.Error("User not authorized")
		log.Println("User not authorized")
		responses.SendOkResponse(writer, NewAuthOkResponse(models.User{}, "", false))

		return
	}

	user, _ := usersList.GetUserBySession(ctx, session.Value)
	profile, _ := profileList.GetProfileByUserID(ctx, user.ID)

	authHandler.UsersList.Logger.Info("User authorized")
	log.Println("User authorized")
	responses.SendOkResponse(writer, NewAuthOkResponse(*user, profile.AvatarIMG, true))
}

func (authHandler *AuthHandler) EditUserEmail(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	ctx := request.Context()

	usersList := authHandler.UsersList

	session, err := request.Cookie("session_id")

	if err != nil {
		authHandler.UsersList.Logger.Error(err, responses.StatusUnauthorized)
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := usersList.GetUserBySession(ctx, session.Value)

	var newUser models.UnauthorizedUser

	err = json.NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		authHandler.UsersList.Logger.Error(err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	email := newUser.Email

	errors := utils.ValidateEmail(email)
	if errors != nil {
		authHandler.UsersList.Logger.Error("Email is not valid", responses.StatusBadRequest)
		log.Println("Email is not valid", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			"Email is not valid"))

		return
	}

	user, err = usersList.EditUserEmail(ctx, user.ID, email)

	if err != nil {
		authHandler.UsersList.Logger.Error("This email is already in use", responses.StatusBadRequest)
		log.Println("This email is already in use", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			"This email is already in use"))

		return
	}

	responses.SendOkResponse(writer, NewAuthOkResponse(*user, session.Value, true))

	authHandler.UsersList.Logger.Info("User", user.Email, "successfully changed his authorization data.")
	log.Println("User", user.Email, "successfully changed his authorization data.")
}
