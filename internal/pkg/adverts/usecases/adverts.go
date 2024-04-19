package usecases

import (
	"context"
	"mime/multipart"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type AdvertsInfo interface {
	GetAdvert(ctx context.Context, advertID uint, city, category string) (*models.ReturningAdvert, error)
	GetAdvertsByCity(ctx context.Context, city string, startID, num uint) ([]*models.ReturningAdInList, error)
	GetAdvertsByCategory(ctx context.Context, category, city string, startID, num uint) ([]*models.ReturningAdInList, error)
	GetAdvertByOnlyByID(ctx context.Context, advertID uint) (*models.ReturningAdvert, error)

	CreateAdvert(ctx context.Context, files []*multipart.FileHeader, data models.ReceivedAdData) (*models.ReturningAdvert, error)

	GetCityID(city string) (uint, error)
	GetCategoryID(city string) (uint, error)

	GetLastAdvertID() uint
	GetLastLocationID() uint
	GetLastCategoryID() uint
}
