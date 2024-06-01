package storage

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

type OrderStorage struct {
	pool    *pgxpool.Pool
	metrics *mymetrics.DatabaseMetrics
}

func NewOrderStorage(pool *pgxpool.Pool, metrics *mymetrics.DatabaseMetrics) *OrderStorage {
	return &OrderStorage{
		pool:    pool,
		metrics: metrics,
	}
}

func (ol *OrderStorage) getBoughtOrdersByUserID(ctx context.Context, tx pgx.Tx,
	userID uint) ([]*models.ReturningOrder, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLGetBoughtOrdersByUserID := `
		SELECT 
		ord.id AS order_id,
		ord.order_status, 
		ord.created_time AS order_created_time, 
		ord.updated_time AS order_updated_time,
		ord.closed_time AS order_closed_time, 
		ord.phone AS order_phone, 
		ord.name AS order_name, 
		ord.email AS order_email, 
		ord.delivery_price AS order_delivery_price,
		ord.delivery_address AS order_delivery_address,
		a.id AS advert_id, 
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
		(SELECT COALESCE(r.rating, 0) AS rating FROM public.advert a LEFT JOIN public.review r ON r.advert_id = a.id WHERE a.id = ord.advert_id) 
	FROM 
		public.advert a
	LEFT JOIN 
		public.city c ON a.city_id = c.id
	LEFT JOIN 
		public.category cat ON a.category_id = cat.id
	LEFT JOIN 
		public.order ord ON a.id = ord.advert_id
	WHERE ord.user_id = $1;`

	logging.LogInfo(logger, "SELECT FROM advert, cart, category, city, advert_image")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLGetBoughtOrdersByUserID, userID)

	ol.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while executing select adverts from the cart, err=%w", err))
		ol.metrics.IncreaseErrors(funcName)

		return nil, err
	}

	defer rows.Close()

	var orderList []*models.ReturningOrder

	for rows.Next() {
		categoryModel := models.Category{}
		cityModel := models.City{}
		advertModel := models.Advert{}
		photoPad := models.PhotoPadSoloImage{}
		orderItem := models.OrderItem{}

		if err := rows.Scan(&orderItem.ID, &orderItem.Status, &orderItem.Created, &orderItem.Updated, &orderItem.Closed,
			&orderItem.Phone, &orderItem.Name, &orderItem.Email, &orderItem.DeliveryPrice, &orderItem.Address,
			&advertModel.ID, &advertModel.UserID, &cityModel.ID, &cityModel.CityName, &cityModel.Translation,
			&categoryModel.ID, &categoryModel.Name, &categoryModel.Translation, &advertModel.Title,
			&advertModel.Description, &advertModel.Price, &advertModel.CreatedTime, &advertModel.ClosedTime,
			&advertModel.IsUsed, &photoPad.Photo, &orderItem.Rating); err != nil {
			logging.LogError(logger,
				fmt.Errorf("something went wrong while scanning adverts from the cart, err=%w", err))

			return nil, err
		}

		advertModel.CityID = cityModel.ID
		advertModel.CategoryID = categoryModel.ID

		photoURLToInsert := ""
		if photoPad.Photo != nil {
			photoURLToInsert = *photoPad.Photo
		}

		returningAdvert := models.ReturningAdvert{
			Advert:   advertModel,
			City:     cityModel,
			Category: categoryModel,
		}

		decodedImage, err := utils.DecodeImage(photoURLToInsert)

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while decoding image, err=%w", err))

			return nil, err
		}

		returningAdvert.Photos = append(returningAdvert.Photos, photoURLToInsert)
		returningAdvert.PhotosIMG = append(returningAdvert.PhotosIMG, decodedImage)

		ReturningOrder := models.ReturningOrder{
			OrderItem:       orderItem,
			ReturningAdvert: returningAdvert,
		}

		orderList = append(orderList, &ReturningOrder)

	}

	return orderList, nil
}

func (ol *OrderStorage) GetBoughtOrdersByUserID(ctx context.Context, userID uint) ([]*models.ReturningOrder, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	orderList := []*models.ReturningOrder{}

	err := pgx.BeginFunc(ctx, ol.pool, func(tx pgx.Tx) error {
		orderListInner, err := ol.getBoughtOrdersByUserID(ctx, tx, userID)
		orderList = orderListInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting orders list, err=%w", err))

		return nil, err
	}

	if orderList == nil {
		orderList = []*models.ReturningOrder{}
	}

	return orderList, nil
}

func (ol *OrderStorage) getSoldOrdersByUserID(ctx context.Context, tx pgx.Tx,
	userID uint) ([]*models.ReturningOrder, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLGetSoldOrdersByUserID := `
		SELECT 
		ord.id AS order_id,
		ord.order_status, 
		ord.created_time AS order_created_time, 
		ord.updated_time AS order_updated_time,
		ord.closed_time AS order_closed_time, 
		ord.phone AS order_phone, 
		ord.name AS order_name, 
		ord.email AS order_email, 
		ord.delivery_price AS order_delivery_price,
		ord.delivery_address AS order_delivery_address,
		a.id AS advert_id, 
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
		(SELECT COALESCE(r.rating, 0) AS rating FROM public.advert a LEFT JOIN public.review r ON r.advert_id = a.id WHERE a.id = ord.advert_id)	
	FROM 
		public.advert a
	LEFT JOIN 
		public.city c ON a.city_id = c.id
	LEFT JOIN 
		public.category cat ON a.category_id = cat.id
	LEFT JOIN 
		public.order ord ON a.id = ord.advert_id
	WHERE a.user_id = $1 AND a.id = ord.advert_id;`

	logging.LogInfo(logger, "SELECT FROM advert, cart, category, city, advert_image")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLGetSoldOrdersByUserID, userID)

	ol.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while executing select adverts from the cart, err=%w", err))
		ol.metrics.IncreaseErrors(funcName)

		return nil, err
	}

	defer rows.Close()

	var orderList []*models.ReturningOrder

	for rows.Next() {
		categoryModel := models.Category{}
		cityModel := models.City{}
		advertModel := models.Advert{}
		photoPad := models.PhotoPadSoloImage{}
		orderItem := models.OrderItem{}

		if err := rows.Scan(&orderItem.ID, &orderItem.Status, &orderItem.Created, &orderItem.Updated, &orderItem.Closed,
			&orderItem.Phone, &orderItem.Name, &orderItem.Email, &orderItem.DeliveryPrice, &orderItem.Address,
			&advertModel.ID, &advertModel.UserID, &cityModel.ID, &cityModel.CityName, &cityModel.Translation,
			&categoryModel.ID, &categoryModel.Name, &categoryModel.Translation, &advertModel.Title,
			&advertModel.Description, &advertModel.Price, &advertModel.CreatedTime, &advertModel.ClosedTime,
			&advertModel.IsUsed, &photoPad.Photo, &orderItem.Rating); err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while scanning adverts from the cart, err=%w", err))

			return nil, err
		}

		advertModel.CityID = cityModel.ID
		advertModel.CategoryID = categoryModel.ID

		photoURLToInsert := ""
		if photoPad.Photo != nil {
			photoURLToInsert = *photoPad.Photo
		}

		returningAdvert := models.ReturningAdvert{
			Advert:   advertModel,
			City:     cityModel,
			Category: categoryModel,
		}

		decodedImage, err := utils.DecodeImage(photoURLToInsert)

		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while decoding image, err=%w", err))

			return nil, err
		}

		returningAdvert.Photos = append(returningAdvert.Photos, photoURLToInsert)
		returningAdvert.PhotosIMG = append(returningAdvert.PhotosIMG, decodedImage)

		ReturningOrder := models.ReturningOrder{
			OrderItem:       orderItem,
			ReturningAdvert: returningAdvert,
		}

		orderList = append(orderList, &ReturningOrder)
	}

	return orderList, nil
}

func (ol *OrderStorage) GetSoldOrdersByUserID(ctx context.Context, userID uint) ([]*models.ReturningOrder, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	orderList := []*models.ReturningOrder{}

	err := pgx.BeginFunc(ctx, ol.pool, func(tx pgx.Tx) error {
		orderListInner, err := ol.getSoldOrdersByUserID(ctx, tx, userID)
		orderList = orderListInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting orders list, err=%w", err))

		return nil, err
	}

	if orderList == nil {
		orderList = []*models.ReturningOrder{}
	}

	return orderList, nil
}

func (ol *OrderStorage) createOrderByID(ctx context.Context, tx pgx.Tx, userID uint,
	data *models.ReceivedOrderItem) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCreateOrder :=
		`INSERT INTO public."order"(
			user_id, advert_id, order_status, phone, name, surname, patronymic, email, delivery_price, delivery_address)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	logging.LogInfo(logger, "INSERT INTO user")

	var err error

	const (
		paidStatus     string = "В обработке"
		surnamePlug    string = "Фамилия"
		patronymicPlug string = "Отчество"
	)

	start := time.Now()

	_, err = tx.Exec(ctx, SQLCreateOrder, userID, data.AdvertID, paidStatus, data.Phone, data.Name, surnamePlug,
		patronymicPlug, data.Email, data.DeliveryPrice, data.DeliveryAddress)

	ol.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing create order query, err=%w",
			err))
		ol.metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

func (ol *OrderStorage) CreateOrderByID(ctx context.Context, userID uint, data *models.ReceivedOrderItem) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, ol.pool, func(tx pgx.Tx) error {
		err := ol.createOrderByID(ctx, tx, userID, data)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong while creating order, err=%w", err))

			return err
		}

		return nil
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while creating user, err=%w", err))

		return err
	}

	return nil
}

func (ol *OrderStorage) orderExists(ctx context.Context, tx pgx.Tx, userID, advertID uint) (bool, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLUserExists := `SELECT EXISTS(SELECT 1 FROM public."order" WHERE user_id = $1 and advert_id = $2);`

	logging.LogInfo(logger, "SELECT FROM user")

	start := time.Now()

	userLine := tx.QueryRow(ctx, SQLUserExists, userID, advertID)

	ol.metrics.AddDuration(funcName, time.Since(start))

	var exists bool

	if err := userLine.Scan(&exists); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning order exists, err=%w", err))

		return false, err
	}

	return exists, nil
}

func (ol *OrderStorage) OrderExists(ctx context.Context, userID, advertID uint) bool {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var exists bool

	err := pgx.BeginFunc(ctx, ol.pool, func(tx pgx.Tx) error {
		userExists, err := ol.orderExists(ctx, tx, userID, advertID)
		exists = userExists

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing advert exists query, err=%w",
			err))
	}

	return exists
}

func (ol *OrderStorage) createReview(ctx context.Context, tx pgx.Tx, userID, advertID, rating uint) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCreateReview :=
		`INSERT INTO public.review(
			user_id, advert_id, rating)
			VALUES ($1, $2, $3);`

	logging.LogInfo(logger, "INSERT INTO review")

	var err error

	start := time.Now()

	_, err = tx.Exec(ctx, SQLCreateReview, userID, advertID, rating)

	ol.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing create review query, err=%w",
			err))
		ol.metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

func (ol *OrderStorage) CreateReview(ctx context.Context, userID, advertID, rating uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, ol.pool, func(tx pgx.Tx) error {
		err := ol.createReview(ctx, tx, userID, advertID, rating)

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while creating review, err=%w", err))

		return err
	}

	return nil
}
