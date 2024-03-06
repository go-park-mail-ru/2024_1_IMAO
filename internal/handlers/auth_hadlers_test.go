package myhandlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "github.com/go-park-mail-ru/2024_1_IMAO/internal/handlers"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
)

func TestLoginHandlerSuccessful(t *testing.T) { //nolint:funlen
	t.Parallel()

	type TestCase struct {
		name             string
		inputUser        *storage.UnauthorizedUser
		expectedResponse *responses.AuthOkResponse
	}

	testCases := [...]TestCase{
		{
			name: "Base Test",
			inputUser: &storage.UnauthorizedUser{
				Email:          `example@mail.ru`,
				Password:       "123456",
				PasswordRepeat: "",
			},
			expectedResponse: responses.NewAuthOkResponse(storage.User{
				ID:           1,
				Name:         "Barak",
				Surname:      "Obama",
				Email:        "example@mail.ru",
				PasswordHash: "111-222-333",
			}, "111-222-333-444", true),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			if err := writer.WriteField("email", testCase.inputUser.Email); err != nil {
				t.Fatalf("Failed to write 'email' field: %v", err)
			}

			if err := writer.WriteField("password", testCase.inputUser.Password); err != nil {
				t.Fatalf("Failed to write 'password' field: %v", err)
			}

			if err := writer.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/login", body)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()

			usersList := storage.NewActiveUser()
			authHandler := &handler.AuthHandler{
				List: usersList,
			}
			authHandler.Login(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)

			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse responses.AuthOkResponse

			err = json.Unmarshal(receivedResponse, &resultResponse)
			if err != nil {
				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
			}

			codeIsEqual := resultResponse.Code == testCase.expectedResponse.Code
			userIDIsEqual := resultResponse.User.ID == testCase.expectedResponse.User.ID
			emailIsEqual := resultResponse.User.Email == testCase.expectedResponse.User.Email
			isAuthIsEqual := resultResponse.IsAuth == testCase.expectedResponse.IsAuth

			if !codeIsEqual || !userIDIsEqual || !emailIsEqual || !isAuthIsEqual {
				t.Errorf("wrong Response: got %+v, expected %+v",
					resultResponse, testCase.expectedResponse)
			}
		})
	}
}

func TestSignUpHandlerSuccessful(t *testing.T) { //nolint:funlen
	t.Parallel()

	type TestCase struct {
		name             string
		inputUser        *storage.UnauthorizedUser
		expectedResponse *responses.AuthOkResponse
	}

	testCases := [...]TestCase{
		{
			name: "Base Test",
			inputUser: &storage.UnauthorizedUser{
				Email:          "bigbob@mail.ru",
				Password:       "BigBob-123456",
				PasswordRepeat: "BigBob-123456",
			},
			expectedResponse: responses.NewAuthOkResponse(storage.User{
				ID:           2,
				Name:         "Barak",
				Surname:      "Obama",
				Email:        "bigbob@mail.ru",
				PasswordHash: "111-222-333",
			}, "111-222-333-444", true),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			if err := writer.WriteField("email", testCase.inputUser.Email); err != nil {
				t.Fatalf("Failed to write 'email' field: %v", err)
			}

			if err := writer.WriteField("password", testCase.inputUser.Password); err != nil {
				t.Fatalf("Failed to write 'password' field: %v", err)
			}

			if err := writer.WriteField("passwordRepeat", testCase.inputUser.PasswordRepeat); err != nil {
				t.Fatalf("Failed to write 'passwordRepeat' field: %v", err)
			}

			if err := writer.Close(); err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/signup", body)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())

			w := httptest.NewRecorder()

			usersList := storage.NewActiveUser()
			authHandler := &handler.AuthHandler{
				List: usersList,
			}
			authHandler.Signup(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)

			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse responses.AuthOkResponse

			err = json.Unmarshal(receivedResponse, &resultResponse)
			if err != nil {
				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
			}

			codeIsEqual := resultResponse.Code == testCase.expectedResponse.Code
			userIDIsEqual := resultResponse.User.ID == testCase.expectedResponse.User.ID
			emailIsEqual := resultResponse.User.Email == testCase.expectedResponse.User.Email
			isAuthIsEqual := resultResponse.IsAuth == testCase.expectedResponse.IsAuth

			if !codeIsEqual || !userIDIsEqual || !emailIsEqual || !isAuthIsEqual {
				t.Errorf("wrong Response: got %+v, expected %+v",
					resultResponse, testCase.expectedResponse)
			}
		})
	}
}
