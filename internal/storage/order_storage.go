package storage

import (
	"sync"
	"time"
)

const (
	_ = iota
	StatusCreated
	StatusInDelivery
	StatusDelivered
)

var (
// errNotInCart = errors.New("there is no advert in the cart")
)

type ReceivedOrderItem struct {
	AdvertID      uint   `json:"advertID"`
	Phone         string `json:"phone"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Adress        string `json:"adress"`
	DeliveryPrice uint   `json:"deliveryPrice"`
}

type ReceivedOrderItems struct {
	Adverts []*ReceivedOrderItem `json:"adverts"`
}

// type ReturningOrderItem struct {
// 	AdvertID      uint
// 	StatusID      uint
// 	Adress        string
// 	DeliveryPrice uint
// }

type ReturningOrder struct {
	OrderItem OrderItem `json:"orderItem"`
	Advert    Advert    `json:"advert"`
}

type OrderItem struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"userId"`
	AdvertID      uint      `json:"advertId"`
	StatusID      uint      `json:"statusId"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
	Closed        time.Time `json:"closed"`
	Phone         string    `json:"phone"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Adress        string    `json:"adress"`
	DeliveryPrice uint      `json:"deliveryPrice"`
}

type OrderList struct {
	Items []*OrderItem
	mu    sync.RWMutex
}

type OrderInfo interface {
	GetOrdersByUserID(userID uint, userList UsersInfo, advertsList AdvertsInfo) ([]*OrderItem, error)
	GetReturningOrderByUserID(userID uint, userList UsersInfo, advertsList AdvertsInfo) ([]*ReturningOrder, error)
	//DeleteAdvByIDs(userID uint, advertID uint, userList UsersInfo, advertsList AdvertsInfo) error
	CreateOrderByIDs(userID uint, orderItem ReceivedOrderItem, userList UsersInfo, advertsList AdvertsInfo) bool
}

func (ol *OrderList) GetOrdersByUserID(userID uint, userList UsersInfo, advertsList AdvertsInfo) ([]*OrderItem, error) {
	cart := []*OrderItem{}

	for i := range ol.Items {
		ol.mu.Lock()
		item := ol.Items[i]
		ol.mu.Unlock()

		if item.UserID != userID {
			continue
		}

		cart = append(cart, item)
	}

	return cart, nil
}

func (ol *OrderList) GetReturningOrderByUserID(userID uint, userList UsersInfo, advertsList AdvertsInfo) ([]*ReturningOrder, error) {
	cart := []*ReturningOrder{}

	for i := range ol.Items {
		ol.mu.Lock()
		item := ol.Items[i]
		ol.mu.Unlock()

		if item.UserID != userID {
			continue
		}
		advert, err := advertsList.GetAdvert(item.AdvertID)

		if err != nil {
			return cart, err
		}

		cart = append(cart, &ReturningOrder{
			OrderItem: *item,
			Advert:    advert.Advert,
		})
	}

	return cart, nil
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

func (ol *OrderList) CreateOrderByIDs(userID uint, orderItem ReceivedOrderItem, userList UsersInfo, advertsList AdvertsInfo) bool {

	newOrderItem := OrderItem{
		ID:            0,
		UserID:        userID,
		AdvertID:      orderItem.AdvertID,
		StatusID:      StatusCreated,
		Created:       time.Now(),
		Updated:       time.Now(),
		Closed:        time.Now(),
		Phone:         orderItem.Phone,
		Name:          orderItem.Name,
		Email:         orderItem.Email,
		Adress:        orderItem.Adress,
		DeliveryPrice: orderItem.DeliveryPrice,
	}
	ol.mu.Lock()
	ol.Items = append(ol.Items, &newOrderItem)
	ol.mu.Unlock()
	return true
}

func NewOrderList() *OrderList {
	return &OrderList{
		Items: make([]*OrderItem, 0),
		mu:    sync.RWMutex{},
	}
}
