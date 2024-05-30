//nolint:tagliatelle
package models

import (
	"sync"
	"time"
)

type ReceivedOrderItem struct {
	AdvertID        uint   `json:"advertID"`
	Phone           string `json:"phone"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Adress          string `json:"adress"`
	DeliveryPrice   uint   `json:"deliveryPrice"`
	DeliveryAddress string `json:"address"`
}

type ReceivedOrderItems struct {
	Adverts []*ReceivedOrderItem `json:"adverts"`
}

type OrderCreated struct {
	IsCreated bool `json:"isCreated"`
}

type ReturningOrder struct {
	OrderItem       OrderItem       `json:"orderItem"`
	ReturningAdvert ReturningAdvert `json:"advert"`
}

type OrderItem struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"userId"`
	AdvertID      uint      `json:"advertId"`
	Status        string    `json:"status"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
	Closed        time.Time `json:"closed"`
	Phone         string    `json:"phone"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Address       string    `json:"address"`
	DeliveryPrice uint      `json:"deliveryPrice"`
}

type ReviewItem struct {
	AdvertID uint `json:"advertId"`
	Rating   uint `json:"rating"`
}

type OrderList struct {
	Items []*OrderItem
	Mux   sync.RWMutex
}
