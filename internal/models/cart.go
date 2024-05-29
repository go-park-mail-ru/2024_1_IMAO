package models

import (
	"sync"
)

type ReceivedCartItem struct {
	// UserID   uint `json:"userID"`
	AdvertID uint `json:"advertId"`
}

type ReceivedCartItems struct {
	// UserID    uint    `json:"userID"`
	AdvertIDs []uint `json:"advertIDs"`
}

type CartItem struct {
	UserID   uint
	AdvertID uint
}

type CartList struct {
	Items []*CartItem
	Mux   sync.RWMutex
}

type Appended struct {
	IsAppended bool `json:"isAppended"`
}
