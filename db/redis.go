package db

import (
	"github.com/go-redis/redis"
	"net/url"
	"os"
	"strconv"
)

var R *redis.Client

func getOptions() (*redis.Options, error) {
	redisUrl, err := url.Parse(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}
	redisUrlPassword, ok := redisUrl.User.Password()
	redisPassword := ""
	if ok {
		redisPassword = redisUrlPassword
	}
	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}
	return &redis.Options{
		Addr:     redisUrl.Host,
		Password: redisPassword,
		DB:       redisDb,
	}, nil
}

func init() {
	opt, err := getOptions()
	if err != nil {
		panic(err)
	}

	R = redis.NewClient(opt)
}
