package delivery

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	profileproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery/protobuf"

	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	csrf "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/middleware/csrf"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
)

const (
	sessionTime   = 24 * time.Hour
	tokenDuration = 24 * time.Hour
	sameSite      = 4
)

type AuthHandler struct {
	authClient    authproto.AuthClient
	profileClient profileproto.ProfileClient
}

func NewAuthHandler(authClient authproto.AuthClient,
	profileClient profileproto.ProfileClient) *AuthHandler {
	return &AuthHandler{
		authClient:    authClient,
		profileClient: profileClient,
	}
}

func createSession(sessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(sessionTime),
		HttpOnly: true,
		SameSite: sameSite,
		Secure:   true,
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

	client := authHandler.authClient

	user := new(models.UnauthorizedUser)

	data, _ := io.ReadAll(request.Body)

	err := user.UnmarshalJSON(data)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	authUser, err := client.Login(ctx, &authproto.ExistedUserData{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusUnauthorized)
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	cookie := createSession(authUser.SessionID)

	http.SetCookie(writer, cookie)

	userData := models.AuthResponse{
		User: models.User{
			ID:    uint(authUser.ID),
			Email: authUser.Email,
		},
		IsAuth: true,
	}
	responses.SendOkResponse(writer, responses.NewOkResponse(userData))

	logging.LogHandlerInfo(logger, fmt.Sprintf("User %s have been authorized with session ID: %s ",
		user.Email, authUser.SessionID), responses.StatusOk)
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

	client := authHandler.authClient

	session, err := request.Cookie("session_id")
	_, clientErr := client.Logout(ctx, &authproto.SessionData{SessionID: session.Value})

	if err != nil || clientErr != nil {
		logging.LogHandlerError(logger, err, responses.StatusUnauthorized)
		logging.LogHandlerError(logger, clientErr, responses.StatusUnauthorized)
		log.Println(err, responses.StatusUnauthorized)
		log.Println(clientErr, responses.StatusUnauthorized)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(writer, session)

	userData := models.AuthResponse{
		IsAuth: false,
	}
	responses.SendOkResponse(writer, responses.NewOkResponse(userData))

	// ПО-ХОРОШЕМУ НУЖНО ПЕРЕПИСАТЬ ХЭНДЛЕР, ЧТОБЫ В ЛОГЕ МОЖНО БЫЛО ВЫВОДИТЬ КАКОЙ ИМЕННО ПОЛЬЗОВАТЕЛЬ РАЗЛОГИНИЛСЯ
	logging.LogHandlerInfo(logger, "User have been logged out", responses.StatusOk)
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
func (authHandler *AuthHandler) Signup(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	client := authHandler.authClient
	profileClient := authHandler.profileClient

	var newUser models.UnauthorizedUser

	data, _ := io.ReadAll(request.Body)

	err := newUser.UnmarshalJSON(data)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	user, err := client.Signup(ctx, &authproto.NewUserData{
		Email:          newUser.Email,
		Password:       newUser.Password,
		PasswordRepeat: newUser.PasswordRepeat,
	})
	if err != nil {
		logging.LogHandlerError(logger, responses.ErrWrongCredentials, responses.StatusBadRequest)
		log.Println(err)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			responses.ErrWrongCredentials))

		return
	}

	cookie := createSession(user.SessionID)

	http.SetCookie(writer, cookie)

	_, _ = profileClient.CreateProfile(ctx, &profileproto.ProfileIDRequest{ID: user.ID})

	userData := models.AuthResponse{
		User: models.User{
			ID:    uint(user.ID),
			Email: user.Email,
		},
		IsAuth: true,
	}
	responses.SendOkResponse(writer, responses.NewOkResponse(userData))

	logging.LogHandlerInfo(logger, fmt.Sprintf("User %s have been authorized with session ID: %s ",
		user.Email, user.SessionID), responses.StatusOk)
	log.Println("User", user.Email, "have been authorized with session ID:", user.SessionID)
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

	client := authHandler.authClient
	profileClient := authHandler.profileClient

	session, err := request.Cookie("session_id")

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusUnauthorized)
		responses.SendOkResponse(writer,
			responses.NewOkResponse(models.AuthResponse{IsAuth: false}))

		return
	}

	user, _ := client.GetCurrentUser(ctx, &authproto.SessionData{
		SessionID: session.Value,
	})

	var responseData models.AdditionalUserData

	if user.IsAuth {
		profile, _ := profileClient.GetProfile(ctx, &profileproto.ProfileIDRequest{ID: user.ID})

		logging.LogHandlerInfo(logger, fmt.Sprintf("User %s is authorized", user.Email), responses.StatusOk)

		responseData.Avatar = profile.AvatarIMG
		responseData.PhoneNumber = profile.Phone
		responseData.FavNum = uint(profile.FavNum)
		responseData.CartNum = uint(profile.CartNum)
	} else {
		logging.LogHandlerInfo(logger, "User not authorized", responses.StatusOk)
	}

	responseData.User = models.User{
		ID:    uint(user.ID),
		Email: user.Email,
	}
	responseData.IsAuth = user.IsAuth

	responses.SendOkResponse(writer, responses.NewOkResponse(responseData))
}

func (authHandler *AuthHandler) EditUserEmail(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	client := authHandler.authClient

	session, err := request.Cookie("session_id")
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusUnauthorized)
		log.Println(err, responses.StatusUnauthorized)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusUnauthorized,
			responses.ErrUnauthorized))

		return
	}

	var newUser models.UnauthorizedUser

	data, _ := io.ReadAll(request.Body)

	err = newUser.UnmarshalJSON(data)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	user, err := client.EditEmail(ctx, &authproto.EditEmailRequest{
		Email:     newUser.Email,
		SessionID: session.Value,
	})
	if err != nil {
		logging.LogHandlerInfo(logger, responses.ErrUserAlreadyExists, responses.StatusBadRequest)
		log.Println("This email is already in use", responses.StatusBadRequest)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusBadRequest,
			"This email is already in use"))

		return
	}

	logging.LogHandlerInfo(logger, fmt.Sprintf("User %s successfully changed his authorization data",
		user.Email), responses.StatusOk)

	userData := models.AuthResponse{
		User: models.User{
			ID:    uint(user.ID),
			Email: user.Email,
		},
		IsAuth: true,
	}
	responses.SendOkResponse(writer, responses.NewOkResponse(userData))
}

func (authHandler *AuthHandler) GetCSRFToken(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	sessionInstance, ok := ctx.Value(config.SessionContextKey).(models.Session)
	if !ok {
		err := errors.New("error while getting sessionInstance from context")
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	secret := "Vol4okSecretKey"
	hashToken, err := csrf.NewHMACHashToken(secret)

	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	fmt.Printf("ID сессии: %s\n", sessionInstance.Value)

	tokenExpTime := time.Now().Add(tokenDuration).Unix()

	token, err := hashToken.Create(&sessionInstance, tokenExpTime)
	if err != nil {
		logging.LogHandlerError(logger, err, responses.StatusInternalServerError)
		log.Println(err, responses.StatusInternalServerError)
		responses.SendErrResponse(request, writer, responses.NewErrResponse(responses.StatusInternalServerError,
			responses.ErrInternalServer))
	}

	log.Printf("Сгенерированный токен: %s\n", token)

	logging.LogHandlerInfo(logger, "success", responses.StatusOk)
	responses.SendOkResponse(writer, responses.NewOkResponse(token))
}
