package repository

import (
	"context"
	"fmt"
	"time"

	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type FavouritesStorage struct {
	pool    *pgxpool.Pool
	metrics *mymetrics.DatabaseMetrics
}

func NewFavouritesStorage(pool *pgxpool.Pool, metrics *mymetrics.DatabaseMetrics) *FavouritesStorage {
	return &FavouritesStorage{
		pool:    pool,
		metrics: metrics,
	}
}

func (favouritesStorage *FavouritesStorage) getFavouritesByUserID(ctx context.Context, tx pgx.Tx,
	userID uint) ([]*models.ReturningAdInList, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLGetFavouritesByUserID := `SELECT a.id, c.translation, category.translation, a.title, a.price,
	(SELECT array_agg(url) FROM 
	                           (SELECT url 
	                            FROM advert_image 
	                            WHERE advert_id = a.id 
	                            ORDER BY id) AS ordered_images) AS image_urls,
	CAST(CASE WHEN EXISTS (SELECT 1 FROM cart c WHERE c.user_id = $1 AND c.advert_id = a.id)
		THEN 1 ELSE 0 END AS bool) AS in_cart
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	INNER JOIN favourite f ON a.id = f.advert_id
	WHERE f.user_id = $1
	ORDER BY a.id;
	`

	logging.LogInfo(logger, "SELECT FROM advert, city, category, advert_image, favourite, cart")

	start := time.Now()
	rows, err := tx.Query(ctx, SQLGetFavouritesByUserID, userID)
	favouritesStorage.metrics.AddDuration(funcName, time.Since(start))
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select adverts query, err=%v",
			err))
		favouritesStorage.metrics.IncreaseErrors(funcName)

		return nil, err
	}
	defer rows.Close()

	var adsList []*models.ReturningAdInList
	for rows.Next() {
		returningAdInList := models.ReturningAdInList{}

		photoPad := models.PhotoPad{}

		if err := rows.Scan(&returningAdInList.ID, &returningAdInList.City, &returningAdInList.Category,
			&returningAdInList.Title, &returningAdInList.Price, &photoPad.Photo,
			&returningAdInList.InCart); err != nil {
			return nil, err
		}

		returningAdInList.InFavourites = true
		returningAdInList.IsActive = true

		if photoPad.Photo != nil {
			for _, ptr := range photoPad.Photo {
				returningAdInList.Photos = append(returningAdInList.Photos, *ptr)
			}
		}

		for i := 0; i < len(returningAdInList.Photos); i++ {

			image, err := utils.DecodeImage(returningAdInList.Photos[i])
			returningAdInList.PhotosIMG = append(returningAdInList.PhotosIMG, image)
			if err != nil {
				logging.LogError(logger, fmt.Errorf("error occurred while decoding advert_image %v, err = %v",
					returningAdInList.Photos[i], err))

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

func (favouritesStorage *FavouritesStorage) GetFavouritesByUserID(ctx context.Context,
	userID uint) ([]*models.ReturningAdInList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var favourites []*models.ReturningAdInList

	err := pgx.BeginFunc(ctx, favouritesStorage.pool, func(tx pgx.Tx) error {
		favouritesInner, err := favouritesStorage.getFavouritesByUserID(ctx, tx, userID)
		favourites = favouritesInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, err=%v", err))

		return nil, err
	}

	if favourites == nil {
		favourites = []*models.ReturningAdInList{}
	}

	return favourites, nil
}

func (favouritesStorage *FavouritesStorage) deleteAdvByIDs(ctx context.Context, tx pgx.Tx, userID uint,
	advertID uint) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLDeleteFromCart := `DELETE FROM public.favourite
		WHERE user_id = $1 AND advert_id = $2;`

	logging.LogInfo(logger, "DELETE FROM favourite")

	var err error

	start := time.Now()
	_, err = tx.Exec(ctx, SQLDeleteFromCart, userID, advertID)
	favouritesStorage.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while executing advert delete from the favourite, err=%v", err))
		favouritesStorage.metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

func (favouritesStorage *FavouritesStorage) DeleteAdvByIDs(ctx context.Context, userID uint, advertID uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, favouritesStorage.pool, func(tx pgx.Tx) error {
		err := favouritesStorage.deleteAdvByIDs(ctx, tx, userID, advertID)

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting adverts list, most likely , err=%v", err))

		return err
	}

	return nil
}

func (favouritesStorage *FavouritesStorage) appendAdvByIDs(ctx context.Context, tx pgx.Tx, userID uint,
	advertID uint) (bool, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLAddToFavourites := `WITH deletion AS (
		DELETE FROM public.favourite
		WHERE user_id = $1 AND advert_id = $2
		RETURNING user_id, advert_id
	)
	INSERT INTO public.favourite (user_id, advert_id)
	SELECT $1, $2
	WHERE NOT EXISTS (
		SELECT 1 FROM deletion
	) RETURNING true;
	`
	logging.LogInfo(logger, "DELETE or SELECT FROM favourite")

	start := time.Now()
	userLine := tx.QueryRow(ctx, SQLAddToFavourites, userID, advertID)
	favouritesStorage.metrics.AddDuration(funcName, time.Since(start))

	added := false

	if err := userLine.Scan(&added); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning advert added, err=%v", err))

		return false, nil
	}

	return added, nil
}

func (favouritesStorage *FavouritesStorage) AppendAdvByIDs(ctx context.Context, userID uint, advertID uint) bool {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var added bool

	err := pgx.BeginFunc(ctx, favouritesStorage.pool, func(tx pgx.Tx) error {
		addedInner, err := favouritesStorage.appendAdvByIDs(ctx, tx, userID, advertID)
		added = addedInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing addvert add to favourites, err=%v", err))
	}

	return added
}
