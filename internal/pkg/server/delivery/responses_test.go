//nolint:all
package responses_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	responses "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
)

func TestSendOkResponse(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder() //nolint:varnamelen

	testVar := []int{1, 2, 3}

	testResponse := models.OkResponse{
		Code:  200,
		Items: []int{1, 2, 3},
	}

	responses.SendOkResponse(w, responses.NewOkResponse(testVar))

	if status := w.Code; status != responses.StatusOk {
		t.Errorf("Expected status %v, got %v", responses.StatusOk, status)
	}

	expected, err := json.Marshal(testResponse)
	if err != nil {
		t.Fatal("Errow while json.Marshal")
	}

	if w.Body.String() != string(expected) {
		t.Errorf("Expected %v, got %v", string(expected), w.Body.String())
	}
}

func TestSendErrResponse(t *testing.T) {
	t.Parallel()

	responseWriter := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "http://example.com/api/handler/", bytes.NewBufferString(""))

	testResponse := models.ErrResponse{
		Code:   400,
		Status: "Bad request",
	}

	code := new(int)
	*code = 400
	ctx := context.WithValue(req.Context(), "code", code)
	req = req.WithContext(ctx)

	responses.SendErrResponse(req, responseWriter, responses.NewErrResponse(responses.StatusBadRequest,
		responses.ErrBadRequest))

	if status := responseWriter.Code; status != responses.StatusOk {
		t.Errorf("Expected status %v, got %v", responses.StatusOk, status)
	}

	expected, err := json.Marshal(testResponse)
	if err != nil {
		t.Fatal("Errow while json.Marshal")
	}

	if responseWriter.Body.String() != string(expected) {
		t.Errorf("Expected %v, got %v", string(expected), responseWriter.Body.String())
	}
}
