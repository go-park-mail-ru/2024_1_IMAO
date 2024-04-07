package responses_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestSendOkResponse(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()

	testResponse := map[string]string{"message": "success"}
	responses.SendOkResponse(w, testResponse)

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

	testResponse := map[string]string{"error": responses.ErrUserAlreadyExists}
	responses.SendErrResponse(responseWriter, testResponse)

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
