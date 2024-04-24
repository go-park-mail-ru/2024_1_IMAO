package repository

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	errProfileNotExists = errors.New("profile does not exist")
	NameSeqProfile      = pgx.Identifier{"public", "city_id_seq"} //nolint:gochecknoglobals
)

type CityStorage struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func NewCityStorage(pool *pgxpool.Pool, logger *zap.SugaredLogger) *CityStorage {
	return &CityStorage{
		pool:   pool,
		logger: logger,
	}
}

func (cl *CityStorage) getCityList(ctx context.Context, tx pgx.Tx) (*models.CityList, error) {
	SQLCityList := `SELECT id, name, translation FROM public.city;`
	cl.logger.Infof(`SELECT id, name, translation FROM public.city;`)
	rows, err := tx.Query(ctx, SQLCityList)
	if err != nil {
		cl.logger.Errorf("Something went wrong while executing select city query, err=%v", err)

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
		cl.logger.Errorf("Something went wrong while scanning city rows, err=%v", err)

		return nil, err
	}

	return &cityList, nil
}

func (cl *CityStorage) GetCityList(ctx context.Context) (*models.CityList, error) {
	var cityList *models.CityList

	err := pgx.BeginFunc(ctx, cl.pool, func(tx pgx.Tx) error {
		cityListInner, err := cl.getCityList(ctx, tx)
		cityList = cityListInner

		return err
	})

	if err != nil {
		cl.logger.Errorf("Something went wrong while getting city list, err=%v", err)

		return nil, err
	}

	return cityList, nil
}
