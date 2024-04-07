package storage

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/mdigger/translit"
)

var (
	errWrongAdvertID      = errors.New("wrong advert ID")
	errWrongCityName      = errors.New("wrong city name")
	errWrongCategoryName  = errors.New("wrong category name")
	errWrongIDinCategory  = errors.New("there is no ad with such id in category")
	errWrongIDinCity      = errors.New("there is no ad with such id in city")
	errWrongAdvertsAmount = errors.New("too many elements specified")
	errAlreadyClosed      = errors.New("advert already closed")
)

type AdvertsListWrapper struct {
	AdvertsList *models.AdvertsList
}

func (ads *AdvertsListWrapper) GetAdvertByOnlyByID(advertID uint) (*models.ReturningAdvert, error) {
	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	if advertID > ads.AdvertsList.AdvertsCounter {
		return nil, errWrongAdvertID
	}

	cityID := ads.AdvertsList.Adverts[advertID-1].CityID
	categoryID := ads.AdvertsList.Adverts[advertID-1].CategoryID

	return &models.ReturningAdvert{
		Advert:   *ads.AdvertsList.Adverts[advertID-1],
		City:     *ads.AdvertsList.Cities[cityID-1],
		Category: *ads.AdvertsList.Categories[categoryID-1],
	}, nil
}

func (ads *AdvertsListWrapper) GetAdvert(advertID uint, city, category string) (*models.ReturningAdvert, error) {
	cityID, err := ads.GetCityID(city)
	if err != nil {
		return nil, err
	}

	categoryID, err := ads.GetCategoryID(category)
	if err != nil {
		return nil, err
	}

	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	if advertID > ads.AdvertsList.AdvertsCounter || ads.AdvertsList.Adverts[advertID-1].Deleted {
		return nil, errWrongAdvertID
	}

	if ads.AdvertsList.Adverts[advertID-1].CityID != cityID {
		return nil, errWrongIDinCity
	}

	if ads.AdvertsList.Adverts[advertID-1].CategoryID != categoryID {
		return nil, errWrongIDinCategory
	}

	return &models.ReturningAdvert{
		Advert:   *ads.AdvertsList.Adverts[advertID-1],
		City:     *ads.AdvertsList.Cities[cityID-1],
		Category: *ads.AdvertsList.Categories[categoryID-1],
	}, nil
}

func (ads *AdvertsListWrapper) GetAdvertsByCity(city string, startID, num uint) ([]*models.ReturningAdInList, error) {
	if num > ads.AdvertsList.AdvertsCounter {
		return nil, errWrongAdvertsAmount
	}

	cityID, err := ads.GetCityID(city)
	if err != nil {
		return nil, err
	}

	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	var returningAds []*models.ReturningAdInList
	var counter uint = 0

	for counter != num && counter+startID-1 != ads.AdvertsList.AdvertsCounter {
		ad := ads.AdvertsList.Adverts[startID+counter-1]
		exists := ad.Active && !ad.Deleted

		if exists && ad.CityID == cityID {
			returningAds = append(returningAds, &models.ReturningAdInList{
				ID:       ad.ID,
				Title:    ad.Title,
				Price:    ad.Price,
				City:     ads.AdvertsList.Cities[ad.CityID-1].Translation,
				Category: ads.AdvertsList.Categories[ad.CategoryID-1].Translation,
			})
		}

		counter++
	}

	return returningAds, nil
}

func (ads *AdvertsListWrapper) GetAdvertsByCategory(category, city string, startID, num uint) ([]*models.ReturningAdInList, error) {
	if num > ads.AdvertsList.AdvertsCounter {
		return nil, errWrongAdvertsAmount
	}

	cityID, err := ads.GetCityID(city)
	if err != nil {
		return nil, err
	}

	categoryID, err := ads.GetCategoryID(category)
	if err != nil {
		return nil, err
	}

	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	var returningAds []*models.ReturningAdInList
	var counter uint = 0

	for counter != num && counter+startID-1 != ads.AdvertsList.AdvertsCounter {
		ad := ads.AdvertsList.Adverts[startID+counter-1]
		exists := ad.Active && !ad.Deleted

		if exists && ad.CityID == cityID && ad.CategoryID == categoryID {
			returningAds = append(returningAds, &models.ReturningAdInList{
				ID:       ad.ID,
				Title:    ad.Title,
				Price:    ad.Price,
				City:     ads.AdvertsList.Cities[ad.CityID-1].Translation,
				Category: ads.AdvertsList.Categories[ad.CategoryID-1].Translation,
			})
		}

		counter++
	}

	return returningAds, nil
}

func (ads *AdvertsListWrapper) GetAdvertsByUserIDFiltered(userID uint, filter func(*models.Advert) bool) ([]*models.ReturningAdvert, error) {
	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()
	var returningAds []*models.ReturningAdvert
	for _, ad := range ads.AdvertsList.Adverts {
		if ad.UserID == userID && filter(ad) {
			returningAds = append(returningAds, &models.ReturningAdvert{
				Advert:   *ad,
				City:     *ads.AdvertsList.Cities[ad.CityID-1],
				Category: *ads.AdvertsList.Categories[ad.CategoryID-1],
			})
		}
	}
	return returningAds, nil
}

func (ads *AdvertsListWrapper) CreateAdvert(data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	cityID, err := ads.GetCityID(data.City)
	if err != nil {
		return nil, err
	}

	categoryID, err := ads.GetCategoryID(data.Category)
	if err != nil {
		return nil, err
	}

	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	newAd := &models.Advert{
		ID:          ads.GetLastAdvertID(),
		UserID:      data.UserID,
		CityID:      cityID,
		CategoryID:  categoryID,
		Title:       data.Title,
		Description: data.Description,
		Price:       data.Price,
		CreatedTime: time.Now(),
		Active:      true,
		IsUsed:      data.IsUsed,
	}

	ads.AdvertsList.Adverts = append(ads.AdvertsList.Adverts, newAd)

	return &models.ReturningAdvert{
		Advert:   *newAd,
		City:     *ads.AdvertsList.Cities[cityID-1],
		Category: *ads.AdvertsList.Categories[categoryID-1],
	}, nil
}

func (ads *AdvertsListWrapper) EditAdvert(data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	id := data.ID
	if id > ads.AdvertsList.AdvertsCounter || ads.AdvertsList.Adverts[id-1].Deleted {
		return nil, errWrongAdvertID
	}

	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	ads.AdvertsList.Adverts[id-1] = &models.Advert{
		ID:          id,
		UserID:      data.UserID,
		Title:       data.Title,
		Description: data.Description,
		Price:       data.Price,
		CityID:      ads.AdvertsList.Adverts[id-1].CityID,
		CategoryID:  ads.AdvertsList.Adverts[id-1].CategoryID,
		CreatedTime: ads.AdvertsList.Adverts[id-1].CreatedTime,
		Active:      true,
		IsUsed:      data.IsUsed,
		Deleted:     false,
	}

	return &models.ReturningAdvert{
		Advert:   *ads.AdvertsList.Adverts[id-1],
		Category: *ads.AdvertsList.Categories[ads.AdvertsList.Adverts[id-1].CategoryID-1],
		City:     *ads.AdvertsList.Cities[ads.AdvertsList.Adverts[id-1].CityID-1],
	}, nil
}

func (ads *AdvertsListWrapper) CloseAdvert(advertID uint) error {
	if advertID > ads.AdvertsList.AdvertsCounter || ads.AdvertsList.Adverts[advertID-1].Deleted {
		return errWrongAdvertID
	}

	if !ads.AdvertsList.Adverts[advertID-1].Active {
		return errAlreadyClosed
	}

	ads.AdvertsList.Adverts[advertID-1].Active = false

	return nil
}

func (ads *AdvertsListWrapper) DeleteAdvert(advertID uint) error {
	if advertID > ads.AdvertsList.AdvertsCounter || ads.AdvertsList.Adverts[advertID-1].Deleted {
		return errWrongAdvertID
	}

	ads.AdvertsList.Adverts[advertID-1].Deleted = true

	return nil
}

func (ads *AdvertsListWrapper) GetCityID(city string) (uint, error) {
	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	for _, val := range ads.AdvertsList.Cities {
		if val.CityName == city || val.Translation == city {
			return val.ID, nil
		}
	}

	return 0, errWrongCityName
}

func (ads *AdvertsListWrapper) GetCategoryID(category string) (uint, error) {
	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	for _, val := range ads.AdvertsList.Categories {
		if val.Name == category || val.Translation == category {
			return val.ID, nil
		}
	}

	return 0, errWrongCategoryName
}

func (ads *AdvertsListWrapper) GetLastAdvertID() uint {
	ads.AdvertsList.AdvertsCounter++

	return ads.AdvertsList.AdvertsCounter
}

func (ads *AdvertsListWrapper) GetLastLocationID() uint {
	ads.AdvertsList.CitiesCounter++

	return ads.AdvertsList.CitiesCounter
}

func (ads *AdvertsListWrapper) GetLastCategoryID() uint {
	ads.AdvertsList.CategoriesCounter++

	return ads.AdvertsList.CategoriesCounter
}

func NewAdvertsList() *AdvertsListWrapper {
	return &AdvertsListWrapper{
		AdvertsList: &models.AdvertsList{
			AdvertsCounter:    0,
			CitiesCounter:     0,
			CategoriesCounter: 0,
			Adverts:           make([]*models.Advert, 0),
			Cities:            make([]*models.City, 0),
			Categories:        make([]*models.Category, 0),
			Mux:               sync.RWMutex{},
		},
	}
}

func FillAdvertsList(ads *AdvertsListWrapper) {
	locationID := ads.GetLastLocationID()
	ads.AdvertsList.Cities = append(ads.AdvertsList.Cities, &models.City{
		ID:          locationID,
		CityName:    "Москва",
		Translation: translit.Ru("Москва"),
	})

	categoryID := ads.GetLastCategoryID()
	ads.AdvertsList.Categories = append(ads.AdvertsList.Categories, &models.Category{
		ID:          categoryID,
		Name:        "Тест",
		Translation: translit.Ru("Тест"),
	})

	for i := 1; i <= 100; i++ {
		price, _ := rand.Int(rand.Reader, big.NewInt(int64(models.MaxPrice)))
		advertID := ads.GetLastAdvertID()
		ads.AdvertsList.Adverts = append(ads.AdvertsList.Adverts, &models.Advert{
			ID:          advertID,
			UserID:      1,
			Title:       fmt.Sprintf("Объявление № %d", advertID),
			Description: fmt.Sprintf("Текст в объявлениии № %d", advertID),
			Price:       uint(price.Uint64()) * advertID,
			CityID:      1,
			CategoryID:  1,
			CreatedTime: time.Now(),
			Active:      true,
			IsUsed:      true,
			Deleted:     false,
		})
	}
}

func AddAdvert(ads *models.AdvertsList, advert *models.Advert) {

	ads.Adverts = append(ads.Adverts, advert)

}
