//nolint:forcetypeassert,cyclop,prealloc
package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
)

const (
	waitingMinutes = 10
	activeStatus   = "Активно"
)

type AdvertStorage struct {
	pool    *pgxpool.Pool
	metrics *mymetrics.DatabaseMetrics
}

func NewAdvertStorage(pool *pgxpool.Pool, metrics *mymetrics.DatabaseMetrics) *AdvertStorage {
	return &AdvertStorage{
		pool:    pool,
		metrics: metrics,
	}
}

func (ads *AdvertStorage) getAdvertOnlyByID(ctx context.Context, tx pgx.Tx,
	advertID uint) (*models.ReturningAdvert, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAdvertByID := `
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

	start := time.Now()

	advertLine := tx.QueryRow(ctx, SQLAdvertByID, advertID)

	ads.metrics.AddDuration(funcName, time.Since(start))

	var (
		categoryModel models.Category
		cityModel     models.City
		advertModel   models.Advert
		advertStatus  string
	)

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &cityModel.CityName,
		&cityModel.Translation, &categoryModel.ID, &categoryModel.Name, &categoryModel.Translation,
		&advertModel.Title, &advertModel.Description, &advertModel.Price, &advertModel.CreatedTime,
		&advertModel.ClosedTime, &advertModel.IsUsed, &advertModel.Views, &advertStatus, &advertModel.FavouritesNum,
		&advertModel.Phone); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert, err=%w", err))

		return nil, err
	}

	if advertStatus == activeStatus {
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
		logging.LogError(logger, fmt.Errorf("something went wrong while getting advert by id only, err=%w", err))

		return nil, err
	}

	advertsList.Photos, err = ads.GetAdvertImagesURLs(ctx, advertsList.Advert.ID)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting advert images urls , err=%w",
			err))

		return nil, err
	}

	for i := 0; i < len(advertsList.Photos); i++ {
		image, err := utils.DecodeImage(advertsList.Photos[i])
		if err != nil {
			logging.LogError(logger, fmt.Errorf("error occurred while decoding advert_image %s, err = %w",
				advertsList.Photos[i], err))
		}

		advertsList.PhotosIMG = append(advertsList.PhotosIMG, image)
	}

	// advertsList.Advert.Sanitize()

	return advertsList, nil
}

func (ads *AdvertStorage) getAdvertImagesURLs(ctx context.Context, tx pgx.Tx, advertID uint) ([]string, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAdvertImagesURLs := `
	SELECT url
	FROM public.advert_image
	WHERE advert_id = $1;`

	logging.LogInfo(logger, "SELECT FROM advert_image")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLAdvertImagesURLs, advertID)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts urls, err=%w",
			err))
		ads.metrics.IncreaseErrors(funcName)

		return nil, err
	}

	defer rows.Close()

	var urlArray []string

	for rows.Next() {
		var returningURL string

		if err := rows.Scan(&returningURL); err != nil {
			logging.LogError(logger,
				fmt.Errorf("something went wrong while scanning rows of advert images for advert %v, err=%w",
					advertID, err))

			return nil, err
		}

		urlArray = append(urlArray, returningURL)
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
		logging.LogError(logger, fmt.Errorf("something went wrong getting image urls for advertID=%v , err=%w",
			advertID, err))

		return nil, err
	}

	return urlArray, nil
}

func (ads *AdvertStorage) insertView(ctx context.Context, tx pgx.Tx, userID, advertID uint) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLInsertView := `
	INSERT INTO public.view(
		user_id, advert_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, advert_id) DO NOTHING;`

	logging.LogInfo(logger, "INSERT INTO view")

	var err error

	start := time.Now()

	_, err = tx.Exec(ctx, SQLInsertView, userID, advertID)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while executing insertInsertView query, err=%w", err))
		ads.metrics.IncreaseErrors(funcName)

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
		logging.LogError(logger, fmt.Errorf("something went wrong while inserting view, err=%w", err))

		return err
	}

	return nil
}

func (ads *AdvertStorage) getAdvert(ctx context.Context, tx pgx.Tx, advertID uint) (*models.ReturningAdvert, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAdvertByID := `
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

	start := time.Now()

	advertLine := tx.QueryRow(ctx, SQLAdvertByID, advertID)

	ads.metrics.AddDuration(funcName, time.Since(start))

	var (
		categoryModel models.Category
		cityModel     models.City
		advertModel   models.Advert
		advertStatus  string
	)

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &cityModel.CityName,
		&cityModel.Translation, &categoryModel.ID, &categoryModel.Name, &categoryModel.Translation, &advertModel.Title,
		&advertModel.Description, &advertModel.Price, &advertModel.CreatedTime, &advertModel.ClosedTime,
		&advertModel.IsUsed, &advertModel.Views, &advertStatus, &advertModel.FavouritesNum); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert, err=%w", err))

		return nil, err
	}

	if advertStatus == activeStatus {
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

func (ads *AdvertStorage) getAdvertAuth(ctx context.Context, tx pgx.Tx, userID,
	advertID uint) (*models.ReturningAdvert, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAdvertByID := `
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
		a.is_promoted,
		a.promotion_start,
		a.promotion_duration,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $1 AND f.advert_id = a.id)
         	THEN 1 ELSE 0 END AS bool) AS in_favourites,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $1 AND c.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_cart,
		a.favourites_number,
		ARRAY_AGG(pay.created_time ORDER BY pay.created_time DESC) FILTER (WHERE pay.payment_status = 'pending')
		FROM 
		public.advert a
		LEFT JOIN 
		public.city c ON a.city_id = c.id
		LEFT JOIN 
		public.category cat ON a.category_id = cat.id
		LEFT JOIN 
		public.payments pay ON pay.advert_id = a.id
		WHERE a.id = $2
		GROUP BY 
		     a.id, a.user_id, a.city_id, c.name, c.translation, a.category_id, cat.name, cat.translation, a.title, 
		     a.description, a.price, a.created_time, a.closed_time, a.is_used, a.views, a.advert_status, a.is_promoted, 
		     a.promotion_start, a.promotion_duration, a.favourites_number;`

	logging.LogInfo(logger, "SELECT FROM advert, city, category")

	start := time.Now()

	advertLine := tx.QueryRow(ctx, SQLAdvertByID, userID, advertID)

	ads.metrics.AddDuration(funcName, time.Since(start))

	categoryModel := models.Category{}
	cityModel := models.City{}
	advertModel := models.Advert{}
	promotionModel := models.Promotion{}
	paymentsDates := models.PaymentsDatesList{}

	var advertStatus string

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &cityModel.CityName,
		&cityModel.Translation, &categoryModel.ID, &categoryModel.Name, &categoryModel.Translation, &advertModel.Title,
		&advertModel.Description, &advertModel.Price, &advertModel.CreatedTime, &advertModel.ClosedTime,
		&advertModel.IsUsed, &advertModel.Views, &advertStatus, &promotionModel.IsPromoted,
		&promotionModel.PromotionStart, &promotionModel.PromotionDuration, &advertModel.InFavourites,
		&advertModel.InCart, &advertModel.FavouritesNum, &paymentsDates.List); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert, err=%w", err))

		return nil, err
	}

	for _, date := range paymentsDates.List {
		if time.Since(*date).Minutes() < waitingMinutes {
			promotionModel.NeedPing = true

			break
		}
	}

	if advertStatus == activeStatus {
		advertModel.Active = true
	} else {
		advertModel.Active = false
	}

	advertModel.CityID = cityModel.ID
	advertModel.CategoryID = categoryModel.ID

	return &models.ReturningAdvert{
		Advert:    advertModel,
		City:      cityModel,
		Category:  categoryModel,
		Promotion: promotionModel,
	}, nil
}

func (ads *AdvertStorage) GetAdvert(ctx context.Context, userID, advertID uint) (*models.ReturningAdvert, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var (
		advertsList *models.ReturningAdvert
		err         error
	)

	if userID == 0 {
		err = pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvert(ctx, tx, advertID)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting advert, err=%w", err))

			return nil, err
		}
	} else {
		err = pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertAuth(ctx, tx, userID, advertID)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting advert, err=%w", err))

			return nil, err
		}
	}

	advertsList.Photos, err = ads.GetAdvertImagesURLs(ctx, advertsList.Advert.ID)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting advert images urls , err=%w",
			err))

		return nil, err
	}

	for i := 0; i < len(advertsList.Photos); i++ {
		image, err := utils.DecodeImage(advertsList.Photos[i])
		if err != nil {
			logging.LogError(logger, fmt.Errorf("error occurred while decoding advert_image %v, err = %w",
				advertsList.Photos[i], err))
		}

		advertsList.PhotosIMG = append(advertsList.PhotosIMG, image)
	}

	// advertsList.Advert.Sanitize()

	return advertsList, nil
}

func (ads *AdvertStorage) getAdvertsByCity(ctx context.Context, tx pgx.Tx, city string,
	startID, num uint) ([]*models.ReturningAdInList, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var param uint

	if (startID-1)%num != 0 {
		return nil, nil
	}

	param = (startID - 1) / num

	SQLAdvertsByCityPromoted := `
	WITH promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted, a.promotion_start,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = TRUE AND a.advert_status = 'Активно' AND c.translation = $2
		ORDER BY promotion_start DESC, id
		OFFSET 5 * $1
		LIMIT 5
	), non_promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted, a.promotion_start,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = FALSE AND a.advert_status = 'Активно' AND c.translation = $2
		ORDER BY id
		OFFSET 15 * $1 + 5 * div((SELECT 
		CASE 
			WHEN ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) < 0 THEN 0
			ELSE ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE)
		END), 5) + (SELECT
			CASE
				WHEN (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) - $1 * 5 < 0 
						THEN 5 - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) % 5
				ELSE 0
			END)
	)
	SELECT * FROM promoted_adverts
	UNION ALL
	SELECT * FROM non_promoted_adverts
	ORDER BY is_promoted DESC, promotion_start DESC, id ASC
	LIMIT 20;
	`

	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLAdvertsByCityPromoted, param, city)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts query, err=%w",
			err))
		ads.metrics.IncreaseErrors(funcName)

		return nil, err
	}

	defer rows.Close()

	var adsList []*models.ReturningAdInList

	for rows.Next() {
		returningAdInList := models.ReturningAdInList{}

		var (
			photoPad       models.PhotoPad
			isPromotedPlug interface{}
		)

		if err := rows.Scan(&returningAdInList.ID, &returningAdInList.City, &returningAdInList.Category,
			&returningAdInList.Title, &returningAdInList.Price, &returningAdInList.IsPromoted, &isPromotedPlug,
			&photoPad.Photo); err != nil {
			ads.metrics.IncreaseErrors(funcName)

			return nil, err
		}

		returningAdInList.InFavourites = false
		returningAdInList.InCart = false
		returningAdInList.IsActive = true

		if photoPad.Photo != nil {
			for _, ptr := range photoPad.Photo {
				returningAdInList.Photos = append(returningAdInList.Photos, *ptr)
			}
		}
		// returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%w", err))
		ads.metrics.IncreaseErrors(funcName)

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) getAdvertsByCityAuth(ctx context.Context, tx pgx.Tx, city string, userID, startID,
	num uint) ([]*models.ReturningAdInList, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var param uint

	if (startID-1)%num != 0 {
		return nil, nil
	}

	param = (startID - 1) / num

	SQLAdvertsByCityAuth := `
	WITH promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted, a.promotion_start,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $3 AND f.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_favourites,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $3 AND c.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_cart								
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = TRUE AND a.advert_status = 'Активно' AND c.translation = $2
		ORDER BY promotion_start DESC, id
		OFFSET 5 * $1
		LIMIT 5
	), non_promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted, a.promotion_start,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $3 AND f.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_favourites,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $3 AND c.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_cart								
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = FALSE AND a.advert_status = 'Активно' AND c.translation = $2
		ORDER BY id
		OFFSET 15 * $1 + 5 * div((SELECT 
		CASE 
			WHEN ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) < 0 THEN 0
			ELSE ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE)
		END), 5) + (SELECT
			CASE
				WHEN (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) - $1 * 5 < 0 
						THEN 5 - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) % 5
				ELSE 0
			END)
	)
	SELECT * FROM promoted_adverts
	UNION ALL
	SELECT * FROM non_promoted_adverts
	ORDER BY is_promoted DESC, promotion_start DESC, id ASC
	LIMIT 20;
	`

	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image, favourite, cart")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLAdvertsByCityAuth, param, city, userID)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts query, err=%w",
			err))
		ads.metrics.IncreaseErrors(funcName)

		return nil, err
	}

	defer rows.Close()

	var adsList []*models.ReturningAdInList

	for rows.Next() {
		returningAdInList := models.ReturningAdInList{}

		var (
			photoPad       models.PhotoPad
			isPromotedPlug interface{}
		)

		if err := rows.Scan(&returningAdInList.ID, &returningAdInList.City, &returningAdInList.Category,
			&returningAdInList.Title, &returningAdInList.Price, &returningAdInList.IsPromoted, &isPromotedPlug,
			&photoPad.Photo, &returningAdInList.InFavourites, &returningAdInList.InCart); err != nil {
			return nil, err
		}

		returningAdInList.IsActive = true

		if photoPad.Photo != nil {
			for _, ptr := range photoPad.Photo {
				returningAdInList.Photos = append(returningAdInList.Photos, *ptr)
			}
		}

		// returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%w", err))

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) GetAdvertsByCity(ctx context.Context, city string, userID, startID,
	num uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList []*models.ReturningAdInList

	if userID == 0 {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertsByCity(ctx, tx, city, startID, num)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%w", err))

			return nil, err
		}
	} else {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertsByCityAuth(ctx, tx, city, userID, startID, num)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%w", err))

			return nil, err
		}
	}

	return advertsList, nil
}

func (ads *AdvertStorage) getAdvertsByCategory(ctx context.Context, tx pgx.Tx, category, city string,
	startID, num uint) ([]*models.ReturningAdInList, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var param uint

	if (startID-1)%num != 0 {
		return nil, nil
	}

	param = (startID - 1) / num

	SQLAdvertsByCity := `
	WITH promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted, a.promotion_start,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = TRUE AND a.advert_status = 'Активно' AND c.translation = $2 AND category.translation = $3
		ORDER BY promotion_start DESC, id
		OFFSET 5 * $1
		LIMIT 5
	), non_promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted, a.promotion_start,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = FALSE AND a.advert_status = 'Активно' AND c.translation = $2 AND category.translation = $3
		ORDER BY id
		OFFSET 15 * $1 + 5 * div((SELECT 
		CASE 
			WHEN ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) < 0 THEN 0
			ELSE ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE)
		END), 5) + (SELECT
			CASE
				WHEN (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) - $1 * 5 < 0 
						THEN 5 - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) % 5
				ELSE 0
			END)
	)
	SELECT * FROM promoted_adverts
	UNION ALL
	SELECT * FROM non_promoted_adverts
	ORDER BY is_promoted DESC, promotion_start DESC, id ASC
	LIMIT 20;
	`

	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLAdvertsByCity, param, city, category)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts query, err=%w",
			err))
		ads.metrics.IncreaseErrors(funcName)

		return nil, err
	}

	defer rows.Close()

	var adsList []*models.ReturningAdInList

	for rows.Next() {
		returningAdInList := models.ReturningAdInList{}

		var (
			photoPad       models.PhotoPad
			isPromotedPlug interface{}
		)

		if err := rows.Scan(&returningAdInList.ID, &returningAdInList.City, &returningAdInList.Category,
			&returningAdInList.Title, &returningAdInList.Price, &returningAdInList.IsPromoted, &isPromotedPlug,
			&photoPad.Photo); err != nil {
			return nil, err
		}

		returningAdInList.IsActive = true

		if photoPad.Photo != nil {
			for _, ptr := range photoPad.Photo {
				returningAdInList.Photos = append(returningAdInList.Photos, *ptr)
			}
		}

		returningAdInList.InFavourites = false
		returningAdInList.InCart = false

		// returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%w", err))

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) getAdvertsByCategoryAuth(ctx context.Context, tx pgx.Tx, category, city string, userID,
	startID, num uint) ([]*models.ReturningAdInList, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var param uint

	if (startID-1)%num != 0 {
		return nil, nil
	}

	param = (startID - 1) / num

	SQLAdvertsByCityAuth := `
	WITH promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted, a.promotion_start,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $4 AND f.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_favourites,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $4 AND c.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_cart								
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = TRUE AND a.advert_status = 'Активно' AND c.translation = $2 AND category.translation = $3
		ORDER BY promotion_start DESC, id
		OFFSET 5 * $1
		LIMIT 5
	), non_promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted, a.promotion_start,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $4 AND f.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_favourites,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $4 AND c.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_cart								
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = FALSE AND a.advert_status = 'Активно' AND c.translation = $2 AND category.translation = $3
		ORDER BY id
		OFFSET 15 * $1 + 5 * div((SELECT 
		CASE 
			WHEN ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) < 0 THEN 0
			ELSE ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE)
		END), 5) + (SELECT
			CASE
				WHEN (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) - $1 * 5 < 0 
						THEN 5 - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) % 5
				ELSE 0
			END)
	)
	SELECT * FROM promoted_adverts
	UNION ALL
	SELECT * FROM non_promoted_adverts
	ORDER BY is_promoted DESC, promotion_start DESC, id ASC
	LIMIT 20;
	`

	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image, favourite, cart")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLAdvertsByCityAuth, param, city, category, userID)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts query, err=%w",
			err))
		ads.metrics.IncreaseErrors(funcName)

		return nil, err
	}

	defer rows.Close()

	var adsList []*models.ReturningAdInList

	for rows.Next() {
		returningAdInList := models.ReturningAdInList{}

		var (
			photoPad       models.PhotoPad
			isPromotedPlug interface{}
		)

		if err := rows.Scan(&returningAdInList.ID, &returningAdInList.City, &returningAdInList.Category,
			&returningAdInList.Title, &returningAdInList.Price, &returningAdInList.IsPromoted, &isPromotedPlug,
			&photoPad.Photo, &returningAdInList.InFavourites, &returningAdInList.InCart); err != nil {
			return nil, err
		}

		returningAdInList.IsActive = true

		if photoPad.Photo != nil {
			for _, ptr := range photoPad.Photo {
				returningAdInList.Photos = append(returningAdInList.Photos, *ptr)
			}
		}

		// returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%w", err))

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) GetAdvertsByCategory(ctx context.Context, category, city string, userID, startID,
	num uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList []*models.ReturningAdInList

	if userID == 0 {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertsByCategory(ctx, tx, category, city, startID, num)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%w", err))

			return nil, err
		}
	} else {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertsByCategoryAuth(ctx, tx, category, city, userID, startID, num)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%w", err))

			return nil, err
		}
	}

	return advertsList, nil
}

func (ads *AdvertStorage) getAdvertsForUserWhereStatusIs(ctx context.Context, tx pgx.Tx, userID, deleted,
	advertNum uint) ([]*models.ReturningAdInList, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	advertStatus := activeStatus
	if deleted == 1 {
		advertStatus = "Продано"
	}

	SQLGetAdvertsForUserWhereStatusIs := `SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.user_id = $1 AND a.advert_status = $2 
	ORDER BY id
	LIMIT $3;
	`

	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLGetAdvertsForUserWhereStatusIs, userID, advertStatus, advertNum)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts query, err=%w",
			err))
		ads.metrics.IncreaseErrors(funcName)

		return nil, err
	}

	defer rows.Close()

	var adsList []*models.ReturningAdInList

	for rows.Next() {
		returningAdInList := models.ReturningAdInList{}

		photoPad := models.PhotoPad{}

		if err := rows.Scan(&returningAdInList.ID, &returningAdInList.City, &returningAdInList.Category,
			&returningAdInList.Title, &returningAdInList.Price, &photoPad.Photo); err != nil {
			return nil, err
		}

		returningAdInList.IsActive = true
		if deleted == 1 {
			returningAdInList.IsActive = false
		}

		if photoPad.Photo != nil {
			for _, ptr := range photoPad.Photo {
				returningAdInList.Photos = append(returningAdInList.Photos, *ptr)
			}
		}

		returningAdInList.InFavourites = false
		returningAdInList.InCart = false

		// returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%w", err))

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) getAdvertsForUserWhereStatusIsAuth(ctx context.Context, tx pgx.Tx, userID, authorID,
	deleted, advertNum uint) ([]*models.ReturningAdInList, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	advertStatus := activeStatus
	if deleted == 1 {
		advertStatus = "Продано"
	}

	SQLGetAdvertsForUserWhereStatusIsAuth := `SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls,
	CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $1 AND f.advert_id = a.id)
         THEN 1 ELSE 0 END AS bool) AS in_favourites,
	CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $1 AND c.advert_id = a.id)
		THEN 1 ELSE 0 END AS bool) AS in_cart
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.user_id = $2 AND a.advert_status = $3 
	ORDER BY id
	LIMIT $4;
	`

	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image, favourite, cart")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLGetAdvertsForUserWhereStatusIsAuth, userID, authorID, advertStatus, advertNum)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts query, err=%w",
			err))
		ads.metrics.IncreaseErrors(funcName)

		return nil, err
	}

	defer rows.Close()

	var adsList []*models.ReturningAdInList

	for rows.Next() {
		returningAdInList := models.ReturningAdInList{}

		photoPad := models.PhotoPad{}

		if err := rows.Scan(&returningAdInList.ID, &returningAdInList.City, &returningAdInList.Category,
			&returningAdInList.Title, &returningAdInList.Price, &photoPad.Photo, &returningAdInList.InFavourites,
			&returningAdInList.InCart); err != nil {
			return nil, err
		}

		returningAdInList.IsActive = true
		if deleted == 1 {
			returningAdInList.IsActive = false
		}

		if photoPad.Photo != nil {
			for _, ptr := range photoPad.Photo {
				returningAdInList.Photos = append(returningAdInList.Photos, *ptr)
			}
		}
		// returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%w", err))

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) GetAdvertsForUserWhereStatusIs(ctx context.Context, userID, authorID,
	deleted, advertNum uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList []*models.ReturningAdInList

	if userID == 0 {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertsForUserWhereStatusIs(ctx, tx, authorID, deleted, advertNum)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%w", err))

			return nil, err
		}
	} else {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.getAdvertsForUserWhereStatusIsAuth(ctx, tx, userID, authorID,
				deleted, advertNum)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%w", err))

			return nil, err
		}
	}

	return advertsList, nil
}

func (ads *AdvertStorage) createAdvert(ctx context.Context, tx pgx.Tx,
	data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCreateAdvert :=
		`WITH ins AS (
		INSERT INTO advert (user_id, city_id, category_id, title, description, price, is_used, phone, price_history)
		SELECT
			$1,
			city.id,
			category.id,
			$2,
			$3,
			$4,
			$5,
			$8,
			ARRAY['{"updated_time":"' || $9 || '", "new_price":' || $10 || '}']::jsonb[]
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
	SELECT ins.*, c.name AS city_name, c.translation AS city_translation, cat.name AS category_name, 
			cat.translation AS category_translation
	FROM ins
	LEFT JOIN public.city c ON ins.city_id = c.id
	LEFT JOIN public.category cat ON ins.category_id = cat.id;`

	logging.LogInfo(logger, "INSERT INTO advert")

	start := time.Now()

	advertLine := tx.QueryRow(ctx, SQLCreateAdvert, data.UserID, data.Title, data.Description, data.Price, data.IsUsed,
		data.City, data.Category, data.Phone, time.Now().Format("2006-01-02 15:04:05"),
		strconv.Itoa(int(data.Price)))

	ads.metrics.AddDuration(funcName, time.Since(start))

	categoryModel := models.Category{}
	cityModel := models.City{}
	advertModel := models.Advert{}

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &categoryModel.ID,
		&advertModel.Title, &advertModel.Description, &advertModel.CreatedTime, &advertModel.ClosedTime,
		&advertModel.Price, &advertModel.IsUsed, &cityModel.CityName, &cityModel.Translation, &categoryModel.Name,
		&categoryModel.Translation); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert, err=%w", err))
		ads.metrics.IncreaseErrors(funcName)

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
		logging.LogError(logger, fmt.Errorf("something went wrong while getting advert, err=%w", err))

		return nil, err
	}

	advertsList.Photos, err = ads.SetAdvertImages(ctx, files, "advert_images",
		"advert_images_resized", advertsList.Advert.ID)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while updating advert image url , err=%w",
			err))

		return nil, err
	}

	for i := 0; i < len(advertsList.Photos); i++ {
		image, err := utils.DecodeImage(advertsList.Photos[i])
		if err != nil {
			logging.LogError(logger, fmt.Errorf("error occurred while decoding advert_image %v, err = %w",
				advertsList.Photos[i], err))

			return nil, err
		}

		advertsList.PhotosIMG = append(advertsList.PhotosIMG, image)
	}

	// advertsList.Advert.Sanitize()

	return advertsList, nil
}

func (ads *AdvertStorage) setAdvertImage(ctx context.Context, tx pgx.Tx, advertID uint, originalImageURL,
	resizedImage string) (string, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLUpdateProfileAvatarURL := `
	INSERT INTO advert_image (url, advert_id, url_resized)
	VALUES 
    ($1, $2, $3)
	RETURNING url;`

	logging.LogInfo(logger, "INSERT INTO advert_image")

	var returningURL string

	start := time.Now()

	urlLine := tx.QueryRow(ctx, SQLUpdateProfileAvatarURL, originalImageURL, advertID, resizedImage)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err := urlLine.Scan(&returningURL); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert image , err=%w", err))

		return "", err
	}

	return returningURL, nil
}

func (ads *AdvertStorage) deleteAllImagesForAdvertFromLocalStorage(ctx context.Context, tx pgx.Tx,
	advertID uint) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAdvertImagesURLs := `
	SELECT url
	FROM public.advert_image
	WHERE advert_id = $1;`

	logging.LogInfo(logger, "SELECT FROM advert_image")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLAdvertImagesURLs, advertID)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts urls, err=%w",
			err))
		ads.metrics.IncreaseErrors(funcName)

		return err
	}

	defer rows.Close()

	var oldURL interface{}

	for rows.Next() {
		if err := rows.Scan(&oldURL); err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while scanning rows "+
				"for deleting advert images for advert %d, err=%w", advertID, err))

			return err
		}

		if oldURL != nil {
			err := os.Remove(oldURL.(string))
			if err != nil {
				logging.LogError(logger, fmt.Errorf("something went wrong while deleting image %s, err=%w",
					oldURL.(string), err))

				return err
			}
		}
	}

	return nil
}

func (ads *AdvertStorage) deleteAllImagesForAdvertFromDatabase(ctx context.Context, tx pgx.Tx, advertID uint) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCreateProfile := `DELETE FROM public.advert_image WHERE advert_id = $1;`

	logging.LogInfo(logger, "DELETE FROM advert_image")

	var err error

	start := time.Now()

	_, err = tx.Exec(ctx, SQLCreateProfile, advertID)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing delete advert_image query, "+
			"err=%w", err))
		ads.metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

func (ads *AdvertStorage) SetAdvertImages(ctx context.Context, files []*multipart.FileHeader, originalImageFolderName,
	resizedImageFolderName string, advertID uint) ([]string, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		return ads.deleteAllImagesForAdvertFromLocalStorage(ctx, tx, advertID)
	})

	if err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong deleting AllImagesForAdvertFromLocalStorage , err=%w", err))

		return nil, err
	}

	err = pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		return ads.deleteAllImagesForAdvertFromDatabase(ctx, tx, advertID)
	})

	if err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while deleting AllImagesForAdvertFromDatabase , err=%w", err))

		return nil, err
	}

	var urlArray []string

	for i := 0; i < len(files); i++ {
		var url string

		originalImageFullPath, err := utils.WriteFile(files[i], originalImageFolderName)

		if err != nil {
			logging.LogError(logger,
				fmt.Errorf("something went wrong while writing file of the original image , err=%w", err))

			return nil, err
		}

		resizedImageFullPath, err := utils.WriteResizedFile(files[i], resizedImageFolderName)

		if err != nil {
			logging.LogError(logger,
				fmt.Errorf("something went wrong while writing file of the resized image , err=%w", err))

			return nil, err
		}

		err = pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			urlInner, err := ads.setAdvertImage(ctx, tx, advertID, originalImageFullPath, resizedImageFullPath)
			url = urlInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while updating profile url , err=%w", err))

			return nil, err
		}

		urlArray = append(urlArray, url)
	}

	return urlArray, nil
}

func (ads *AdvertStorage) editAdvert(ctx context.Context, tx pgx.Tx,
	data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	funcName := logging.GetOnlyFunctionName()
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
				phone = $8,
				price_history = price_history || 
					ARRAY['{"updated_time":"' || $9 || '", "new_price":' || $10 || '}']::jsonb[]
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

	start := time.Now()

	advertLine := tx.QueryRow(ctx, SQLUpdateAdvert, data.Title, data.Description, data.Price, data.IsUsed,
		data.City, data.Category, data.ID, data.Phone, time.Now().Format("2006-01-02 15:04:05"),
		strconv.Itoa(int(data.Price)))

	ads.metrics.AddDuration(funcName, time.Since(start))

	categoryModel := models.Category{}
	cityModel := models.City{}
	advertModel := models.Advert{}

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &categoryModel.ID,
		&advertModel.Title, &advertModel.Description, &advertModel.CreatedTime, &advertModel.ClosedTime,
		&advertModel.Price, &advertModel.IsUsed, &cityModel.CityName, &cityModel.Translation,
		&categoryModel.Name, &categoryModel.Translation); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert, err=%w", err))

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

func (ads *AdvertStorage) EditAdvert(ctx context.Context, files []*multipart.FileHeader,
	data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList *models.ReturningAdvert

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.editAdvert(ctx, tx, data)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while updating advert, err=%w", err))

		return nil, err
	}

	advertsList.Photos, err = ads.SetAdvertImages(ctx, files, "advert_images",
		"advert_images_resized", data.ID)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while updating advert image url , err=%w",
			err))

		return nil, err
	}

	for i := 0; i < len(advertsList.Photos); i++ {
		image, err := utils.DecodeImage(advertsList.Photos[i])
		if err != nil {
			logging.LogError(logger, fmt.Errorf("error occurred while decoding advert_image %s, err = %w",
				advertsList.Photos[i], err))

			return nil, err
		}

		advertsList.PhotosIMG = append(advertsList.PhotosIMG, image)
	}

	// advertsList.Advert.Sanitize()

	return advertsList, nil
}

func (ads *AdvertStorage) closeAdvert(ctx context.Context, tx pgx.Tx, advertID uint) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCloseAdvert := `UPDATE public.advert	SET  advert_status='Скрыто'	WHERE id = $1;`

	logging.LogInfo(logger, "UPDATE advert")

	var err error

	start := time.Now()

	_, err = tx.Exec(ctx, SQLCloseAdvert, advertID)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing close advert query, err=%w",
			err))
		ads.metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

func (ads *AdvertStorage) CloseAdvert(ctx context.Context, advertID uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		err := ads.closeAdvert(ctx, tx, advertID)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while closing advert, err=%w", err))

			return err
		}

		return nil
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while closing advert, err=%w", err))

		return err
	}

	return nil
}

func (ads *AdvertStorage) searchAdvertByTitle(ctx context.Context, tx pgx.Tx, title, city string, startID,
	num uint) ([]*models.ReturningAdInList, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var param uint

	if (startID-1)%num != 0 {
		return nil, nil
	}

	param = (startID - 1) / num

	SQLSearchAdvertByTitle := `
	WITH promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted, a.promotion_start,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id 
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = TRUE AND a.advert_status = 'Активно' AND c.translation = $3
			AND (to_tsvector(a.title) @@ to_tsquery(replace($2 || ':*', ' ', ' | ')))
		ORDER BY ts_rank(to_tsvector(a.title), to_tsquery(replace($2 || ':*', ' ', ' | '))) DESC, promotion_start DESC, id
		OFFSET 5 * $1
		LIMIT 5
	), non_promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted, a.promotion_start,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id 
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = FALSE AND a.advert_status = 'Активно' AND c.translation = $3
			AND (to_tsvector(a.title) @@ to_tsquery(replace($2 || ':*', ' ', ' | ')))
		ORDER BY ts_rank(to_tsvector(a.title), to_tsquery(replace($2 || ':*', ' ', ' | '))) DESC, id
		OFFSET 15 * $1 + 5 * div((SELECT 
		CASE 
			WHEN ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) < 0 THEN 0
			ELSE ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE)
		END), 5) + (SELECT
			CASE
				WHEN (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) - $1 * 5 < 0 
						THEN 5 - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) % 5
				ELSE 0
			END)
	)
	SELECT * FROM promoted_adverts
	UNION ALL
	SELECT * FROM non_promoted_adverts
	ORDER BY is_promoted DESC, promotion_start DESC, id ASC
	LIMIT 20;
	`

	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image using index")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLSearchAdvertByTitle, param, title, city)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts query, err=%w",
			err))
		ads.metrics.IncreaseErrors(funcName)

		return nil, err
	}

	defer rows.Close()

	var adsList []*models.ReturningAdInList

	for rows.Next() {
		returningAdInList := models.ReturningAdInList{}

		var (
			photoPad       models.PhotoPad
			isPromotedPlug interface{}
		)

		if err := rows.Scan(&returningAdInList.ID, &returningAdInList.City, &returningAdInList.Category,
			&returningAdInList.Title, &returningAdInList.Price, &returningAdInList.IsPromoted, &isPromotedPlug,
			&photoPad.Photo); err != nil {
			return nil, err
		}

		returningAdInList.IsActive = true

		if photoPad.Photo != nil {
			for _, ptr := range photoPad.Photo {
				returningAdInList.Photos = append(returningAdInList.Photos, *ptr)
			}
		}

		returningAdInList.InFavourites = false
		returningAdInList.InCart = false

		// returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%w", err))

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) searchAdvertByTitleAuth(ctx context.Context, tx pgx.Tx, title, city string, userID, startID,
	num uint) ([]*models.ReturningAdInList, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var param uint

	if (startID-1)%num != 0 {
		return nil, nil
	}

	SQLSearchAdvertByTitle := `
	WITH promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted, a.promotion_start,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $3 AND f.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_favourites,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $3 AND c.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_cart								
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id 
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = TRUE AND a.advert_status = 'Активно' AND c.translation = $4
			AND (to_tsvector(a.title) @@ to_tsquery(replace($2 || ':*', ' ', ' | ')))
		ORDER BY ts_rank(to_tsvector(a.title), to_tsquery(replace($2 || ':*', ' ', ' | '))) DESC, 
				promotion_start DESC, id
		OFFSET 5 * $1
		LIMIT 5
	), non_promoted_adverts AS (
		SELECT a.id, c.translation, category.translation, a.title, a.price, a.is_promoted, a.promotion_start,
		(SELECT array_agg(url_resized) FROM 
	                                   (SELECT url_resized 
	                                    FROM advert_image 
	                                    WHERE advert_id = a.id 
	                                    ORDER BY id) AS ordered_images) AS image_urls,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $3 AND f.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_favourites,
		CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $3 AND c.advert_id = a.id)
			THEN 1 ELSE 0 END AS bool) AS in_cart
		FROM public.advert a
		INNER JOIN city c ON a.city_id = c.id 
		INNER JOIN category ON a.category_id = category.id
		WHERE is_promoted = FALSE AND a.advert_status = 'Активно' AND c.translation = $4
			AND (to_tsvector(a.title) @@ to_tsquery(replace($2 || ':*', ' ', ' | ')))
		ORDER BY ts_rank(to_tsvector(a.title), to_tsquery(replace($2 || ':*', ' ', ' | '))) DESC, id
		OFFSET 15 * $1 + 5 * div((SELECT 
		CASE 
			WHEN ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) < 0 THEN 0
			ELSE ($1 * 5) - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE)
		END), 5) + (SELECT
			CASE
				WHEN (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) - $1 * 5 < 0 
						THEN 5 - (SELECT COUNT(*) FROM public.advert WHERE is_promoted = TRUE) % 5
				ELSE 0
			END)
	)
	SELECT * FROM promoted_adverts
	UNION ALL
	SELECT * FROM non_promoted_adverts
	ORDER BY is_promoted DESC, promotion_start DESC, id ASC
	LIMIT 20;
	`

	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image using index")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLSearchAdvertByTitle, param, title, userID, city)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts query, err=%w",
			err))
		ads.metrics.IncreaseErrors(funcName)

		return nil, err
	}

	defer rows.Close()

	var adsList []*models.ReturningAdInList

	for rows.Next() {
		var (
			returningAdInList models.ReturningAdInList
			photoPad          models.PhotoPad
			isPromotedPlug    interface{}
		)

		if err := rows.Scan(&returningAdInList.ID, &returningAdInList.City, &returningAdInList.Category,
			&returningAdInList.Title, &returningAdInList.Price, &returningAdInList.IsPromoted, &isPromotedPlug,
			&photoPad.Photo, &returningAdInList.InFavourites, &returningAdInList.InCart); err != nil {
			return nil, err
		}

		returningAdInList.IsActive = true

		if photoPad.Photo != nil {
			for _, ptr := range photoPad.Photo {
				returningAdInList.Photos = append(returningAdInList.Photos, *ptr)
			}
		}

		// returningAdInList.Sanitize()

		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts rows, err=%w", err))

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertStorage) SearchAdvertByTitle(ctx context.Context, title, city string, userID, startID,
	num uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var advertsList []*models.ReturningAdInList

	if userID == 0 {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.searchAdvertByTitle(ctx, tx, title, city, startID, num)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%w", err))

			return nil, err
		}
	} else {
		err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
			advertsListInner, err := ads.searchAdvertByTitleAuth(ctx, tx, title, city, userID, startID, num)
			advertsList = advertsListInner

			return err
		})

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%w", err))

			return nil, err
		}
	}

	return advertsList, nil
}

func (ads *AdvertStorage) getSuggestions(ctx context.Context, tx pgx.Tx, title, city string, num uint) ([]string, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	words := strings.Fields(title)
	wordsCount := len(words)

	SQLSelectSuggestionsOneWord := `
	SELECT DISTINCT LOWER(regexp_replace(ts_headline(a.title, to_tsquery(replace($1 || ':*', ' ', ' | ')), 
								'MaxFragments=1,' || 'FragmentDelimiter=...,MaxWords=2,MinWords=1'), 
								'<b>|</b>', '', 'g')) AS title
	FROM public.advert a
	JOIN public.city c on a.city_id = c.id 
	WHERE (to_tsvector(a.title) @@ to_tsquery(replace($1 || ':*', ' ', ' | '))) AND a.advert_status = 'Активно' AND c.translation = $3
	ORDER BY title
	LIMIT $2;
	`

	SQLSelectSuggestionsManyWords := `WITH one_word_titles AS (
		SELECT DISTINCT LOWER(regexp_replace(ts_headline(a.title, to_tsquery(replace($1 || ':*', ' ', ' | ')), 
									  'MaxFragments=1,' || 'FragmentDelimiter=...,MaxWords=2,MinWords=1'), 
										'<b>|</b>', '', 'g')) AS title
		FROM public.advert a
		JOIN public.city c on a.city_id = c.id 
		WHERE (to_tsvector(a.title) @@ to_tsquery(replace($1 || ':*', ' ', ' | '))) AND a.advert_status = 'Активно' AND c.translation = $3
	),
	two_word_titles AS (
		SELECT DISTINCT LOWER(regexp_replace(ts_headline(a.title, to_tsquery(replace($1 || ':*', ' ', ' | ')), 
									  'MaxFragments=2,' || 'FragmentDelimiter=...,MaxWords=3,MinWords=2'), 
										'<b>|</b>', '', 'g')) AS title
		FROM public.advert a
		JOIN public.city c on a.city_id = c.id 
		WHERE (to_tsvector(a.title) @@ to_tsquery(replace($1 || ':*', ' ', ' | '))) AND a.advert_status = 'Активно' AND c.translation = $3
	)
	SELECT * FROM one_word_titles
	UNION
	SELECT * FROM two_word_titles
	ORDER BY title
	LIMIT $2;
	`

	logging.LogInfo(logger, "SELECT FROM advert")

	var (
		rows pgx.Rows
		err  error
	)

	start := time.Now()

	if wordsCount > 1 {
		rows, err = tx.Query(ctx, SQLSelectSuggestionsManyWords, title, num, city)
	} else {
		rows, err = tx.Query(ctx, SQLSelectSuggestionsOneWord, title, num, city)
	}

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while executing select advert title query, err=%w", err))

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
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning advert titles rows, err=%w",
			err))

		return nil, err
	}

	return suggestions, nil
}

func (ads *AdvertStorage) GetSuggestions(ctx context.Context, title, city string, num uint) ([]string, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var suggestions []string

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		suggestionsInner, err := ads.getSuggestions(ctx, tx, title, city, num)
		suggestions = suggestionsInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%w", err))

		return nil, err
	}

	return suggestions, nil
}

func (ads *AdvertStorage) getPriceHistory(ctx context.Context, tx pgx.Tx, id uint) ([]*models.PriceHistoryItem, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLSelectPriceHistory := `SELECT price_history
	FROM public.advert
	WHERE id = $1; `

	logging.LogInfo(logger, "SELECT FROM advert")

	start := time.Now()

	priceHistoryLine := tx.QueryRow(ctx, SQLSelectPriceHistory, id)

	ads.metrics.AddDuration(funcName, time.Since(start))

	var (
		priceHistory []*models.PriceHistoryItem
		myJsonbArray []interface{}
	)

	if err := priceHistoryLine.Scan(&myJsonbArray); err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while scanning priceHistory into []interface{} , err=%w", err))

		return nil, err
	}

	for _, item := range myJsonbArray {
		var priceHistoryItem models.PriceHistoryItem

		if abobaMap, ok := item.(map[string]interface{}); ok {
			priceHistoryItem.NewPrice = abobaMap["new_price"].(float64)
			priceHistoryItem.UpdatedTime = abobaMap["updated_time"].(string)
		} else {
			err := fmt.Errorf("item %v is not of type map[string]interface{}", item)
			logging.LogError(logger, err)

			return nil, err
		}

		priceHistory = append(priceHistory, &priceHistoryItem)
	}

	return priceHistory, nil
}

func (ads *AdvertStorage) GetPriceHistory(ctx context.Context, userID uint) ([]*models.PriceHistoryItem, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var priceHistoryItem []*models.PriceHistoryItem

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		priceHistoryItemInner, err := ads.getPriceHistory(ctx, tx, userID)
		priceHistoryItem = priceHistoryItemInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting priceHistory, err=%w", err))

		return nil, err
	}

	return priceHistoryItem, nil
}

func (ads *AdvertStorage) checkAdvertOwnership(ctx context.Context, tx pgx.Tx, advertID, userID uint) (bool, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCheckAdvertOwnership :=
		`SELECT EXISTS(SELECT 1 FROM public.advert WHERE id=$1 AND user_id=$2 AND is_promoted=false);`

	logging.LogInfo(logger, "SELECT FROM advert")

	start := time.Now()

	userLine := tx.QueryRow(ctx, SQLCheckAdvertOwnership, advertID, userID)

	ads.metrics.AddDuration(funcName, time.Since(start))

	var ownership bool

	if err := userLine.Scan(&ownership); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning advert ownership exists, err=%w", err))

		return false, err
	}

	return ownership, nil
}

func (ads *AdvertStorage) CheckAdvertOwnership(ctx context.Context, advertID, userID uint) bool {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var ownership bool

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		ownershipExists, err := ads.checkAdvertOwnership(ctx, tx, advertID, userID)
		ownership = ownershipExists

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing ownership exists query, err=%w", err))
	}

	return ownership
}

func (ads *AdvertStorage) getPaymnetUUIDList(ctx context.Context, tx pgx.Tx,
	advertID uint) (*models.PaymnetUUIDList, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLSelectPaymnetUUIDList := `SELECT 
	(SELECT array_agg(payment_uuid) 
	 FROM (SELECT payment_uuid FROM payments 
		   WHERE advert_id = $1 AND payment_status='pending' ORDER BY id)
	 AS ordered_uuids)
	 AS uuid_list;`

	logging.LogInfo(logger, "SELECT FROM payments")

	start := time.Now()

	userLine := tx.QueryRow(ctx, SQLSelectPaymnetUUIDList, advertID)

	ads.metrics.AddDuration(funcName, time.Since(start))

	PaymnetUUIDList := models.PaymnetUUIDList{}
	UUIDListPad := models.PaymnetUUIDListPad{}

	if err := userLine.Scan(&UUIDListPad.Pad); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning payments uuids, err=%w", err))

		return nil, err
	}

	if UUIDListPad.Pad != nil {
		for _, ptr := range UUIDListPad.Pad {
			PaymnetUUIDList.UUIDList = append(PaymnetUUIDList.UUIDList, *ptr)
		}
	}

	return &PaymnetUUIDList, nil
}

func (ads *AdvertStorage) GetPaymnetUUIDList(ctx context.Context, advertID uint) (*models.PaymnetUUIDList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var PaymnetUUIDList *models.PaymnetUUIDList

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		PaymnetUUIDListInner, err := ads.getPaymnetUUIDList(ctx, tx, advertID)
		PaymnetUUIDList = PaymnetUUIDListInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing ownership exists query, err=%w", err))

		return nil, err
	}

	return PaymnetUUIDList, nil
}

// ПЕРЕПИСАТЬ ЧЕРЕЗ ПЕРЕСЕЧЕНИЕ МНОЖЕСТВ И BULK UPDATE
func (ads *AdvertStorage) yuKassaUpdateOneRecord(ctx context.Context, tx pgx.Tx, uuid string) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCloseAdvert := `UPDATE public.payments
	SET  payment_status='waiting_for_capture'
	WHERE payment_uuid=$1;`

	logging.LogInfo(logger, "UPDATE advert")

	var err error

	start := time.Now()

	_, err = tx.Exec(ctx, SQLCloseAdvert, uuid)

	ads.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while executing updating payment query, err=%w", err))
		ads.metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

// ПЕРЕПИСАТЬ ЧЕРЕЗ ПЕРЕСЕЧЕНИЕ МНОЖЕСТВ И BULK UPDATE
func (ads *AdvertStorage) YuKassaUpdateOneRecord(ctx context.Context, uuid string) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		err := ads.yuKassaUpdateOneRecord(ctx, tx, uuid)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while updating payment status, err=%w", err))

			return err
		}

		return nil
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while closing advert, err=%w", err))

		return err
	}

	return nil
}

// ПЕРЕПИСАТЬ ЧЕРЕЗ ПЕРЕСЕЧЕНИЕ МНОЖЕСТВ И BULK UPDATE
func (ads *AdvertStorage) YuKassaUpdateDB(ctx context.Context, paymentList *models.PaymentList, advertID uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	dbUUIDList, err := ads.GetPaymnetUUIDList(ctx, advertID)

	if err != nil {
		logging.LogError(logger, fmt.Errorf("no payments for advert with id=%d or error occurred, err=%w",
			advertID, err))

		return err
	}

	yukassaUUIDArray := []string{}

	for i := 0; i < len(paymentList.Items); i++ {
		yukassaUUIDArray = append(yukassaUUIDArray, paymentList.Items[i].ID)
	}

	resultArray := utils.FindIntersection(dbUUIDList.UUIDList, yukassaUUIDArray)

	for i := 0; i < len(resultArray); i++ {
		err = ads.YuKassaUpdateOneRecord(ctx, resultArray[i])
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while updating payment status, err=%w",
				err))

			return err
		}
	}

	return nil
}

func (ads *AdvertStorage) getPromotionData(ctx context.Context, tx pgx.Tx, advertID uint) (*models.Promotion, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	getPromotionDataQuery := `
		SELECT 
		a.is_promoted,
		a.promotion_start,
		a.promotion_duration
		FROM 
		public.advert a
		WHERE a.id = $1;`

	logging.LogInfo(logger, "SELECT promotion FROM advert")

	start := time.Now()

	line := tx.QueryRow(ctx, getPromotionDataQuery, advertID)

	ads.metrics.AddDuration(funcName, time.Since(start))

	var promotionData models.Promotion

	if err := line.Scan(&promotionData.IsPromoted, &promotionData.PromotionStart,
		&promotionData.PromotionDuration); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning promotion data, err=%w", err))

		return nil, err
	}

	return &promotionData, nil
}

func (ads *AdvertStorage) GetPromotionData(ctx context.Context, advertID uint) (*models.Promotion, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var promotionData *models.Promotion

	err := pgx.BeginFunc(ctx, ads.pool, func(tx pgx.Tx) error {
		promotionDataInner, err := ads.getPromotionData(ctx, tx, advertID)
		promotionData = promotionDataInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while promotion data advert, err=%w", err))

		return nil, err
	}

	return promotionData, nil
}
