package repository

import (
	"context"
	"fmt"
	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"
	"time"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type CityStorage struct {
	pool    *pgxpool.Pool
	metrics *mymetrics.DatabaseMetrics
}

func NewCityStorage(pool *pgxpool.Pool, metrics *mymetrics.DatabaseMetrics) *CityStorage {
	return &CityStorage{
		pool:    pool,
		metrics: metrics,
	}
}

func (cl *CityStorage) getCityList(ctx context.Context, tx pgx.Tx) (*models.CityList, error) {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCityList := `SELECT id, name, translation FROM public.city;`

	logging.LogInfo(logger, "SELECT FROM city")

	start := time.Now()

	rows, err := tx.Query(ctx, SQLCityList)

	cl.metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger,
			fmt.Errorf("something went wrong while executing select city query, err=%w", err))
		cl.metrics.IncreaseErrors(funcName)

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
		logging.LogError(logger, fmt.Errorf("something went wrong while scanning city rows, err=%w", err))

		return nil, err
	}

	return &cityList, nil
}

func (cl *CityStorage) GetCityList(ctx context.Context) (*models.CityList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var cityList *models.CityList

	err := pgx.BeginFunc(ctx, cl.pool, func(tx pgx.Tx) error {
		cityListInner, err := cl.getCityList(ctx, tx)
		cityList = cityListInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting city list, err=%w", err))

		return nil, err
	}

	return cityList, nil
}
