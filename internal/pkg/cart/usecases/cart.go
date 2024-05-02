package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type CartStorageInterface interface {
	GetCartByUserID(ctx context.Context, userID uint) ([]*models.ReturningAdvert, error)
	DeleteAdvByIDs(ctx context.Context, userID uint, advertID uint) error
	AppendAdvByIDs(ctx context.Context, userID uint, advertID uint) bool
}
