package storage

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mdigger/translit"
	"go.uber.org/zap"
)

var (
	errWrongAdvertID      = errors.New("wrong advert ID")
	errWrongCityName      = errors.New("wrong city name")
	errWrongCategoryName  = errors.New("wrong category name")
	errWrongIDinCategory  = errors.New("there is no ad with such id in category")
	errWrongIDinCity      = errors.New("there is no ad with such id in city")
	errWrongAdvertsAmount = errors.New("too many elements specified")
	errAlreadyClosed      = errors.New("advert already closed")
)

type AdvertsListWrapper struct {
	AdvertsList *models.AdvertsList
	Pool        *pgxpool.Pool
	Logger      *zap.SugaredLogger
}

func (ads *AdvertsListWrapper) GetAdvertByOnlyByID(advertID uint) (*models.ReturningAdvert, error) {
	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	if advertID > ads.AdvertsList.AdvertsCounter {
		return nil, errWrongAdvertID
	}

	cityID := ads.AdvertsList.Adverts[advertID-1].CityID
	categoryID := ads.AdvertsList.Adverts[advertID-1].CategoryID

	return &models.ReturningAdvert{
		Advert:   *ads.AdvertsList.Adverts[advertID-1],
		City:     *ads.AdvertsList.Cities[cityID-1],
		Category: *ads.AdvertsList.Categories[categoryID-1],
	}, nil
}

func (ads *AdvertsListWrapper) getAdvert(ctx context.Context, tx pgx.Tx, advertID uint, city, category string) (*models.ReturningAdvert, error) {
	SQLAdvertById := `
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
		WHERE a.id = $1;`

	ads.Logger.Infof(`
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
		WHERE a.id = %s;`, advertID)
	advertLine := tx.QueryRow(ctx, SQLAdvertById, advertID)

	categoryModel := models.Category{}
	cityModel := models.City{}
	advertModel := models.Advert{}

	if err := advertLine.Scan(&advertModel.ID, &advertModel.UserID, &cityModel.ID, &cityModel.CityName, &cityModel.Translation,
		&categoryModel.ID, &categoryModel.Name, &categoryModel.Translation, &advertModel.Title, &advertModel.Description, &advertModel.Price,
		&advertModel.CreatedTime, &advertModel.ClosedTime, &advertModel.IsUsed); err != nil {

		ads.Logger.Errorf("Something went wrong while scanning advert, err=%v", err)

		return nil, err
	}

	advertModel.CityID = cityModel.ID
	advertModel.CategoryID = categoryModel.ID

	return &models.ReturningAdvert{
		Advert:   advertModel,
		City:     cityModel,
		Category: categoryModel,
	}, nil
}

func (ads *AdvertsListWrapper) GetAdvert(ctx context.Context, advertID uint, city, category string) (*models.ReturningAdvert, error) {
	var advertsList *models.ReturningAdvert

	err := pgx.BeginFunc(ctx, ads.Pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.getAdvert(ctx, tx, advertID, city, category)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		ads.Logger.Errorf("Something went wrong while getting adverts list, err=%v", err)

		return nil, err
	}

	return advertsList, nil
}

// func (ads *AdvertsListWrapper) GetAdvertsByCity(city string, startID, num uint) ([]*models.ReturningAdInList, error) {
// 	if num > ads.AdvertsList.AdvertsCounter {
// 		return nil, errWrongAdvertsAmount
// 	}

// 	cityID, err := ads.GetCityID(city)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ads.AdvertsList.Mux.Lock()
// 	defer ads.AdvertsList.Mux.Unlock()

// 	var returningAds []*models.ReturningAdInList
// 	var counter uint = 0

// 	for counter != num && counter+startID-1 != ads.AdvertsList.AdvertsCounter {
// 		ad := ads.AdvertsList.Adverts[startID+counter-1]
// 		exists := ad.Active && !ad.Deleted

// 		if exists && ad.CityID == cityID {
// 			returningAds = append(returningAds, &models.ReturningAdInList{
// 				ID:       ad.ID,
// 				Title:    ad.Title,
// 				Price:    ad.Price,
// 				City:     ads.AdvertsList.Cities[ad.CityID-1].Translation,
// 				Category: ads.AdvertsList.Categories[ad.CategoryID-1].Translation,
// 			})
// 		}

// 		counter++
// 	}

// 	return returningAds, nil
// }

func (ads *AdvertsListWrapper) getAdvertsByCity(ctx context.Context, tx pgx.Tx, city string, startID, num uint) ([]*models.ReturningAdInList, error) {
	SQLAdvertsByCity := `SELECT a.id, c.translation, category.translation, a.title, a.price
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.id >= $1 AND a.advert_status = 'Активно' AND c.translation = $2
	LIMIT $3;
	`
	ads.Logger.Infof(`SELECT a.id, c.translation, category.translation, a.title, a.price
	FROM public.advert a
	INNER JOIN city c ON a.city_id = c.id
	INNER JOIN category ON a.category_id = category.id
	WHERE a.id >= %s AND a.advert_status = 'Активно' AND c.translation = %s
	LIMIT %s`, startID, city, num)
	rows, err := tx.Query(ctx, SQLAdvertsByCity, startID, city, num)
	if err != nil {
		ads.Logger.Errorf("Something went wrong while executing select adverts query, err=%v", err)

		return nil, err
	}
	defer rows.Close()

	var adsList []*models.ReturningAdInList
	for rows.Next() {
		returningAdInList := models.ReturningAdInList{}
		if err := rows.Scan(&returningAdInList.ID, &returningAdInList.City, &returningAdInList.Category, &returningAdInList.Title, &returningAdInList.Price); err != nil {
			return nil, err
		}
		adsList = append(adsList, &returningAdInList)
	}

	if err := rows.Err(); err != nil {
		ads.Logger.Errorf("Something went wrong while scanning adverts rows, err=%v", err)

		return nil, err
	}

	return adsList, nil
}

func (ads *AdvertsListWrapper) GetAdvertsByCity(ctx context.Context, city string, startID, num uint) ([]*models.ReturningAdInList, error) {
	var advertsList []*models.ReturningAdInList

	err := pgx.BeginFunc(ctx, ads.Pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.getAdvertsByCity(ctx, tx, city, startID, num)
		advertsList = advertsListInner

		return err
	})

	if err != nil {
		ads.Logger.Errorf("Something went wrong while getting adverts list, err=%v", err)

		return nil, err
	}

	return advertsList, nil
}

func (ads *AdvertsListWrapper) GetAdvertsByCategory(category, city string, startID, num uint) ([]*models.ReturningAdInList, error) {
	if num > ads.AdvertsList.AdvertsCounter {
		return nil, errWrongAdvertsAmount
	}

	cityID, err := ads.GetCityID(city)
	if err != nil {
		return nil, err
	}

	categoryID, err := ads.GetCategoryID(category)
	if err != nil {
		return nil, err
	}

	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	var returningAds []*models.ReturningAdInList
	var counter uint = 0

	for counter != num && counter+startID-1 != ads.AdvertsList.AdvertsCounter {
		ad := ads.AdvertsList.Adverts[startID+counter-1]
		exists := ad.Active && !ad.Deleted

		if exists && ad.CityID == cityID && ad.CategoryID == categoryID {
			returningAds = append(returningAds, &models.ReturningAdInList{
				ID:       ad.ID,
				Title:    ad.Title,
				Price:    ad.Price,
				City:     ads.AdvertsList.Cities[ad.CityID-1].Translation,
				Category: ads.AdvertsList.Categories[ad.CategoryID-1].Translation,
			})
		}

		counter++
	}

	return returningAds, nil
}

func (ads *AdvertsListWrapper) getAdvertsForUserWhereStatusIs(ctx context.Context, tx pgx.Tx, userId, deleted uint) (*models.ReturningAdvertList, error) {

	statusID := 23
	if deleted == 1 {
		statusID = 24
	}

	SQLGetAdvertsForUserWhereStatusIs := `SELECT 
		a.id, 
		a.user_id, 
		a.city_id, 
		a.category_id, 
		a.title, 
		a.description, 
		a.price, 
		a.created_time, 
		a.closed_time, 
		a.is_used, 
		a.status_id,
		c.name AS city_name,
		c.translation AS city_translation,
		cat.name AS category_name,
		cat.translation AS category_translation
	FROM 
		public.advert a
	INNER JOIN 
		city c ON a.city_id = c.id
	INNER JOIN 
		category cat ON a.category_id = cat.id
	WHERE 
		a.user_id = $1 AND a.status_id = $2;
	`
	ads.Logger.Infof(`SELECT a.id,	a.user_id, 	a.city_id, 	a.category_id, 	a.title, a.description, a.price, a.created_time, a.closed_time,	a.is_used, 	a.status_id, c.name AS city_name,
		c.translation AS city_translation,	cat.name AS category_name,	cat.translation AS category_translation FROM public.advert a INNER JOIN city c ON a.city_id = c.id 
		INNER JOIN category cat ON a.category_id = cat.id WHERE 	a.user_id = %s AND a.status_id = %s; `, userId, statusID)

	rows, err := tx.Query(ctx, SQLGetAdvertsForUserWhereStatusIs, userId, statusID)
	if err != nil {
		ads.Logger.Errorf("Something went wrong while executing select adverts for user where status is, err=%v", err)

		return nil, err
	}
	defer rows.Close()

	returningAdvertList := models.ReturningAdvertList{}
	for rows.Next() {
		advert := models.Advert{}
		city := models.City{}
		category := models.Category{}
		var status uint // ЗАГЛУШКА

		if err := rows.Scan(&advert.ID, &advert.UserID, &advert.CityID, &advert.CategoryID, &advert.Title, &advert.Description, &advert.Price, &advert.CreatedTime, &advert.ClosedTime, &advert.IsUsed, &status,
			&city.CityName, &city.Translation, &category.Name, &category.Translation); err != nil {

			ads.Logger.Errorf("Something went wrong while scanning adverts rows, err=%v", err)

			return nil, err
		}
		advert.Deleted = false
		if status == 24 {
			advert.Deleted = true
		}
		city.ID = advert.CityID
		category.ID = advert.CategoryID
		returningAdvert := models.ReturningAdvert{
			Advert:   advert,
			City:     city,
			Category: category,
		}

		returningAdvertList.AdvertItems = append(returningAdvertList.AdvertItems, &returningAdvert)
	}

	if err := rows.Err(); err != nil {
		ads.Logger.Errorf("Something went wrong while scanning adverts rows, err=%v", err)

		return nil, err
	}

	return &returningAdvertList, nil
}

func (ads *AdvertsListWrapper) GetAdvertsForUserWhereStatusIs(ctx context.Context, userId, deleted uint) ([]*models.ReturningAdInList, error) {

	var advertsList []*models.ReturningAdInList

	err := pgx.BeginFunc(ctx, ads.Pool, func(tx pgx.Tx) error {
		advertsListInner, err := ads.getAdvertsForUserWhereStatusIs(ctx, tx, userId, deleted)
		for _, num := range advertsListInner.AdvertItems {
			returningAdInList := models.ReturningAdInList{
				ID:       num.Advert.ID,
				Title:    num.Advert.Title,
				Price:    num.Advert.Price,
				City:     num.City.Translation,
				Category: num.Category.Translation,
			}
			advertsList = append(advertsList, &returningAdInList)
		}

		return err
	})

	if err != nil {
		ads.Logger.Errorf("Something went wrong while getting adverts list, err=%v", err)

		return nil, err
	}

	return advertsList, nil
}

func (ads *AdvertsListWrapper) GetAdvertsByUserIDFiltered(userID uint, filter func(*models.Advert) bool) ([]*models.ReturningAdvert, error) {
	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()
	var returningAds []*models.ReturningAdvert
	for _, ad := range ads.AdvertsList.Adverts {
		if ad.UserID == userID && filter(ad) {
			returningAds = append(returningAds, &models.ReturningAdvert{
				Advert:   *ad,
				City:     *ads.AdvertsList.Cities[ad.CityID-1],
				Category: *ads.AdvertsList.Categories[ad.CategoryID-1],
			})
		}
	}
	return returningAds, nil
}

func (ads *AdvertsListWrapper) CreateAdvert(data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	cityID, err := ads.GetCityID(data.City)
	if err != nil {
		return nil, err
	}

	categoryID, err := ads.GetCategoryID(data.Category)
	if err != nil {
		return nil, err
	}

	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	newAd := &models.Advert{
		ID:          ads.GetLastAdvertID(),
		UserID:      data.UserID,
		CityID:      cityID,
		CategoryID:  categoryID,
		Title:       data.Title,
		Description: data.Description,
		Price:       data.Price,
		CreatedTime: time.Now(),
		Active:      true,
		IsUsed:      data.IsUsed,
	}

	ads.AdvertsList.Adverts = append(ads.AdvertsList.Adverts, newAd)

	return &models.ReturningAdvert{
		Advert:   *newAd,
		City:     *ads.AdvertsList.Cities[cityID-1],
		Category: *ads.AdvertsList.Categories[categoryID-1],
	}, nil
}

func (ads *AdvertsListWrapper) EditAdvert(data models.ReceivedAdData) (*models.ReturningAdvert, error) {
	id := data.ID
	if id > ads.AdvertsList.AdvertsCounter || ads.AdvertsList.Adverts[id-1].Deleted {
		return nil, errWrongAdvertID
	}

	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	ads.AdvertsList.Adverts[id-1] = &models.Advert{
		ID:          id,
		UserID:      data.UserID,
		Title:       data.Title,
		Description: data.Description,
		Price:       data.Price,
		CityID:      ads.AdvertsList.Adverts[id-1].CityID,
		CategoryID:  ads.AdvertsList.Adverts[id-1].CategoryID,
		CreatedTime: ads.AdvertsList.Adverts[id-1].CreatedTime,
		Active:      true,
		IsUsed:      data.IsUsed,
		Deleted:     false,
	}

	return &models.ReturningAdvert{
		Advert:   *ads.AdvertsList.Adverts[id-1],
		Category: *ads.AdvertsList.Categories[ads.AdvertsList.Adverts[id-1].CategoryID-1],
		City:     *ads.AdvertsList.Cities[ads.AdvertsList.Adverts[id-1].CityID-1],
	}, nil
}

func (ads *AdvertsListWrapper) CloseAdvert(advertID uint) error {
	if advertID > ads.AdvertsList.AdvertsCounter || ads.AdvertsList.Adverts[advertID-1].Deleted {
		return errWrongAdvertID
	}

	if !ads.AdvertsList.Adverts[advertID-1].Active {
		return errAlreadyClosed
	}

	ads.AdvertsList.Adverts[advertID-1].Active = false

	return nil
}

func (ads *AdvertsListWrapper) DeleteAdvert(advertID uint) error {
	if advertID > ads.AdvertsList.AdvertsCounter || ads.AdvertsList.Adverts[advertID-1].Deleted {
		return errWrongAdvertID
	}

	ads.AdvertsList.Adverts[advertID-1].Deleted = true

	return nil
}

func (ads *AdvertsListWrapper) GetCityID(city string) (uint, error) {
	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	for _, val := range ads.AdvertsList.Cities {
		if val.CityName == city || val.Translation == city {
			return val.ID, nil
		}
	}

	return 0, errWrongCityName
}

func (ads *AdvertsListWrapper) GetCategoryID(category string) (uint, error) {
	ads.AdvertsList.Mux.Lock()
	defer ads.AdvertsList.Mux.Unlock()

	for _, val := range ads.AdvertsList.Categories {
		if val.Name == category || val.Translation == category {
			return val.ID, nil
		}
	}

	return 0, errWrongCategoryName
}

func (ads *AdvertsListWrapper) GetLastAdvertID() uint {
	ads.AdvertsList.AdvertsCounter++

	return ads.AdvertsList.AdvertsCounter
}

func (ads *AdvertsListWrapper) GetLastLocationID() uint {
	ads.AdvertsList.CitiesCounter++

	return ads.AdvertsList.CitiesCounter
}

func (ads *AdvertsListWrapper) GetLastCategoryID() uint {
	ads.AdvertsList.CategoriesCounter++

	return ads.AdvertsList.CategoriesCounter
}

func NewAdvertsList(pool *pgxpool.Pool, logger *zap.SugaredLogger) *AdvertsListWrapper {
	return &AdvertsListWrapper{
		AdvertsList: &models.AdvertsList{
			AdvertsCounter:    0,
			CitiesCounter:     0,
			CategoriesCounter: 0,
			Adverts:           make([]*models.Advert, 0),
			Cities:            make([]*models.City, 0),
			Categories:        make([]*models.Category, 0),
			Mux:               sync.RWMutex{},
		},
		Pool:   pool,
		Logger: logger,
	}
}

func FillAdvertsList(ads *AdvertsListWrapper) {
	locationID := ads.GetLastLocationID()
	ads.AdvertsList.Cities = append(ads.AdvertsList.Cities, &models.City{
		ID:          locationID,
		CityName:    "Москва",
		Translation: "Moscow",
	})

	categoryID := ads.GetLastCategoryID()
	ads.AdvertsList.Categories = append(ads.AdvertsList.Categories, &models.Category{
		ID:          categoryID,
		Name:        "Тест",
		Translation: translit.Ru("Тест"),
	})

	for i := 1; i <= 100; i++ {
		price, _ := rand.Int(rand.Reader, big.NewInt(int64(models.MaxPrice)))
		advertID := ads.GetLastAdvertID()
		ads.AdvertsList.Adverts = append(ads.AdvertsList.Adverts, &models.Advert{
			ID:          advertID,
			UserID:      1,
			Title:       fmt.Sprintf("Объявление № %d", advertID),
			Description: fmt.Sprintf("Текст в объявлениии № %d", advertID),
			Price:       uint(price.Uint64()) * advertID,
			CityID:      1,
			CategoryID:  1,
			CreatedTime: time.Now(),
			Active:      true,
			IsUsed:      true,
			Deleted:     false,
		})
	}
}

func AddAdvert(ads *models.AdvertsList, advert *models.Advert) {

	ads.Adverts = append(ads.Adverts, advert)

}
