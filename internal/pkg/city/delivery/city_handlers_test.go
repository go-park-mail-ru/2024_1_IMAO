package delivery_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	delivery "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/delivery"
	mock_city "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/city/mocks"
)

func TestGetCityList(t *testing.T) {

	tests := []struct {
		name           string
		siMocker       func(ctx context.Context, si *mock_city.MockCityStorageInterface)
		expectedStatus int
	}{
		{
			name: "Bad_Request",
			siMocker: func(ctx context.Context, si *mock_city.MockCityStorageInterface) {
				si.EXPECT().GetCityList(ctx).Return(nil, errors.New("error"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Test_Success_1",
			siMocker: func(ctx context.Context, si *mock_city.MockCityStorageInterface) {
				si.EXPECT().GetCityList(ctx).Return(&models.CityList{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			storageInterface := mock_city.NewMockCityStorageInterface(ctrl)
			defer ctrl.Finish()
			req := httptest.NewRequest("GET", "http://example.com/api/handler/", bytes.NewBufferString(""))

			code := new(int)
			*code = 400

			ctx := context.WithValue(req.Context(), "code", code)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			tt.siMocker(ctx, storageInterface)

			h := delivery.NewCityHandler(storageInterface)

			// if tt.name == "Bad_Request" {

			// 	ctx = context.WithValue(req.Context(), )

			// 	req = req.WithContext(ctx)
			// }

			h.GetCityList(w, req)

			var testResponse models.ErrResponse
			_ = json.NewDecoder(w.Body).Decode(&testResponse)

			fmt.Println(testResponse)

			assert.Equal(t, tt.expectedStatus, testResponse.Code)
		})
	}
}
