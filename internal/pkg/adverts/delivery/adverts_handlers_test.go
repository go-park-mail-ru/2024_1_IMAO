package delivery_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

//nolint:funlen
func TestGetPostsListHandlerSuccessful(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name             string
		inputParamCount  int
		handler          *handler.AdvertsHandler
		postsForStorage  []storage.Advert
		expectedResponse responses.AdvertsOkResponse
	}

	testCases := [...]TestCase{
		{
			name:            "test basic work",
			inputParamCount: 1,
			handler:         &handler.AdvertsHandler{List: storage.NewAdvertsList()},
			postsForStorage: []storage.Advert{{
				ID:          1,
				UserID:      1,
				Title:       "Test Title",
				Description: "Test Description",
				Price:       1337,
				Image:       storage.Image{},
				Location:    "Moscow",
			}},
			expectedResponse: responses.AdvertsOkResponse{
				Code: responses.StatusOk,
				Adverts: []*storage.Advert{{
					ID:          1,
					UserID:      1,
					Title:       "Test Title",
					Description: "Test Description",
					Price:       1337,
					Image:       storage.Image{},
					Location:    "Moscow",
				}},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			storage.FillAdvertsList(testCase.handler.List)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			testCase.handler.Root(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse responses.AdvertsOkResponse

			err = json.Unmarshal(receivedResponse, &resultResponse)
			if err != nil {
				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
			}

			if resultResponse.Code != responses.StatusOk {
				t.Errorf("wrong Response code: got %+v, expected %+v",
					resultResponse.Code, responses.StatusOk)
			}
		})
	}
}

//nolint:funlen
func TestGetPostsListHandlerUnuccessful(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name             string
		inputParamCount  int
		handler          *handler.AdvertsHandler
		postsForStorage  []storage.Advert
		expectedResponse responses.AdvertsOkResponse
	}

	testCases := [...]TestCase{
		{
			name:            "test basic work",
			inputParamCount: 1,
			handler:         &handler.AdvertsHandler{List: storage.NewAdvertsList()},
			postsForStorage: []storage.Advert{{
				ID:          1,
				UserID:      1,
				Title:       "Test Title",
				Description: "Test Description",
				Price:       1337,
				Image:       storage.Image{},
				Location:    "Moscow",
			}},
			expectedResponse: responses.AdvertsOkResponse{
				Code: responses.StatusOk,
				Adverts: []*storage.Advert{{
					ID:          1,
					UserID:      1,
					Title:       "Test Title",
					Description: "Test Description",
					Price:       1337,
					Image:       storage.Image{},
					Location:    "Moscow",
				}},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			storage.AddAdvert(testCase.handler.List, &testCase.postsForStorage[0])

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			testCase.handler.Root(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			receivedResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to ReadAll resp.Body: %v", err)
			}

			var resultResponse responses.AdvertsOkResponse

			err = json.Unmarshal(receivedResponse, &resultResponse)
			if err != nil {
				t.Fatalf("Failed to Unmarshal(receivedResponse): %v", err)
			}

			if resultResponse.Code != responses.StatusBadRequest {
				t.Errorf("wrong Response code: got %+v, expected %+v",
					resultResponse.Code, responses.StatusBadRequest)
			}
		})
	}
}
