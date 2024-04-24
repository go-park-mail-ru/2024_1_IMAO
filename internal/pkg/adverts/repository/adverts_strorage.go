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
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var advertsList *models.ReturningAdvert

	city := ""     // ЗАГЛУШКИ
	category := "" // ЗАГЛУШКИ

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.getAdvert(ctx, tx, advertID, city, category)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while getting adverts list, err=%v", err)

		return nil, err
	}

	return advertsList, nil
}

func (ads *AdvertStorage) getAdvertImagesURLs(ctx context.Context, tx pgx.Tx, advertID uint) ([]string, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLAdvertImagesURLs := `
	SELECT url
	FROM public.advert_image
	WHERE advert_id = $1;`

	childLogger.Infof(`
	SELECT url
	FROM public.advert_image
	WHERE advert_id = %s;`, advertID)

	rows, err := tx.Query(ctx, SQLAdvertImagesURLs, advertID)
	if err != nil {
		childLogger.Errorf("Something went wrong while executing select adverts urls, err=%v", err)

		return nil, err
	}
	defer rows.Close()

	var urlArray []string

	for rows.Next() {
		var returningUrl string

		if err := rows.Scan(&returningUrl); err != nil {
			childLogger.Errorf("Something went wrong while scanning rows of advert images for advert %v, err=%v", advertID, err)

			return nil, err
		}

		urlArray = append(urlArray, returningUrl)
	}

	return urlArray, nil
}

func (ads *AdvertStorage) GetAdvertImagesURLs(ctx context.Context, advertID uint) ([]string, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var urlArray []string

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		urlArrayInner, err := ads.getAdvertImagesURLs(ctx, tx, advertID)
		urlArray = urlArrayInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong getting image urls for advertID=%v , err=%v", advertID, err)

		return nil, err
	}

	return urlArray, nil
}

func (ads *AdvertStorage) getAdvert(ctx context.Context, tx pgx.Tx, advertID uint, city, category string) (*models.ReturningAdvert, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

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
		a.is_used
		FROM 
		public.advert a
		LEFT JOIN 
		public.city c ON a.city_id = c.id
		LEFT JOIN 
		public.category cat ON a.category_id = cat.id
		WHERE a.id = $1;`

	childLogger.Infof(`
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
		a.is_used
		FROM 
		public.advert a
		LEFT JOIN 
		public.city c ON a.city_id = c.id
		LEFT JOIN 
		public.category cat ON a.category_id = cat.id
		WHERE a.id = %s;`, advertID)
	advertLine := tx.QueryRow(ctx, SQLAdvertById, advertID)

	categoryModel := models.Category{}
	cityModel := models.City{}
	advertModel := models.Advert{}

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &cityModel.CityName, &cityModel.Translation,
		&categoryModel.ID, &categoryModel.Name, &categoryModel.Translation, &advertModel.Title, &advertModel.Description, &advertModel.Price,
		&advertModel.CreatedTime, &advertModel.ClosedTime, &advertModel.IsUsed); err != nil {

		childLogger.Errorf("Something went wrong while scanning advert, err=%v", err)

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
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var advertsList *models.ReturningAdvert

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.getAdvert(ctx, tx, advertID, city, category)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while getting adverts list, err=%v", err)

		return nil, err
	}

	advertsList.Photos, err = ads.GetAdvertImagesURLs(ctx, advertsList.Advert.ID)
	if err != nil {
		childLogger.Errorf("Something went wrong while getting advert images urls , err=%v", err)

		return nil, err
	}

	for i := 0; i < len(advertsList.Photos); i++ {

		image, err := utils.DecodeImage(advertsList.Photos[i])
		advertsList.PhotosIMG = append(advertsList.PhotosIMG, image)
		if err != nil {
			childLogger.Errorf("Error occurred while decoding advert_image %v, err = %v", advertsList.Photos[i], err)
		}
	}

	return advertsList, nil
}

func (ads *AdvertStorage) getAdvertsByCity(ctx context.Context, tx pgx.Tx, city string, startID, num uint) ([]*models.ReturningAdInList, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLAdvertsByCity := `SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT url FROM advert_image WHERE advert_id = a.id ORDER BY id LIMIT 1) AS first_image_url
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.id >= $1 AND a.advert_status = 'Активно' AND c.translation = $2
	LIMIT $3;
	`
	childLogger.Infof(`SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT url FROM advert_image WHERE advert_id = a.id ORDER BY id LIMIT 1) AS first_image_url
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.id >= %s AND a.advert_status = 'Активно' AND c.translation = %s
	LIMIT %s`, startID, city, num)
	rows, err := tx.Query(ctx, SQLAdvertsByCity, startID, city, num)
	if err != nil {
		childLogger.Errorf("Something went wrong while executing select adverts query, err=%v", err)

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

		photoURLToInsert := ""
		if photoPad.Photo != nil {
			photoURLToInsert = *photoPad.Photo
		}

		returningAdInList.PhotoIMG, err = utils.DecodeImage(photoURLToInsert)

		if err != nil {
			childLogger.Errorf("Something went wrong while decoding image, err=%v", err)

			return nil, err
		}

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		childLogger.Errorf("Something went wrong while scanning adverts rows, err=%v", err)

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) GetAdvertsByCity(ctx context.Context, city string, startID, num uint) ([]*models.ReturningAdInList, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var advertsList []*models.ReturningAdInList

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.getAdvertsByCity(ctx, tx, city, startID, num)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while getting adverts list, err=%v", err)

		return nil, err
	}

	return advertsList, nil
}

func (ads *AdvertStorage) getAdvertsByCategory(ctx context.Context, tx pgx.Tx, category, city string, startID, num uint) ([]*models.ReturningAdInList, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLAdvertsByCity := `SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT url FROM advert_image WHERE advert_id = a.id ORDER BY id LIMIT 1) AS first_image_url
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.id >= $1 AND a.advert_status = 'Активно' AND c.translation = $2 AND category.translation = $3
	LIMIT $4;
	`
	childLogger.Infof(`SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT url FROM advert_image WHERE advert_id = a.id ORDER BY id LIMIT 1) AS first_image_url
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.id >= %s AND a.advert_status = 'Активно' AND c.translation = %s AND category.translation = %s
	LIMIT %s`, startID, city, category, num)
	rows, err := tx.Query(ctx, SQLAdvertsByCity, startID, city, category, num)
	if err != nil {
		childLogger.Errorf("Something went wrong while executing select adverts query, err=%v", err)

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

		photoURLToInsert := ""
		if photoPad.Photo != nil {
			photoURLToInsert = *photoPad.Photo
		}

		returningAdInList.PhotoIMG, err = utils.DecodeImage(photoURLToInsert)

		if err != nil {
			childLogger.Errorf("Something went wrong while decoding image, err=%v", err)

			return nil, err
		}

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		childLogger.Errorf("Something went wrong while scanning adverts rows, err=%v", err)

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) GetAdvertsByCategory(ctx context.Context, category, city string, startID, num uint) ([]*models.ReturningAdInList, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var advertsList []*models.ReturningAdInList

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.getAdvertsByCategory(ctx, tx, category, city, startID, num)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while getting adverts list, err=%v", err)

		return nil, err
	}

	return advertsList, nil
}

func (ads *AdvertStorage) getAdvertsForUserWhereStatusIs(ctx context.Context, tx pgx.Tx, userId, deleted uint) (*models.ReturningAdvertList, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

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
		(SELECT url FROM advert_image WHERE advert_id = a.id ORDER BY id LIMIT 1) AS first_image_url
	FROM 
		public.advert a
	INNER JOIN 
		city c ON a.city_id = c.id
	INNER JOIN 
		category cat ON a.category_id = cat.id
	WHERE 
		a.user_id = $1 AND a.advert_status = $2;
	`
	childLogger.Infof(`SELECT a.id,	a.user_id, 	a.city_id, 	a.category_id, 	a.title, a.description, a.price, a.created_time, a.closed_time,	a.is_used, 	a.advert_status, c.name AS city_name,
		c.translation AS city_translation,	cat.name AS category_name,	cat.translation AS category_translation,
		(SELECT url FROM advert_image WHERE advert_id = a.id ORDER BY id LIMIT 1) AS first_image_url
		FROM public.advert a INNER JOIN city c ON a.city_id = c.id 
		INNER JOIN category cat ON a.category_id = cat.id WHERE 	a.user_id = %s AND a.advert_status = %s; `, userId, advertStatus)

	rows, err := tx.Query(ctx, SQLGetAdvertsForUserWhereStatusIs, userId, advertStatus)
	if err != nil {
		childLogger.Errorf("Something went wrong while executing select adverts for user where status is, err=%v", err)

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

			childLogger.Errorf("Something went wrong while scanning adverts rows, err=%v", err)

			return nil, err
		}

		photoURLToInsert := ""
		if photoPad.Photo != nil {
			photoURLToInsert = *photoPad.Photo
		}

		var photoArray []string

		photoArray = append(photoArray, photoURLToInsert)

		photoIMG, err := utils.DecodeImage(photoURLToInsert)

		if err != nil {
			childLogger.Errorf("Something went wrong while decoding image, err=%v", err)

			return nil, err
		}

		var photoIMGArray []string

		photoIMGArray = append(photoIMGArray, photoIMG)

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
		childLogger.Errorf("Something went wrong while scanning adverts rows, err=%v", err)

		return nil, err
	}

	return &returningAdvertList, nil
}

func (ads *AdvertStorage) GetAdvertsForUserWhereStatusIs(ctx context.Context, userId, deleted uint) ([]*models.ReturningAdInList, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var advertsList []*models.ReturningAdInList

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.getAdvertsForUserWhereStatusIs(ctx, tx, userId, deleted)
		for _, num := range advertsListInner.AdvertItems {
			returningAdInList := models.ReturningAdInList{
				ID:       num.Advert.ID,
				Title:    num.Advert.Title,
				Price:    num.Advert.Price,
				City:     num.City.Translation,
				Category: num.Category.Translation,
				Photo:    num.Photos[0],
				PhotoIMG: num.PhotosIMG[0],
			}
			advertsList = append(advertsList, &returningAdInList)
		}

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while getting adverts list, err=%v", err)

		return nil, err
	}

	return advertsList, nil
}

func (ads *AdvertStorage) createAdvert(ctx context.Context, tx pgx.Tx, data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

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

	childLogger.Infof(
		`WITH ins AS (
			INSERT INTO advert (user_id, city_id, category_id, title, description, price, is_used)
			SELECT
				%s,
				city.id,
				category.id,
				%s,
				%s,
				%s,
				%s
			FROM
				city
			JOIN
				category ON city.name = %s AND category.translation = %s
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
		LEFT JOIN public.category cat ON ins.category_id = cat.id;`,
		data.UserID, data.Title, data.Description, data.Price, data.IsUsed, data.City, data.Category)

	advertLine := tx.QueryRow(ctx, SQLCreateAdvert, data.UserID, data.Title, data.Description, data.Price, data.IsUsed, data.City, data.Category)

	categoryModel := models.Category{}
	cityModel := models.City{}
	advertModel := models.Advert{}

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &categoryModel.ID, &advertModel.Title, &advertModel.Description, &advertModel.CreatedTime,
		&advertModel.ClosedTime, &advertModel.Price, &advertModel.IsUsed, &cityModel.CityName, &cityModel.Translation, &categoryModel.Name, &categoryModel.Translation); err != nil {

		childLogger.Errorf("Something went wrong while scanning advert, err=%v", err)

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
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var advertsList *models.ReturningAdvert

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.createAdvert(ctx, tx, data)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while getting advert, err=%v", err)

		return nil, err
	}

	advertsList.Photos, err = ads.SetAdvertImages(ctx, files, "advert_images", advertsList.Advert.ID)
	if err != nil {
		childLogger.Errorf("Something went wrong while updating advert image url , err=%v", err)

		return nil, err
	}

	for i := 0; i < len(advertsList.Photos); i++ {
		image, err := utils.DecodeImage(advertsList.Photos[i])
		advertsList.PhotosIMG = append(advertsList.PhotosIMG, image)
		if err != nil {
			childLogger.Errorf("Error occurred while decoding advert_image %v, err = %v", advertsList.Photos[i], err)
		}
	}

	return advertsList, nil
}

func (ads *AdvertStorage) setAdvertImage(ctx context.Context, tx pgx.Tx, advertID uint, imageUrl string) (string, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLUpdateProfileAvatarURL := `
	INSERT INTO advert_image (url, advert_id)
	VALUES 
    ($1, $2)
	RETURNING url;`

	childLogger.Infof(`
	INSERT INTO advert_image (url, advert_id)
	VALUES 
    %1, %2)
	RETURNING url;`, imageUrl, advertID)

	var returningUrl string

	urlLine := tx.QueryRow(ctx, SQLUpdateProfileAvatarURL, imageUrl, advertID)

	if err := urlLine.Scan(&returningUrl); err != nil {

		childLogger.Errorf("Something went wrong while scanning advert image , err=%v", err)

		return "", err
	}

	return returningUrl, nil
}

func (ads *AdvertStorage) deleteAllImagesForAdvertFromLocalStorage(ctx context.Context, tx pgx.Tx, advertID uint) error {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLAdvertImagesURLs := `
	SELECT url
	FROM public.advert_image
	WHERE advert_id = $1;`

	childLogger.Infof(`
	SELECT url
	FROM public.advert_image
	WHERE advert_id = %s;`, advertID)

	rows, err := tx.Query(ctx, SQLAdvertImagesURLs, advertID)
	if err != nil {
		childLogger.Errorf("Something went wrong while executing select adverts urls, err=%v", err)

		return err
	}
	defer rows.Close()

	var oldUrl interface{}

	for rows.Next() {
		if err := rows.Scan(&oldUrl); err != nil {
			childLogger.Errorf("Something went wrong while scanning rows for deleting advert images for advert %v, err=%v", advertID, err)

			return err
		}

		if oldUrl != nil {
			err := os.Remove(oldUrl.(string))
			if err != nil {
				childLogger.Errorf("Something went wrong while deleting image %v, err=%v", oldUrl.(string), err)

				return err
			}
		}
	}

	return nil
}

func (ads *AdvertStorage) deleteAllImagesForAdvertFromDatabase(ctx context.Context, tx pgx.Tx, advertID uint) error {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLCreateProfile := `DELETE FROM public.advert_image WHERE advert_id = $1;`
	childLogger.Infof(`DELETE FROM public.advert_image WHERE advert_id = %s;`, advertID)

	var err error

	_, err = tx.Exec(ctx, SQLCreateProfile, advertID)

	if err != nil {
		childLogger.Errorf("Something went wrong while executing delete advert_image query, err=%v", err)

		return fmt.Errorf("Something went wrong while executing delete advert_image profile query", err)
	}

	return nil
}

func (ads *AdvertStorage) SetAdvertImages(ctx context.Context, files []*multipart.FileHeader, folderName string, advertID uint) ([]string, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		return ads.deleteAllImagesForAdvertFromLocalStorage(ctx, tx, advertID)
	})

	if err != nil {
		childLogger.Errorf("Something went wrong deleting AllImagesForAdvertFromLocalStorage , err=%v", err)

		return nil, err
	}

	err = pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		return ads.deleteAllImagesForAdvertFromDatabase(ctx, tx, advertID)
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while deleting AllImagesForAdvertFromDatabase , err=%v", err)

		return nil, err
	}

	var urlArray []string

	for i := 0; i < len(files); i++ {
		var url string

		fullPath, err := utils.WriteFile(files[i], folderName)

		if err != nil {
			childLogger.Errorf("Something went wrong while writing file of the image , err=%v", err)

			return nil, err
		}

		err = pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			urlInner, err := ads.setAdvertImage(ctx, tx, advertID, fullPath)
			url = urlInner

			return err
		})

		if err != nil {
			childLogger.Errorf("Something went wrong while updating profile url , err=%v", err)

			return nil, err
		}

		urlArray = append(urlArray, url)
	}

	return urlArray, nil
}

func (ads *AdvertStorage) editAdvert(ctx context.Context, tx pgx.Tx, data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

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

	childLogger.Infof(
		`WITH upd AS (
			UPDATE advert
			SET 
				city_id = city.id,
				category_id = category.id,
				title = %s,
				description = %s,
				price = %s,
				is_used = %s
			 FROM
				city
			JOIN
				category ON city.name = %s AND category.translation = %s
			WHERE 
				advert.id = %s
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
			public.category cat ON upd.category_id = cat.id;`,
		data.Title, data.Description, data.Price, data.IsUsed, data.City, data.Category, data.ID)

	advertLine := tx.QueryRow(ctx, SQLUpdateAdvert, data.Title, data.Description, data.Price, data.IsUsed, data.City, data.Category, data.ID)

	categoryModel := models.Category{}
	cityModel := models.City{}
	advertModel := models.Advert{}

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &categoryModel.ID, &advertModel.Title, &advertModel.Description, &advertModel.CreatedTime,
		&advertModel.ClosedTime, &advertModel.Price, &advertModel.IsUsed, &cityModel.CityName, &cityModel.Translation, &categoryModel.Name, &categoryModel.Translation); err != nil {

		childLogger.Errorf("Something went wrong while scanning advert, err=%v", err)

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
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	var advertsList *models.ReturningAdvert

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.editAdvert(ctx, tx, data)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		childLogger.Errorf("Something went wrong while updating advert, err=%v", err)

		return nil, err
	}

	advertsList.Photos, err = ads.SetAdvertImages(ctx, files, "advert_images", data.ID)
	if err != nil {
		childLogger.Errorf("Something went wrong while updating advert image url , err=%v", err)

		return nil, err
	}

	for i := 0; i < len(advertsList.Photos); i++ {

		image, err := utils.DecodeImage(advertsList.Photos[i])
		advertsList.PhotosIMG = append(advertsList.PhotosIMG, image)
		if err != nil {
			childLogger.Errorf("Error occurred while decoding advert_image %v, err = %v", advertsList.Photos[i], err)
		}
	}

	return advertsList, nil
}

func (ads *AdvertStorage) closeAdvert(ctx context.Context, tx pgx.Tx, advertID uint) error {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLCloseAdvert := `UPDATE public.advert	SET  advert_status='Скрыто'	WHERE id = $1;`
	childLogger.Infof(`UPDATE public.advert	SET  advert_status='Скрыто'	WHERE id = %s;`, advertID)
	var err error

	_, err = tx.Exec(ctx, SQLCloseAdvert, advertID)

	if err != nil {
		childLogger.Errorf("Something went wrong while executing close advert query, err=%v", err)
		return fmt.Errorf("Something went wrong while executing close advert query", err)
	}

	return nil
}

func (ads *AdvertStorage) CloseAdvert(ctx context.Context, advertID uint) error {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		err := ads.closeAdvert(ctx, tx, advertID)
		if err != nil {
			childLogger.Errorf("Something went wrong while closing advert, err=%v", err)
			return fmt.Errorf("Something went wrong while closing advert", err)
		}

		return nil
	})

	if err != nil {

		childLogger.Errorf("Error while closing advert, err=%v", err)
		return err
	}

	return nil
}

func (ads *AdvertStorage) deleteAdvert(ctx context.Context, tx pgx.Tx, user *models.User) error {
	requestUUID, ok := ctx.Value("requestUUID").(string)
	if !ok {
		requestUUID = "unknow"
	}

	childLogger := ads.logger.With(
		zap.String("requestUUID", requestUUID),
	)

	SQLCreateUser := `INSERT INTO public."user"(email, password_hash) VALUES ($1, $2);`
	childLogger.Infof(`INSERT INTO public."user"(email, password_hash) VALUES (%s, %s)`, user.Email, user.PasswordHash)
	var err error

	_, err = tx.Exec(ctx, SQLCreateUser, user.Email, user.PasswordHash)

	if err != nil {
		childLogger.Errorf("Something went wrong while executing create user query, err=%v", err)
		return fmt.Errorf("Something went wrong while executing create user query in func createUser", err)
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
