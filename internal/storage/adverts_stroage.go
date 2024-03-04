package storage

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

var (
	errWrongID            = errors.New("wrong adverts ID")
	errWrongAdvertsAmount = errors.New("too many elements specified")
)

type Image struct {
}

type Advert struct {
	ID          uint   `json:"ID"`
	UserID      uint   `json:"userID"`
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

	for i := 0; i < int(number); i++ {
		returningAds[i] = &Advert{
			ID:          ads.Adverts[i].ID,
			UserID:      ads.Adverts[i].UserID,
			Title:       ads.Adverts[i].Title,
			Description: ads.Adverts[i].Description,
			Price:       ads.Adverts[i].Price,
			Image:       ads.Adverts[i].Image,
			Location:    ads.Adverts[i].Location,
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
		advertID := ads.GetLastID()
		ads.Adverts = append(ads.Adverts, &Advert{
			ID:          advertID,
			UserID:      1,
			Title:       fmt.Sprintf("Объявление № %d", advertID),
			Image:       Image{},
			Description: fmt.Sprintf("Текст в объявлениии № %d", advertID),
			Price:       uint(rand.Intn(1000)) * advertID,
			Location:    "Москва",
		})
	}
}
