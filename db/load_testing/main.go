package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	pgxpoolconfig "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/server/repository"
	"github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

func main() {
	connPool, err := pgxpool.NewWithConfig(context.Background(), pgxpoolconfig.PGXPoolConfig())
	if err != nil {
		log.Fatal("Error while creating connection pool to the database!!")
	}

	startTime := time.Now() // Запоминаем время начала

	for i := 0; i < 10000; i++ {
		fmt.Println("Iteration ", i)
		_ = CreateAdvert(context.Background(), connPool)
	}

	endTime := time.Now() // Запоминаем время окончания

	elapsedTime := endTime.Sub(startTime) // Вычисляем затраченное время

	fmt.Printf("Цикл выполнен за %v\n", elapsedTime)

}
