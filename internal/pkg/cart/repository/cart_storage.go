package storage

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type CartStorage struct {
	pool    *pgxpool.Pool
	metrics *mymetrics.DatabaseMetrics
}

func NewCartStorage(pool *pgxpool.Pool, metrics *mymetrics.DatabaseMetrics) *CartStorage {
	return &CartStorage{
		pool:    pool,
		metrics: metrics,
	}
}

func (cl *CartStorage) getCartByUserID(ctx context.Context, tx pgx.Tx, userID uint) ([]*models.ReturningAdvert, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAdvertByUserId := `
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
			(SELECT url FROM advert_image WHERE advert_id = a.id ORDER BY id LIMIT 1) AS first_image_url,
			CAST(CASE WHEN EXISTS (SELECT 1 FROM favourite f WHERE f.user_id = $1 AND f.advert_id = a.id)
         		THEN 1 ELSE 0 END AS bool) AS in_favourites
		FROM 
			public.advert a
		LEFT JOIN 
			public.city c ON a.city_id = c.id
		LEFT JOIN 
			public.category cat ON a.category_id = cat.id
		LEFT JOIN 
			public.cart cart ON a.id = cart.advert_id
		WHERE cart.user_id = $1;`

	logging.LogInfo(logger, "SELECT FROM advert, cart, category, city, advert_image")

	start := time.Now()
	rows, err := tx.Query(ctx, SQLAdvertByUserId, userID)
	cl.metrics.AddDuration(funcName, time.Since(start))
	if err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while executing select adverts from the cart, err=%v", err))
		cl.metrics.IncreaseErrors(funcName)

		return nil, err
	}
	defer rows.Close()

	var adsList []*models.ReturningAdvert
	for rows.Next() {

		categoryModel := models.Category{}
		cityModel := models.City{}
		advertModel := models.Advert{}
		photoPad := models.PhotoPadSoloImage{}

		if err := rows.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &cityModel.CityName,
			&cityModel.Translation, &categoryModel.ID, &categoryModel.Name, &categoryModel.Translation,
			&advertModel.Title, &advertModel.Description, &advertModel.Price, &advertModel.CreatedTime,
			&advertModel.ClosedTime, &advertModel.IsUsed, &photoPad.Photo, &advertModel.InFavourites); err != nil {

			logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts from the cart, err=%v", err))

			return nil, err
		}

		advertModel.CityID = cityModel.ID
		advertModel.CategoryID = categoryModel.ID

		photoURLToInsert := ""
		if photoPad.Photo != nil {
			photoURLToInsert = *photoPad.Photo
		}

		returningAdvertList := models.ReturningAdvert{
			Advert:   advertModel,
			City:     cityModel,
			Category: categoryModel,
		}

		decodedImage, err := utils.DecodeImage(photoURLToInsert)

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while decoding image, err=%v", err))

			return nil, err
		}

		returningAdvertList.Photos = append(returningAdvertList.Photos, photoURLToInsert)
		returningAdvertList.PhotosIMG = append(returningAdvertList.PhotosIMG, decodedImage)

		adsList = append(adsList, &returningAdvertList)
	}

	return adsList, nil
}

func (cl *CartStorage) GetCartByUserID(ctx context.Context, userID uint) ([]*models.ReturningAdvert, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	cart := []*models.ReturningAdvert{}

	err := pgx.BeginFunc(ctx, cl.pool, func(tx pgx.Tx) error {
		cartInner, err := cl.getCartByUserID(ctx, tx, userID)
		cart = cartInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%v", err))

		return nil, err
	}

	if cart == nil {
		cart = []*models.ReturningAdvert{}
	}

	return cart, nil
}

func (cl *CartStorage) deleteAdvByIDs(ctx context.Context, tx pgx.Tx, userID uint, advertID uint) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLDeleteFromCart := `DELETE FROM public.cart
		WHERE user_id = $1 AND advert_id = $2;`

	logging.LogInfo(logger, "DELETE FROM cart")

	var err error

	start := time.Now()
	_, err = tx.Exec(ctx, SQLDeleteFromCart, userID, advertID)
	cl.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while executing advert delete from the cart, err=%v", err))
		cl.metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

func (cl *CartStorage) DeleteAdvByIDs(ctx context.Context, userID uint, advertID uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, cl.pool, func(tx pgx.Tx) error {
		err := cl.deleteAdvByIDs(ctx, tx, userID, advertID)

		return err
	})

	if err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while getting adverts list, most likely , err=%v", err))

		return err
	}

	return nil
}

func (cl *CartStorage) appendAdvByIDs(ctx context.Context, tx pgx.Tx, userID uint, advertID uint) (bool, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAddToCart := `WITH deletion AS (
		DELETE FROM public.cart
		WHERE user_id = $1 AND advert_id = $2
		RETURNING user_id, advert_id
	)
	INSERT INTO public.cart (user_id, advert_id)
	SELECT $1, $2
	WHERE NOT EXISTS (
		SELECT 1 FROM deletion
	) RETURNING true;
	`
	logging.LogInfo(logger, "DELETE or SELECT FROM cart")

	start := time.Now()
	userLine := tx.QueryRow(ctx, SQLAddToCart, userID, advertID)
	cl.metrics.AddDuration(funcName, time.Since(start))

	added := false

	if err := userLine.Scan(&added); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning advert added, err=%v", err))
		cl.metrics.IncreaseErrors(funcName)

		return false, nil
	}

	return added, nil
}

func (cl *CartStorage) AppendAdvByIDs(ctx context.Context, userID uint, advertID uint) bool {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var added bool

	err := pgx.BeginFunc(ctx, cl.pool, func(tx pgx.Tx) error {
		addedInner, err := cl.appendAdvByIDs(ctx, tx, userID, advertID)
		added = addedInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing addvert add to cart, err=%v", err))
	}

	return added
}
