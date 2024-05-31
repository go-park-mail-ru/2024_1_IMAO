//nolint:all
package storage

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	errProfileNotExists = errors.New("profile does not exist")
	NameSeqProfile      = pgx.Identifier{"public", "profile_id_seq"} //nolint:gochecknoglobals
)

type ProfileStorage struct {
	pool    *pgxpool.Pool
	metrics *mymetrics.DatabaseMetrics
}

func NewProfileStorage(pool *pgxpool.Pool, metrics *mymetrics.DatabaseMetrics) *ProfileStorage {
	return &ProfileStorage{
		pool:    pool,
		metrics: metrics,
	}
}

func (pl *ProfileStorage) createProfile(ctx context.Context, tx pgx.Tx, profile *models.Profile) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCreateProfile := `INSERT INTO public.profile(user_id) VALUES ($1);`

	logging.LogInfo(logger, "INSERT INTO profile")

	var err error

	start := time.Now()
	_, err = tx.Exec(ctx, SQLCreateProfile, profile.UserID)
	pl.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing create profile query, err=%w", err))
		pl.metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

func (pl *ProfileStorage) CreateProfile(ctx context.Context, userID uint) *models.Profile {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	profile := models.Profile{
		UserID: userID,
	}

	err := pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		err := pl.createProfile(ctx, tx, &profile)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while creating profile, err=%w", err))

			return err
		}

		id, err := repository.GetLastValSeq(ctx, tx, NameSeqProfile)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong getting user id from seq, err=%w", err))

			return err
		}

		profile.ID = uint(id)

		return nil
	})

	if err != nil {

		return nil
	}

	profile.Sanitize()

	return &profile
}

func (pl *ProfileStorage) getProfileByUserID(ctx context.Context, tx pgx.Tx, id uint) (*models.Profile, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

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
			(SELECT COUNT(*) FROM subscription WHERE user_id_merchant = $1) AS subscriber_count,
			(SELECT COUNT(*) FROM subscription WHERE user_id_subscriber = $1) AS subscription_count,
			(SELECT COUNT(*) FROM public.review JOIN public.advert ON review.advert_id = advert.id JOIN public.user ON advert.user_id = "user".id WHERE "user".id = $1) AS review_count,
			(SELECT COUNT(*) FROM advert WHERE user_id = $1 AND advert_status = 'Активно') AS active_ads_count,
			(SELECT COUNT(*) FROM advert WHERE user_id = $1 AND advert_status = 'Продано') AS sold_ads_count,
			p.rating,
			p.cart_adverts_number,
			p.fav_adverts_number
		FROM 
			public.profile p
		INNER JOIN 
			public.city c
		ON 
			p.city_id = c.id
		WHERE 
			p.user_id = $1`

	logging.LogInfo(logger, "SELECT FROM profile, city, subscription, review")

	start := time.Now()
	profileLine := tx.QueryRow(ctx, SQLUserById, id)
	pl.metrics.AddDuration(funcName, time.Since(start))

	profile := models.Profile{}
	city := models.City{}
	profilePad := models.ProfilePad{}

	if err := profileLine.Scan(&profile.ID, &profile.UserID, &city.ID, &profilePad.Phone, &profilePad.Name,
		&profilePad.Surname, &profile.RegisterTime, &profile.Approved, &profilePad.Avatar, &city.CityName, &city.Translation, &profile.SubersCount,
		&profile.SubonsCount, &profile.ReactionsCount, &profile.ActiveAddsCount, &profile.SoldAddsCount, &profile.Rating,
		&profile.CartNum, &profile.FavNum); err != nil {

		logging.LogError(logger, fmt.Errorf("something went wrong while scanning profile, err=%w", err))

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
	profile.MerchantsName = nameToInsert
	profile.IsSubscribed = false

	return &profile, nil
}

func (pl *ProfileStorage) getProfileByUserIDAuth(ctx context.Context, tx pgx.Tx, profileId, userId uint) (*models.Profile, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	fmt.Println("userId", userId, " profileId", profileId)

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
			(SELECT COUNT(*) FROM subscription WHERE user_id_merchant = $1) AS subscriber_count,
			(SELECT COUNT(*) FROM subscription WHERE user_id_subscriber = $1) AS subscription_count,
			(SELECT COUNT(*) FROM public.review JOIN public.advert ON review.advert_id = advert.id JOIN public.user ON advert.user_id = "user".id WHERE "user".id = $1) AS review_count,
			(SELECT COUNT(*) FROM advert WHERE user_id = $1 AND advert_status = 'Активно') AS active_ads_count,
			(SELECT COUNT(*) FROM advert WHERE user_id = $1 AND advert_status = 'Продано') AS sold_ads_count,
			CAST(CASE WHEN EXISTS (SELECT 1 FROM subscription s WHERE s.user_id_merchant = $1 AND s.user_id_subscriber = $2)
			THEN 1 ELSE 0 END AS bool) AS is_subscribed,
			p.rating,
			p.cart_adverts_number,
			p.fav_adverts_number
		FROM 
			public.profile p
		INNER JOIN 
			public.city c
		ON 
			p.city_id = c.id
		WHERE 
			p.user_id = $1`

	logging.LogInfo(logger, "SELECT FROM profile, city, subscription, review")

	start := time.Now()
	profileLine := tx.QueryRow(ctx, SQLUserById, profileId, userId)
	pl.metrics.AddDuration(funcName, time.Since(start))

	profile := models.Profile{}
	city := models.City{}
	profilePad := models.ProfilePad{}

	if err := profileLine.Scan(&profile.ID, &profile.UserID, &city.ID, &profilePad.Phone, &profilePad.Name,
		&profilePad.Surname, &profile.RegisterTime, &profile.Approved, &profilePad.Avatar, &city.CityName, &city.Translation, &profile.SubersCount,
		&profile.SubonsCount, &profile.ReactionsCount, &profile.ActiveAddsCount, &profile.SoldAddsCount, &profile.IsSubscribed, &profile.Rating,
		&profile.CartNum, &profile.FavNum); err != nil {

		logging.LogError(logger, fmt.Errorf("something went wrong while scanning profile, err=%w", err))

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
	profile.MerchantsName = nameToInsert

	fmt.Println("profile.IsSubscribed", profile.IsSubscribed)

	return &profile, nil
}

func (pl *ProfileStorage) GetProfileByUserID(ctx context.Context, profileId, userId uint) (*models.Profile, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var profile *models.Profile

	if userId == 0 {

		err := pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
			profileInner, err := pl.getProfileByUserID(ctx, tx, profileId)
			profile = profileInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting profile by UserID , err=%w", errProfileNotExists))

			return nil, errProfileNotExists
		}

		profile.AvatarIMG, err = utils.DecodeImage(profile.Avatar)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("error occurred while decoding avatar image, err = %w", err))

			return nil, err
		}
	} else {

		err := pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
			profileInner, err := pl.getProfileByUserIDAuth(ctx, tx, profileId, userId)
			profile = profileInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting profile by UserID , err=%w", errProfileNotExists))

			return nil, errProfileNotExists
		}

		profile.AvatarIMG, err = utils.DecodeImage(profile.Avatar)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("error occurred while decoding avatar image, err = %w", err))

			return nil, err
		}

	}

	profile.Sanitize()

	return profile, nil
}

func (pl *ProfileStorage) setProfileCity(ctx context.Context, tx pgx.Tx, userID uint, data models.City) (*models.Profile, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLUpdateProfileCity := `UPDATE public.profile p
	SET city_id = $1
	FROM public.city c
	WHERE c.id = $1 AND p.user_id = $2
	RETURNING p.id, p.user_id, p.city_id, p.phone, p.name, p.surname, p.regtime, p.verified, p.avatar_url, p.rating, c.name AS city_name, c.translation AS city_translation;`

	logging.LogInfo(logger, "UPDATE profile")

	start := time.Now()
	profileLine := tx.QueryRow(ctx, SQLUpdateProfileCity, data.ID, userID)
	pl.metrics.AddDuration(funcName, time.Since(start))

	profile := models.Profile{}
	city := models.City{}
	profilePad := models.ProfilePad{}

	if err := profileLine.Scan(&profile.ID, &profile.UserID, &city.ID, &profilePad.Phone, &profilePad.Name,
		&profilePad.Surname, &profile.RegisterTime, &profile.Approved, &profilePad.Avatar, &profile.Rating, &city.CityName, &city.Translation); err != nil {

		logging.LogError(logger, fmt.Errorf("something went wrong while scanning profile lines, err=%w", err))

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

	//rand.Seed(time.Now().UnixNano()) //nolint:staticcheck
	//profile.Rating = float32(math.Round((rand.Float64()*4+1)*100) / 100)
	profile.ReactionsCount = 10
	//profile.Approved = true
	profile.MerchantsName = nameToInsert
	profile.SubersCount = rand.Intn(10)
	profile.SubonsCount = rand.Intn(10)

	return &profile, nil
}

func (pl *ProfileStorage) SetProfileCity(ctx context.Context, userID uint, data models.City) (*models.Profile, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var profile *models.Profile

	err := pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		profileInner, err := pl.setProfileCity(ctx, tx, userID, data)
		profile = profileInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while updating profile city, err=%w", errProfileNotExists))

		return nil, errProfileNotExists
	}

	profile.AvatarIMG, err = utils.DecodeImage(profile.Avatar)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("error occurred while decoding avatar image, err = %w", err))

		return nil, err
	}

	profile.Sanitize()

	return profile, nil
}

func (pl *ProfileStorage) setProfilePhone(ctx context.Context, tx pgx.Tx, userID uint, data models.SetProfilePhoneNec) (*models.Profile, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLUpdateProfilePhone := `UPDATE public.profile p
	SET phone = $1
	FROM public.city c
	WHERE c.id = p.city_id AND p.user_id = $2
	RETURNING p.id, p.user_id, p.city_id, p.phone, p.name, p.surname, p.regtime, p.verified,  p.avatar_url, p.rating, c.name AS city_name, c.translation AS city_translation;`

	logging.LogInfo(logger, "UPDATE profile")

	start := time.Now()
	profileLine := tx.QueryRow(ctx, SQLUpdateProfilePhone, data.Phone, userID)
	pl.metrics.AddDuration(funcName, time.Since(start))

	profile := models.Profile{}
	city := models.City{}
	profilePad := models.ProfilePad{}

	if err := profileLine.Scan(&profile.ID, &profile.UserID, &city.ID, &profilePad.Phone, &profilePad.Name,
		&profilePad.Surname, &profile.RegisterTime, &profile.Approved, &profilePad.Avatar, &profile.Rating, &city.CityName, &city.Translation); err != nil {

		logging.LogError(logger, fmt.Errorf("something went wrong while scanning profile lines , err=%w", err))

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

	//rand.Seed(time.Now().UnixNano()) //nolint:staticcheck
	//profile.Rating = float32(math.Round((rand.Float64()*4+1)*100) / 100)
	profile.ReactionsCount = 10
	//profile.Approved = true
	profile.MerchantsName = nameToInsert
	profile.SubersCount = rand.Intn(10)
	profile.SubonsCount = rand.Intn(10)

	return &profile, nil
}

func (pl *ProfileStorage) SetProfilePhone(ctx context.Context, userID uint, data models.SetProfilePhoneNec) (*models.Profile, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var profile *models.Profile

	err := pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		profileInner, err := pl.setProfilePhone(ctx, tx, userID, data)
		profile = profileInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while updating profile phone , err=%w", errProfileNotExists))

		return nil, errProfileNotExists
	}

	profile.AvatarIMG, err = utils.DecodeImage(profile.Avatar)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("error occurred while decoding avatar image, err = %w", err))

		return nil, err
	}

	profile.Sanitize()

	return profile, nil
}

func (pl *ProfileStorage) setProfileAvatarUrl(ctx context.Context, tx pgx.Tx, userID uint, avatar string) (string, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLUpdateProfileAvatarURL := `
	UPDATE public.profile p
	SET avatar_url = $1
	WHERE p.user_id = $2
	RETURNING avatar_url;`

	logging.LogInfo(logger, "UPDATE profile")

	var url string

	start := time.Now()
	urlLine := tx.QueryRow(ctx, SQLUpdateProfileAvatarURL, avatar, userID)
	pl.metrics.AddDuration(funcName, time.Since(start))

	if err := urlLine.Scan(&url); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning url line , err=%w", err))

		return "", err
	}

	return url, nil
}

func (pl *ProfileStorage) deleteAvatar(ctx context.Context, tx pgx.Tx, userID uint) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLGetAvatarURL := `
	SELECT p.avatar_url
	FROM public.profile p
	WHERE p.user_id = $1`

	logging.LogInfo(logger, "SELECT FROM profile")

	var oldUrl interface{}

	start := time.Now()
	urlLine := tx.QueryRow(ctx, SQLGetAvatarURL, userID)
	pl.metrics.AddDuration(funcName, time.Since(start))

	if err := urlLine.Scan(&oldUrl); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while deleting url , err=%w", err))

		return err
	}

	if oldUrl != nil {
		os.Remove(oldUrl.(string))
	}
	return nil
}

func (pl *ProfileStorage) SetProfileAvatarUrl(ctx context.Context, fullPath string, userID uint) (string, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		return pl.deleteAvatar(ctx, tx, userID)
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while updating profile url , err=%w", err))

		return "", err
	}

	var url string

	err = pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		urlInner, err := pl.setProfileAvatarUrl(ctx, tx, userID, fullPath)
		url = urlInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while updating profile url , err=%w", errProfileNotExists))

		return "", err
	}

	return url, nil
}

func (pl *ProfileStorage) setProfileInfo(ctx context.Context, tx pgx.Tx, userID uint,
	data models.EditProfileNec) (*models.Profile, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLUpdateProfileInfo := `
		UPDATE public.profile p
		SET 
			name = $1,
			surname = $2
		FROM public.city c
		WHERE c.id = p.city_id AND p.user_id = $3
		RETURNING p.id, p.user_id, p.city_id, p.phone, p.name, p.surname, p.regtime, p.verified, p.avatar_url, 
			c.name AS city_name, c.translation AS city_translation;`

	logging.LogInfo(logger, "UPDATE profile")

	start := time.Now()
	profileLine := tx.QueryRow(ctx, SQLUpdateProfileInfo, data.Name, data.Surname, userID)
	pl.metrics.AddDuration(funcName, time.Since(start))

	profile := models.Profile{}
	city := models.City{}
	profilePad := models.ProfilePad{}

	if err := profileLine.Scan(&profile.ID, &profile.UserID, &city.ID, &profilePad.Phone, &profilePad.Name,
		&profilePad.Surname, &profile.RegisterTime, &profile.Approved, &profilePad.Avatar,
		&city.CityName, &city.Translation); err != nil {

		logging.LogError(logger, fmt.Errorf("something went wrong while scanning profile lines , err=%w", err))

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
	profile.Rating = float32(math.Round((rand.Float64()*4+1)*100) / 100)
	profile.ReactionsCount = 10
	//profile.Approved = true
	profile.MerchantsName = nameToInsert
	profile.SubersCount = rand.Intn(10)
	profile.SubonsCount = rand.Intn(10)

	return &profile, nil
}

func (pl *ProfileStorage) SetProfileInfo(ctx context.Context, userID uint,
	data models.EditProfileNec) (*models.Profile, error) {

	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var (
		profile *models.Profile
		err     error
	)

	if data.Avatar != "" {
		data.Avatar, err = pl.SetProfileAvatarUrl(ctx, data.Avatar, userID)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while updating profile url , err=%w", err))

			return nil, err
		}
	}

	err = pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		profileInner, err := pl.setProfileInfo(ctx, tx, userID, data)
		profile = profileInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while updating profile url , err=%w", errProfileNotExists))

		return nil, errProfileNotExists
	}

	profile.AvatarIMG, err = utils.DecodeImage(profile.Avatar)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("error occurred while decoding avatar image, err = %w", err))

		return nil, err
	}

	profile.Sanitize()

	return profile, nil
}

func (pl *ProfileStorage) appendSubByIDs(ctx context.Context, tx pgx.Tx, userID uint, merchantID uint) bool {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAddToCart := `WITH deletion AS (
		DELETE FROM public.subscription
		WHERE user_id_subscriber = $1 AND user_id_merchant = $2
		RETURNING user_id_subscriber, user_id_merchant
	)
	INSERT INTO public.subscription (user_id_subscriber, user_id_merchant)
	SELECT $1, $2
	WHERE NOT EXISTS (
		SELECT 1 FROM deletion
	) RETURNING true;
	`
	logging.LogInfo(logger, "DELETE or SELECT FROM subscription")

	start := time.Now()
	userLine := tx.QueryRow(ctx, SQLAddToCart, userID, merchantID)
	pl.metrics.AddDuration(funcName, time.Since(start))

	added := false

	if err := userLine.Scan(&added); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning subscriber added, err=%w", err))
		pl.metrics.IncreaseErrors(funcName)

		return false
	}

	return added
}

func (pl *ProfileStorage) AppendSubByIDs(ctx context.Context, userID uint, advertID uint) bool {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var (
		added bool
		err   error
	)

	err = pgx.BeginFunc(ctx, pl.pool, func(tx pgx.Tx) error {
		addedInner := pl.appendSubByIDs(ctx, tx, userID, advertID)
		added = addedInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing subscriber add to subscriber list, err=%w",
			err))
	}

	return added
}
