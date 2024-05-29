package delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

	models "github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/config"
	profileproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery/protobuf"
	mock_user_profile "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/delivery/protobuf/mocks"
	delivery "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"
	authproto "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	mock_user_client "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf/mocks"
	utils "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	authClient := mock_user_client.NewMockAuthClient(ctrl)
	profileClient := mock_user_profile.NewMockProfileClient(ctrl)
	defer ctrl.Finish()

	authHandler := delivery.NewAuthHandler(authClient, profileClient)
	mockResponseWriter := httptest.NewRecorder()

	testUser := models.UnauthorizedUser{
		Email:    "test@example.com",
		Password: "password",
	}
	testUserProto := authproto.ExistedUserData{
		Email:    "test@example.com",
		Password: "password",
	}
	testUserJSON, _ := json.Marshal(testUser)
	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(testUserJSON))
	request.Header.Set("Content-Type", "application/json")

	authClient.EXPECT().Login(gomock.Any(), &testUserProto).Return(&authproto.LoggedUser{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: utils.HashPassword("password"),
		SessionID:    "123456",
		IsAuth:       true,
	}, nil)

	authHandler.Login(mockResponseWriter, request)
	assert.Equal(t, http.StatusOK, mockResponseWriter.Code)
}

func TestLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	authClient := mock_user_client.NewMockAuthClient(ctrl)
	profileClient := mock_user_profile.NewMockProfileClient(ctrl)
	defer ctrl.Finish()

	authHandler := delivery.NewAuthHandler(authClient, profileClient)
	mockResponseWriter := httptest.NewRecorder()

	testUser := models.UnauthorizedUser{
		Email:    "test@example.com",
		Password: "password",
	}
	testSessionProto := authproto.SessionData{SessionID: "123456"}
	empty := emptypb.Empty{}
	testUserJSON, _ := json.Marshal(testUser)
	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(testUserJSON))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "session_id", Value: "123456"})

	authClient.EXPECT().Logout(gomock.Any(), &testSessionProto).Return(&empty, nil)

	authHandler.Logout(mockResponseWriter, request)
	assert.Equal(t, http.StatusOK, mockResponseWriter.Code)
}

func TestSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	authClient := mock_user_client.NewMockAuthClient(ctrl)
	profileClient := mock_user_profile.NewMockProfileClient(ctrl)
	defer ctrl.Finish()

	authHandler := delivery.NewAuthHandler(authClient, profileClient)
	mockResponseWriter := httptest.NewRecorder()

	testUser := models.UnauthorizedUser{
		Email:          "test@example.com",
		Password:       "password",
		PasswordRepeat: "password",
	}
	testNewUserDataProto := authproto.NewUserData{
		Email:          "test@example.com",
		Password:       "password",
		PasswordRepeat: "password",
	}
	testUserJSON, _ := json.Marshal(testUser)
	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(testUserJSON))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "session_id", Value: "123456"})
	loggedUser := authproto.LoggedUser{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: utils.HashPassword("password"),
		SessionID:    "123456",
		IsAuth:       true,
	}

	testNewProfileProto := profileproto.ProfileIDRequest{ID: 1}

	authClient.EXPECT().Signup(gomock.Any(), &testNewUserDataProto).Return(&loggedUser, nil)
	profileClient.EXPECT().CreateProfile(gomock.Any(), &testNewProfileProto)
	authHandler.Signup(mockResponseWriter, request)
	assert.Equal(t, http.StatusOK, mockResponseWriter.Code)
}

func TestCheckAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	authClient := mock_user_client.NewMockAuthClient(ctrl)
	profileClient := mock_user_profile.NewMockProfileClient(ctrl)
	defer ctrl.Finish()

	authHandler := delivery.NewAuthHandler(authClient, profileClient)
	mockResponseWriter := httptest.NewRecorder()

	testUser := models.UnauthorizedUser{
		Email:          "test@example.com",
		Password:       "password",
		PasswordRepeat: "password",
	}

	testAuthUser := authproto.AuthUser{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: utils.HashPassword("password"),
		IsAuth:       true,
	}

	testUserJSON, _ := json.Marshal(testUser)
	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(testUserJSON))
	request.Header.Set("Content-Type", "application/json")
	testSessionProto := authproto.SessionData{SessionID: "123456"}
	request.AddCookie(&http.Cookie{Name: "session_id", Value: "123456"})

	testNewProfileProto := profileproto.ProfileIDRequest{ID: 1}
	testNewProfileDataProto := profileproto.ProfileData{
		ID:              1,
		UserID:          1,
		Name:            "name",
		Surname:         "name",
		CityID:          1,
		CityName:        "name",
		Translation:     "name",
		Phone:           "name",
		Avatar:          "name",
		Rating:          1,
		ReactionsCount:  1,
		Approved:        true,
		MerchantsName:   "name",
		SubersCount:     1,
		SubonsCount:     1,
		AvatarIMG:       "name",
		ActiveAddsCount: 1,
		SoldAddsCount:   1,
		CartNum:         1,
		FavNum:          1,
	}

	authClient.EXPECT().GetCurrentUser(gomock.Any(), &testSessionProto).Return(&testAuthUser, nil)
	profileClient.EXPECT().GetProfile(gomock.Any(), &testNewProfileProto).Return(&testNewProfileDataProto, nil)
	authHandler.CheckAuth(mockResponseWriter, request)
	assert.Equal(t, http.StatusOK, mockResponseWriter.Code)
}

func TestEditUserEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	authClient := mock_user_client.NewMockAuthClient(ctrl)
	profileClient := mock_user_profile.NewMockProfileClient(ctrl)
	defer ctrl.Finish()

	authHandler := delivery.NewAuthHandler(authClient, profileClient)
	mockResponseWriter := httptest.NewRecorder()

	testUser := models.UnauthorizedUser{
		Email:          "test@example.com",
		Password:       "password",
		PasswordRepeat: "password",
	}

	testEditEmail := authproto.EditEmailRequest{
		Email:     "test@example.com",
		SessionID: "123456",
	}

	testUserJSON, _ := json.Marshal(testUser)
	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(testUserJSON))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "session_id", Value: "123456"})

	testUserProto := authproto.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: utils.HashPassword("password"),
	}

	authClient.EXPECT().EditEmail(gomock.Any(), &testEditEmail).Return(&testUserProto, nil)

	authHandler.EditUserEmail(mockResponseWriter, request)
	assert.Equal(t, http.StatusOK, mockResponseWriter.Code)
}

func TestGetCSRFToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	authClient := mock_user_client.NewMockAuthClient(ctrl)
	profileClient := mock_user_profile.NewMockProfileClient(ctrl)
	defer ctrl.Finish()

	authHandler := delivery.NewAuthHandler(authClient, profileClient)
	mockResponseWriter := httptest.NewRecorder()

	testUser := models.UnauthorizedUser{
		Email:          "test@example.com",
		Password:       "password",
		PasswordRepeat: "password",
	}

	testUserJSON, _ := json.Marshal(testUser)
	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(testUserJSON))
	request.AddCookie(&http.Cookie{Name: "session_id", Value: "123456"})
	request.Header.Set("Content-Type", "application/json")

	sessionInstance := models.Session{
		UserID: uint32(1),
		Value:  "123456",
	}

	ctx := context.WithValue(request.Context(), config.SessionContextKey, sessionInstance)
	request = request.WithContext(ctx)

	authHandler.GetCSRFToken(mockResponseWriter, request)
	assert.Equal(t, http.StatusOK, mockResponseWriter.Code)
}
