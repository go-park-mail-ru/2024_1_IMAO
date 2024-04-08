package storage

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	errProfileNotExists = errors.New("profile does not exist")
	NameSeqProfile      = pgx.Identifier{"public", "profile_id_seq"} //nolint:gochecknoglobals
)

type ProfileListWrapper struct {
	ProfileList *models.ProfileList
	Pool        *pgxpool.Pool
}

func (pl *ProfileListWrapper) createProfile(ctx context.Context, tx pgx.Tx, profile *models.Profile) error {
	SQLCreateProfile := `INSERT INTO public.profile(user_id) VALUES ($1);`

	var err error

	_, err = tx.Exec(ctx, SQLCreateProfile, profile.UserID)

	if err != nil {
		return fmt.Errorf("Something went wrong while executing create profile query", err)
	}

	return nil
}

func (pl *ProfileListWrapper) CreateProfile(ctx context.Context, userID uint) *models.Profile {
	// 	НАВЕРНОЕ, СЮДА НУЖНО ДОПИСАТЬ ПРОВЕРКУ НА ТО, ЧТО ЮЗЕР С ТАКИМ ID ДЕЙСТВИТЕЛЬНО СУЩЕСТВУЕТ

	profile := models.Profile{
		UserID: userID,
	}

	err := pgx.BeginFunc(ctx, pl.Pool, func(tx pgx.Tx) error {
		err := pl.createProfile(ctx, tx, &profile)
		if err != nil {

			return fmt.Errorf("Something went wrong while creating profile", err)
		}
		id, err := repository.GetLastValSeq(ctx, tx, NameSeqProfile)
		if err != nil {

			return fmt.Errorf("Something went wrong getting user id from seq", err)
		}
		profile.ID = uint(id)

		return nil
	})

	if err != nil {

		return nil
	}
	fmt.Println("profile", profile)
	return &profile
}

func (pl *ProfileListWrapper) getProfileByUserID(ctx context.Context, tx pgx.Tx, id uint) (*models.Profile, error) {
	SQLUserById := `
		SELECT 
			p.id, 
			p.user_id, 
			p.city_id, 
			p.phone, 
			p.name, 
			p.surname, 
			p.regtime, 
			p.verified, 
			p.avatar_url,
			c.name AS city_name,
			c.translation AS city_translation
		FROM 
			public.profile p
		INNER JOIN 
			public.city c
		ON 
			p.city_id = c.id
		WHERE 
			p.user_id = $1`
	profileLine := tx.QueryRow(ctx, SQLUserById, id)

	profile := models.Profile{}
	city := models.City{}
	profilePad := models.ProfilePad{}
	var avatar_url *string // ЗАГЛУШКА

	if err := profileLine.Scan(&profile.ID, &profile.UserID, &city.ID, &profilePad.Phone, &profilePad.Name,
		&profilePad.Surname, &profile.RegisterTime, &profile.Approved, &avatar_url, &city.CityName, &city.Translation); err != nil {

		return nil, err
	}

	phoneToInsert := ""
	if profilePad.Phone != nil {
		phoneToInsert = *profilePad.Phone
	}

	nameToInsert := ""
	if profilePad.Name != nil {
		nameToInsert = *profilePad.Name
	}

	surnameToInsert := ""
	if profilePad.Surname != nil {
		surnameToInsert = *profilePad.Surname
	}

	profile.Phone = phoneToInsert
	profile.Name = nameToInsert
	profile.Surname = surnameToInsert
	profile.Avatar = models.Image{}
	profile.City = city

	rand.Seed(time.Now().UnixNano())
	profile.Rating = math.Round((rand.Float64()*4+1)*100) / 100
	profile.ReactionsCount = 10
	profile.Approved = true
	profile.MerchantsName = nameToInsert
	profile.SubersCount = rand.Intn(10)
	profile.SubonsCount = rand.Intn(10)

	return &profile, nil
}

func (pl *ProfileListWrapper) GetProfileByUserID(ctx context.Context, userID uint) (*models.Profile, error) {
	// pl.ProfileList.Mux.Lock()
	// defer pl.ProfileList.Mux.Unlock()

	// p, ok := pl.ProfileList.Profiles[userID]
	// if !ok {
	// 	return nil, errProfileNotExists
	// }

	var profile *models.Profile

	err := pgx.BeginFunc(ctx, pl.Pool, func(tx pgx.Tx) error {
		profileInner, err := pl.getProfileByUserID(ctx, tx, userID)
		profile = profileInner

		return err
	})

	fmt.Println("err", err)

	if err != nil {
		return nil, errProfileNotExists
	}

	fmt.Println("profile", profile)

	return profile, nil

	//return p, nil
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

func NewProfileList(pool *pgxpool.Pool) *ProfileListWrapper {
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
		Pool: pool,
	}
}
