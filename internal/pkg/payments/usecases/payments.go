package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type PaymentsStorageInterface interface {
	CheckAdvertOwnership(ctx context.Context, advertId, userId uint) bool
	GetPriceAndDescription(ctx context.Context, advertId, rateCode uint) (*models.PriceAndDescription, error)
	CreatePayment(ctx context.Context, payment *models.Payment, idempotencyKey string, advertId uint, duration string) error
	//GetPaymentForm(ctx context.Context) (*models.CityList, error)
}
