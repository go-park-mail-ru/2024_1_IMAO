package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	advuc "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/adverts/usecases"
	useruc "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/usecases"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	errNotInCart = errors.New("there is no advert in the cart")
)

type CartListWrapper struct {
	CartList *models.CartList
	Pool     *pgxpool.Pool
	Logger   *zap.SugaredLogger
}

func (cl *CartListWrapper) getCartByUserID(ctx context.Context, tx pgx.Tx, userID uint) ([]*models.ReturningAdvert, error) {
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
			a.is_used
		FROM 
			public.advert a
		LEFT JOIN 
			public.city c ON a.city_id = c.id
		LEFT JOIN 
			public.category cat ON a.category_id = cat.id
		LEFT JOIN 
			public.cart cart ON a.id = cart.advert_id
		WHERE cart.user_id = $1;`

	cl.Logger.Infof(`
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
			a.is_used
		FROM 
			public.advert a
		LEFT JOIN 
			public.city c ON a.city_id = c.id
		LEFT JOIN 
			public.category cat ON a.category_id = cat.id
		LEFT JOIN 
			public.cart cart ON a.id = cart.advert_id
		WHERE cart.user_id = %1;`, userID)

	rows, err := tx.Query(ctx, SQLAdvertByUserId, userID)
	if err != nil {
		cl.Logger.Errorf("Something went wrong while executing select adverts from the cart, err=%v", err)

		return nil, err
	}
	defer rows.Close()

	var adsList []*models.ReturningAdvert
	for rows.Next() {

		categoryModel := models.Category{}
		cityModel := models.City{}
		advertModel := models.Advert{}

		if err := rows.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &cityModel.CityName, &cityModel.Translation,
			&categoryModel.ID, &categoryModel.Name, &categoryModel.Translation, &advertModel.Title, &advertModel.Description, &advertModel.Price,
			&advertModel.CreatedTime, &advertModel.ClosedTime, &advertModel.IsUsed); err != nil {

			cl.Logger.Errorf("Something went wrong while scanning adverts from the cart, err=%v", err)

			return nil, err
		}

		advertModel.CityID = cityModel.ID
		advertModel.CategoryID = categoryModel.ID

		returningAdvertList := models.ReturningAdvert{
			Advert:   advertModel,
			City:     cityModel,
			Category: categoryModel,
		}

		adsList = append(adsList, &returningAdvertList)
	}

	return adsList, nil
}

func (cl *CartListWrapper) GetCartByUserID(ctx context.Context, userID uint, userList useruc.UsersInfo, advertsList advuc.AdvertsInfo) ([]*models.ReturningAdvert, error) {
	cart := []*models.ReturningAdvert{}

	err := pgx.BeginFunc(ctx, cl.Pool, func(tx pgx.Tx) error {
		cartInner, err := cl.getCartByUserID(ctx, tx, userID)
		cart = cartInner

		return err
	})

	if err != nil {
		cl.Logger.Errorf("Something went wrong while getting adverts list, err=%v", err)

		return nil, err
	}

	if cart == nil {
		cart = []*models.ReturningAdvert{}
	}

	return cart, nil
}

func (cl *CartListWrapper) deleteAdvByIDs(ctx context.Context, tx pgx.Tx, userID uint, advertID uint) error {
	SQLDeleteFromCart := `DELETE FROM public.cart
		WHERE user_id = $1 AND advert_id = $2;`

	cl.Logger.Infof(`DELETE FROM public.cart
		WHERE user_id = $1 AND advert_id = $2;`, userID, advertID, userID, advertID)

	var err error

	_, err = tx.Exec(ctx, SQLDeleteFromCart, userID, advertID)

	if err != nil {
		cl.Logger.Errorf("Something went wrong while executing advert delete from the cart, err=%v", err)
		return fmt.Errorf("Something went wrong while executing advert delete from the cart", err)
	}

	return nil
}

func (cl *CartListWrapper) DeleteAdvByIDs(ctx context.Context, userID uint, advertID uint, userList useruc.UsersInfo, advertsList advuc.AdvertsInfo) error {

	err := pgx.BeginFunc(ctx, cl.Pool, func(tx pgx.Tx) error {
		err := cl.deleteAdvByIDs(ctx, tx, userID, advertID)

		return err
	})

	if err != nil {
		cl.Logger.Errorf("Something went wrong while getting adverts list, most likely , err=%v", err)

		return err
	}

	return nil
}

func (cl *CartListWrapper) appendAdvByIDs(ctx context.Context, tx pgx.Tx, userID uint, advertID uint) (bool, error) {
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
	cl.Logger.Infof(`WITH deletion AS (
		DELETE FROM public.cart
		WHERE user_id = %s AND advert_id = %s
		RETURNING user_id, advert_id
	)
	INSERT INTO public.cart (user_id, advert_id)
	SELECT %s, %s
	WHERE NOT EXISTS (
		SELECT 1 FROM deletion
	) RETURNING true;
	`, userID, advertID, userID, advertID)

	userLine := tx.QueryRow(ctx, SQLAddToCart, userID, advertID)

	added := false

	if err := userLine.Scan(&added); err != nil {
		//cl.Logger.Errorf("Error while scanning advert added, err=%v", err)
		//return false, err
	}

	return added, nil
}

func (cl *CartListWrapper) AppendAdvByIDs(ctx context.Context, userID uint, advertID uint, userList useruc.UsersInfo, advertsList advuc.AdvertsInfo) bool {
	var added bool

	err := pgx.BeginFunc(ctx, cl.Pool, func(tx pgx.Tx) error {
		addedInner, err := cl.appendAdvByIDs(ctx, tx, userID, advertID)
		added = addedInner

		return err
	})

	if err != nil {
		cl.Logger.Errorf("Error while executing addvert add to cart, err=%v", err)
	}

	return added
}

func NewCartList(pool *pgxpool.Pool, logger *zap.SugaredLogger) *CartListWrapper {
	return &CartListWrapper{
		CartList: &models.CartList{
			Items: make([]*models.CartItem, 0),
			Mux:   sync.RWMutex{},
		},
		Pool:   pool,
		Logger: logger,
	}
}
