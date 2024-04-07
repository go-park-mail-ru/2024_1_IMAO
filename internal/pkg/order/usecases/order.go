package usecases

import (
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
)

type OrderInfo interface {
	GetOrdersByUserID(userID uint, advertsList usecases.AdvertsInfo) ([]*models.OrderItem, error)
	GetReturningOrderByUserID(userID uint, advertsList usecases.AdvertsInfo) ([]*models.ReturningOrder, error)
	//DeleteAdvByIDs(userID uint, advertID uint, userList UsersInfo, advertsList AdvertsInfo) error
	CreateOrderByID(userID uint, orderItem *models.ReceivedOrderItem, advertsList usecases.AdvertsInfo) bool
}
