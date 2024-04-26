package delivery_test

import (
	"testing"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	delivery "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery"
)

func TestNewAuthOkResponse(t *testing.T) {
	t.Parallel()

	user := models.User{
		ID:           2,
		Email:        "bigbob@mail.ru",
		PasswordHash: "111-222-333",
	}
	avatar := "imageBytes"
	isAuth := true

	response := delivery.NewAuthOkResponse(user, avatar, isAuth)

	if response.User.Email != user.Email {
		t.Errorf("Expected Code to be %s, got %s", user.Email, response.User.Email)
	}

	if response.Code != responses.StatusOk {
		t.Errorf("Expected Code to be %d, got %d", responses.StatusOk, response.Code)
	}

	if response.Avatar != avatar {
		t.Errorf("Expected SessionID to be %s, got %s", avatar, response.Avatar)
	}

	if response.IsAuth != isAuth {
		t.Errorf("Expected IsAuth to be %v, got %v", isAuth, response.IsAuth)
	}
}

func TestNewAuthErrResponse(t *testing.T) {
	t.Parallel()

	code := responses.StatusUnauthorized
	status := responses.ErrUnauthorized

	response := delivery.NewAuthErrResponse(code, status)

	if response.Code != code {
		t.Errorf("Expected Code to be %d, got %d", code, response.Code)
	}

	if response.Status != status {
		t.Errorf("Expected Status to be %s, got %s", status, response.Status)
	}
}
