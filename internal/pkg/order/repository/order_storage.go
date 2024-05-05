package storage

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
// errNotInCart = errors.New("there is no advert in the cart")
)

type OrderStorage struct {
	pool *pgxpool.Pool
}

func NewOrderStorage(pool *pgxpool.Pool) *OrderStorage {
	return &OrderStorage{
		pool: pool,
	}
}

// func (ol *OrderStorage) GetOrdersByUserID(userID uint, advertsList advuc.AdvertsStorageInterface) ([]*models.OrderItem, error) {
// 	cart := []*models.OrderItem{}

// 	for i := range ol.OrderList.Items {
// 		ol.OrderList.Mux.Lock()
// 		item := ol.OrderList.Items[i]
// 		ol.OrderList.Mux.Unlock()

// 		if item.UserID != userID {
// 			continue
// 		}

// 		cart = append(cart, item)
// 	}

// 	return cart, nil
// }

// func (ol *OrderStorage) GetReturningOrderByUserID(ctx context.Context, userID uint, advertsList advuc.AdvertsStorageInterface) ([]*models.ReturningOrder, error) {
// 	order := []*models.ReturningOrder{}

// 	for i := range ol.OrderList.Items {
// 		ol.OrderList.Mux.Lock()
// 		item := ol.OrderList.Items[i]
// 		ol.OrderList.Mux.Unlock()

// 		if item.UserID != userID {
// 			continue
// 		}
// 		advert, err := advertsList.GetAdvertOnlyByID(ctx, item.AdvertID)

// 		if err != nil {
// 			return order, err
// 		}

// 		order = append(order, &models.ReturningOrder{
// 			OrderItem: *item,
// 			Advert:    advert.Advert,
// 		})
// 	}

// 	return order, nil
// }

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

func (ol *OrderStorage) сreateOrderByID(ctx context.Context, tx pgx.Tx, userID uint, data *models.ReceivedOrderItem) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCreateOrder :=
		`INSERT INTO public."order"(
			user_id, advert_id, order_status, phone, name, surname, patronymic, email, delivery_price)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

	logging.LogInfo(logger, "INSERT INTO user")

	var err error

	const paidStatus string = "Оплачено"
	const surnamePlug string = "Фамилия"
	const patronymicPlug string = "Отчество"

	_, err = tx.Exec(ctx, SQLCreateOrder, userID, data.AdvertID, paidStatus, data.Phone, data.Name, surnamePlug, patronymicPlug, data.Email, data.DeliveryPrice)

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing create order query, err=%v", err))

		return err
	}

	return nil
}

func (ol *OrderStorage) CreateOrderByID(ctx context.Context, userID uint, data *models.ReceivedOrderItem) error {

	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, ol.pool, func(tx pgx.Tx) error {
		err := ol.сreateOrderByID(ctx, tx, userID, data)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while creating order, err=%v", err))

			return err
		}

		return nil
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while creating user, err=%v", err))

		return err
	}

	return nil
}

// func NewOrderList(pool *pgxpool.Pool, logger *zap.SugaredLogger) *models.OrderList {
// 	return &models.OrderList{
// 		Items: make([]*models.OrderItem, 0),
// 		Mux:   sync.RWMutex{},
// 	}
// }
