package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"mime/multipart"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	errProfileNotExists = errors.New("profile does not exist")
	NameSeqProfile      = pgx.Identifier{"public", "profile_id_seq"} //nolint:gochecknoglobals
)

type ProfileListWrapper struct {
	ProfileList *models.ProfileList
	Pool        *pgxpool.Pool
	Logger      *zap.SugaredLogger
}

func (pl *ProfileListWrapper) createProfile(ctx context.Context, tx pgx.Tx, profile *models.Profile) error {
	SQLCreateProfile := `INSERT INTO public.profile(user_id) VALUES ($1);`
	pl.Logger.Infof(`INSERT INTO public.profile(user_id) VALUES (%s);`, profile)

	var err error

	_, err = tx.Exec(ctx, SQLCreateProfile, profile.UserID)

	if err != nil {
		pl.Logger.Errorf("Something went wrong while executing create profile query, err=%v", err)

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
			pl.Logger.Errorf("Something went wrong while creating profile, err=%v", err)
			return fmt.Errorf("Something went wrong while creating profile", err)
		}
		id, err := repository.GetLastValSeq(ctx, tx, NameSeqProfile)
		if err != nil {
			pl.Logger.Errorf("Something went wrong getting user id from seq, err=%v", err)

			return fmt.Errorf("Something went wrong getting user id from seq", err)
		}
		profile.ID = uint(id)

		return nil
	})

	if err != nil {

		return nil
	}

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
	pl.Logger.Infof(`
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
		p.user_id = %s`, id)
	profileLine := tx.QueryRow(ctx, SQLUserById, id)

	profile := models.Profile{}
	city := models.City{}
	profilePad := models.ProfilePad{}

	if err := profileLine.Scan(&profile.ID, &profile.UserID, &city.ID, &profilePad.Phone, &profilePad.Name,
		&profilePad.Surname, &profile.RegisterTime, &profile.Approved, &profilePad.Avatar, &city.CityName, &city.Translation); err != nil {

		pl.Logger.Errorf("Something went wrong while scanning profile, err=%v", err)

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

	avatartToInsert := ""
	if profilePad.Avatar != nil {
		avatartToInsert = *profilePad.Avatar
	}

	profile.Phone = phoneToInsert
	profile.Name = nameToInsert
	profile.Surname = surnameToInsert
	profile.Avatar = avatartToInsert
	profile.City = city

	rand.Seed(time.Now().UnixNano())
	profile.Rating = math.Round((rand.Float64()*4+1)*100) / 100
	profile.ReactionsCount = 10
	//profile.Approved = true
	profile.MerchantsName = nameToInsert
	profile.SubersCount = rand.Intn(10)
	profile.SubonsCount = rand.Intn(10)

	return &profile, nil
}

func (pl *ProfileListWrapper) GetProfileByUserID(ctx context.Context, userID uint) (*models.Profile, error) {
	var profile *models.Profile

	err := pgx.BeginFunc(ctx, pl.Pool, func(tx pgx.Tx) error {
		profileInner, err := pl.getProfileByUserID(ctx, tx, userID)
		profile = profileInner

		return err
	})

	if err != nil {
		pl.Logger.Errorf("Something went wrong while getting profile by UserID , err=%v", errProfileNotExists)

		return nil, errProfileNotExists
	}

	return profile, nil
}

func (pl *ProfileListWrapper) setProfileCity(ctx context.Context, tx pgx.Tx, userID uint, data models.City) (*models.Profile, error) {

	SQLUpdateProfileCity := `UPDATE public.profile p
	SET city_id = $1
	FROM public.city c
	WHERE c.id = $1 AND p.user_id = $2
	RETURNING p.id, p.user_id, p.city_id, p.phone, p.name, p.surname, p.regtime, p.verified, p.avatar_url, c.name AS city_name, c.translation AS city_translation;`

	pl.Logger.Infof(`UPDATE public.profile p
	SET city_id = %s
	FROM public.city c
	WHERE c.id = %s AND p.user_id = %s
	RETURNING p.id, p.user_id, p.city_id, p.phone, p.name, p.surname, p.regtime, p.verified, p.avatar_url, c.name AS city_name, c.translation AS city_translation;`, data.ID, data.ID, userID)

	profileLine := tx.QueryRow(ctx, SQLUpdateProfileCity, data.ID, userID)

	profile := models.Profile{}
	city := models.City{}
	profilePad := models.ProfilePad{}

	if err := profileLine.Scan(&profile.ID, &profile.UserID, &city.ID, &profilePad.Phone, &profilePad.Name,
		&profilePad.Surname, &profile.RegisterTime, &profile.Approved, &profilePad.Avatar, &city.CityName, &city.Translation); err != nil {

		pl.Logger.Errorf("Something went wrong while scanning profile lines, err=%v", err)

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

	avatartToInsert := ""
	if profilePad.Avatar != nil {
		avatartToInsert = *profilePad.Avatar
	}

	profile.Phone = phoneToInsert
	profile.Name = nameToInsert
	profile.Surname = surnameToInsert
	profile.Avatar = avatartToInsert
	profile.City = city

	rand.Seed(time.Now().UnixNano())
	profile.Rating = math.Round((rand.Float64()*4+1)*100) / 100
	profile.ReactionsCount = 10
	//profile.Approved = true
	profile.MerchantsName = nameToInsert
	profile.SubersCount = rand.Intn(10)
	profile.SubonsCount = rand.Intn(10)

	return &profile, nil
}

func (pl *ProfileListWrapper) SetProfileCity(ctx context.Context, userID uint, data models.City) (*models.Profile, error) {
	// ЭТО НУЖНО РЕАЛИЗОВАТЬ (НО ЭТО НЕ ТОЧНО)
	// if pl.ProfileExists(ctx, id) {
	// 	return nil, errProfileDoNotExists
	// }

	var profile *models.Profile

	err := pgx.BeginFunc(ctx, pl.Pool, func(tx pgx.Tx) error {
		profileInner, err := pl.setProfileCity(ctx, tx, userID, data)
		profile = profileInner

		return err
	})

	if err != nil {
		pl.Logger.Errorf("Something went wrong while updating profile city, err=%v", errProfileNotExists)

		return nil, errProfileNotExists
	}

	return profile, nil
}

func (pl *ProfileListWrapper) setProfilePhone(ctx context.Context, tx pgx.Tx, userID uint, data models.SetProfilePhoneNec) (*models.Profile, error) {

	SQLUpdateProfilePhone := `UPDATE public.profile p
	SET phone = $1
	FROM public.city c
	WHERE c.id = p.city_id AND p.user_id = $2
	RETURNING p.id, p.user_id, p.city_id, p.phone, p.name, p.surname, p.regtime, p.verified, p.avatar_url, c.name AS city_name, c.translation AS city_translation;`

	pl.Logger.Infof(`UPDATE public.profile p
	SET phone = %s
	FROM public.city c
	WHERE c.id = p.city_id AND p.user_id = %s
	RETURNING p.id, p.user_id, p.city_id, p.phone, p.name, p.surname, p.regtime, p.verified, p.avatar_url, c.name AS city_name, c.translation AS city_translation;`, data.Phone, userID)

	profileLine := tx.QueryRow(ctx, SQLUpdateProfilePhone, data.Phone, userID)

	profile := models.Profile{}
	city := models.City{}
	profilePad := models.ProfilePad{}

	if err := profileLine.Scan(&profile.ID, &profile.UserID, &city.ID, &profilePad.Phone, &profilePad.Name,
		&profilePad.Surname, &profile.RegisterTime, &profile.Approved, &profilePad.Avatar, &city.CityName, &city.Translation); err != nil {

		pl.Logger.Errorf("Something went wrong while scanning profile lines , err=%v", err)

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

	avatartToInsert := ""
	if profilePad.Avatar != nil {
		avatartToInsert = *profilePad.Avatar
	}

	profile.Phone = phoneToInsert
	profile.Name = nameToInsert
	profile.Surname = surnameToInsert
	profile.Avatar = avatartToInsert
	profile.City = city

	rand.Seed(time.Now().UnixNano())
	profile.Rating = math.Round((rand.Float64()*4+1)*100) / 100
	profile.ReactionsCount = 10
	//profile.Approved = true
	profile.MerchantsName = nameToInsert
	profile.SubersCount = rand.Intn(10)
	profile.SubonsCount = rand.Intn(10)

	return &profile, nil
}

func (pl *ProfileListWrapper) SetProfilePhone(ctx context.Context, userID uint, data models.SetProfilePhoneNec) (*models.Profile, error) {
	// ЭТО НУЖНО РЕАЛИЗОВАТЬ (НО ЭТО НЕ ТОЧНО)
	// if pl.ProfileExists(ctx, id) {
	// 	return nil, errProfileDoNotExists
	// }

	var profile *models.Profile

	err := pgx.BeginFunc(ctx, pl.Pool, func(tx pgx.Tx) error {
		profileInner, err := pl.setProfilePhone(ctx, tx, userID, data)
		profile = profileInner

		return err
	})

	if err != nil {
		pl.Logger.Errorf("Something went wrong while updating profile phone , err=%v", errProfileNotExists)

		return nil, errProfileNotExists
	}

	return profile, nil
}

func (pl *ProfileListWrapper) setProfileAvatarUrl(ctx context.Context, tx pgx.Tx, userID uint, url string) (*models.Profile, error) {

	SQLUpdateProfileAvatarURL := `
	UPDATE public.profile p
	SET avatar_url= $1
	FROM public.city c
	WHERE c.id = p.city_id AND p.user_id = $2
	RETURNING p.id, p.user_id, p.city_id, p.phone, p.name, p.surname, p.regtime, p.verified, p.avatar_url, c.name AS city_name, c.translation AS city_translation;`

	pl.Logger.Infof(`UPDATE public.profile p
	SET avatar_url= %s
	FROM public.city c
	WHERE c.id = p.city_id AND p.user_id = %s
	RETURNING p.id, p.user_id, p.city_id, p.phone, p.name, p.surname, p.regtime, p.verified, p.avatar_url, c.name AS city_name, c.translation AS city_translation;`, url, userID)

	profileLine := tx.QueryRow(ctx, SQLUpdateProfileAvatarURL, url, userID)

	profile := models.Profile{}
	city := models.City{}
	profilePad := models.ProfilePad{}

	if err := profileLine.Scan(&profile.ID, &profile.UserID, &city.ID, &profilePad.Phone, &profilePad.Name,
		&profilePad.Surname, &profile.RegisterTime, &profile.Approved, &profilePad.Avatar, &city.CityName, &city.Translation); err != nil {

		pl.Logger.Errorf("Something went wrong while scanning profile lines , err=%v", err)

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

	avatartToInsert := ""
	if profilePad.Avatar != nil {
		avatartToInsert = *profilePad.Avatar
	}

	profile.Phone = phoneToInsert
	profile.Name = nameToInsert
	profile.Surname = surnameToInsert
	profile.Avatar = avatartToInsert
	profile.City = city

	rand.Seed(time.Now().UnixNano())
	profile.Rating = math.Round((rand.Float64()*4+1)*100) / 100
	profile.ReactionsCount = 10
	//profile.Approved = true
	profile.MerchantsName = nameToInsert
	profile.SubersCount = rand.Intn(10)
	profile.SubonsCount = rand.Intn(10)

	return &profile, nil
}

func (pl *ProfileListWrapper) SetProfileAvatarUrl(ctx context.Context, file *multipart.FileHeader, folderName string, userID uint) (*models.Profile, error) {
	// ЭТО НУЖНО РЕАЛИЗОВАТЬ (НО ЭТО НЕ ТОЧНО)
	// if pl.ProfileExists(ctx, id) {
	// 	return nil, errProfileDoNotExists
	// }

	var profile *models.Profile

	fullPath, err := utils.WriteFile(file, folderName)

	if err != nil {
		pl.Logger.Errorf("Something went wrong while writing file of the image , err=%v", err)

		return nil, errProfileNotExists
	}

	err = pgx.BeginFunc(ctx, pl.Pool, func(tx pgx.Tx) error {
		profileInner, err := pl.setProfileAvatarUrl(ctx, tx, userID, fullPath)
		profile = profileInner

		return err
	})

	if err != nil {
		pl.Logger.Errorf("Something went wrong while updating profile url , err=%v", errProfileNotExists)

		return nil, errProfileNotExists
	}

	return profile, nil
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
	p.Avatar = "" // ОПАСНОСТЬ

	return p, nil
}

func (pl *ProfileListWrapper) EditProfile(userID uint, data models.EditProfileNec) (*models.Profile, error) {
	pl.ProfileList.Mux.Lock()
	defer pl.ProfileList.Mux.Unlock()
	log.Println(pl.ProfileList.Profiles[userID])
	log.Println(pl.ProfileList.Profiles[userID-1])
	_, ok := pl.ProfileList.Profiles[userID]
	log.Println(ok)
	if !ok {
		return nil, errProfileNotExists
	}

	old := pl.ProfileList.Profiles[userID]

	old.Name = data.Name
	old.Surname = data.Surname
	old.Avatar = data.Avatar // ОПАСНОСТЬ

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

func NewProfileList(pool *pgxpool.Pool, logger *zap.SugaredLogger) *ProfileListWrapper {
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
					Avatar:         "", // Предполагается, что Image имеет конструктор по умолчанию
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
					Avatar:         "", // Предполагается, что Image имеет конструктор по умолчанию
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
		Pool:   pool,
		Logger: logger,
	}
}
