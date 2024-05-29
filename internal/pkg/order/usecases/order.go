package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type OrderStorageInterface interface {
	CreateOrderByID(ctx context.Context, userID uint, data *models.ReceivedOrderItem) error
	GetBoughtOrdersByUserID(ctx context.Context, userID uint) ([]*models.ReturningOrder, error)
	GetSoldOrdersByUserID(ctx context.Context, userID uint) ([]*models.ReturningOrder, error)
}
