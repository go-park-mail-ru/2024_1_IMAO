package delivery_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	models "github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	mock_profile "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/profile/mocks"
	delivery "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"
	mock_user "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/mocks"
	utils "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
)

// func TestHandler_Login(t *testing.T) {

// }

const plug = ""

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	usersStorageInterface := mock_user.NewMockUsersStorageInterface(ctrl)
	profileStorageInterface := mock_profile.NewMockProfileStorageInterface(ctrl)
	defer ctrl.Finish()

	authHandler := delivery.NewAuthHandler(usersStorageInterface, profileStorageInterface, plug, plug)
	mockResponseWriter := httptest.NewRecorder()

	// Создаем экземпляр AuthHandler с моком storage

	// Подготовка тестовых данных
	testUser := models.UnauthorizedUser{
		Email:    "test@example.com",
		Password: "password",
	}
	testUserJSON, _ := json.Marshal(testUser)
	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(testUserJSON))
	request.Header.Set("Content-Type", "application/json")

	// Настройка мока для GetUserByEmail
	usersStorageInterface.EXPECT().GetUserByEmail(gomock.Any(), testUser.Email).Return(&models.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: utils.HashPassword("password"), //
	}, nil)
	usersStorageInterface.EXPECT().AddSession(uint(1)).Return("session_id")
	//usersStorageInterface.EXPECT().AddSession(100).Return()

	// Настройка мока для AddSession
	//mockStorage.EXPECT().AddSession(gomock.Any()).Return("session_id", nil)

	// Вызов функции Login
	authHandler.Login(mockResponseWriter, request)

	// Проверка ответа
	assert.Equal(t, http.StatusOK, mockResponseWriter.Code)
	// Здесь можно добавить дополнительные проверки, например, содержимое ответа
}

// func TestLogout(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	usersStorageInterface := mock_user.NewMockUsersStorageInterface(ctrl)
// 	//profileStorageInterface := mock_profile.NewMockProfileStorageInterface(ctrl)
// 	defer ctrl.Finish()

// 	//mockResponseWriter := httptest.NewRecorder()

// 	// Создаем экземпляр AuthHandler с моком storage
// 	//authHandler := delivery.NewAuthHandler(usersStorageInterface, profileStorageInterface, plug, plug)

// 	// Подготовка тестового запроса
// 	//request, _ := http.NewRequest("POST", "/logout", nil)
// 	//request.AddCookie(&http.Cookie{Name: "session_id", Value: "test_session_id"})

// 	// Настройка мока для RemoveSession
// 	usersStorageInterface.EXPECT().RemoveSession(gomock.Any(), "test_session_id").Return(nil)

// 	// Вызов функции Logout
// 	//authHandler.Logout(mockResponseWriter, request)

// 	// Проверка ответа
// 	//assert.Equal(t, http.StatusOK, mockResponseWriter.Code)
// 	// Проверка, что cookie был установлен с датой истечения
// 	//cookies := mockResponseWriter.Header().Get("Set-Cookie")
// 	//assert.Contains(t, cookies, "session_id=; expires=")
// 	// Здесь можно добавить дополнительные проверки, например, содержимое ответа
// }

// func TestLoginHandlerSuccessful(t *testing.T) { //nolint:funlen
// 	t.Parallel()

// 	type TestCase struct {
// 		name             string
// 		inputUser        *storage.UnauthorizedUser
// 		expectedResponse *responses.AuthOkResponse
// 	}

// 	testCases := [...]TestCase{
// 		{
// 			name: "Base Test",
// 			inputUser: &storage.UnauthorizedUser{
// 				Email:          `example@mail.ru`,
// 				Password:       "123456",
// 				PasswordRepeat: "",
// 			},
// 			expectedResponse: responses.NewAuthOkResponse(storage.User{
// 				ID:           1,
// 				Name:         "Barak",
// 				Surname:      "Obama",
// 				Email:        "example@mail.ru",
// 				PasswordHash: "111-222-333",
// 			}, "111-222-333-444", true),
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		testCase := testCase

// 		t.Run(testCase.name, func(t *testing.T) {
// 			t.Parallel()

// 			body := &bytes.Buffer{}
// 			writer := multipart.NewWriter(body)

// 			if err := writer.WriteField("email", testCase.inputUser.Email); err != nil {
// 				t.Fatalf("Failed to write 'email' field: %v", err)
// 			}

// 			if err := writer.WriteField("password", testCase.inputUser.Password); err != nil {
// 				t.Fatalf("Failed to write 'password' field: %v", err)
// 			}

// 			if err := writer.Close(); err != nil {
// 				t.Fatalf("Failed to close writer: %v", err)
// 			}

// 			req, err := http.NewRequest(http.MethodPost, "/login", body)
// 			if err != nil {
// 				t.Fatalf("Failed to create request: %v", err)
// 			}

// 			req.Header.Set("Content-Type", writer.FormDataContentType())

// 			w := httptest.NewRecorder()

// 			usersList := storage.NewActiveUser()
// 			authHandler := &handler.AuthHandler{
// 				List: usersList,
// 			}
// 			authHandler.Login(w, req)

// 			resp := w.Result()
// 			defer resp.Body.Close()

// 			receivedResponse, err := io.ReadAll(resp.Body)

// 			if err != nil {
// 				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
// 			}

// 			var resultResponse responses.AuthOkResponse

// 			err = json.Unmarshal(receivedResponse, &resultResponse)
// 			if err != nil {
// 				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
// 			}

// 			codeIsEqual := resultResponse.Code == testCase.expectedResponse.Code
// 			userIDIsEqual := resultResponse.User.ID == testCase.expectedResponse.User.ID
// 			emailIsEqual := resultResponse.User.Email == testCase.expectedResponse.User.Email
// 			isAuthIsEqual := resultResponse.IsAuth == testCase.expectedResponse.IsAuth

// 			if !codeIsEqual || !userIDIsEqual || !emailIsEqual || !isAuthIsEqual {
// 				t.Errorf("wrong Response: got %+v, expected %+v",
// 					resultResponse, testCase.expectedResponse)
// 			}
// 		})
// 	}
// }

// func TestSignUpHandlerSuccessful(t *testing.T) { //nolint:funlen
// 	t.Parallel()

// 	type TestCase struct {
// 		name             string
// 		inputUser        *storage.UnauthorizedUser
// 		expectedResponse *responses.AuthOkResponse
// 	}

// 	testCases := [...]TestCase{
// 		{
// 			name: "Base Test",
// 			inputUser: &storage.UnauthorizedUser{
// 				Email:          "bigbob@mail.ru",
// 				Password:       "BigBob-123456",
// 				PasswordRepeat: "BigBob-123456",
// 			},
// 			expectedResponse: responses.NewAuthOkResponse(storage.User{
// 				ID:           2,
// 				Name:         "Barak",
// 				Surname:      "Obama",
// 				Email:        "bigbob@mail.ru",
// 				PasswordHash: "111-222-333",
// 			}, "111-222-333-444", true),
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		testCase := testCase

// 		t.Run(testCase.name, func(t *testing.T) {
// 			t.Parallel()

// 			body := &bytes.Buffer{}
// 			writer := multipart.NewWriter(body)

// 			if err := writer.WriteField("email", testCase.inputUser.Email); err != nil {
// 				t.Fatalf("Failed to write 'email' field: %v", err)
// 			}

// 			if err := writer.WriteField("password", testCase.inputUser.Password); err != nil {
// 				t.Fatalf("Failed to write 'password' field: %v", err)
// 			}

// 			if err := writer.WriteField("passwordRepeat", testCase.inputUser.PasswordRepeat); err != nil {
// 				t.Fatalf("Failed to write 'passwordRepeat' field: %v", err)
// 			}

// 			if err := writer.Close(); err != nil {
// 				t.Fatalf("Failed to close writer: %v", err)
// 			}

// 			req, err := http.NewRequest(http.MethodPost, "/signup", body)
// 			if err != nil {
// 				t.Fatalf("Failed to create request: %v", err)
// 			}

// 			req.Header.Set("Content-Type", writer.FormDataContentType())

// 			w := httptest.NewRecorder()

// 			usersList := storage.NewActiveUser()
// 			authHandler := &handler.AuthHandler{
// 				List: usersList,
// 			}
// 			authHandler.Signup(w, req)

// 			resp := w.Result()
// 			defer resp.Body.Close()

// 			receivedResponse, err := io.ReadAll(resp.Body)

// 			if err != nil {
// 				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
// 			}

// 			var resultResponse responses.AuthOkResponse

// 			err = json.Unmarshal(receivedResponse, &resultResponse)
// 			if err != nil {
// 				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
// 			}

// 			codeIsEqual := resultResponse.Code == testCase.expectedResponse.Code
// 			userIDIsEqual := resultResponse.User.ID == testCase.expectedResponse.User.ID
// 			emailIsEqual := resultResponse.User.Email == testCase.expectedResponse.User.Email
// 			isAuthIsEqual := resultResponse.IsAuth == testCase.expectedResponse.IsAuth

// 			if !codeIsEqual || !userIDIsEqual || !emailIsEqual || !isAuthIsEqual {
// 				t.Errorf("wrong Response: got %+v, expected %+v",
// 					resultResponse, testCase.expectedResponse)
// 			}
// 		})
// 	}
// }

// func TestLogoutHandlerSuccessful(t *testing.T) { //nolint:funlen
// 	t.Parallel()

// 	type TestCase struct {
// 		name             string
// 		inputUser        *storage.UnauthorizedUser
// 		expectedResponse *responses.AuthOkResponse
// 	}

// 	testCases := [...]TestCase{
// 		{
// 			name: "Base Test",
// 			inputUser: &storage.UnauthorizedUser{
// 				Email:          "example@mail.ru",
// 				Password:       "123456",
// 				PasswordRepeat: "",
// 			},
// 			expectedResponse: responses.NewAuthOkResponse(storage.User{
// 				ID:           0,
// 				Name:         "",
// 				Surname:      "",
// 				Email:        "",
// 				PasswordHash: "",
// 			}, "", false),
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		testCase := testCase
// 		t.Run(testCase.name, func(t *testing.T) {
// 			t.Parallel()

// 			body, writer := createLoginRequestBody(t, testCase.inputUser)
// 			req, err := http.NewRequest(http.MethodPost, "/login", body)

// 			if err != nil {
// 				t.Fatalf("Failed to create request: %v", err)
// 			}

// 			req.Header.Set("Content-Type", writer.FormDataContentType())

// 			responseWriter1 := httptest.NewRecorder()
// 			usersList := storage.NewActiveUser()
// 			authHandler := &handler.AuthHandler{
// 				List: usersList,
// 			}
// 			authHandler.Login(responseWriter1, req)

// 			sessionID := extractSessionID(t, responseWriter1.Result())

// 			req2, err := createLogoutRequest(sessionID)
// 			if err != nil {
// 				t.Fatalf("Failed to create logout request: %v", err)
// 			}

// 			responseWriter2 := httptest.NewRecorder()
// 			authHandler.Logout(responseWriter2, req2)

// 			resp := responseWriter2.Result()

// 			defer resp.Body.Close()

// 			checkLogoutResponse(t, resp, testCase.expectedResponse)
// 		})
// 	}
// }

// func TestCheckAuthHandlerSuccessful(t *testing.T) { //nolint:funlen
// 	t.Parallel()

// 	type TestCase struct {
// 		name             string
// 		inputUser        *storage.UnauthorizedUser
// 		expectedResponse *responses.AuthOkResponse
// 	}

// 	testCases := [...]TestCase{
// 		{
// 			name: "AuthorizedUser",
// 			inputUser: &storage.UnauthorizedUser{
// 				Email:          "example@mail.ru",
// 				Password:       "111111",
// 				PasswordRepeat: "",
// 			},
// 			expectedResponse: responses.NewAuthOkResponse(storage.User{
// 				ID:           0,
// 				Name:         "",
// 				Surname:      "",
// 				Email:        "",
// 				PasswordHash: "",
// 			}, "", false),
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		testCase := testCase
// 		t.Run(testCase.name, func(t *testing.T) {
// 			t.Parallel()

// 			body, writer := createLoginRequestBody(t, testCase.inputUser)
// 			req, err := http.NewRequest(http.MethodPost, "/login", body)

// 			if err != nil {
// 				t.Fatalf("Failed to create request: %v", err)
// 			}

// 			req.Header.Set("Content-Type", writer.FormDataContentType())

// 			responseWriter1 := httptest.NewRecorder()
// 			usersList := storage.NewActiveUser()
// 			authHandler := &handler.AuthHandler{
// 				List: usersList,
// 			}
// 			authHandler.Login(responseWriter1, req)

// 			sessionID := extractSessionID(t, responseWriter1.Result())

// 			req2, err := createCheckAuthRequest(sessionID)
// 			if err != nil {
// 				t.Fatalf("Failed to create logout request: %v", err)
// 			}

// 			responseWriter2 := httptest.NewRecorder()
// 			authHandler.CheckAuth(responseWriter2, req2)

// 			resp := responseWriter2.Result()

// 			defer resp.Body.Close()

// 			checkLogoutResponse(t, resp, testCase.expectedResponse)
// 		})
// 	}
// }

// func createLoginRequestBody(t *testing.T, user *storage.UnauthorizedUser) (*bytes.Buffer, *multipart.Writer) {
// 	t.Helper()

// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)

// 	for _, field := range []struct {
// 		Name  string
// 		Value string
// 	}{
// 		{"email", user.Email},
// 		{"password", user.Password},
// 		{"passwordRepeat", user.PasswordRepeat},
// 	} {
// 		if err := writer.WriteField(field.Name, field.Value); err != nil {
// 			t.Fatalf("Failed to write '%s' field: %v", field.Name, err)
// 		}
// 	}

// 	if err := writer.Close(); err != nil {
// 		t.Fatalf("Failed to close writer: %v", err)
// 	}

// 	return body, writer
// }

// func createLogoutRequest(sessionID string) (*http.Request, error) {
// 	req, err := http.NewRequest(http.MethodPost, "/logout", nil)

// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Add("Cookie", fmt.Sprintf("session_id=%s", sessionID))

// 	return req, nil
// }

// func createCheckAuthRequest(sessionID string) (*http.Request, error) {
// 	req, err := http.NewRequest(http.MethodGet, "/check_auth", nil)

// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Add("Cookie", fmt.Sprintf("session_id=%s", sessionID))

// 	return req, nil
// }

// func extractSessionID(t *testing.T, resp *http.Response) string {
// 	t.Helper()

// 	for _, cookie := range resp.Cookies() {
// 		if cookie.Name == "session_id" {
// 			return cookie.Value
// 		}
// 	}

// 	return ""
// }

// func checkLogoutResponse(t *testing.T, resp *http.Response, expectedResponse *responses.AuthOkResponse) {
// 	t.Helper()

// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("Expected status OK, got %v", resp.StatusCode)
// 	}

// 	var resultResponse responses.AuthOkResponse

// 	bodyBytes, _ := io.ReadAll(resp.Body)
// 	err := json.Unmarshal(bodyBytes, &resultResponse)

// 	if err != nil {
// 		t.Fatalf("Failed to unmarshal response body: %v", err)
// 	}

// 	if !reflect.DeepEqual(resultResponse, *expectedResponse) {
// 		t.Errorf("Expected response %+v, got %+v", *expectedResponse, resultResponse)
// 	}
// }

// func TestCheckAuthAllowedMethods(t *testing.T) {
// 	t.Parallel()

// 	authHandler := &handler.AuthHandler{}

// 	t.Run("GET request", func(t *testing.T) {
// 		t.Parallel()

// 		req := httptest.NewRequest(http.MethodGet, "/check_auth", nil)
// 		w := httptest.NewRecorder()

// 		authHandler.CheckAuth(w, req)

// 		if w.Code != http.StatusOK {
// 			t.Errorf("Expected status OK, got %v", w.Code)
// 		}
// 	})

// 	t.Run("POST request", func(t *testing.T) {
// 		t.Parallel()

// 		req := httptest.NewRequest(http.MethodPost, "/check_auth", nil)
// 		w := httptest.NewRecorder()

// 		authHandler.CheckAuth(w, req)

// 		if w.Code != http.StatusMethodNotAllowed {
// 			t.Errorf("Expected status MethodNotAllowed, got %v", w.Code)
// 		}
// 	})
// }

// func TestLoginAllowedMethods(t *testing.T) {
// 	t.Parallel()

// 	authHandler := &handler.AuthHandler{}

// 	t.Run("GET request", func(t *testing.T) {
// 		t.Parallel()

// 		req := httptest.NewRequest(http.MethodGet, "/login", nil)
// 		w := httptest.NewRecorder()

// 		authHandler.Login(w, req)

// 		if w.Code != http.StatusMethodNotAllowed {
// 			t.Errorf("Expected status MethodNotAllowed, got %v", w.Code)
// 		}
// 	})

// 	t.Run("POST request", func(t *testing.T) {
// 		t.Parallel()

// 		body := &bytes.Buffer{}
// 		writer := multipart.NewWriter(body)

// 		if err := writer.WriteField("email", "example@mail.ru"); err != nil {
// 			t.Fatalf("Failed to write 'email' field: %v", err)
// 		}

// 		if err := writer.Close(); err != nil {
// 			t.Fatalf("Failed to close writer: %v", err)
// 		}

// 		req, err := http.NewRequest(http.MethodPost, "/login", body)
// 		if err != nil {
// 			t.Fatalf("Failed to create request: %v", err)
// 		}

// 		req.Header.Set("Content-Type", writer.FormDataContentType())

// 		w := httptest.NewRecorder()

// 		usersList := storage.NewActiveUser()
// 		authHandler := &handler.AuthHandler{
// 			List: usersList,
// 		}

// 		authHandler.Login(w, req)

// 		if w.Code != http.StatusOK {
// 			t.Errorf("Expected status OK, got %v", w.Code)
// 		}
// 	})
// }

// func TestSignUpAllowedMethods(t *testing.T) {
// 	t.Parallel()

// 	authHandler := &handler.AuthHandler{}

// 	t.Run("GET request", func(t *testing.T) {
// 		t.Parallel()

// 		req := httptest.NewRequest(http.MethodGet, "/signup", nil)
// 		w := httptest.NewRecorder()

// 		authHandler.Signup(w, req)

// 		if w.Code != http.StatusMethodNotAllowed {
// 			t.Errorf("Expected status MethodNotAllowed, got %v", w.Code)
// 		}
// 	})

// 	t.Run("POST request", func(t *testing.T) {
// 		t.Parallel()

// 		body := &bytes.Buffer{}
// 		writer := multipart.NewWriter(body)

// 		if err := writer.WriteField("email", "example@mail.ru"); err != nil {
// 			t.Fatalf("Failed to write 'email' field: %v", err)
// 		}

// 		if err := writer.Close(); err != nil {
// 			t.Fatalf("Failed to close writer: %v", err)
// 		}

// 		req, err := http.NewRequest(http.MethodPost, "/signup", body)
// 		if err != nil {
// 			t.Fatalf("Failed to create request: %v", err)
// 		}

// 		req.Header.Set("Content-Type", writer.FormDataContentType())

// 		w := httptest.NewRecorder()

// 		usersList := storage.NewActiveUser()
// 		authHandler := &handler.AuthHandler{
// 			List: usersList,
// 		}

// 		authHandler.Signup(w, req)

// 		if w.Code != http.StatusOK {
// 			t.Errorf("Expected status OK, got %v", w.Code)
// 		}
// 	})
// }

// func TestLogoutAllowedMethods(t *testing.T) {
// 	t.Parallel()

// 	authHandler := &handler.AuthHandler{}

// 	t.Run("GET request", func(t *testing.T) {
// 		t.Parallel()

// 		req := httptest.NewRequest(http.MethodGet, "/logout", nil)
// 		w := httptest.NewRecorder()

// 		authHandler.Logout(w, req)

// 		if w.Code != http.StatusMethodNotAllowed {
// 			t.Errorf("Expected status MethodNotAllowed, got %v", w.Code)
// 		}
// 	})

// 	t.Run("POST request", func(t *testing.T) {
// 		t.Parallel()

// 		body := &bytes.Buffer{}
// 		writer := multipart.NewWriter(body)

// 		if err := writer.WriteField("email", "example@mail.ru"); err != nil {
// 			t.Fatalf("Failed to write 'email' field: %v", err)
// 		}

// 		if err := writer.Close(); err != nil {
// 			t.Fatalf("Failed to close writer: %v", err)
// 		}

// 		req, err := http.NewRequest(http.MethodPost, "/logout", body)
// 		if err != nil {
// 			t.Fatalf("Failed to create request: %v", err)
// 		}

// 		req.Header.Set("Content-Type", writer.FormDataContentType())

// 		w := httptest.NewRecorder()

// 		usersList := storage.NewActiveUser()
// 		authHandler := &handler.AuthHandler{
// 			List: usersList,
// 		}

// 		authHandler.Logout(w, req)

// 		if w.Code != http.StatusOK {
// 			t.Errorf("Expected status OK, got %v", w.Code)
// 		}
// 	})
// }
