package storage

import (
	"errors"
	"sync"
	"time"
)

var errProfileNotExists = errors.New("profile does not exist")

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
	mu       sync.RWMutex
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

type ProfileInfo interface {
	CreateProfile(userID uint) *Profile
	GetProfileByUserID(userID uint) (*Profile, error)

	SetProfileCity(userID uint, data SetProfileCityNec)
	SetProfilePhone(userID uint, data SetProfilePhoneNec)
	SetProfileRating(userID uint, data SetProfileRatingNec)
	SetProfile(userID uint, data SetProfileNec)
	EditProfile(userID uint, data EditProfileNec)
	SetProfileApproved(userID uint)
}

func (pl *ProfileList) CreateProfile(userID uint) *Profile {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.Profiles[userID] = &Profile{
		UserID:       userID,
		RegisterTime: time.Now(),
	}

	return pl.Profiles[userID]
}

func (pl *ProfileList) GetProfileByUserID(userID uint) (*Profile, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	p, ok := pl.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	return p, nil
}

func (pl *ProfileList) SetProfileCity(userID uint, data SetProfileCityNec) (*Profile, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	p, ok := pl.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.City = data.City

	return p, nil
}

func (pl *ProfileList) SetProfilePhone(userID uint, data SetProfilePhoneNec) (*Profile, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	p, ok := pl.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.Phone = data.Phone

	return p, nil
}

func (pl *ProfileList) SetProfileRating(userID uint, data SetProfileRatingNec) (*Profile, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	p, ok := pl.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.Rating = (p.Rating*p.ReactionsCount + data.Reaction) / (p.ReactionsCount + 1)

	return p, nil
}

func (pl *ProfileList) SetProfileApproved(userID uint) (*Profile, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	p, ok := pl.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.Approved = true

	return p, nil
}

func (pl *ProfileList) SetProfile(userID uint, data SetProfileNec) (*Profile, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	p, ok := pl.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.Name = data.Name
	p.Surname = data.Surname
	p.Avatar = data.Avatar

	return p, nil
}

func (pl *ProfileList) EditProfile(userID uint, data EditProfileNec) (*Profile, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	_, ok := pl.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	old := pl.Profiles[userID]

	old.Name = data.Name
	old.Surname = data.Surname
	old.City = data.City
	old.Phone = data.Phone
	old.Avatar = data.Avatar
	old.MerchantsName = data.MerchantsName
	old.SubersCount = data.SubersCount
	old.SubonsCount = data.SubonsCount

	// pl.Profiles[userID] = &Profile{
	// 	UserID:       old.UserID,
	// 	RegisterTime: old.RegisterTime,
	// 	Name:         data.Name,
	// 	Surname:      data.Surname,
	// 	City:         data.City,
	// 	Phone:        data.Phone,
	// 	Avatar:       data.Avatar,
	// }

	return pl.Profiles[userID], nil
}

// func NewProfileList() *ProfileList {
// 	return &ProfileList{
// 		mu: sync.RWMutex{},
// 	}
// }

func NewProfileList() *ProfileList {
	return &ProfileList{
		Profiles: map[uint]*Profile{
			1: {
				UserID:  1,
				Name:    "Vladimir",
				Surname: "Vasilievich",
				City: City{
					ID:          1,
					Name:        "Moscow",
					Translation: "Москва",
				},
				Phone:          "1234567890",
				Avatar:         Image{}, // Предполагается, что Image имеет конструктор по умолчанию
				RegisterTime:   time.Now(),
				Rating:         5.0,
				ReactionsCount: 10,
				Approved:       true,
				MerchantsName:  "Vova",
				SubersCount:    10,
				SubonsCount:    100,
			},
			2: {
				Name:    "Petr",
				Surname: "Andreevich",
				UserID:  2,
				City: City{
					ID:          1,
					Name:        "Kaluga",
					Translation: "Калуга",
				},
				Phone:          "1234567890",
				Avatar:         Image{}, // Предполагается, что Image имеет конструктор по умолчанию
				RegisterTime:   time.Now(),
				Rating:         4.4,
				ReactionsCount: 10,
				Approved:       false,
				MerchantsName:  "Petya",
				SubersCount:    100,
				SubonsCount:    10,
			},
		},
		mu: sync.RWMutex{},
	}
}
