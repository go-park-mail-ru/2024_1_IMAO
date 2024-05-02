package storage

import (
	"context"
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

type OrderStorage struct {
	pool      *pgxpool.Pool
	OrderList *models.OrderList
}

func NewOrderStorage(pool *pgxpool.Pool) *OrderStorage {
	return &OrderStorage{
		pool: pool,
	}
}

func (ol *OrderStorage) GetOrdersByUserID(userID uint) ([]*models.OrderItem, error) {
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

func (ol *OrderStorage) GetReturningOrderByUserID(ctx context.Context, userID uint,
	advertsList advuc.AdvertsStorageInterface) ([]*models.ReturningOrder, error) {
	order := []*models.ReturningOrder{}

	for i := range ol.OrderList.Items {
		ol.OrderList.Mux.Lock()
		item := ol.OrderList.Items[i]
		ol.OrderList.Mux.Unlock()

		if item.UserID != userID {
			continue
		}
		advert, err := advertsList.GetAdvertByOnlyByID(ctx, item.AdvertID)

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

func (ol *OrderStorage) CreateOrderByID(userID uint, orderItem *models.ReceivedOrderItem) error {

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

func NewOrderList(pool *pgxpool.Pool, logger *zap.SugaredLogger) *models.OrderList {
	return &models.OrderList{
		Items: make([]*models.OrderItem, 0),
		Mux:   sync.RWMutex{},
	}
}
