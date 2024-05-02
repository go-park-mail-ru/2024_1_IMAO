package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
)

type OrderStorageInterface interface {
	GetOrdersByUserID(userID uint) ([]*models.OrderItem, error)
	GetReturningOrderByUserID(ctx context.Context, userID uint,
		advertsList usecases.AdvertsStorageInterface) ([]*models.ReturningOrder, error)
	// DeleteAdvByIDs(userID uint, advertID uint, userList UsersInfo, advertsList AdvertsInfo) error
	CreateOrderByID(userID uint, orderItem *models.ReceivedOrderItem) error
}
