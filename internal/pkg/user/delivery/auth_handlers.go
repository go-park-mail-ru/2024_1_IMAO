package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	csrf "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/csrf"
	profusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/usecases"
	userusecases "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"

	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	utils "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"

	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
)

const (
	sessionTime = 24 * time.Hour
)

type AuthHandler struct {
	storage        userusecases.UsersStorageInterface
	profileStorage profusecases.ProfileStorageInterface
	addrOrigin     string
	schema         string
}

func NewAuthHandler(storage userusecases.UsersStorageInterface, profileStorage profusecases.ProfileStorageInterface,
	addrOrigin string, schema string) *AuthHandler {
	return &AuthHandler{
		storage:        storage,
		profileStorage: profileStorage,
		addrOrigin:     addrOrigin,
		schema:         schema,
	}
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
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	storage := authHandler.storage

	session, cookieErr := request.Cookie("session_id")

	if cookieErr == nil && storage.SessionExists(session.Value) {
		logging.LogHandlerInfo(logger, responses.ErrAuthorized, responses.StatusBadRequest)
		log.Println(responses.ErrAuthorized, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrAuthorized))

		return
	}

	var user models.UnauthorizedUser

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	email := user.Email
	password := user.Password

	expectedUser, err := storage.GetUserByEmail(ctx, email)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusBadRequest)
		log.Println(err, responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrWrongCredentials))

		return
	}

	if !utils.CheckPassword(password, expectedUser.PasswordHash) {
		logging.LogHandlerInfo(logger, responses.ErrDifferentPasswords, responses.StatusBadRequest)
		log.Println("Passwords do not match", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrWrongCredentials))

		return
	}

	sessionID := storage.AddSession(expectedUser.ID)

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

	logging.LogHandlerInfo(logger, fmt.Sprintf("User %s have been authorized with session ID: %s ", user.Email, sessionID), responses.StatusOk)
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
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	storage := authHandler.storage

	session, err := request.Cookie("session_id")
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusUnauthorized)
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	err = storage.RemoveSession(ctx, session.Value)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusUnauthorized)
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(writer, session)

	userData := NewAuthOkResponse(models.User{}, "", false)
	responses.SendOkResponse(writer, userData)

	// ПО-ХОРОШЕМУ НУЖНО ПЕРЕПИСАТЬ ХЭНДЛЕР, ЧТОБЫ В ЛОГЕ МОЖНО БЫЛО ВЫВОДИТЬ КАКОЙ ИМЕННО ПОЛЬЗОВАТЕЛЬ РАЗЛОГИНИЛСЯ
	logging.LogHandlerInfo(logger, "User have been logged out", responses.StatusOk)
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
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	storage := authHandler.storage
	profileStorage := authHandler.profileStorage

	session, cookieErr := request.Cookie("session_id")

	if cookieErr == nil && storage.SessionExists(session.Value) {
		logging.LogHandlerInfo(logger, responses.ErrAuthorized, responses.StatusBadRequest)
		log.Println("User already authorized", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrAuthorized))

		return
	}

	var newUser models.UnauthorizedUser

	err := json.NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	email := newUser.Email
	password := newUser.Password
	passwordRepeat := newUser.PasswordRepeat

	if storage.UserExists(ctx, email) {
		logging.LogHandlerInfo(logger, responses.ErrUserAlreadyExists, responses.StatusBadRequest)
		log.Println("User already exists", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrUserAlreadyExists))

		return
	}

	errors := utils.Validate(email, password)
	if errors != nil {
		logging.LogHandlerInfo(logger, responses.ErrNotValidData, responses.StatusBadRequest)
		log.Println("Bad user data", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewValidationErrResponse(responses.StatusBadRequest,
			errors))

		return
	}

	if password != passwordRepeat {
		logging.LogHandlerInfo(logger, responses.ErrDifferentPasswords, responses.StatusBadRequest)
		log.Println("Passwords do not match")
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			responses.ErrDifferentPasswords))

		return
	}

	user, _ := storage.CreateUser(ctx, email, utils.HashPassword(password))

	profileStorage.CreateProfile(ctx, user.ID)

	sessionID := storage.AddSession(user.ID)
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

	logging.LogHandlerInfo(logger, fmt.Sprintf("User %s have been authorized with session ID: %s ", user.Email, sessionID), responses.StatusOk)
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
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	storage := authHandler.storage
	profileStorage := authHandler.profileStorage

	session, err := request.Cookie("session_id")

	if err != nil || !storage.SessionExists(session.Value) {
		if err == nil {
			err = errors.New("no such cookie in userStorage")
		}
		logging.LogHandlerError(logger, err, responses.StatusUnauthorized)
		log.Println("User not authorized")
		responses.SendOkResponse(writer, NewAuthOkResponse(models.User{}, "", false))

		return
	}

	user, _ := storage.GetUserBySession(ctx, session.Value)
	profile, _ := profileStorage.GetProfileByUserID(ctx, user.ID)

	log.Println("User authorized")

	logging.LogHandlerInfo(logger, fmt.Sprintf("User %s is authorized", user.Email), responses.StatusOk)
	responses.SendOkResponse(writer, NewAuthOkResponse(*user, profile.AvatarIMG, true))
}

func (authHandler *AuthHandler) EditUserEmail(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	storage := authHandler.storage

	session, err := request.Cookie("session_id")

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusUnauthorized)
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	user, _ := storage.GetUserBySession(ctx, session.Value)

	var newUser models.UnauthorizedUser

	err = json.NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	email := newUser.Email

	errors := utils.ValidateEmail(email)
	if errors != nil {
		logging.LogHandlerInfo(logger, responses.ErrNotValidData, responses.StatusBadRequest)
		log.Println("Email is not valid", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			"Email is not valid"))

		return
	}

	user, err = storage.EditUserEmail(ctx, user.ID, email)

	if err != nil {
		logging.LogHandlerInfo(logger, responses.ErrUserAlreadyExists, responses.StatusBadRequest)
		log.Println("This email is already in use", responses.StatusBadRequest)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusBadRequest,
			"This email is already in use"))

		return
	}

	logging.LogHandlerInfo(logger, fmt.Sprintf("User %s successfully changed his authorization data", user.Email), responses.StatusOk)
	responses.SendOkResponse(writer, NewAuthOkResponse(*user, "", true))

	log.Println("User", user.Email, "successfully changed his authorization data.")
}

func (authHandler *AuthHandler) GetCSRFToken(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	if request.Method != http.MethodGet {
		http.Error(writer, responses.ErrNotAllowed, responses.StatusNotAllowed)

		return
	}

	sessionInstance, ok := ctx.Value(config.SessionContextKey).(models.Session)
	if !ok {
		err := errors.New("error while getting sessionInstance from context")
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	secret := "Vol4okSecretKey"
	hashToken, err := csrf.NewHMACHashToken(secret)

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	fmt.Printf("ID сессии: %s\n", sessionInstance.Value)

	tokenExpTime := time.Now().Add(24 * time.Hour).Unix()
	token, err := hashToken.Create(&sessionInstance, tokenExpTime)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(writer, NewAuthErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}
	fmt.Printf("Сгенерированный токен: %s\n", token)

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, NewSessionOkResponse(token))
}
