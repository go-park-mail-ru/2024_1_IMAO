package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
)

var errProfileNotExists = errors.New("profile does not exist")

type ProfileListWrapper struct {
	ProfileList *models.ProfileList
}

func (pl *ProfileListWrapper) CreateProfile(userID uint) *models.Profile {
	pl.ProfileList.Mux.Lock()
	defer pl.ProfileList.Mux.Unlock()

	pl.ProfileList.Profiles[userID] = &models.Profile{
		UserID:       userID,
		RegisterTime: time.Now(),
	}

	return pl.ProfileList.Profiles[userID]
}

func (pl *ProfileListWrapper) GetProfileByUserID(userID uint) (*models.Profile, error) {
	pl.ProfileList.Mux.Lock()
	defer pl.ProfileList.Mux.Unlock()

	p, ok := pl.ProfileList.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	return p, nil
}

func (pl *ProfileListWrapper) SetProfileCity(userID uint, data models.SetProfileCityNec) (*models.Profile, error) {
	pl.ProfileList.Mux.Lock()
	defer pl.ProfileList.Mux.Unlock()

	p, ok := pl.ProfileList.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.City = data.City

	return p, nil
}

func (pl *ProfileListWrapper) SetProfilePhone(userID uint, data models.SetProfilePhoneNec) (*models.Profile, error) {
	pl.ProfileList.Mux.Lock()
	defer pl.ProfileList.Mux.Unlock()

	p, ok := pl.ProfileList.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.Phone = data.Phone

	return p, nil
}

func (pl *ProfileListWrapper) SetProfileRating(userID uint, data models.SetProfileRatingNec) (*models.Profile, error) {
	pl.ProfileList.Mux.Lock()
	defer pl.ProfileList.Mux.Unlock()

	p, ok := pl.ProfileList.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.Rating = (p.Rating*p.ReactionsCount + data.Reaction) / (p.ReactionsCount + 1)

	return p, nil
}

func (pl *ProfileListWrapper) SetProfileApproved(userID uint) (*models.Profile, error) {
	pl.ProfileList.Mux.Lock()
	defer pl.ProfileList.Mux.Unlock()

	p, ok := pl.ProfileList.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.Approved = true

	return p, nil
}

func (pl *ProfileListWrapper) SetProfile(userID uint, data models.SetProfileNec) (*models.Profile, error) {
	pl.ProfileList.Mux.Lock()
	defer pl.ProfileList.Mux.Unlock()

	p, ok := pl.ProfileList.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	p.Name = data.Name
	p.Surname = data.Surname
	p.Avatar = data.Avatar

	return p, nil
}

func (pl *ProfileListWrapper) EditProfile(userID uint, data models.EditProfileNec) (*models.Profile, error) {
	pl.ProfileList.Mux.Lock()
	defer pl.ProfileList.Mux.Unlock()

	_, ok := pl.ProfileList.Profiles[userID]
	if !ok {
		return nil, errProfileNotExists
	}

	old := pl.ProfileList.Profiles[userID]

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

	return pl.ProfileList.Profiles[userID], nil
}

// func NewProfileList() *ProfileList {
// 	return &ProfileList{
// 		mu: sync.RWMutex{},
// 	}
// }

func NewProfileList() *ProfileListWrapper {
	return &ProfileListWrapper{
		ProfileList: &models.ProfileList{
			Profiles: map[uint]*models.Profile{
				1: {
					UserID:  1,
					Name:    "Vladimir",
					Surname: "Vasilievich",
					City: models.City{
						ID:          1,
						CityName:    "Moscow",
						Translation: "Москва",
					},
					Phone:          "1234567890",
					Avatar:         models.Image{}, // Предполагается, что Image имеет конструктор по умолчанию
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
					City: models.City{
						ID:          1,
						CityName:    "Kaluga",
						Translation: "Калуга",
					},
					Phone:          "1234567890",
					Avatar:         models.Image{}, // Предполагается, что Image имеет конструктор по умолчанию
					RegisterTime:   time.Now(),
					Rating:         4.4,
					ReactionsCount: 10,
					Approved:       false,
					MerchantsName:  "Petya",
					SubersCount:    100,
					SubonsCount:    10,
				},
			},
			Mux: sync.RWMutex{},
		},
	}
}
