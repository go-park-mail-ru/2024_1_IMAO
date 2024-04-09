package storage

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	advuc "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	resp "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/delivery"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
// errNotInCart = errors.New("there is no advert in the cart")
)

type OrderListWrapper struct {
	OrderList *models.OrderList
	Pool      *pgxpool.Pool
	Logger    *zap.SugaredLogger
}

func (ol *OrderListWrapper) GetOrdersByUserID(userID uint, advertsList advuc.AdvertsInfo) ([]*models.OrderItem, error) {
	cart := []*models.OrderItem{}

	for i := range ol.OrderList.Items {
		ol.OrderList.Mux.Lock()
		item := ol.OrderList.Items[i]
		ol.OrderList.Mux.Unlock()

		if item.UserID != userID {
			continue
		}

		cart = append(cart, item)
	}

	return cart, nil
}

func (ol *OrderListWrapper) GetReturningOrderByUserID(userID uint, advertsList advuc.AdvertsInfo) ([]*models.ReturningOrder, error) {
	order := []*models.ReturningOrder{}

	for i := range ol.OrderList.Items {
		ol.OrderList.Mux.Lock()
		item := ol.OrderList.Items[i]
		ol.OrderList.Mux.Unlock()

		if item.UserID != userID {
			continue
		}
		advert, err := advertsList.GetAdvertByOnlyByID(item.AdvertID)

		if err != nil {
			return order, err
		}

		order = append(order, &models.ReturningOrder{
			OrderItem: *item,
			Advert:    advert.Advert,
		})
	}

	return order, nil
}

// func (cl *CartList) DeleteAdvByIDs(userID uint, advertID uint, userList UsersInfo, advertsList AdvertsInfo) error {
// 	for i := range cl.Items {
// 		cl.mu.Lock()
// 		item := cl.Items[i]
// 		cl.mu.Unlock()

// 		if item.UserID != userID || item.AdvertID != advertID {
// 			continue
// 		}
// 		cl.mu.Lock()
// 		cl.Items = append(cl.Items[:i], cl.Items[i+1:]...)
// 		cl.mu.Unlock()
// 		return nil
// 	}

// 	return errNotInCart
// }

func (ol *OrderListWrapper) CreateOrderByID(userID uint, orderItem *models.ReceivedOrderItem, advertsList advuc.AdvertsInfo) error {

	newOrderItem := models.OrderItem{
		ID:            0,
		UserID:        userID,
		AdvertID:      orderItem.AdvertID,
		StatusID:      resp.StatusCreated,
		Created:       time.Now(),
		Updated:       time.Now(),
		Closed:        time.Now(),
		Phone:         orderItem.Phone,
		Name:          orderItem.Name,
		Email:         orderItem.Email,
		Adress:        orderItem.Adress,
		DeliveryPrice: orderItem.DeliveryPrice,
	}
	ol.OrderList.Mux.Lock()
	ol.OrderList.Items = append(ol.OrderList.Items, &newOrderItem)
	ol.OrderList.Mux.Unlock()
	return nil
}

func NewOrderList(pool *pgxpool.Pool, logger *zap.SugaredLogger) *OrderListWrapper {
	return &OrderListWrapper{
		OrderList: &models.OrderList{
			Items: make([]*models.OrderItem, 0),
			Mux:   sync.RWMutex{},
		},
		Pool:   pool,
		Logger: logger,
	}
}
