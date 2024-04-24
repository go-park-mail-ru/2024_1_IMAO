package storage

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"mime/multipart"
	"os"
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

type ProfileStorage struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewProfileStorage(pool *pgxpool.Pool, logger *zap.SugaredLogger) *ProfileStorage {
	return &ProfileStorage{
		pool:   pool,
		logger: logger,
	}
}

func (pl *ProfileStorage) createProfile(ctx context.Context, tx pgx.Tx, profile *models.Profile) error {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := pl.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLCreateProfile := `INSERT INTO public.profile(user_id) VALUES ($1);`
	childLogger.Infof(`INSERT INTO public.profile(user_id) VALUES (%s);`, profile.UserID)

	var err error

	_, err = tx.Exec(ctx, SQLCreateProfile, profile.UserID)

	if err != nil {
		childLogger.Errorf("Something went wrong while executing create profile query, err=%v", err)

		return fmt.Errorf("Something went wrong while executing create profile query", err)
	}

	return nil
}

func (pl *ProfileStorage) CreateProfile(ctx context.Context, userID uint) *models.Profile {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := pl.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	profile := models.Profile{
		UserID: userID,
	}

	err := pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		err := pl.createProfile(ctx, tx, &profile)
		if err != nil {
			childLogger.Errorf("Something went wrong while creating profile, err=%v", err)
			return fmt.Errorf("Something went wrong while creating profile", err)
		}
		id, err := repository.GetLastValSeq(ctx, tx, NameSeqProfile)
		if err != nil {
			childLogger.Errorf("Something went wrong getting user id from seq, err=%v", err)

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

func (pl *ProfileStorage) getProfileByUserID(ctx context.Context, tx pgx.Tx, id uint) (*models.Profile, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := pl.logger.With(
		zap.String("requestUUID", requestUUID),
	)

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
			c.translation AS city_translation,
			(SELECT COUNT(*) FROM subscription WHERE user_id_subscriber = $1 ) AS subscriber_count,
			(SELECT COUNT(*) FROM subscription WHERE user_id_merchant = $1 ) AS subscription_count,
			(SELECT COUNT(*) FROM public.review JOIN public.advert ON review.advert_id = advert.id JOIN public.user ON advert.user_id = "user".id WHERE "user".id = $1) AS review_count,
			(SELECT COUNT(*) FROM advert WHERE user_id = $1 AND advert_status = 'Активно') AS active_ads_count,
			(SELECT COUNT(*) FROM advert WHERE user_id = $1 AND advert_status = 'Продано') AS sold_ads_count
		FROM 
			public.profile p
		INNER JOIN 
			public.city c
		ON 
			p.city_id = c.id
		WHERE 
			p.user_id = $1`
	childLogger.Infof(`
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
		c.translation AS city_translation,
		(SELECT COUNT(*) FROM subscription WHERE user_id_subscriber = $1 ) AS subscriber_count,
		(SELECT COUNT(*) FROM subscription WHERE user_id_merchant = $1 ) AS subscription_count,
		(SELECT COUNT(*) FROM public.review JOIN public.advert ON review.advert_id = advert.id JOIN public.user ON advert.user_id = "user".id WHERE "user".id = $1) AS review_count,
		(SELECT COUNT(*) FROM advert WHERE user_id = $1 AND advert_status = 'Активно') AS active_ads_count,
		(SELECT COUNT(*) FROM advert WHERE user_id = $1 AND advert_status = 'Продано') AS sold_ads_count
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
		&profilePad.Surname, &profile.RegisterTime, &profile.Approved, &profilePad.Avatar, &city.CityName, &city.Translation, &profile.SubersCount,
		&profile.SubonsCount, &profile.ReactionsCount, &profile.ActiveAddsCount, &profile.SoldAddsCount); err != nil {

		childLogger.Errorf("Something went wrong while scanning profile, err=%v", err)

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
	//profile.ReactionsCount = 10
	//profile.Approved = true
	profile.MerchantsName = nameToInsert
	//profile.SubersCount = rand.Intn(10)
	//profile.SubonsCount = rand.Intn(10)

	return &profile, nil
}

func (pl *ProfileStorage) GetProfileByUserID(ctx context.Context, userID uint) (*models.Profile, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := pl.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var profile *models.Profile

	err := pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		profileInner, err := pl.getProfileByUserID(ctx, tx, userID)
		profile = profileInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while getting profile by UserID , err=%v", errProfileNotExists)

		return nil, errProfileNotExists
	}

	profile.AvatarIMG, err = utils.DecodeImage(profile.Avatar)
	if err != nil {
		childLogger.Errorf("Error occurred while decoding avatar image, err = %v", err)
	}

	return profile, nil
}

func (pl *ProfileStorage) setProfileCity(ctx context.Context, tx pgx.Tx, userID uint, data models.City) (*models.Profile, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := pl.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLUpdateProfileCity := `UPDATE public.profile p
	SET city_id = $1
	FROM public.city c
	WHERE c.id = $1 AND p.user_id = $2
	RETURNING p.id, p.user_id, p.city_id, p.phone, p.name, p.surname, p.regtime, p.verified, p.avatar_url, c.name AS city_name, c.translation AS city_translation;`

	childLogger.Infof(`UPDATE public.profile p
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

		childLogger.Errorf("Something went wrong while scanning profile lines, err=%v", err)

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

func (pl *ProfileStorage) SetProfileCity(ctx context.Context, userID uint, data models.City) (*models.Profile, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := pl.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var profile *models.Profile

	err := pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		profileInner, err := pl.setProfileCity(ctx, tx, userID, data)
		profile = profileInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while updating profile city, err=%v", errProfileNotExists)

		return nil, errProfileNotExists
	}

	profile.AvatarIMG, err = utils.DecodeImage(profile.Avatar)
	if err != nil {
		childLogger.Errorf("Error occurred while decoding avatar image, err = %v", err)
	}

	return profile, nil
}

func (pl *ProfileStorage) setProfilePhone(ctx context.Context, tx pgx.Tx, userID uint, data models.SetProfilePhoneNec) (*models.Profile, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := pl.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLUpdateProfilePhone := `UPDATE public.profile p
	SET phone = $1
	FROM public.city c
	WHERE c.id = p.city_id AND p.user_id = $2
	RETURNING p.id, p.user_id, p.city_id, p.phone, p.name, p.surname, p.regtime, p.verified, p.avatar_url, c.name AS city_name, c.translation AS city_translation;`

	childLogger.Infof(`UPDATE public.profile p
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

		childLogger.Errorf("Something went wrong while scanning profile lines , err=%v", err)

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

func (pl *ProfileStorage) SetProfilePhone(ctx context.Context, userID uint, data models.SetProfilePhoneNec) (*models.Profile, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := pl.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var profile *models.Profile

	err := pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		profileInner, err := pl.setProfilePhone(ctx, tx, userID, data)
		profile = profileInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while updating profile phone , err=%v", errProfileNotExists)

		return nil, errProfileNotExists
	}

	profile.AvatarIMG, err = utils.DecodeImage(profile.Avatar)
	if err != nil {
		childLogger.Errorf("Error occurred while decoding avatar image, err = %v", err)
	}

	return profile, nil
}

func (pl *ProfileStorage) setProfileAvatarUrl(ctx context.Context, tx pgx.Tx, userID uint, avatar string) (string, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := pl.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLUpdateProfileAvatarURL := `
	UPDATE public.profile p
	SET avatar_url = $1
	WHERE p.user_id = $2
	RETURNING avatar_url;`

	childLogger.Infof(`
	UPDATE public.profile p
	SET avatar_url = %s
	WHERE p.user_id = %s
	RETURNING avatar_url;`, avatar, userID)

	var url string

	urlLine := tx.QueryRow(ctx, SQLUpdateProfileAvatarURL, avatar, userID)

	if err := urlLine.Scan(&url); err != nil {

		childLogger.Errorf("Something went wrong while scanning url line , err=%v", err)

		return "", err
	}

	return url, nil
}

func (pl *ProfileStorage) deleteAvatar(ctx context.Context, tx pgx.Tx, userID uint) error {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := pl.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLGetAvatarURL := `
	SELECT p.avatar_url
	FROM public.profile p
	WHERE p.user_id = $1`

	childLogger.Infof(`
	SELECT p.avatar_url
	FROM public.profile p
	WHERE p.user_id = %s`, userID)

	var oldUrl interface{}

	urlLine := tx.QueryRow(ctx, SQLGetAvatarURL, userID)
	if err := urlLine.Scan(&oldUrl); err != nil {
		childLogger.Errorf("Something went wrong while deleting url , err=%v", err)

		return err
	}

	if oldUrl != nil {
		os.Remove(oldUrl.(string))
	}
	return nil
}

func (pl *ProfileStorage) SetProfileAvatarUrl(ctx context.Context, file *multipart.FileHeader, folderName string,
	userID uint) (string, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := pl.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	err := pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		return pl.deleteAvatar(ctx, tx, userID)
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while updating profile url , err=%v", err)

		return "", err
	}

	var url string
	fullPath, err := utils.WriteFile(file, folderName)

	if err != nil {
		childLogger.Errorf("Something went wrong while writing file of the image , err=%v", err)

		return "", errProfileNotExists
	}

	err = pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		urlInner, err := pl.setProfileAvatarUrl(ctx, tx, userID, fullPath)
		url = urlInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while updating profile url , err=%v", errProfileNotExists)

		return "", err
	}

	return url, nil
}

func (pl *ProfileStorage) setProfileInfo(ctx context.Context, tx pgx.Tx, userID uint,
	data models.EditProfileNec) (*models.Profile, error) {

	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := pl.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLUpdateProfileInfo := `
		UPDATE public.profile p
		SET 
			name = $1,
			surname = $2
		FROM public.city c
		WHERE c.id = p.city_id AND p.user_id = $3
		RETURNING p.id, p.user_id, p.city_id, p.phone, p.name, p.surname, p.regtime, p.verified, p.avatar_url, 
			c.name AS city_name, c.translation AS city_translation;`

	childLogger.Infof(`UPDATE public.profile p
		SET 
			name = %s,
			surname = %s
		FROM public.city c
		WHERE c.id = p.city_id AND p.user_id = %s
		RETURNING p.id, p.user_id, p.city_id, p.phone, p.name, p.surname, p.regtime, p.verified, p.avatar_url, 
			c.name AS city_name, c.translation AS city_translation;`, data.Name, data.Surname, userID)

	profileLine := tx.QueryRow(ctx, SQLUpdateProfileInfo, data.Name, data.Surname, userID)

	profile := models.Profile{}
	city := models.City{}
	profilePad := models.ProfilePad{}

	if err := profileLine.Scan(&profile.ID, &profile.UserID, &city.ID, &profilePad.Phone, &profilePad.Name,
		&profilePad.Surname, &profile.RegisterTime, &profile.Approved, &profilePad.Avatar,
		&city.CityName, &city.Translation); err != nil {

		childLogger.Errorf("Something went wrong while scanning profile lines , err=%v", err)

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

func (pl *ProfileStorage) SetProfileInfo(ctx context.Context, userID uint, file *multipart.FileHeader,
	data models.EditProfileNec) (*models.Profile, error) {

	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := pl.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var profile *models.Profile
	var err error

	if file != nil {
		data.Avatar, err = pl.SetProfileAvatarUrl(ctx, file, "avatars", userID)
		if err != nil {
			childLogger.Errorf("Something went wrong while updating profile url , err=%v", err)

			return nil, err
		}
	}

	err = pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		profileInner, err := pl.setProfileInfo(ctx, tx, userID, data)
		profile = profileInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while updating profile url , err=%v", errProfileNotExists)

		return nil, errProfileNotExists
	}

	profile.AvatarIMG, err = utils.DecodeImage(profile.Avatar)
	if err != nil {
		childLogger.Errorf("Error occurred while decoding avatar image, err = %v", err)
	}

	return profile, nil
}
