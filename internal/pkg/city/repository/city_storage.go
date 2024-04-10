package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	errProfileNotExists = errors.New("profile does not exist")
	NameSeqProfile      = pgx.Identifier{"public", "city_id_seq"} //nolint:gochecknoglobals
)

type CityListWrapper struct {
	CityList *models.CityList
	Pool     *pgxpool.Pool
	Logger   *zap.SugaredLogger
}

func (cl *CityListWrapper) getCityList(ctx context.Context, tx pgx.Tx) (*models.CityList, error) {
	SQLCityList := `SELECT id, name, translation FROM public.city;`
	cl.Logger.Infof(`SELECT id, name, translation FROM public.city;`)
	rows, err := tx.Query(ctx, SQLCityList)
	if err != nil {
		cl.Logger.Errorf("Something went wrong while executing select city query, err=%v", err)

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
		cl.Logger.Errorf("Something went wrong while scanning city rows, err=%v", err)

		return nil, err
	}

	return &cityList, nil
}

func (cl *CityListWrapper) GetCityList(ctx context.Context) (*models.CityList, error) {
	var cityList *models.CityList

	err := pgx.BeginFunc(ctx, cl.Pool, func(tx pgx.Tx) error {
		cityListInner, err := cl.getCityList(ctx, tx)
		cityList = cityListInner

		return err
	})

	if err != nil {
		cl.Logger.Errorf("Something went wrong while getting city list, err=%v", err)

		return nil, err
	}

	return cityList, nil
}

func NewCityList(pool *pgxpool.Pool, logger *zap.SugaredLogger) *CityListWrapper {
	return &CityListWrapper{
		CityList: &models.CityList{
			CityItems: make([]*models.City, 0),
			Mux:       sync.RWMutex{},
		},
		Pool:   pool,
		Logger: logger,
	}
}
