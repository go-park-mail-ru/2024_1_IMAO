package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PaymentsStorage struct {
	pool    *pgxpool.Pool
	metrics *mymetrics.DatabaseMetrics
}

func NewPaymentsStorage(pool *pgxpool.Pool, metrics *mymetrics.DatabaseMetrics) *PaymentsStorage {
	return &PaymentsStorage{
		pool:    pool,
		metrics: metrics,
	}
}

func (paymentsStorage *PaymentsStorage) getCityList(ctx context.Context, tx pgx.Tx) (*models.CityList, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCityList := `SELECT id, name, translation FROM public.city;`

	logging.LogInfo(logger, "SELECT FROM city")

	start := time.Now()
	rows, err := tx.Query(ctx, SQLCityList)
	paymentsStorage.metrics.AddDuration(funcName, time.Since(start))
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select city query, err=%v",
			err))
		paymentsStorage.metrics.IncreaseErrors(funcName)

		return nil, err
	}
	defer rows.Close()

	cityList := models.CityList{}
	for rows.Next() {
		city := models.City{}
		if err := rows.Scan(&city.ID, &city.CityName, &city.Translation); err != nil {
			return nil, err
		}
		cityList.CityItems = append(cityList.CityItems, &city)
	}

	if err := rows.Err(); err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning city rows, err=%v", err))

		return nil, err
	}

	return &cityList, nil
}

func (paymentsStorage *PaymentsStorage) GetCityList(ctx context.Context) (*models.CityList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var cityList *models.CityList

	err := pgx.BeginFunc(ctx, paymentsStorage.pool, func(tx pgx.Tx) error {
		cityListInner, err := paymentsStorage.getCityList(ctx, tx)
		cityList = cityListInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting city list, err=%v", err))

		return nil, err
	}

	return cityList, nil
}

func (paymentsStorage *PaymentsStorage) paymentFormExists(ctx context.Context, tx pgx.Tx, email string) (bool, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLUserExists := `SELECT EXISTS(SELECT 1 FROM public."user" WHERE email=$1 );`

	logging.LogInfo(logger, "SELECT FROM payments")

	start := time.Now()
	userLine := tx.QueryRow(ctx, SQLUserExists, email)
	paymentsStorage.metrics.AddDuration(funcName, time.Since(start))

	var exists bool

	if err := userLine.Scan(&exists); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning user exists, err=%v", err))

		return false, err
	}

	return exists, nil
}

func (paymentsStorage *PaymentsStorage) PaymentFormExists(ctx context.Context, email string) bool {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var exists bool

	err := pgx.BeginFunc(ctx, paymentsStorage.pool, func(tx pgx.Tx) error {
		userExists, err := paymentsStorage.paymentFormExists(ctx, tx, email)
		exists = userExists

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing user exists query, err=%v", err))

	}

	return exists
}

func (paymentsStorage *PaymentsStorage) getPriceAndDescription(ctx context.Context, tx pgx.Tx, advertId, rateCode uint) (*models.PriceAndDescription, error) {
	priceMap := map[uint]string{
		1: "100",
		2: "250",
		3: "450",
	}

	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCheckAdvertOwnership := `SELECT title, c.translation, cat.translation
	FROM public.advert a
	LEFT JOIN 
	public.city c ON a.city_id = c.id
	LEFT JOIN 
	public.category cat ON a.category_id = cat.id
	WHERE a.id=$1;
	`

	logging.LogInfo(logger, "SELECT FROM advert")

	start := time.Now()
	userLine := tx.QueryRow(ctx, SQLCheckAdvertOwnership, advertId)
	paymentsStorage.metrics.AddDuration(funcName, time.Since(start))

	var (
		title                string
		city_translation     string
		category_translation string
	)
	priceAndDescription := models.PriceAndDescription{}

	if err := userLine.Scan(&title, &city_translation, &category_translation); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning advert title, err=%v", err))

		return nil, err
	}

	if val, ok := priceMap[rateCode]; ok {
		priceAndDescription.Price = val
		priceAndDescription.UrlEnding = city_translation + "/" + category_translation + "/" + strconv.FormatUint(uint64(advertId), 10)

		switch key := rateCode; key {
		case 1:
			priceAndDescription.Description = `Платное продвижение товара "` + title + `" на 1 день`
		case 2:
			priceAndDescription.Description = `Платное продвижение товара "` + title + `" на 3 дня`
		case 3:
			priceAndDescription.Description = `Платное продвижение товара "` + title + `" на 7 дней`
		default:
			fmt.Println("no such key in the map")
		}
	} else {
		logging.LogError(logger, fmt.Errorf("no such key in the map"))

		return nil, fmt.Errorf("no such key in the map")
	}

	return &priceAndDescription, nil

}

func (paymentsStorage *PaymentsStorage) GetPriceAndDescription(ctx context.Context, advertId, rateCode uint) (*models.PriceAndDescription, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var priceAndDescription *models.PriceAndDescription

	err := pgx.BeginFunc(ctx, paymentsStorage.pool, func(tx pgx.Tx) error {
		priceAndDescriptionInner, err := paymentsStorage.getPriceAndDescription(ctx, tx, advertId, rateCode)
		priceAndDescription = priceAndDescriptionInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing get payment's price and description, err=%v", err))

		return nil, err
	}

	return priceAndDescription, nil
}

func (paymentsStorage *PaymentsStorage) checkAdvertOwnership(ctx context.Context, tx pgx.Tx, advertId, userId uint) (bool, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCheckAdvertOwnership := `SELECT EXISTS(SELECT 1 FROM public.advert WHERE id=$1 AND user_id=$2 );`

	logging.LogInfo(logger, "SELECT FROM advert")

	start := time.Now()
	userLine := tx.QueryRow(ctx, SQLCheckAdvertOwnership, advertId, userId)
	paymentsStorage.metrics.AddDuration(funcName, time.Since(start))

	var ownership bool

	if err := userLine.Scan(&ownership); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning advert ownership exists, err=%v", err))

		return false, err
	}

	return ownership, nil
}

func (paymentsStorage *PaymentsStorage) CheckAdvertOwnership(ctx context.Context, advertId, userId uint) bool {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var ownership bool

	err := pgx.BeginFunc(ctx, paymentsStorage.pool, func(tx pgx.Tx) error {
		ownershipExists, err := paymentsStorage.checkAdvertOwnership(ctx, tx, advertId, userId)
		ownership = ownershipExists

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing user exists query, err=%v", err))

	}

	return ownership
}

func (paymentsStorage *PaymentsStorage) createPayment(ctx context.Context, tx pgx.Tx, payment *models.Payment,
	idempotencyKey string, advertId uint) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCreateProfile := `INSERT INTO public.payments(
		advert_id, payment_uuid, payment_value, payment_description, payment_status, payment_form_url, created_time, idempotency_key)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	logging.LogInfo(logger, "INSERT INTO payments")

	var err error

	float64Value, err := strconv.ParseFloat(payment.Amount.Value, 64)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing insert payment query, err=%v", err))
		paymentsStorage.metrics.IncreaseErrors(funcName)

		return err
	}

	uintValue := uint(float64Value)

	start := time.Now()
	_, err = tx.Exec(ctx, SQLCreateProfile, advertId, payment.ID, uintValue, payment.Description,
		payment.Status, payment.Confirmation.ConfirmationURL, payment.CreatedAt, idempotencyKey)
	paymentsStorage.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing insert payment query, err=%v", err))
		paymentsStorage.metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

func (paymentsStorage *PaymentsStorage) CreatePayment(ctx context.Context, payment *models.Payment, idempotencyKey string,
	advertId uint) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, paymentsStorage.pool, func(tx pgx.Tx) error {
		err := paymentsStorage.createPayment(ctx, tx, payment, idempotencyKey, advertId)

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while inserting payment in db, err=%v", err))

		return err
	}

	return nil
}
