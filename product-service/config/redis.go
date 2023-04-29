package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Redis *redis.Client
)

func ConnectToRedis() error {
	var counts int64
	for {
		rdb := redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
		_, err := rdb.Ping(context.Background()).Result()
		if err == nil {
			log.Println("Connected to Redis!")
			Redis = rdb
			return nil
		}

		log.Println("Redis not yet ready ...")
		counts++

		if counts > 20 {
			return err
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
