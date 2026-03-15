package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        os.Getenv("REDIS_CONN"),
		Password:    "",
		DB:          0,
		ReadTimeout: 2 * time.Second,
	})

	fmt.Println(rdb.Ping(context.Background()).Result())
	return rdb, nil
}
