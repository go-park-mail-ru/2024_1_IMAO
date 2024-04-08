package models

import (
	"sync"
	"time"
)

type Profile struct {
	UserID         uint      `json:"id"`
	Name           string    `json:"name"`
	Surname        string    `json:"surname"`
	City           City      `json:"city"`
	Phone          string    `json:"phoneNumber"`
	Avatar         Image     `json:"avatar"`
	RegisterTime   time.Time `json:"regTime"`
	Rating         float64   `json:"rating"`
	ReactionsCount float64   `json:"reactionsCount"`
	Approved       bool      `json:"approved"`
	MerchantsName  string    `json:"merchantsName"`
	SubersCount    int       `json:"subersCount"`
	SubonsCount    int       `json:"subonsCount"`
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
	Avatar        Image  `json:"avatar"`
	City          City   `json:"city"`
	Phone         string `json:"phone"`
	MerchantsName string `json:"merchantsName"`
	SubersCount   int    `json:"subersCount"`
	SubonsCount   int    `json:"subonsCount"`
}