package storage

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"strings"

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
	pool *pgxpool.Pool
}

func NewAdvertStorage(pool *pgxpool.Pool) *AdvertStorage {
	return &AdvertStorage{
		pool: pool,
	}
}

func (ads *AdvertStorage) getAdvertOnlyByID(ctx context.Context, tx pgx.Tx, advertID uint) (*models.ReturningAdvert, error) {
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
		a.views,
		a.advert_status,
		a.favourites_number,
		a.phone
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
	var advertStatus string

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &cityModel.CityName, &cityModel.Translation,
		&categoryModel.ID, &categoryModel.Name, &categoryModel.Translation, &advertModel.Title, &advertModel.Description, &advertModel.Price,
		&advertModel.CreatedTime, &advertModel.ClosedTime, &advertModel.IsUsed, &advertModel.Views, &advertStatus, &advertModel.FavouritesNum,
		&advertModel.Phone); err != nil {

		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert, err=%v", err))

		return nil, err
	}

	if advertStatus == "Активно" {
		advertModel.Active = true
	} else {
		advertModel.Active = false
	}
	advertModel.CityID = cityModel.ID
	advertModel.CategoryID = categoryModel.ID

	return &models.ReturningAdvert{
		Advert:   advertModel,
		City:     cityModel,
		Category: categoryModel,
	}, nil
}

func (ads *AdvertStorage) GetAdvertOnlyByID(ctx context.Context, advertID uint) (*models.ReturningAdvert, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList *models.ReturningAdvert

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.getAdvertOnlyByID(ctx, tx, advertID)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting advert by id only, err=%v", err))

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
		a.views,
		a.advert_status,
		a.favourites_number
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
	var advertStatus string

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &cityModel.CityName, &cityModel.Translation,
		&categoryModel.ID, &categoryModel.Name, &categoryModel.Translation, &advertModel.Title, &advertModel.Description, &advertModel.Price,
		&advertModel.CreatedTime, &advertModel.ClosedTime, &advertModel.IsUsed, &advertModel.Views, &advertStatus, &advertModel.FavouritesNum); err != nil {

		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert, err=%v", err))

		return nil, err
	}

	if advertStatus == "Активно" {
		advertModel.Active = true
	} else {
		advertModel.Active = false
	}
	advertModel.CityID = cityModel.ID
	advertModel.CategoryID = categoryModel.ID

	return &models.ReturningAdvert{
		Advert:   advertModel,
		City:     cityModel,
		Category: categoryModel,
	}, nil
}

func (ads *AdvertStorage) getAdvertAuth(ctx context.Context, tx pgx.Tx, userID, advertID uint, city, category string) (*models.ReturningAdvert, error) {
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
		a.views,
		a.advert_status,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $1 AND f.advert_id = a.id)
         	THEN 1 ELSE 0 END AS bool) AS in_favourites,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $1 AND c.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_cart,
		a.favourites_number	
		FROM 
		public.advert a
		LEFT JOIN 
		public.city c ON a.city_id = c.id
		LEFT JOIN 
		public.category cat ON a.category_id = cat.id
		WHERE a.id = $2;`

	logging.LogInfo(logger, "SELECT FROM advert, city, category")

	advertLine := tx.QueryRow(ctx, SQLAdvertById, userID, advertID)

	categoryModel := models.Category{}
	cityModel := models.City{}
	advertModel := models.Advert{}
	var advertStatus string

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &cityModel.CityName, &cityModel.Translation,
		&categoryModel.ID, &categoryModel.Name, &categoryModel.Translation, &advertModel.Title, &advertModel.Description, &advertModel.Price,
		&advertModel.CreatedTime, &advertModel.ClosedTime, &advertModel.IsUsed, &advertModel.Views, &advertStatus,
		&advertModel.InFavourites, &advertModel.InCart, &advertModel.FavouritesNum); err != nil {

		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert, err=%v", err))

		return nil, err
	}

	if advertStatus == "Активно" {
		advertModel.Active = true
	} else {
		advertModel.Active = false
	}
	advertModel.CityID = cityModel.ID
	advertModel.CategoryID = categoryModel.ID

	return &models.ReturningAdvert{
		Advert:   advertModel,
		City:     cityModel,
		Category: categoryModel,
	}, nil
}

func (ads *AdvertStorage) GetAdvert(ctx context.Context, userID, advertID uint, city, category string) (*models.ReturningAdvert, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList *models.ReturningAdvert
	var err error

	if userID == 0 {
		err = pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvert(ctx, tx, advertID, city, category)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting advert, err=%v", err))

			return nil, err
		}

	} else {
		err = pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertAuth(ctx, tx, userID, advertID, city, category)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting advert, err=%v", err))

			return nil, err
		}
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
	(SELECT array_agg(url_resized) FROM (SELECT url_resized FROM advert_image WHERE advert_id = a.id ORDER BY id) AS ordered_images) AS image_urls
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
		returningAdInList.InCart = false

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
    (SELECT array_agg(url_resized) FROM (SELECT url_resized FROM advert_image WHERE advert_id = a.id ORDER BY id) AS ordered_images) AS image_urls,
    CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $1 AND f.advert_id = a.id)
         THEN 1 ELSE 0 END AS bool) AS in_favourites,
	CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $1 AND c.advert_id = a.id)
	THEN 1 ELSE 0 END AS bool) AS in_cart	 
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.id >= $2 AND a.advert_status = 'Активно' AND c.translation = $3
	ORDER BY id
	LIMIT $4;
	`
	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image, favourite, cart")

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
			&returningAdInList.Price, &photoPad.Photo, &returningAdInList.InFavourites, &returningAdInList.InCart); err != nil {
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
	(SELECT array_agg(url_resized) FROM (SELECT url_resized FROM advert_image WHERE advert_id = a.id ORDER BY id) AS ordered_images) AS image_urls
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.id >= $1 AND a.advert_status = 'Активно' AND c.translation = $2 AND category.translation = $3
	ORDER BY id
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

		returningAdInList.InFavourites = false
		returningAdInList.InCart = false

		returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%v", err))

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) getAdvertsByCategoryAuth(ctx context.Context, tx pgx.Tx, category, city string, userID, startID, num uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAdvertsByCityAuth := `SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT array_agg(url_resized) FROM (SELECT url_resized FROM advert_image WHERE advert_id = a.id ORDER BY id) AS ordered_images) AS image_urls,
	CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $1 AND f.advert_id = a.id)
         THEN 1 ELSE 0 END AS bool) AS in_favourites,
	CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $1 AND c.advert_id = a.id)
	THEN 1 ELSE 0 END AS bool) AS in_cart
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.id >= $2 AND a.advert_status = 'Активно' AND c.translation = $3 AND category.translation = $4
	ORDER BY id
	LIMIT $5;
	`
	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image, favourite, cart")

	rows, err := tx.Query(ctx, SQLAdvertsByCityAuth, userID, startID, city, category, num)
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
			&returningAdInList.Price, &photoPad.Photo, &returningAdInList.InFavourites, &returningAdInList.InCart); err != nil {
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

func (ads *AdvertStorage) GetAdvertsByCategory(ctx context.Context, category, city string, userID, startID, num uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList []*models.ReturningAdInList

	if userID == 0 {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertsByCategory(ctx, tx, category, city, startID, num)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%v", err))

			return nil, err
		}

	} else {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertsByCategoryAuth(ctx, tx, category, city, userID, startID, num)
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

func (ads *AdvertStorage) getAdvertsForUserWhereStatusIs(ctx context.Context, tx pgx.Tx, userId, deleted uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	advertStatus := "Активно"
	if deleted == 1 {
		advertStatus = "Продано"
	}

	SQLGetAdvertsForUserWhereStatusIs := `SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT array_agg(url_resized) FROM (SELECT url_resized FROM advert_image WHERE advert_id = a.id ORDER BY id) AS ordered_images) AS image_urls
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.user_id = $1 AND a.advert_status = $2 
	ORDER BY id;
	`

	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image")

	rows, err := tx.Query(ctx, SQLGetAdvertsForUserWhereStatusIs, userId, advertStatus)
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
			&returningAdInList.Price, &photoPad.Photo); err != nil {
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

		returningAdInList.InFavourites = false
		returningAdInList.InCart = false

		returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%v", err))

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) getAdvertsForUserWhereStatusIsAuth(ctx context.Context, tx pgx.Tx, userID, userId, deleted uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	advertStatus := "Активно"
	if deleted == 1 {
		advertStatus = "Продано"
	}

	SQLGetAdvertsForUserWhereStatusIsAuth := `SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT array_agg(url_resized) FROM (SELECT url_resized FROM advert_image WHERE advert_id = a.id ORDER BY id) AS ordered_images) AS image_urls,
	CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $1 AND f.advert_id = a.id)
         THEN 1 ELSE 0 END AS bool) AS in_favourites,
	CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $1 AND c.advert_id = a.id)
		THEN 1 ELSE 0 END AS bool) AS in_cart
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.user_id = $2 AND a.advert_status = $3 
	ORDER BY id;
	`

	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image, favourite, cart")

	rows, err := tx.Query(ctx, SQLGetAdvertsForUserWhereStatusIsAuth, userID, userId, advertStatus)
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
			&returningAdInList.Price, &photoPad.Photo, &returningAdInList.InFavourites, &returningAdInList.InCart); err != nil {
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

func (ads *AdvertStorage) GetAdvertsForUserWhereStatusIs(ctx context.Context, userID, userId, deleted uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList []*models.ReturningAdInList

	if userID == 0 {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertsForUserWhereStatusIs(ctx, tx, userId, deleted)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%v", err))

			return nil, err
		}

	} else {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertsForUserWhereStatusIsAuth(ctx, tx, userID, userId, deleted)
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

func (ads *AdvertStorage) createAdvert(ctx context.Context, tx pgx.Tx, data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCreateAdvert :=
		`WITH ins AS (
		INSERT INTO advert (user_id, city_id, category_id, title, description, price, is_used, phone)
		SELECT
			$1,
			city.id,
			category.id,
			$2,
			$3,
			$4,
			$5,
			$8
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

	advertLine := tx.QueryRow(ctx, SQLCreateAdvert, data.UserID, data.Title, data.Description, data.Price, data.IsUsed, data.City, data.Category, data.Phone)

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
				is_used = $4,
				phone = $8
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

	advertLine := tx.QueryRow(ctx, SQLUpdateAdvert, data.Title, data.Description, data.Price, data.IsUsed, data.City, data.Category, data.ID, data.Phone)

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

func (ads *AdvertStorage) searchAdvertByTitle(ctx context.Context, tx pgx.Tx, title string, startID, num uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLSearchAdvertByTitle := `SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT array_agg(url_resized) FROM (SELECT url_resized FROM advert_image WHERE advert_id = a.id ORDER BY id) AS ordered_images) AS image_urls
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE (to_tsvector(a.title) @@ to_tsquery(replace($1 || ':*', ' ', ' | '))) AND a.advert_status = 'Активно' AND a.id >= $2
	ORDER BY ts_rank(to_tsvector(a.title), to_tsquery(replace($1 || ':*', ' ', ' | '))) DESC	 
	LIMIT $3;
	`
	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image using index")

	const (
		startAdvertID uint = 1
		advertLimit   uint = 30
	)

	rows, err := tx.Query(ctx, SQLSearchAdvertByTitle, title, startID, num)
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

		returningAdInList.InFavourites = false
		returningAdInList.InCart = false

		returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%v", err))

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) searchAdvertByTitleAuth(ctx context.Context, tx pgx.Tx, title string, userID, startID, num uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLSearchAdvertByTitle := `SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT array_agg(url_resized) FROM (SELECT url_resized FROM advert_image WHERE advert_id = a.id ORDER BY id) AS ordered_images) AS image_urls,
	CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $4 AND f.advert_id = a.id)
         THEN 1 ELSE 0 END AS bool) AS in_favourites,
	CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $4 AND c.advert_id = a.id)
	THEN 1 ELSE 0 END AS bool) AS in_cart
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE (to_tsvector(a.title) @@ to_tsquery(replace($1 || ':*', ' ', ' | '))) AND a.advert_status = 'Активно' AND a.id >= $2
	ORDER BY ts_rank(to_tsvector(a.title), to_tsquery(replace($1 || ':*', ' ', ' | '))) DESC	 
	LIMIT $3;
	`
	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image using index")

	const (
		startAdvertID uint = 1
		advertLimit   uint = 30
	)

	rows, err := tx.Query(ctx, SQLSearchAdvertByTitle, title, startID, num, userID)
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
			&returningAdInList.Price, &photoPad.Photo, &returningAdInList.InFavourites, &returningAdInList.InCart); err != nil {
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

func (ads *AdvertStorage) SearchAdvertByTitle(ctx context.Context, title string, userID, startID, num uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList []*models.ReturningAdInList

	if userID == 0 {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.searchAdvertByTitle(ctx, tx, title, startID, num)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%v", err))

			return nil, err
		}

	} else {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.searchAdvertByTitleAuth(ctx, tx, title, userID, startID, num)
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

func (ads *AdvertStorage) getSuggestions(ctx context.Context, tx pgx.Tx, title string, num uint) ([]string, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	words := strings.Fields(title)
	wordsCount := len(words)

	SQLSelectSuggestionsOneWord := `
	SELECT DISTINCT LOWER(regexp_replace(ts_headline(a.title, to_tsquery(replace($1 || ':*', ' ', ' | ')), 
								'MaxFragments=1,' || 'FragmentDelimiter=...,MaxWords=2,MinWords=1'), '<b>|</b>', '', 'g')) AS title
	FROM public.advert a
	WHERE (to_tsvector(a.title) @@ to_tsquery(replace($1 || ':*', ' ', ' | '))) AND a.advert_status = 'Активно'
	ORDER BY title
	LIMIT $2;
	`

	SQLSelectSuggestionsManyWords := `WITH one_word_titles AS (
		SELECT DISTINCT LOWER(regexp_replace(ts_headline(a.title, to_tsquery(replace($1 || ':*', ' ', ' | ')), 
									  'MaxFragments=1,' || 'FragmentDelimiter=...,MaxWords=2,MinWords=1'), '<b>|</b>', '', 'g')) AS title
		FROM public.advert a
		WHERE (to_tsvector(a.title) @@ to_tsquery(replace($1 || ':*', ' ', ' | '))) AND a.advert_status = 'Активно'
	),
	two_word_titles AS (
		SELECT DISTINCT LOWER(regexp_replace(ts_headline(a.title, to_tsquery(replace($1 || ':*', ' ', ' | ')), 
									  'MaxFragments=2,' || 'FragmentDelimiter=...,MaxWords=3,MinWords=2'), '<b>|</b>', '', 'g')) AS title
		FROM public.advert a
		WHERE (to_tsvector(a.title) @@ to_tsquery(replace($1 || ':*', ' ', ' | '))) AND a.advert_status = 'Активно'
	)
	SELECT * FROM one_word_titles
	UNION
	SELECT * FROM two_word_titles
	ORDER BY title
	LIMIT $2;
	`
	logging.LogInfo(logger, "SELECT FROM advert")
	var rows pgx.Rows
	var err error
	fmt.Println(wordsCount)
	if wordsCount > 1 {
		rows, err = tx.Query(ctx, SQLSelectSuggestionsManyWords, title, num)
	} else {
		rows, err = tx.Query(ctx, SQLSelectSuggestionsOneWord, title, num)
	}

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select advert title query, err=%v", err))

		return nil, err
	}
	defer rows.Close()

	var suggestions []string
	for rows.Next() {
		var suggestionPad *string

		if err := rows.Scan(&suggestionPad); err != nil {
			return nil, err
		}

		suggestionToInsert := ""
		if suggestionPad != nil {
			suggestionToInsert = *suggestionPad
		}

		// ПО-ХОРОШЕМУ СЮДА НУЖНО ПРИКРУТИТЬ САНИТАЙЗЕР
		suggestions = append(suggestions, suggestionToInsert)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert titles rows, err=%v", err))

		return nil, err
	}

	return suggestions, nil
}

func (ads *AdvertStorage) GetSuggestions(ctx context.Context, title string, num uint) ([]string, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var suggestions []string

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		suggestionsInner, err := ads.getSuggestions(ctx, tx, title, num)
		suggestions = suggestionsInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%v", err))

		return nil, err
	}

	return suggestions, nil
}
