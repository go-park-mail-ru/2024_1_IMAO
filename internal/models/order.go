package models

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
	Adress        string    `json:"adress"`
	DeliveryPrice uint      `json:"deliveryPrice"`
}

type OrderList struct {
	Items []*OrderItem
	Mux   sync.RWMutex
}
