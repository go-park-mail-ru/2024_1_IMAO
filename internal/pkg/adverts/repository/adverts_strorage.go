package storage

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
)

var (
	errWrongAdvertID      = errors.New("wrong advert ID")
	errWrongCityName      = errors.New("wrong city name")
	errWrongCategoryName  = errors.New("wrong category name")
	errWrongIDinCategory  = errors.New("there is no ad with such id in category")
	errWrongIDinCity      = errors.New("there is no ad with such id in city")
	errWrongAdvertsAmount = errors.New("too many elements specified")
	errAlreadyClosed      = errors.New("advert already closed")
)

type AdvertStorage struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewAdvertStorage(pool *pgxpool.Pool, logger *zap.SugaredLogger) *AdvertStorage {
	return &AdvertStorage{
		pool:   pool,
		logger: logger,
	}
}

func (ads *AdvertStorage) GetAdvertByOnlyByID(ctx context.Context, advertID uint) (*models.ReturningAdvert, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList *models.ReturningAdvert

	city := ""     // ЗАГЛУШКИ
	category := "" // ЗАГЛУШКИ

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.getAdvert(ctx, tx, advertID, city, category)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%v", err))

		return nil, err
	}

	return advertsList, nil
}

func (ads *AdvertStorage) getAdvertImagesURLs(ctx context.Context, tx pgx.Tx, advertID uint) ([]string, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAdvertImagesURLs := `
	SELECT url
	FROM public.advert_image
	WHERE advert_id = $1;`

	logging.LogInfo(logger, "SELECT FROM advert_image")

	rows, err := tx.Query(ctx, SQLAdvertImagesURLs, advertID)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts urls, err=%v", err))

		return nil, err
	}
	defer rows.Close()

	var urlArray []string

	for rows.Next() {
		var returningUrl string

		if err := rows.Scan(&returningUrl); err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while scanning rows of advert images for advert %v, err=%v", advertID, err))

			return nil, err
		}

		urlArray = append(urlArray, returningUrl)
	}

	return urlArray, nil
}

func (ads *AdvertStorage) GetAdvertImagesURLs(ctx context.Context, advertID uint) ([]string, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var urlArray []string

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		urlArrayInner, err := ads.getAdvertImagesURLs(ctx, tx, advertID)
		urlArray = urlArrayInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong getting image urls for advertID=%v , err=%v", advertID, err))

		return nil, err
	}

	return urlArray, nil
}

func (ads *AdvertStorage) insertView(ctx context.Context, tx pgx.Tx, userID, advertID uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLInsertView := `
	INSERT INTO public.view(
		user_id, advert_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, advert_id) DO NOTHING;`

	logging.LogInfo(logger, "INSERT INTO view")

	var err error

	_, err = tx.Exec(ctx, SQLInsertView, userID, advertID)

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing insertInsertView query, err=%v", err))

		return err
	}

	return nil
}

func (ads *AdvertStorage) InsertView(ctx context.Context, userID, advertID uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		err := ads.insertView(ctx, tx, userID, advertID)

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while inserting view, err=%v", err))

		return err
	}

	return nil
}

func (ads *AdvertStorage) getAdvert(ctx context.Context, tx pgx.Tx, advertID uint, city, category string) (*models.ReturningAdvert, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAdvertById := `
		SELECT 
		a.id, 
		a.user_id,
		a.city_id, 
		c.name AS city_name, 
		c.translation AS city_translation, 
		a.category_id, 
		cat.name AS category_name, 
		cat.translation AS category_translation, 
		a.title, 
		a.description, 
		a.price, 
		a.created_time, 
		a.closed_time, 
		a.is_used,
		a.views
		FROM 
		public.advert a
		LEFT JOIN 
		public.city c ON a.city_id = c.id
		LEFT JOIN 
		public.category cat ON a.category_id = cat.id
		WHERE a.id = $1;`

	logging.LogInfo(logger, "SELECT FROM advert, city, category")

	advertLine := tx.QueryRow(ctx, SQLAdvertById, advertID)

	categoryModel := models.Category{}
	cityModel := models.City{}
	advertModel := models.Advert{}

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &cityModel.CityName, &cityModel.Translation,
		&categoryModel.ID, &categoryModel.Name, &categoryModel.Translation, &advertModel.Title, &advertModel.Description, &advertModel.Price,
		&advertModel.CreatedTime, &advertModel.ClosedTime, &advertModel.IsUsed, &advertModel.Views); err != nil {

		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert, err=%v", err))

		return nil, err
	}

	advertModel.CityID = cityModel.ID
	advertModel.CategoryID = categoryModel.ID

	return &models.ReturningAdvert{
		Advert:   advertModel,
		City:     cityModel,
		Category: categoryModel,
	}, nil
}

func (ads *AdvertStorage) GetAdvert(ctx context.Context, advertID uint, city, category string) (*models.ReturningAdvert, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList *models.ReturningAdvert

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.getAdvert(ctx, tx, advertID, city, category)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%v", err))

		return nil, err
	}

	advertsList.Photos, err = ads.GetAdvertImagesURLs(ctx, advertsList.Advert.ID)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting advert images urls , err=%v", err))

		return nil, err
	}

	for i := 0; i < len(advertsList.Photos); i++ {

		image, err := utils.DecodeImage(advertsList.Photos[i])
		advertsList.PhotosIMG = append(advertsList.PhotosIMG, image)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("error occurred while decoding advert_image %v, err = %v", advertsList.Photos[i], err))
		}
	}

	advertsList.Advert.Sanitize()

	return advertsList, nil
}

func (ads *AdvertStorage) getAdvertsByCity(ctx context.Context, tx pgx.Tx, city string, startID, num uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAdvertsByCity := `SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT array_agg(url) FROM (SELECT url FROM advert_image WHERE advert_id = a.id ORDER BY id) AS ordered_images) AS image_urls
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.id >= $1 AND a.advert_status = 'Активно' AND c.translation = $2
	ORDER BY id
	LIMIT $3;
	`
	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image")

	rows, err := tx.Query(ctx, SQLAdvertsByCity, startID, city, num)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts query, err=%v", err))

		return nil, err
	}
	defer rows.Close()

	var adsList []*models.ReturningAdInList
	for rows.Next() {
		returningAdInList := models.ReturningAdInList{}

		photoPad := models.PhotoPad{}

		if err := rows.Scan(&returningAdInList.ID, &returningAdInList.City, &returningAdInList.Category, &returningAdInList.Title, &returningAdInList.Price, &photoPad.Photo); err != nil {
			return nil, err
		}

		returningAdInList.InFavourites = false

		if photoPad.Photo != nil {
			for _, ptr := range photoPad.Photo {
				returningAdInList.Photos = append(returningAdInList.Photos, *ptr)
			}
		}

		for i := 0; i < len(returningAdInList.Photos); i++ {

			image, err := utils.DecodeImage(returningAdInList.Photos[i])
			returningAdInList.PhotosIMG = append(returningAdInList.PhotosIMG, image)
			if err != nil {
				logging.LogError(logger, fmt.Errorf("error occurred while decoding advert_image %v, err = %v", returningAdInList.Photos[i], err))

				return nil, err
			}
		}

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while decoding image, err=%v", err))

			return nil, err
		}

		returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%v", err))

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) getAdvertsByCityAuth(ctx context.Context, tx pgx.Tx, city string, userID, startID, num uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAdvertsByCityAuth := `SELECT a.id, c.translation, category.translation, a.title, a.price,
    (SELECT array_agg(url) FROM (SELECT url FROM advert_image WHERE advert_id = a.id ORDER BY id) AS ordered_images) AS image_urls,
    CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $1 AND f.advert_id = a.id)
         THEN 1 ELSE 0 END AS bool) AS in_favourites
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.id >= $2 AND a.advert_status = 'Активно' AND c.translation = $3
	ORDER BY id
	LIMIT $4;
	`
	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image")

	rows, err := tx.Query(ctx, SQLAdvertsByCityAuth, userID, startID, city, num)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts query, err=%v", err))

		return nil, err
	}
	defer rows.Close()

	var adsList []*models.ReturningAdInList
	for rows.Next() {
		returningAdInList := models.ReturningAdInList{}

		photoPad := models.PhotoPad{}

		if err := rows.Scan(&returningAdInList.ID, &returningAdInList.City, &returningAdInList.Category, &returningAdInList.Title,
			&returningAdInList.Price, &photoPad.Photo, &returningAdInList.InFavourites); err != nil {
			return nil, err
		}

		if photoPad.Photo != nil {
			for _, ptr := range photoPad.Photo {
				returningAdInList.Photos = append(returningAdInList.Photos, *ptr)
			}
		}

		for i := 0; i < len(returningAdInList.Photos); i++ {

			image, err := utils.DecodeImage(returningAdInList.Photos[i])
			returningAdInList.PhotosIMG = append(returningAdInList.PhotosIMG, image)
			if err != nil {
				logging.LogError(logger, fmt.Errorf("error occurred while decoding advert_image %v, err = %v", returningAdInList.Photos[i], err))

				return nil, err
			}
		}

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while decoding image, err=%v", err))

			return nil, err
		}

		returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%v", err))

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) GetAdvertsByCity(ctx context.Context, city string, userID, startID, num uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList []*models.ReturningAdInList

	if userID == 0 {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertsByCity(ctx, tx, city, startID, num)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%v", err))

			return nil, err
		}

	} else {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertsByCityAuth(ctx, tx, city, userID, startID, num)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%v", err))

			return nil, err
		}
	}

	return advertsList, nil
}

func (ads *AdvertStorage) getAdvertsByCategory(ctx context.Context, tx pgx.Tx, category, city string, startID, num uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAdvertsByCity := `SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT array_agg(url) FROM (SELECT url FROM advert_image WHERE advert_id = a.id ORDER BY id) AS ordered_images) AS image_urls
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.id >= $1 AND a.advert_status = 'Активно' AND c.translation = $2 AND category.translation = $3
	LIMIT $4;
	`
	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image")

	rows, err := tx.Query(ctx, SQLAdvertsByCity, startID, city, category, num)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts query, err=%v", err))

		return nil, err
	}
	defer rows.Close()

	var adsList []*models.ReturningAdInList
	for rows.Next() {
		returningAdInList := models.ReturningAdInList{}

		photoPad := models.PhotoPad{}

		if err := rows.Scan(&returningAdInList.ID, &returningAdInList.City, &returningAdInList.Category, &returningAdInList.Title, &returningAdInList.Price, &photoPad.Photo); err != nil {
			return nil, err
		}

		if photoPad.Photo != nil {
			for _, ptr := range photoPad.Photo {
				returningAdInList.Photos = append(returningAdInList.Photos, *ptr)
			}
		}

		for i := 0; i < len(returningAdInList.Photos); i++ {

			image, err := utils.DecodeImage(returningAdInList.Photos[i])
			returningAdInList.PhotosIMG = append(returningAdInList.PhotosIMG, image)
			if err != nil {
				logging.LogError(logger, fmt.Errorf("error occurred while decoding advert_image %v, err = %v", returningAdInList.Photos[i], err))
			}
		}

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while decoding image, err=%v", err))

			return nil, err
		}

		returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%v", err))

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) GetAdvertsByCategory(ctx context.Context, category, city string, startID, num uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList []*models.ReturningAdInList

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.getAdvertsByCategory(ctx, tx, category, city, startID, num)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%v", err))

		return nil, err
	}

	return advertsList, nil
}

func (ads *AdvertStorage) getAdvertsForUserWhereStatusIs(ctx context.Context, tx pgx.Tx, userId, deleted uint) (*models.ReturningAdvertList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	advertStatus := "Активно"
	if deleted == 1 {
		advertStatus = "Продано"
	}

	SQLGetAdvertsForUserWhereStatusIs := `SELECT 
		a.id, 
		a.user_id, 
		a.city_id, 
		a.category_id, 
		a.title, 
		a.description, 
		a.price, 
		a.created_time, 
		a.closed_time, 
		a.is_used, 
		a.advert_status,
		c.name AS city_name,
		c.translation AS city_translation,
		cat.name AS category_name,
		cat.translation AS category_translation,
		(SELECT array_agg(url) FROM (SELECT url FROM advert_image WHERE advert_id = a.id ORDER BY id) AS ordered_images) AS image_urls
	FROM 
		public.advert a
	INNER JOIN 
		city c ON a.city_id = c.id
	INNER JOIN 
		category cat ON a.category_id = cat.id
	WHERE 
		a.user_id = $1 AND a.advert_status = $2;
	`

	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image")

	rows, err := tx.Query(ctx, SQLGetAdvertsForUserWhereStatusIs, userId, advertStatus)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts for user where status is, err=%v", err))

		return nil, err
	}
	defer rows.Close()

	returningAdvertList := models.ReturningAdvertList{}
	for rows.Next() {
		advert := models.Advert{}
		city := models.City{}
		category := models.Category{}
		var status string // ЗАГЛУШКА

		photoPad := models.PhotoPad{}

		if err := rows.Scan(&advert.ID, &advert.UserID, &advert.CityID, &advert.CategoryID, &advert.Title, &advert.Description, &advert.Price, &advert.CreatedTime, &advert.ClosedTime, &advert.IsUsed, &status,
			&city.CityName, &city.Translation, &category.Name, &category.Translation, &photoPad.Photo); err != nil {

			logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%v", err))

			return nil, err
		}

		var photoArray []string
		if photoPad.Photo != nil {
			for _, ptr := range photoPad.Photo {
				photoArray = append(photoArray, *ptr)
			}
		}

		var photoIMGArray []string

		for i := 0; i < len(photoArray); i++ {

			image, err := utils.DecodeImage(photoArray[i])
			photoIMGArray = append(photoIMGArray, image)
			if err != nil {
				logging.LogError(logger, fmt.Errorf("error occurred while decoding advert_image %v, err = %v", photoArray[i], err))
			}
		}

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while decoding image, err=%v", err))

			return nil, err
		}

		advert.Deleted = false

		if status == "Продано" {
			advert.Deleted = true
		}

		city.ID = advert.CityID
		category.ID = advert.CategoryID
		returningAdvert := models.ReturningAdvert{
			Advert:    advert,
			City:      city,
			Category:  category,
			Photos:    photoArray,
			PhotosIMG: photoIMGArray,
		}

		returningAdvertList.AdvertItems = append(returningAdvertList.AdvertItems, &returningAdvert)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%v", err))

		return nil, err
	}

	return &returningAdvertList, nil
}

func (ads *AdvertStorage) GetAdvertsForUserWhereStatusIs(ctx context.Context, userId, deleted uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList []*models.ReturningAdInList

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.getAdvertsForUserWhereStatusIs(ctx, tx, userId, deleted)
		for _, num := range advertsListInner.AdvertItems {
			returningAdInList := models.ReturningAdInList{
				ID:        num.Advert.ID,
				Title:     num.Advert.Title,
				Price:     num.Advert.Price,
				City:      num.City.Translation,
				Category:  num.Category.Translation,
				Photos:    num.Photos,
				PhotosIMG: num.PhotosIMG,
			}

			returningAdInList.Sanitize()

			advertsList = append(advertsList, &returningAdInList)
		}

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%v", err))

		return nil, err
	}

	return advertsList, nil
}

func (ads *AdvertStorage) createAdvert(ctx context.Context, tx pgx.Tx, data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCreateAdvert :=
		`WITH ins AS (
		INSERT INTO advert (user_id, city_id, category_id, title, description, price, is_used)
		SELECT
			$1,
			city.id,
			category.id,
			$2,
			$3,
			$4,
			$5
		FROM
			city
		JOIN
			category ON city.name = $6 AND category.translation = $7
		RETURNING 
			advert.id, 
			advert.user_id,
			advert.city_id, 
			advert.category_id, 
			advert.title, 
			advert.description,
			advert.created_time,
			advert.closed_time, 
			advert.price, 
			advert.is_used
	)
	SELECT ins.*, c.name AS city_name, c.translation AS city_translation, cat.name AS category_name, cat.translation AS category_translation
	FROM ins
	LEFT JOIN public.city c ON ins.city_id = c.id
	LEFT JOIN public.category cat ON ins.category_id = cat.id;`

	logging.LogInfo(logger, "INSERT INTO advert")

	advertLine := tx.QueryRow(ctx, SQLCreateAdvert, data.UserID, data.Title, data.Description, data.Price, data.IsUsed, data.City, data.Category)

	categoryModel := models.Category{}
	cityModel := models.City{}
	advertModel := models.Advert{}

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &categoryModel.ID, &advertModel.Title, &advertModel.Description, &advertModel.CreatedTime,
		&advertModel.ClosedTime, &advertModel.Price, &advertModel.IsUsed, &cityModel.CityName, &cityModel.Translation, &categoryModel.Name, &categoryModel.Translation); err != nil {

		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert, err=%v", err))

		return nil, err
	}

	advertModel.CityID = cityModel.ID
	advertModel.CategoryID = categoryModel.ID

	return &models.ReturningAdvert{
		Advert:   advertModel,
		City:     cityModel,
		Category: categoryModel,
	}, nil
}

func (ads *AdvertStorage) CreateAdvert(ctx context.Context, files []*multipart.FileHeader,
	data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList *models.ReturningAdvert

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.createAdvert(ctx, tx, data)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting advert, err=%v", err))

		return nil, err
	}

	advertsList.Photos, err = ads.SetAdvertImages(ctx, files, "advert_images", advertsList.Advert.ID)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while updating advert image url , err=%v", err))

		return nil, err
	}

	for i := 0; i < len(advertsList.Photos); i++ {
		image, err := utils.DecodeImage(advertsList.Photos[i])
		advertsList.PhotosIMG = append(advertsList.PhotosIMG, image)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("error occurred while decoding advert_image %v, err = %v", advertsList.Photos[i], err))

			return nil, err
		}
	}

	advertsList.Advert.Sanitize()

	return advertsList, nil
}

func (ads *AdvertStorage) setAdvertImage(ctx context.Context, tx pgx.Tx, advertID uint, imageUrl string) (string, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLUpdateProfileAvatarURL := `
	INSERT INTO advert_image (url, advert_id)
	VALUES 
    ($1, $2)
	RETURNING url;`

	logging.LogInfo(logger, "INSERT INTO advert_image")

	var returningUrl string

	urlLine := tx.QueryRow(ctx, SQLUpdateProfileAvatarURL, imageUrl, advertID)

	if err := urlLine.Scan(&returningUrl); err != nil {

		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert image , err=%v", err))

		return "", err
	}

	return returningUrl, nil
}

func (ads *AdvertStorage) deleteAllImagesForAdvertFromLocalStorage(ctx context.Context, tx pgx.Tx, advertID uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAdvertImagesURLs := `
	SELECT url
	FROM public.advert_image
	WHERE advert_id = $1;`

	logging.LogInfo(logger, "SELECT FROM advert_image")

	rows, err := tx.Query(ctx, SQLAdvertImagesURLs, advertID)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts urls, err=%v", err))

		return err
	}
	defer rows.Close()

	var oldUrl interface{}

	for rows.Next() {
		if err := rows.Scan(&oldUrl); err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while scanning rows for deleting advert images for advert %v, err=%v", advertID, err))

			return err
		}

		if oldUrl != nil {
			err := os.Remove(oldUrl.(string))
			if err != nil {
				logging.LogError(logger, fmt.Errorf("something went wrong while deleting image %v, err=%v", oldUrl.(string), err))

				return err
			}
		}
	}

	return nil
}

func (ads *AdvertStorage) deleteAllImagesForAdvertFromDatabase(ctx context.Context, tx pgx.Tx, advertID uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCreateProfile := `DELETE FROM public.advert_image WHERE advert_id = $1;`

	logging.LogInfo(logger, "DELETE FROM advert_image")

	var err error

	_, err = tx.Exec(ctx, SQLCreateProfile, advertID)

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing delete advert_image query, err=%v", err))

		return err
	}

	return nil
}

func (ads *AdvertStorage) SetAdvertImages(ctx context.Context, files []*multipart.FileHeader, folderName string, advertID uint) ([]string, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		return ads.deleteAllImagesForAdvertFromLocalStorage(ctx, tx, advertID)
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong deleting AllImagesForAdvertFromLocalStorage , err=%v", err))

		return nil, err
	}

	err = pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		return ads.deleteAllImagesForAdvertFromDatabase(ctx, tx, advertID)
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while deleting AllImagesForAdvertFromDatabase , err=%v", err))

		return nil, err
	}

	var urlArray []string

	for i := 0; i < len(files); i++ {
		var url string

		fullPath, err := utils.WriteFile(files[i], folderName)

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while writing file of the image , err=%v", err))

			return nil, err
		}

		err = pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			urlInner, err := ads.setAdvertImage(ctx, tx, advertID, fullPath)
			url = urlInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while updating profile url , err=%v", err))

			return nil, err
		}

		urlArray = append(urlArray, url)
	}

	return urlArray, nil
}

func (ads *AdvertStorage) editAdvert(ctx context.Context, tx pgx.Tx, data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLUpdateAdvert :=
		`WITH upd AS (
			UPDATE advert
			SET 
				city_id = city.id,
				category_id = category.id,
				title = $1,
				description = $2,
				price = $3,
				is_used = $4
			 FROM
				city
			JOIN
				category ON city.name = $5 AND category.translation = $6
			WHERE 
				advert.id = $7
			RETURNING 
				advert.id, 
				advert.user_id,
				advert.city_id, 
				advert.category_id, 
				advert.title, 
				advert.description,
				advert.created_time,
				advert.closed_time,
				advert.price, 
				advert.is_used
		)
		SELECT 
			upd.*, 
			c.name AS city_name, 
			c.translation AS city_translation, 
			cat.name AS category_name, 
			cat.translation AS category_translation
		FROM 
			upd
		LEFT JOIN 
			public.city c ON upd.city_id = c.id
		LEFT JOIN 
			public.category cat ON upd.category_id = cat.id;`

	logging.LogInfo(logger, "UPDATE advert")

	advertLine := tx.QueryRow(ctx, SQLUpdateAdvert, data.Title, data.Description, data.Price, data.IsUsed, data.City, data.Category, data.ID)

	categoryModel := models.Category{}
	cityModel := models.City{}
	advertModel := models.Advert{}

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &categoryModel.ID, &advertModel.Title, &advertModel.Description, &advertModel.CreatedTime,
		&advertModel.ClosedTime, &advertModel.Price, &advertModel.IsUsed, &cityModel.CityName, &cityModel.Translation, &categoryModel.Name, &categoryModel.Translation); err != nil {

		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert, err=%v", err))

		return nil, err
	}

	advertModel.CityID = cityModel.ID
	advertModel.CategoryID = categoryModel.ID

	return &models.ReturningAdvert{
		Advert:   advertModel,
		City:     cityModel,
		Category: categoryModel,
	}, nil
}

func (ads *AdvertStorage) EditAdvert(ctx context.Context, files []*multipart.FileHeader, data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList *models.ReturningAdvert

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.editAdvert(ctx, tx, data)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while updating advert, err=%v", err))

		return nil, err
	}

	advertsList.Photos, err = ads.SetAdvertImages(ctx, files, "advert_images", data.ID)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while updating advert image url , err=%v", err))

		return nil, err
	}

	for i := 0; i < len(advertsList.Photos); i++ {

		image, err := utils.DecodeImage(advertsList.Photos[i])
		advertsList.PhotosIMG = append(advertsList.PhotosIMG, image)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("error occurred while decoding advert_image %v, err = %v", advertsList.Photos[i], err))

			return nil, err
		}
	}

	advertsList.Advert.Sanitize()

	return advertsList, nil
}

func (ads *AdvertStorage) closeAdvert(ctx context.Context, tx pgx.Tx, advertID uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCloseAdvert := `UPDATE public.advert	SET  advert_status='Скрыто'	WHERE id = $1;`

	logging.LogInfo(logger, "UPDATE advert")

	var err error

	_, err = tx.Exec(ctx, SQLCloseAdvert, advertID)

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing close advert query, err=%v", err))

		return err
	}

	return nil
}

func (ads *AdvertStorage) CloseAdvert(ctx context.Context, advertID uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		err := ads.closeAdvert(ctx, tx, advertID)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while closing advert, err=%v", err))

			return err
		}

		return nil
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while closing advert, err=%v", err))

		return err
	}

	return nil
}

func (ads *AdvertStorage) deleteAdvert(ctx context.Context, tx pgx.Tx, user *models.User) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCreateUser := `INSERT INTO public."user"(email, password_hash) VALUES ($1, $2);`
	logging.LogInfo(logger, "INSERT INTO public.user")
	var err error

	_, err = tx.Exec(ctx, SQLCreateUser, user.Email, user.PasswordHash)

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing create user query, err=%v", err))

		return err
	}

	return nil
}

// func (ads *AdvertStorage) DeleteAdvert(advertID uint) error {
// 	if advertID > ads.AdvertsList.AdvertsCounter || ads.AdvertsList.Adverts[advertID-1].Deleted {
// 		return errWrongAdvertID
// 	}

// 	ads.AdvertsList.Adverts[advertID-1].Deleted = true

// 	return nil
// }

func AddAdvert(ads *models.AdvertsList, advert *models.Advert) {

	ads.Adverts = append(ads.Adverts, advert)

}
