package storage

import "sync"

type ReceivedCartItem struct {
	UserID   uint `json:"userID"`
	AdvertID uint `json:"advertID"`
}

type CartItem struct {
	UserID   uint
	AdvertID uint
	deleted  bool
}

type CartList struct {
	Items []*CartItem
	mu    sync.RWMutex
}

type CartInfo interface {
	GetCartByUserID(userID uint, userList UsersInfo, advertsList AdvertsInfo) ([]*Advert, error)
	DeleteAdvByIDs(data ReceivedCartItem, userList UsersInfo, advertsList AdvertsInfo) (bool, error)
	AppendAdvByIDs(data ReceivedCartItem, userList UsersInfo, advertsList AdvertsInfo) (bool, error)
}

func (cl *CartList) GetCartByUserID(userID uint, userList UsersInfo, advertsList AdvertsInfo) ([]*Advert, error) {
	cart := []*Advert{}

	for i := range cl.Items {
		cl.mu.Lock()
		item := cl.Items[i]
		cl.mu.Unlock()

		if item.UserID != userID || item.deleted {
			continue
		}
		advert, err := advertsList.GetAdvert(item.AdvertID)

		if err != nil {
			return cart, err
		}

		cart = append(cart, &advert.Advert)
	}

	return cart, nil
}

func (cl *CartList) DeleteAdvByIDs(data ReceivedCartItem, userList UsersInfo, advertsList AdvertsInfo) (bool, error) {
	for i := range cl.Items {
		cl.mu.Lock()
		item := cl.Items[i]
		cl.mu.Unlock()

		if item.UserID != data.UserID || item.AdvertID != data.AdvertID || item.deleted {
			continue
		}

		item.deleted = true

		return true, nil
	}

	return true, nil
}

func (cl *CartList) AppendAdvByIDs(data ReceivedCartItem, userList UsersInfo, advertsList AdvertsInfo) (bool, error) {
	cartItem := CartItem{data.UserID, data.AdvertID, false}
	cl.mu.Lock()
	cl.Items = append(cl.Items, &cartItem)
	cl.mu.Unlock()
	return true, nil
}

func NewCartList() *CartList {
	return &CartList{
		Items: make([]*CartItem, 0),
		mu:    sync.RWMutex{},
	}
}
