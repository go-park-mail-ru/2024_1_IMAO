package repository

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	lifting = 1
	premium = 2
	maximum = 3
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
		logging.LogError(logger, fmt.Errorf("error while scanning user exists, err=%w", err))

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
		logging.LogError(logger, fmt.Errorf("error while executing user exists query, err=%w", err))
	}

	return exists
}

func (paymentsStorage *PaymentsStorage) getPriceAndDescription(ctx context.Context, tx pgx.Tx, advertID,
	rateCode uint) (*models.PriceAndDescription, error) {
	priceMap := map[uint]string{
		1: "35",
		2: "100",
		3: "199",
	}

	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLSelectPriceAndDescription := `SELECT title, c.translation, cat.translation
	FROM public.advert a
	LEFT JOIN 
	public.city c ON a.city_id = c.id
	LEFT JOIN 
	public.category cat ON a.category_id = cat.id
	WHERE a.id=$1;
	`

	logging.LogInfo(logger, "SELECT FROM advert")

	start := time.Now()

	userLine := tx.QueryRow(ctx, SQLSelectPriceAndDescription, advertID)

	paymentsStorage.metrics.AddDuration(funcName, time.Since(start))

	var (
		title               string
		cityTranslation     string
		categoryTranslation string
		priceAndDescription models.PriceAndDescription
	)

	if err := userLine.Scan(&title, &cityTranslation, &categoryTranslation); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning advert title, err=%w", err))

		return nil, err
	}

	if val, ok := priceMap[rateCode]; ok {
		priceAndDescription.Price = val
		priceAndDescription.URLEnding = cityTranslation + "/" + categoryTranslation +
			"/" + strconv.FormatUint(uint64(advertID), 10)

		switch key := rateCode; key {
		case lifting:
			priceAndDescription.Description = fmt.Sprintf("Платное продвижение товара %s на 1 день", title)
			priceAndDescription.Duration = "1 day"
		case premium:
			priceAndDescription.Description = fmt.Sprintf("Платное продвижение товара %s на 3 дня", title)
			priceAndDescription.Duration = "3 days"
		case maximum:
			priceAndDescription.Description = fmt.Sprintf("Платное продвижение товара %s на 7 дней", title)
			priceAndDescription.Duration = "7 days"
		default:
			log.Println("no such key in the map")
		}
	} else {
		logging.LogError(logger, fmt.Errorf("no such key in the map"))

		return nil, fmt.Errorf("no such key in the map")
	}

	return &priceAndDescription, nil
}

func (paymentsStorage *PaymentsStorage) GetPriceAndDescription(ctx context.Context, advertID,
	rateCode uint) (*models.PriceAndDescription, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var priceAndDescription *models.PriceAndDescription

	err := pgx.BeginFunc(ctx, paymentsStorage.pool, func(tx pgx.Tx) error {
		priceAndDescriptionInner, err := paymentsStorage.getPriceAndDescription(ctx, tx, advertID, rateCode)
		priceAndDescription = priceAndDescriptionInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing get payment's price and description, err=%w", err))

		return nil, err
	}

	return priceAndDescription, nil
}

func (paymentsStorage *PaymentsStorage) checkAdvertOwnership(ctx context.Context, tx pgx.Tx,
	advertID, userID uint) (bool, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCheckAdvertOwnership := `SELECT 
    							EXISTS(
    								SELECT 1 
    								FROM public.advert 
    								WHERE id=$1 AND user_id=$2 AND is_promoted=false);`

	logging.LogInfo(logger, "SELECT FROM advert")

	start := time.Now()

	userLine := tx.QueryRow(ctx, SQLCheckAdvertOwnership, advertID, userID)

	paymentsStorage.metrics.AddDuration(funcName, time.Since(start))

	var ownership bool

	if err := userLine.Scan(&ownership); err != nil {
		logging.LogError(logger, fmt.Errorf("error while scanning advert ownership exists, err=%w", err))

		return false, err
	}

	return ownership, nil
}

func (paymentsStorage *PaymentsStorage) CheckAdvertOwnership(ctx context.Context, advertID, userID uint) bool {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var ownership bool

	err := pgx.BeginFunc(ctx, paymentsStorage.pool, func(tx pgx.Tx) error {
		ownershipExists, err := paymentsStorage.checkAdvertOwnership(ctx, tx, advertID, userID)
		ownership = ownershipExists

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while executing user exists query, err=%w", err))
	}

	return ownership
}

func (paymentsStorage *PaymentsStorage) createPayment(ctx context.Context, tx pgx.Tx, payment *models.Payment,
	idempotencyKey string, advertID uint, duration string) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCreateProfile := `INSERT INTO public.payments(
		advert_id, payment_uuid, payment_value, payment_description, payment_status, payment_form_url, 
                            created_time, idempotency_key, promotion_duration)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

	logging.LogInfo(logger, "INSERT INTO payments")

	var err error

	float64Value, err := strconv.ParseFloat(payment.Amount.Value, 64)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing insert payment query, err=%w",
			err))
		paymentsStorage.metrics.IncreaseErrors(funcName)

		return err
	}

	uintValue := uint(float64Value)

	start := time.Now()

	_, err = tx.Exec(ctx, SQLCreateProfile, advertID, payment.ID, uintValue, payment.Description,
		payment.Status, payment.Confirmation.ConfirmationURL, payment.CreatedAt, idempotencyKey, duration)

	paymentsStorage.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing insert payment query, err=%w",
			err))
		paymentsStorage.metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

func (paymentsStorage *PaymentsStorage) CreatePayment(ctx context.Context, payment *models.Payment,
	idempotencyKey string, advertID uint, duration string) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, paymentsStorage.pool, func(tx pgx.Tx) error {
		err := paymentsStorage.createPayment(ctx, tx, payment, idempotencyKey, advertID, duration)

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while inserting payment in db, err=%w", err))

		return err
	}

	return nil
}
