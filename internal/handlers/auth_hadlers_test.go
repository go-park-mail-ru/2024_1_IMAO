package myhandlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	handler "github.com/go-park-mail-ru/2024_1_IMAO/internal/handlers"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
)

func TestLoginHandlerSuccessful(t *testing.T) {
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

			reqBody, err := json.Marshal(&testCase.inputUser)

			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))

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

			if !reflect.DeepEqual(*testCase.expectedResponse, resultResponse) {
				t.Errorf("wrong Response: got %+v, expected %+v",
					resultResponse, testCase.expectedResponse)
			}
		})
	}
}
