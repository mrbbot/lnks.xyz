package db

import (
	"github.com/go-redis/redis"
	"os"
	"strconv"
)

var R *redis.Client

func init() {
	db, err := strconv.Atoi(os.Getenv("S_REDIS_DB"))
	if err != nil {
		panic(err)
	}

	R = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("S_REDIS_ADDRESS"),
		Password: os.Getenv("S_REDIS_PASSWORD"),
		DB:       db,
	})
}
