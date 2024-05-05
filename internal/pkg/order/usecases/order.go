package usecases

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

type OrderStorageInterface interface {
	//GetOrdersByUserID(userID uint, advertsList usecases.AdvertsStorageInterface) ([]*models.OrderItem, error)
	//GetReturningOrderByUserID(ctx context.Context, userID uint, advertsList usecases.AdvertsStorageInterface) ([]*models.ReturningOrder, error)
	//DeleteAdvByIDs(userID uint, advertID uint, userList UsersInfo, advertsList AdvertsInfo) error
	CreateOrderByID(ctx context.Context, userID uint, data *models.ReceivedOrderItem) error
	GetBoughtOrdersByUserID(ctx context.Context, userID uint) ([]*models.ReturningOrder, error)
}
