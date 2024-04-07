package usecases

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type AdvertsInfo interface {
	GetAdvert(advertID uint, city, category string) (*models.ReturningAdvert, error)
	GetAdvertsByCity(city string, startID, num uint) ([]*models.ReturningAdInList, error)
	GetAdvertsByCategory(category, city string, startID, num uint) ([]*models.ReturningAdInList, error)
	GetAdvertByOnlyByID(advertID uint) (*models.ReturningAdvert, error)

	CreateAdvert(data models.ReceivedAdData) (*models.ReturningAdvert, error)

	GetCityID(city string) (uint, error)
	GetCategoryID(city string) (uint, error)

	GetLastAdvertID() uint
	GetLastLocationID() uint
	GetLastCategoryID() uint
}
