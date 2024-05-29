package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type PaymentsStorageInterface interface {
	CheckAdvertOwnership(ctx context.Context, advertID, userID uint) bool
	GetPriceAndDescription(ctx context.Context, advertID, rateCode uint) (*models.PriceAndDescription, error)
	CreatePayment(ctx context.Context, payment *models.Payment, idempotencyKey string,
		advertID uint, duration string) error
}
