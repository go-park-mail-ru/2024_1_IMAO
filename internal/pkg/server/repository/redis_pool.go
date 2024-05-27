package repository

import (
	"github.com/gomodule/redigo/redis"
	"log"
)

func NewRedisPool(host, password string) *redis.Pool {
	return &redis.Pool{
		MaxActive: 100,
		MaxIdle:   100,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host, redis.DialPassword(password))
			if err != nil {
				log.Println("Failed to connect Redis", err)
				return nil, err
			}

			return conn, nil
		},
	}
}
