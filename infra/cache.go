package infra

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func GetRedisClient() *redis.Client {
	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Username: "default",
			Password: "secret",
			DB:       0,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := rdb.Ping(ctx).Err(); err != nil {
			panic(fmt.Sprintf("redis: can't ping to redi - %v", err))
		}

		fmt.Println("redis: connected to redis")
		redisClient = rdb
	})

	return redisClient
}
