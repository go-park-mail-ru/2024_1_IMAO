package delivery_test

import (
	"testing"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
)

func TestNewAuthOkResponse(t *testing.T) {
	t.Parallel()

	user := storage.User{
		ID:           2,
		Name:         "Barak",
		Surname:      "Obama",
		Email:        "bigbob@mail.ru",
		PasswordHash: "111-222-333",
	}
	sessionID := "testSessionID"
	isAuth := true

	response := responses.NewAuthOkResponse(user, sessionID, isAuth)

	if response.Code != responses.StatusOk {
		t.Errorf("Expected Code to be %d, got %d", responses.StatusOk, response.Code)
	}

	if response.SessionID != sessionID {
		t.Errorf("Expected SessionID to be %s, got %s", sessionID, response.SessionID)
	}

	if response.IsAuth != isAuth {
		t.Errorf("Expected IsAuth to be %v, got %v", isAuth, response.IsAuth)
	}
}

func TestNewAuthErrResponse(t *testing.T) {
	t.Parallel()

	code := responses.StatusUnauthorized
	status := responses.ErrUnauthorized

	response := responses.NewAuthErrResponse(code, status)

	if response.Code != code {
		t.Errorf("Expected Code to be %d, got %d", code, response.Code)
	}

	if response.Status != status {
		t.Errorf("Expected Status to be %s, got %s", status, response.Status)
	}
}
