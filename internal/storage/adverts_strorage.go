package storage

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
)

var (
	errWrongID            = errors.New("wrong adverts ID")
	errWrongAdvertsAmount = errors.New("too many elements specified")
)

const (
	maxPrice = 1000
)

type Image struct{}

type Advert struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	Image       Image  `json:"image"`
	Location    string `json:"location"`
}

type AdvertsList struct {
	Adverts        []*Advert
	AdvertsCounter uint
	mu             sync.RWMutex
}

type AdvertsInfo interface {
	GetAdvert(advertID uint) (*Advert, error)
	GetSeveralAdverts(number uint) (*[]Advert, error)
	GetLastID() uint
}

func (ads *AdvertsList) GetAdvert(advertID uint) (*Advert, error) {
	if advertID > ads.AdvertsCounter {
		return nil, errWrongID
	}

	ads.mu.Lock()
	defer ads.mu.Unlock()

	return ads.Adverts[advertID], nil
}

func (ads *AdvertsList) GetSeveralAdverts(number uint) ([]*Advert, error) {
	if number > ads.AdvertsCounter {
		return nil, errWrongAdvertsAmount
	}

	ads.mu.Lock()
	defer ads.mu.Unlock()

	returningAds := make([]*Advert, number)

	for ind := 0; ind < int(number); ind++ {
		returningAds[ind] = &Advert{
			ID:          ads.Adverts[ind].ID,
			UserID:      ads.Adverts[ind].UserID,
			Title:       ads.Adverts[ind].Title,
			Description: ads.Adverts[ind].Description,
			Price:       ads.Adverts[ind].Price,
			Image:       ads.Adverts[ind].Image,
			Location:    ads.Adverts[ind].Location,
		}
	}

	return returningAds, nil
}

func (ads *AdvertsList) GetLastID() uint {
	ads.AdvertsCounter++

	return ads.AdvertsCounter
}

func NewAdvertsList() *AdvertsList {
	return &AdvertsList{
		AdvertsCounter: 0,
		Adverts:        make([]*Advert, 0),
		mu:             sync.RWMutex{},
	}
}

func FillAdvertsList(ads *AdvertsList) {
	for i := 1; i <= 60; i++ {
		price, _ := rand.Int(rand.Reader, big.NewInt(int64(maxPrice)))
		advertID := ads.GetLastID()
		ads.Adverts = append(ads.Adverts, &Advert{
			ID:          advertID,
			UserID:      1,
			Title:       fmt.Sprintf("Объявление № %d", advertID),
			Image:       Image{},
			Description: fmt.Sprintf("Текст в объявлениии № %d", advertID),
			Price:       uint(price.Uint64()) * advertID,
			Location:    "Москва",
		})
	}
}

func AddAdvert(ads *AdvertsList, advert *Advert) {

	ads.Adverts = append(ads.Adverts, advert)

}
