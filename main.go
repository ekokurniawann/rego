package main

import (
	"belajar-redis/infra"
	"belajar-redis/internal/initializers"
)

func main() {
	redisClient := infra.GetRedisClient()

	initializers.InitializeDB(initializers.Config{
		RedisClient: redisClient,
	})

}
