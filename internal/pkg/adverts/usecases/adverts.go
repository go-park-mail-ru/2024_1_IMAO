//nolint:interfacebloat
package usecases

import (
	"context"
	"mime/multipart"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

//go:generate mockgen -source=adverts.go -destination=mocks/mock.go

type AdvertsStorageInterface interface {
	GetAdvert(ctx context.Context, userID, advertID uint) (*models.ReturningAdvert, error)
	GetAdvertsByCity(ctx context.Context, city string, userID, startID,
		num uint) ([]*models.ReturningAdInList, error)
	GetAdvertsByCategory(ctx context.Context, category, city string, userID, startID,
		num uint) ([]*models.ReturningAdInList, error)
	GetAdvertOnlyByID(ctx context.Context, advertID uint) (*models.ReturningAdvert, error)
	SearchAdvertByTitle(ctx context.Context, title string, userID, startID,
		num uint) ([]*models.ReturningAdInList, error)
	GetSuggestions(ctx context.Context, title string, num uint) ([]string, error)
	GetPriceHistory(ctx context.Context, userID uint) ([]*models.PriceHistoryItem, error)
	CheckAdvertOwnership(ctx context.Context, advertID, userID uint) bool
	GetPaymnetUUIDList(ctx context.Context, advertID uint) (*models.PaymnetUUIDList, error)
	YuKassaUpdateDB(ctx context.Context, paymentList *models.PaymentList, advertID uint) error
	GetPromotionData(ctx context.Context, advertID uint) (*models.Promotion, error)

	CreateAdvert(ctx context.Context, files []*multipart.FileHeader,
		data models.ReceivedAdData) (*models.ReturningAdvert, error)
	EditAdvert(ctx context.Context, files []*multipart.FileHeader,
		data models.ReceivedAdData) (*models.ReturningAdvert, error)
	GetAdvertsForUserWhereStatusIs(ctx context.Context, userID, authorID, deleted,
		advertNum uint) ([]*models.ReturningAdInList, error)
	CloseAdvert(ctx context.Context, advertID uint) error
	InsertView(ctx context.Context, userID, advertID uint) error
}
