package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	errProfileNotExists = errors.New("profile does not exist")
	NameSeqProfile      = pgx.Identifier{"public", "city_id_seq"} //nolint:gochecknoglobals
)

type CityStorage struct {
	pool *pgxpool.Pool
}

func NewCityStorage(pool *pgxpool.Pool) *CityStorage {
	return &CityStorage{
		pool: pool,
	}
}

func (cl *CityStorage) getCityList(ctx context.Context, tx pgx.Tx) (*models.CityList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCityList := `SELECT id, name, translation FROM public.city;`

	logging.LogInfo(logger, "SELECT FROM city")

	rows, err := tx.Query(ctx, SQLCityList)
	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while executing select city query, err=%v", err))

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

func (cl *CityStorage) GetCityList(ctx context.Context) (*models.CityList, error) {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	var cityList *models.CityList

	err := pgx.BeginFunc(ctx, cl.pool, func(tx pgx.Tx) error {
		cityListInner, err := cl.getCityList(ctx, tx)
		cityList = cityListInner

		return err
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong while getting city list, err=%v", err))

		return nil, err
	}

	return cityList, nil
}
