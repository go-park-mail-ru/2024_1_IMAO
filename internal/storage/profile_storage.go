package storage

import (
	"errors"
	"sync"
	"time"
)

var errProfileNotExists = errors.New("profile does not exist")

type Profile struct {
	UserID         uint      `json:"id"`
	City           City      `json:"city"`
	Phone          string    `json:"phoneNumber"`
	Avatar         Image     `json:"avatar"`
	RegisterTime   time.Time `json:"regTime"`
	Rating         float64   `json:"rating"`
	ReactionsCount float64   `json:"reactionsCount"`
	Approved       bool      `json:"approved"`
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

type SetProfileAvatarNec struct {
	Avatar Image `json:"avatar"`
}

type SetProfileRatingNec struct {
	Reaction float64 `json:"reaction"`
}

type ProfileAdvertsNec struct {
	Filter AdvertsFilter `json:"filter"`
}

type ProfileEditNec struct {
	Avatar Image  `json:"avatar"`
	City   City   `json:"city"`
	Phone  string `json:"phone"`
}

type ProfileInfo interface {
	CreateProfile(userID uint) *Profile
	GetProfileByUserID(userID uint) (*Profile, error)

	SetProfileCity(userID uint, to City)
	SetProfilePhone(userID uint, to string)
	SetProfileAvatar(userID uint, to Image)
	SetProfileRating(userID uint, to int)
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

func (pl *ProfileList) SetProfileCity(userID uint, to City) (*Profile, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	p, ok := pl.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.City = to

	return p, nil
}

func (pl *ProfileList) SetProfilePhone(userID uint, to string) (*Profile, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	p, ok := pl.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.Phone = to

	return p, nil
}

func (pl *ProfileList) SetProfileRating(userID uint, reactionProp float64) (*Profile, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	p, ok := pl.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.Rating = (p.Rating*p.ReactionsCount + reactionProp) / (p.ReactionsCount + 1)

	return p, nil
}

func (pl *ProfileList) SetProfileAvatar(userID uint, to Image) (*Profile, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	p, ok := pl.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.Avatar = to

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

func (pl *ProfileList) ProfileEdit(userID uint, data ProfileEditNec) (*Profile, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	_, ok := pl.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	pl.Profiles[userID] = &Profile{
		City:   data.City,
		Phone:  data.Phone,
		Avatar: data.Avatar,
	}

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
				UserID: 1,
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
			},
			2: {
				UserID: 2,
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
			},
		},
		mu: sync.RWMutex{},
	}
}
