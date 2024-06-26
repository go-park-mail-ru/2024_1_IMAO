//nolint:tagliatelle
package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

const (
	MaxPrice = 1000
)

type Image struct{}

type ReceivedAdData struct {
	ID          uint   `json:"Id"` //nolint:tagliatelle
	UserID      uint   `json:"userId"`
	City        string `json:"city"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	IsUsed      bool   `json:"isUsed"`
	Phone       string `json:"phone"`
}

type Category struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Translation string `json:"translation"`
}

type Advert struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"userId"`
	CityID        uint      `json:"cityId"`
	CategoryID    uint      `json:"categoryId"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Phone         string    `json:"phone"`
	Price         uint      `json:"price"`
	CreatedTime   time.Time `json:"created"`
	ClosedTime    time.Time `json:"closed"`
	Active        bool      `json:"active"`
	IsUsed        bool      `json:"isUsed"`
	Views         uint      `json:"views"`
	InFavourites  bool      `json:"inFavourites"`
	InCart        bool      `json:"inCart"`
	FavouritesNum uint      `json:"favouritesNum"`
	Deleted       bool      `json:"-"`
}

type Promotion struct {
	NeedPing          bool             `json:"needPing"`
	IsPromoted        bool             `json:"isPromoted"`
	PromotionStart    pgtype.Timestamp `json:"promotionStart"`
	PromotionDuration pgtype.Interval  `json:"promotionDuration"`
}

type ReturningAdvert struct {
	Advert    Advert    `json:"advert"`
	Promotion Promotion `json:"promotion"`
	City      City      `json:"city"`
	Category  Category  `json:"category"`
	Photos    []string  `json:"photos"`
	PhotosIMG []string  `json:"photosIMG"`
}

type PhotoPad struct {
	Photo []*string `json:"photo"`
}

type PhotoPadSoloImage struct {
	Photo *string `json:"photo"`
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
	ID           uint     `json:"id"`
	Title        string   `json:"title"`
	Price        uint     `json:"price"`
	City         string   `json:"city"`
	Category     string   `json:"category"`
	Photos       []string `json:"photos"`
	PhotosIMG    []string `json:"photosIMG"`
	InFavourites bool     `json:"inFavourites"`
	InCart       bool     `json:"inCart"`
	IsPromoted   bool     `json:"isPromoted"`
	IsActive     bool     `json:"isActive"`
}

type ReturningAdvertList struct {
	AdvertItems []*ReturningAdvert
	Mux         sync.RWMutex
}

type PriceHistoryItem struct {
	UpdatedTime string  `json:"updatedTime"`
	NewPrice    float64 `json:"newPrice"` // ПЕРЕПИСАТЬ НА uint
}

type DBInsertionAdvert struct {
	UserID      uint   `json:"userId"`
	CityID      uint   `json:"cityId"`
	CategoryID  uint   `json:"categoryId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

func (adv *Advert) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	adv.Title = sanitizer.Sanitize(adv.Title)
	adv.Description = sanitizer.Sanitize(adv.Description)
}

func (advl *ReturningAdInList) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	advl.Title = sanitizer.Sanitize(advl.Title)
}
