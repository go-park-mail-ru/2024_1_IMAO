package models

import (
	"sync"
	"time"
)

const (
	MaxPrice = 1000
)

type ReceivedAdData struct {
	ID          uint   `json:"Id"`
	UserID      uint   `json:"userId"`
	City        string `json:"city"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	IsUsed      bool   `json:"isUsed"`
}

type Category struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
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
	CreatedTime time.Time `json:"created"`
	ClosedTime  time.Time `json:"closed"`
	Active      bool      `json:"active"`
	IsUsed      bool      `json:"isUsed"`
	Deleted     bool      `json:"-"`
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
	Mux               sync.RWMutex
}

type ReturningAdInList struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Price    uint   `json:"price"`
	City     string `json:"city"`
	Category string `json:"category"`
}

type ReturningAdvertList struct {
	AdvertItems []*ReturningAdvert
	Mux         sync.RWMutex
}

// type AdvertDB struct {
// 	ID          uint      `json:"id"`
// 	UserID      uint      `json:"userId"`
// 	CityID      uint      `json:"cityId"`
// 	CategoryID  uint      `json:"categoryId"`
// 	Title       string    `json:"title"`
// 	Description string    `json:"description"`
// 	Price       uint      `json:"price"`
// 	CreatedTime time.Time `json:"created"`
// 	ClosedTime  time.Time `json:"closed"`
// 	Active      bool      `json:"active"`
// 	IsUsed      bool      `json:"isUsed"`
// 	StatusID    uint      `json:"statusID"`
// }

// type ReturningAdvertDB struct {
// 	Advert   Advert   `json:"advert"`
// 	City     City     `json:"city"`
// 	Category Category `json:"category"`
// 	Status   Status   `json:"status"`
// }

// type ReturningAdvertDBList struct {
// 	AdvertItems []*ReturningAdvertDB
// 	Mux       sync.RWMutex
// }
