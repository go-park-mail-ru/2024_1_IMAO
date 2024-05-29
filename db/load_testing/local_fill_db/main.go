//nolint:gosec
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserIDAdvertID struct {
	UserID   uint `json:"userId"`
	AdvertID uint `json:"advertId"`
}

func createFavouriteDataWithCopy(advertStartID, userNum, advertNum int) []UserIDAdvertID {
	var favourites []UserIDAdvertID

	startTime := time.Now() // Запоминаем время начала

	for i := advertStartID; i < (advertStartID + advertNum); i += 10 {
		for j := 1; j < userNum; j += 180 {
			userID := uint(j)
			advertID := uint(i)

			favourites = append(favourites, UserIDAdvertID{
				UserID:   userID,
				AdvertID: advertID,
			})

			fmt.Printf("Iteration (%d, %d), elapsed time since start: %v\n", i, j, time.Since(startTime))
		}
	}

	endTime := time.Now() // Запоминаем время окончания

	elapsedTime := endTime.Sub(startTime) // Вычисляем затраченное время

	fmt.Printf("Цикл выполнен за %v\n", elapsedTime)

	return favourites
}

func CreateFavouriteWithCopy(ctx context.Context, pool *pgxpool.Pool, advertStartID int) error {
	favourites := createFavouriteDataWithCopy(advertStartID, 1800, 100)

	// Преобразуем данные в формат, подходящий для pgx.CopyFromRows
	rows := make([][]interface{}, len(favourites))
	for i, favourite := range favourites {
		rows[i] = []interface{}{favourite.UserID, favourite.AdvertID}
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

func main() {
	connPool, err := pgxpool.NewWithConfig(context.Background(), pgxpoolconfig.PGXPoolConfig())
	if err != nil {
		log.Fatal("Error while creating connection pool to the database!!")
	}

	startTime := time.Now() // Запоминаем время начала

	for i := 550000; i < (999000); i += 100 {
		_ = CreateFavouriteWithCopy(context.Background(), connPool, i)
	}

	endTime := time.Now() // Запоминаем время окончания

	elapsedTime := endTime.Sub(startTime) // Вычисляем затраченное время

	fmt.Printf("Цикл выполнен за %v\n", elapsedTime)
}
