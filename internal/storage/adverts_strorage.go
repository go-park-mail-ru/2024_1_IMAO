package storage

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/mdigger/translit"
)

var (
	errWrongAdvertID      = errors.New("wrong advert ID")
	errWrongCityName      = errors.New("wrong city name")
	errWrongCategoryName  = errors.New("wrong category name")
	errWrongAdvertsAmount = errors.New("too many elements specified")
)

const (
	maxPrice = 1000
)

type Image struct{}

type GettingAdsData struct {
	Count   uint   `json:"count"`
	City    string `json:"city"`
	StartID uint   `json:"startId"`
}

type ReceivedAdData struct {
	UserID      uint   `json:"userId"`
	City        string `json:"city"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	Image       Image  `json:"image"`
	IsUsed      bool   `json:"isUsed"`
}

type Category struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Translation string `json:"translation"`
}

type City struct {
	ID          uint   `json:"id"`
	CityName    string `json:"cityName"`
	Translation string `json:"translation"`
}

type Advert struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"userId"`
	CityID      uint      `json:"cityId"`
	CategoryID  uint      `json:"categoryId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       uint      `json:"price"`
	Created     time.Time `json:"created"`
	Image       Image     `json:"image"`
	Closed      time.Time `json:"closed"`
	Active      bool      `json:"active"`
	IsUsed      bool      `json:"isUsed"`
}

type ReturningAdvert struct {
	Advert   Advert   `json:"advert"`
	City     City     `json:"city"`
	Category Category `json:"category"`
}

type AdvertsList struct {
	Adverts           []*Advert
	Categories        []*Category
	Cities            []*City
	AdvertsCounter    uint
	CitiesCounter     uint
	CategoriesCounter uint
	mu                sync.RWMutex
}

type AdvertsInfo interface {
	GetAdvert(advertID uint) (*ReturningAdvert, error)
	GetAdvertsByCity(city string, startID, number uint) ([]*ReturningAdvert, error)
	GetAdvertsByCategory(category, city string, startID, number uint) ([]*ReturningAdvert, error)

	CreateAdvert(data ReceivedAdData) ([]*ReturningAdvert, error)

	getCityID(city string) (uint, error)
	getCategoryID(city string) (uint, error)

	getLastAdvertID() uint
	getLastLocationID() uint
	getLastCategoryID() uint
}

func (ads *AdvertsList) GetAdvert(advertID uint) (*ReturningAdvert, error) {
	ads.mu.Lock()
	defer ads.mu.Unlock()

	if advertID > ads.AdvertsCounter {
		return nil, errWrongAdvertID
	}

	cityID := ads.Adverts[advertID-1].CityID
	categoryID := ads.Adverts[advertID-1].CategoryID

	return &ReturningAdvert{
		Advert:   *ads.Adverts[advertID-1],
		City:     *ads.Cities[cityID-1],
		Category: *ads.Categories[categoryID-1],
	}, nil
}

func (ads *AdvertsList) GetAdvertsByCity(city string, number, startID uint) ([]*ReturningAdvert, error) {
	if number > ads.AdvertsCounter {
		return nil, errWrongAdvertsAmount
	}

	cityID, err := ads.getCityID(city)
	if err != nil {
		return nil, err
	}

	ads.mu.Lock()
	defer ads.mu.Unlock()

	var returningAds []*ReturningAdvert
	var counter uint = 0

	for counter != number && counter+startID-1 != ads.AdvertsCounter {
		ad := ads.Adverts[startID+counter-1]

		if ad.Active && ad.CityID == cityID {
			returningAds = append(returningAds, &ReturningAdvert{
				Advert:   *ad,
				City:     *ads.Cities[cityID-1],
				Category: *ads.Categories[ad.CategoryID-1],
			})
		}

		counter++
	}

	return returningAds, nil
}

func (ads *AdvertsList) GetAdvertsByCategory(city, category string, number, startID uint) ([]*ReturningAdvert, error) {
	if number > ads.AdvertsCounter {
		return nil, errWrongAdvertsAmount
	}

	cityID, err := ads.getCityID(city)
	if err != nil {
		return nil, err
	}

	categoryID, err := ads.getCategoryID(category)
	if err != nil {
		return nil, err
	}

	ads.mu.Lock()
	defer ads.mu.Unlock()

	var returningAds []*ReturningAdvert
	var counter uint = 0

	for counter != number && counter+startID-1 != ads.AdvertsCounter {
		ad := ads.Adverts[startID+counter-1]

		if ad.Active && ad.CityID == cityID && ad.CategoryID == categoryID {
			returningAds = append(returningAds, &ReturningAdvert{
				Advert:   *ad,
				City:     *ads.Cities[cityID-1],
				Category: *ads.Categories[categoryID-1],
			})
		}

		counter++
	}

	return returningAds, nil
}

func (ads *AdvertsList) CreateAdvert(data ReceivedAdData) ([]*ReturningAdvert, error) {
	cityID, err := ads.getCityID(data.City)
	if err != nil {
		return nil, err
	}

	categoryID, err := ads.getCategoryID(data.Category)
	if err != nil {
		return nil, err
	}

	ads.mu.Lock()
	defer ads.mu.Unlock()

	newAd := &Advert{
		ID:          ads.getLastAdvertID(),
		UserID:      data.UserID,
		CityID:      cityID,
		CategoryID:  categoryID,
		Title:       data.Title,
		Description: data.Description,
		Price:       data.Price,
		Created:     time.Now(),
		Image:       data.Image,
		Active:      true,
		IsUsed:      data.IsUsed,
	}

	ads.Adverts = append(ads.Adverts, newAd)

	var returningAds []*ReturningAdvert
	returningAds = append(returningAds, &ReturningAdvert{
		Advert:   *newAd,
		City:     *ads.Cities[cityID-1],
		Category: *ads.Categories[categoryID-1],
	})

	return returningAds, nil
}

func (ads *AdvertsList) getCityID(city string) (uint, error) {
	ads.mu.Lock()
	defer ads.mu.Unlock()

	for _, val := range ads.Cities {
		if val.CityName == city || val.Translation == city {
			return val.ID, nil
		}
	}

	return 0, errWrongCityName
}

func (ads *AdvertsList) getCategoryID(category string) (uint, error) {
	ads.mu.Lock()
	defer ads.mu.Unlock()

	for _, val := range ads.Categories {
		if val.Name == category || val.Translation == category {
			return val.ID, nil
		}
	}

	return 0, errWrongCategoryName
}

func (ads *AdvertsList) getLastAdvertID() uint {
	ads.AdvertsCounter++

	return ads.AdvertsCounter
}

func (ads *AdvertsList) getLastLocationID() uint {
	ads.CitiesCounter++

	return ads.CitiesCounter
}

func (ads *AdvertsList) getLastCategoryID() uint {
	ads.CategoriesCounter++

	return ads.CategoriesCounter
}

func NewAdvertsList() *AdvertsList {
	return &AdvertsList{
		AdvertsCounter:    0,
		CitiesCounter:     0,
		CategoriesCounter: 0,
		Adverts:           make([]*Advert, 0),
		Cities:            make([]*City, 0),
		Categories:        make([]*Category, 0),
		mu:                sync.RWMutex{},
	}
}

func FillAdvertsList(ads *AdvertsList) {
	locationID := ads.getLastLocationID()
	ads.Cities = append(ads.Cities, &City{
		ID:          locationID,
		CityName:    "Москва",
		Translation: translit.Ru("Москва"),
	})

	categoryID := ads.getLastCategoryID()
	ads.Categories = append(ads.Categories, &Category{
		ID:          categoryID,
		Name:        "Тест",
		Translation: translit.Ru("Тест"),
	})

	for i := 1; i <= 60; i++ {
		price, _ := rand.Int(rand.Reader, big.NewInt(int64(maxPrice)))
		advertID := ads.getLastAdvertID()
		ads.Adverts = append(ads.Adverts, &Advert{
			ID:          advertID,
			UserID:      1,
			Title:       fmt.Sprintf("Объявление № %d", advertID),
			Image:       Image{},
			Description: fmt.Sprintf("Текст в объявлениии № %d", advertID),
			Price:       uint(price.Uint64()) * advertID,
			CityID:      1,
			CategoryID:  1,
			Created:     time.Now(),
			Active:      true,
			IsUsed:      true,
		})
	}
}

func AddAdvert(ads *AdvertsList, advert *Advert) {

	ads.Adverts = append(ads.Adverts, advert)

}
