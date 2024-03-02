package adverts

import (
	"errors"
	"fmt"
	"math/rand"
)

var (
	errWrongID            = errors.New("wrong adverts ID")
	errWrongAdvertsAmount = errors.New("too many elements specified")
)

type Image struct {
}

type Advert struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	Image       Image  `json:"image"`
	Location    string `json:"location"`
}

type AdvertsStorage struct {
	Adverts        []*Advert
	AdvertsCounter uint
}

type AdvertsInfo interface {
	GetAdvert(advertID uint) (*Advert, error)
	GetSeveralAdverts(number uint) (*[]Advert, error)
	GetLastID() uint
}

func (ads *AdvertsStorage) GetAdvert(advertID uint) (*Advert, error) {
	if advertID > ads.AdvertsCounter {
		return nil, errWrongID
	}

	return ads.Adverts[advertID], nil
}

func (ads *AdvertsStorage) GetSeveralAdverts(number uint) ([]*Advert, error) {
	if number > ads.AdvertsCounter {
		return nil, errWrongAdvertsAmount
	}

	returningAds := make([]*Advert, number)

	for i := 1; i <= int(number); i++ {
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

func (ads *AdvertsStorage) GetLastID() uint {
	ads.AdvertsCounter++

	return ads.AdvertsCounter
}

func NewAdvertsStorage() *AdvertsStorage {
	return &AdvertsStorage{
		AdvertsCounter: 1,
		Adverts:        make([]*Advert, 1),
	}
}

func FillAdvertsStorage(ads *AdvertsStorage) {
	for i := 1; i <= 60; i++ {
		advertID := ads.GetLastID()
		ads.Adverts[advertID] = &Advert{
			ID:          advertID,
			UserID:      1,
			Title:       fmt.Sprintf("Объявление № %d", advertID),
			Image:       Image{},
			Description: fmt.Sprintf("Текст в объявлениии № %d", advertID),
			Price:       uint(rand.Intn(1000)) * advertID,
			Location:    "Москва",
		}
	}
}
