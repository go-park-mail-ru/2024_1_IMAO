package storage

import (
	"errors"
	"sync"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	advuc "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	useruc "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
)

var (
	errNotInCart = errors.New("there is no advert in the cart")
)

type CartListWrapper struct {
	CartList *models.CartList
}

func (cl *CartListWrapper) GetCartByUserID(userID uint, userList useruc.UsersInfo, advertsList advuc.AdvertsInfo) ([]*models.Advert, error) {
	cart := []*models.Advert{}

	for i := range cl.CartList.Items {
		cl.CartList.Mux.Lock()
		item := cl.CartList.Items[i]
		cl.CartList.Mux.Unlock()

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

func (cl *CartListWrapper) DeleteAdvByIDs(userID uint, advertID uint, userList useruc.UsersInfo, advertsList advuc.AdvertsInfo) error {
	for i := range cl.CartList.Items {
		cl.CartList.Mux.Lock()
		item := cl.CartList.Items[i]
		cl.CartList.Mux.Unlock()

		if item.UserID != userID || item.AdvertID != advertID {
			continue
		}
		cl.CartList.Mux.Lock()
		cl.CartList.Items = append(cl.CartList.Items[:i], cl.CartList.Items[i+1:]...)
		cl.CartList.Mux.Unlock()
		return nil
	}

	return errNotInCart
}

func (cl *CartListWrapper) AppendAdvByIDs(userID uint, advertID uint, userList useruc.UsersInfo, advertsList advuc.AdvertsInfo) bool {
	for i := range cl.CartList.Items {
		cl.CartList.Mux.Lock()
		item := cl.CartList.Items[i]
		cl.CartList.Mux.Unlock()

		if item.UserID != userID || item.AdvertID != advertID {
			continue
		}
		cl.CartList.Mux.Lock()
		cl.CartList.Items = append(cl.CartList.Items[:i], cl.CartList.Items[i+1:]...)
		cl.CartList.Mux.Unlock()
		return false
	}
	cartItem := models.CartItem{
		UserID:   userID,
		AdvertID: advertID,
	}
	cl.CartList.Mux.Lock()
	cl.CartList.Items = append(cl.CartList.Items, &cartItem)
	cl.CartList.Mux.Unlock()
	return true
}

func NewCartList() *CartListWrapper {
	return &CartListWrapper{
		CartList: &models.CartList{
			Items: make([]*models.CartItem, 0),
			Mux:   sync.RWMutex{},
		},
	}
}
