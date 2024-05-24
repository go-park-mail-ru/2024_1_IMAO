package models

import (
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

type Profile struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	Name            string    `json:"name"`
	Surname         string    `json:"surname"`
	City            City      `json:"city"`
	Phone           string    `json:"phoneNumber"`
	Avatar          string    `json:"avatar"`
	RegisterTime    time.Time `json:"regTime"`
	Rating          float32   `json:"rating"`
	ReactionsCount  float32   `json:"reactionsCount"`
	Approved        bool      `json:"approved"`
	MerchantsName   string    `json:"merchantsName"`
	SubersCount     int       `json:"subersCount"`
	SubonsCount     int       `json:"subonsCount"`
	AvatarIMG       string    `json:"avatarImg"`
	ActiveAddsCount int       `json:"activeAddsCount"`
	SoldAddsCount   int       `json:"soldAddsCount"`
	CartNum         uint      `json:"cartNum"`
	FavNum          uint      `json:"favNum"`
}

type ProfilePad struct {
	Name    *string `json:"name"`
	Surname *string `json:"surname"`
	Phone   *string `json:"phoneNumber"`
	Avatar  *string `json:"avatar"`
}

type AdvertsFilter int

const (
	FilterAll = iota
	FilterActive
	FilterClosed
)

type ProfileList struct {
	Profiles map[uint]*Profile
	Mux      sync.RWMutex
}

type SetProfileCityNec struct {
	City City `json:"city"`
}

type SetProfilePhoneNec struct {
	Phone string `json:"phone"`
}

type SetProfileRatingNec struct {
	Reaction float64 `json:"reaction"`
}

type ProfileAdvertsNec struct {
	Filter AdvertsFilter `json:"filter"`
}

type SetProfileNec struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Avatar  Image  `json:"avatar"`
}

type EditProfileNec struct {
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	Avatar        string `json:"avatar"`
	City          City   `json:"city"`
	Phone         string `json:"phone"`
	MerchantsName string `json:"merchantsName"`
	SubersCount   int    `json:"subersCount"`
	SubonsCount   int    `json:"subonsCount"`
}

type DBInsertionProfile struct {
	UserID  uint   `json:"userId"`
	CityID  uint   `json:"cityId"`
	Phone   string `json:"phone"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func (prof *Profile) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()

	prof.Name = sanitizer.Sanitize(prof.Name)
	prof.Surname = sanitizer.Sanitize(prof.Surname)
	prof.MerchantsName = sanitizer.Sanitize(prof.MerchantsName)
	prof.Phone = sanitizer.Sanitize(prof.Phone)
}

type ReceivedMerchantItem struct {
	MerchantID uint `json:"merchantId"`
}
