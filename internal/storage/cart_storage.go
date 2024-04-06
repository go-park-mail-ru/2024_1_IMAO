package storage

import (
	"errors"
	"sync"
)

var (
	errNotInCart = errors.New("there is no advert in the cart")
)

type ReceivedCartItem struct {
	// UserID   uint `json:"userID"`
	AdvertID uint `json:"advertID"`
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
	mu    sync.RWMutex
}

type CartInfo interface {
	GetCartByUserID(userID uint, userList UsersInfo, advertsList AdvertsInfo) ([]*Advert, error)
	DeleteAdvByIDs(userID uint, advertID uint, userList UsersInfo, advertsList AdvertsInfo) error
	AppendAdvByIDs(userID uint, advertID uint, userList UsersInfo, advertsList AdvertsInfo) bool
}

func (cl *CartList) GetCartByUserID(userID uint, userList UsersInfo, advertsList AdvertsInfo) ([]*Advert, error) {
	cart := []*Advert{}

	for i := range cl.Items {
		cl.mu.Lock()
		item := cl.Items[i]
		cl.mu.Unlock()

		if item.UserID != userID {
			continue
		}
		advert, err := advertsList.GetAdvertByOnlyByID(item.AdvertID)

		if err != nil {
			return cart, err
		}

		cart = append(cart, &advert.Advert)
	}

	return cart, nil
}

func (cl *CartList) DeleteAdvByIDs(userID uint, advertID uint, userList UsersInfo, advertsList AdvertsInfo) error {
	for i := range cl.Items {
		cl.mu.Lock()
		item := cl.Items[i]
		cl.mu.Unlock()

		if item.UserID != userID || item.AdvertID != advertID {
			continue
		}
		cl.mu.Lock()
		cl.Items = append(cl.Items[:i], cl.Items[i+1:]...)
		cl.mu.Unlock()
		return nil
	}

	return errNotInCart
}

func (cl *CartList) AppendAdvByIDs(userID uint, advertID uint, userList UsersInfo, advertsList AdvertsInfo) bool {
	for i := range cl.Items {
		cl.mu.Lock()
		item := cl.Items[i]
		cl.mu.Unlock()

		if item.UserID != userID || item.AdvertID != advertID {
			continue
		}
		cl.mu.Lock()
		cl.Items = append(cl.Items[:i], cl.Items[i+1:]...)
		cl.mu.Unlock()
		return false
	}
	cartItem := CartItem{userID, advertID}
	cl.mu.Lock()
	cl.Items = append(cl.Items, &cartItem)
	cl.mu.Unlock()
	return true
}

func NewCartList() *CartList {
	return &CartList{
		Items: make([]*CartItem, 0),
		mu:    sync.RWMutex{},
	}
}
