package initializers

import (
	"belajar-redis/internal/service"
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	RedisClient *redis.Client
	SQLDB       *sql.DB
	SQLxDB      *sqlx.DB
}

func InitializeDB(config Config) {

	if config.RedisClient != nil {
		redisService := service.NewRedisService(config.RedisClient)

		tokenService := service.NewTokenService(redisService)
		userId := "12345"
		userRole := "admin"
		ctx := context.Background()

		tokenString, expiry, err := tokenService.CreateAccess(&ctx, &userId, &userRole)
		if err != nil {
			log.Fatalf("Error creating access token: %v", err)
		}

		fmt.Printf("Token: %s\n", *tokenString)
		fmt.Printf("Expires At: %v\n", expiry)
	}
}
