package delivery_test

import (
	"testing"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/responses"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/storage"
)

func TestNewAdvertsOkResponse(t *testing.T) {
	t.Parallel()

	adverts := []*storage.Advert{}

	response := responses.NewAdvertsOkResponse(adverts)

	if response.Code != responses.StatusOk {
		t.Errorf("Expected Code to be %d, got %d", responses.StatusOk, response.Code)
	}

	if len(response.Adverts) != len(adverts) {
		t.Errorf("Expected Adverts length to be %d, got %d", len(adverts), len(response.Adverts))
	}
}

func TestNewAdvertsErrResponse(t *testing.T) {
	t.Parallel()

	code := responses.StatusNotFound
	status := responses.ErrAdvertNotExist

	response := responses.NewAdvertsErrResponse(code, status)

	if response.Code != code {
		t.Errorf("Expected Code to be %d, got %d", code, response.Code)
	}

	if response.Status != status {
		t.Errorf("Expected Status to be %s, got %s", status, response.Status)
	}
}
