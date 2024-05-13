package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserIdAdvertId struct {
	UserId   uint `json:"userId"`
	AdvertId uint `json:"advertId"`
}

func createAdvert(ctx context.Context, tx pgx.Tx) error {

	user_id := uint(rand.Intn(10) + 1)
	city_id := uint(521)
	category_id := uint(rand.Intn(16) + 1)
	price := uint(rand.Intn(9000) + 1000)
	title := utils.RandString(10)
	description := utils.RandString(50)

	SQLCreateAdvert :=
		`INSERT INTO public.advert(
			user_id, city_id, category_id, title, description, price)
			VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := tx.Exec(ctx, SQLCreateAdvert, user_id, city_id, category_id, title, description, price)

	if err != nil {
		fmt.Println("Something went wrong 32")

		return err
	}

	return nil
}

func CreateAdvert(ctx context.Context, pool *pgxpool.Pool) error {

	err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		err := createAdvert(ctx, tx)

		return err
	})

	if err != nil {
		fmt.Println("Something went wrong 49 ", err)

		return err
	}

	return nil
}

func createAdvertDataWithCopy(n int) ([]models.DBInsertionAdvert, error) {
	var adverts []models.DBInsertionAdvert

	for i := 0; i < n; i++ {
		user_id := uint(rand.Intn(10) + 1)
		city_id := uint(521)
		category_id := uint(rand.Intn(16) + 1)
		price := uint(rand.Intn(9000) + 1000)
		title := utils.RandString(10)
		description := utils.RandString(50)

		adverts = append(adverts, models.DBInsertionAdvert{
			UserID:      user_id,
			CityID:      city_id,
			CategoryID:  category_id,
			Title:       title,
			Description: description,
			Price:       price,
		})
	}

	return adverts, nil
}

func createFavouriteDataWithCopy(advert_start_id, user_num, advert_num int) ([]UserIdAdvertId, error) {
	var favourites []UserIdAdvertId

	startTime := time.Now() // Запоминаем время начала
	for i := advert_start_id; i < (advert_start_id + advert_num); i += 10 {
		for j := 1; j < user_num; j += 180 {
			user_id := uint(j)
			advert_id := uint(i)

			favourites = append(favourites, UserIdAdvertId{
				UserId:   user_id,
				AdvertId: advert_id,
			})

			fmt.Printf("Iteration (%d, %d), elapsed time since start: %v\n", i, j, time.Since(startTime))
		}
	}
	endTime := time.Now() // Запоминаем время окончания

	elapsedTime := endTime.Sub(startTime) // Вычисляем затраченное время

	fmt.Printf("Цикл выполнен за %v\n", elapsedTime)

	return favourites, nil
}

func CreateFavouriteWithCopy(ctx context.Context, pool *pgxpool.Pool, advert_start_id int) error {
	favourites, err := createFavouriteDataWithCopy(advert_start_id, 1800, 100)
	if err != nil {
		return fmt.Errorf("failed to create advert data: %w", err)
	}

	// Преобразуем данные в формат, подходящий для pgx.CopyFromRows
	rows := make([][]interface{}, len(favourites))
	for i, favourite := range favourites {
		rows[i] = []interface{}{favourite.UserId, favourite.AdvertId}
	}

	fmt.Println(rows)

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	_, err = conn.CopyFrom(ctx,
		pgx.Identifier{"favourite"},
		[]string{"user_id", "advert_id"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("something went wrong while copy rows: %w", err)
	}

	return nil
}

func createProfileDataWithCopy(initial_id, n int) ([]models.DBInsertionProfile, error) {
	var profiles []models.DBInsertionProfile

	for i := initial_id; i < (n + initial_id); i++ {
		user_id := uint(i)
		city_id := uint(521)
		phone := utils.RandString(12)
		name := utils.RandString(15)
		surname := utils.RandString(15)

		profiles = append(profiles, models.DBInsertionProfile{
			UserID:  user_id,
			CityID:  city_id,
			Phone:   phone,
			Name:    name,
			Surname: surname,
		})
	}

	return profiles, nil
}

func CreateProfileWithCopy(ctx context.Context, pool *pgxpool.Pool) error {
	profiles, err := createProfileDataWithCopy(12, 1850)
	if err != nil {
		return fmt.Errorf("failed to create advert data: %w", err)
	}

	// Преобразуем данные в формат, подходящий для pgx.CopyFromRows
	rows := make([][]interface{}, len(profiles))
	for i, profile := range profiles {
		rows[i] = []interface{}{profile.UserID, profile.CityID, profile.Phone, profile.Name, profile.Surname}
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	_, err = conn.CopyFrom(ctx,
		pgx.Identifier{"profile"},
		[]string{"user_id", "city_id", "phone", "name", "surname"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("something went wrong while copy rows: %w", err)
	}

	return nil
}

func createUserDataWithCopy(initial_id, n int) ([]models.DBInsertionUser, error) {
	var users []models.DBInsertionUser

	startTime := time.Now() // Запоминаем время начала
	for i := initial_id; i < (n + initial_id); i++ {
		email := utils.RandString(12)
		passwordHash := utils.HashPassword(utils.RandString(8))
		fmt.Printf("Iteration %d, elapsed time since start: %v\n", i, time.Since(startTime))
		users = append(users, models.DBInsertionUser{
			Email:        email,
			PasswordHash: passwordHash,
		})
	}

	endTime := time.Now() // Запоминаем время окончания

	elapsedTime := endTime.Sub(startTime) // Вычисляем затраченное время

	fmt.Printf("Цикл выполнен за %v\n", elapsedTime)

	return users, nil
}

func CreateUserWithCopy(ctx context.Context, pool *pgxpool.Pool) error {
	users, err := createUserDataWithCopy(11, 889)
	if err != nil {
		return fmt.Errorf("failed to create advert data: %w", err)
	}

	// Преобразуем данные в формат, подходящий для pgx.CopyFromRows
	rows := make([][]interface{}, len(users))
	for i, user := range users {
		rows[i] = []interface{}{user.Email, user.PasswordHash}
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	_, err = conn.CopyFrom(ctx,
		pgx.Identifier{"user"},
		[]string{"email", "password_hash"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("something went wrong while copy rows: %w", err)
	}

	return nil
}

func CreateAdvertWithCopy(ctx context.Context, pool *pgxpool.Pool) error {
	adverts, err := createAdvertDataWithCopy(1000)
	if err != nil {
		return fmt.Errorf("failed to create advert data: %w", err)
	}

	// Преобразуем данные в формат, подходящий для pgx.CopyFromRows
	rows := make([][]interface{}, len(adverts))
	for i, advert := range adverts {
		rows[i] = []interface{}{advert.UserID, advert.CityID, advert.CategoryID, advert.Title, advert.Description, advert.Price}
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	_, err = conn.CopyFrom(ctx,
		pgx.Identifier{"advert"},
		[]string{"user_id", "city_id", "category_id", "title", "description", "price"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("something went wrong while copy rows: %w", err)
	}

	return nil
}

func main() {
	connPool, err := pgxpool.NewWithConfig(context.Background(), pgxpoolconfig.PGXPoolConfig())
	if err != nil {
		log.Fatal("Error while creating connection pool to the database!!")
	}

	startTime := time.Now() // Запоминаем время начала

	for i := 550000; i < (999000); i += 100 {
		_ = CreateFavouriteWithCopy(context.Background(), connPool, i)
	}

	//_ = CreateUserWithCopy(context.Background(), connPool)

	//_ = CreateProfileWithCopy(context.Background(), connPool)

	// for i := 0; i < 10; i++ {
	// 	fmt.Printf("Iteration %d, elapsed time since start: %v\n", i, time.Since(startTime))
	// 	//_ = CreateAdvert(context.Background(), connPool)
	// 	_ = CreateAdvertWithCopy(context.Background(), connPool)
	// }

	endTime := time.Now() // Запоминаем время окончания

	elapsedTime := endTime.Sub(startTime) // Вычисляем затраченное время

	fmt.Printf("Цикл выполнен за %v\n", elapsedTime)

}
