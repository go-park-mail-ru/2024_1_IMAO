package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ExecuteInsertQuery(connPool *pgxpool.Pool, query string) {
	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error while acquiring connection from the database pool!", err)
	}
	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("Could not ping database")
	}

	_, err = connPool.Exec(context.Background(), query)
	if err != nil {
		log.Fatal("Error while executing the query", err)
	}
}

func ExecuteSelectQuery(connPool *pgxpool.Pool, query string) (pgx.Rows, error) {
	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Error while acquiring connection from the database pool: %w", err)
	}
	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Error while connecting the database: %w", err)
	}

	rows, err := connPool.Query(context.Background(), query)

	if err != nil {
		return nil, fmt.Errorf("Error while executing the SELECT query: %w", err)
	}
	defer rows.Close()
	return rows, nil
}

func ExecuteSelectQueryBool(connPool *pgxpool.Pool, query string) ([]bool, error) {
	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error while acquiring connection from the database pool!", err)
	}
	//defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("Could not ping database")
	}

	fmt.Println("Connected to the database!")

	rows, err := connPool.Query(context.Background(), query)
	fmt.Println("rows", rows)
	if err != nil {
		return nil, fmt.Errorf("Error while executing the SELECT query: %w", err)
	}
	//defer rows.Close()

	var results []bool
	for rows.Next() {
		var result bool
		fmt.Println("result")
		if err := rows.Scan(&result); err != nil {
			return nil, fmt.Errorf("Error scanning row: %w", err)
		}
		results = append(results, result)
	}

	fmt.Println("results")

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over rows: %w", err)
	}

	return results, nil
}
