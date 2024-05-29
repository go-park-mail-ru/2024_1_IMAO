package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

//go:generate mockgen -source=city.go -destination=../mocks/city_mocks.go

type CityStorageInterface interface {
	GetCityList(ctx context.Context) (*models.CityList, error)
}
