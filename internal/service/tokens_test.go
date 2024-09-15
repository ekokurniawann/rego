package service

import (
	"belajar-redis/infra"
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var redisClient *redis.Client

func setup() {
	redisClient = infra.GetRedisClient()
}

func teardown() {
	redisClient.FlushDB(context.Background())
	redisClient.Close()
}

var Ctx = context.Background()
var userId = "12345"
var userRole = "admin"

func TestCreateAccess_Valid(t *testing.T) {
	setup()
	defer teardown()

	tokenService := NewTokenService(NewRedisService(redisClient))

	tokenString, expiry, err := tokenService.CreateAccess(&Ctx, &userId, &userRole)

	assert.NoError(t, err)
	assert.NotNil(t, tokenString)
	assert.NotEmpty(t, *tokenString)
	assert.NotNil(t, expiry)
	assert.WithinDuration(t, time.Now().Add(15*time.Minute), expiry.Time, 1*time.Second)

	cachedToken, err := redisClient.Get(context.Background(), "access-token-"+userId).Result()
	assert.NoError(t, err)
	assert.NotEmpty(t, cachedToken)
}

func TestCreateAccess_RedisError(t *testing.T) {
	invalidRedisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   9999,
	})

	tokenService := NewTokenService(NewRedisService(invalidRedisClient))

	tokenString, expiry, err := tokenService.CreateAccess(&Ctx, &userId, &userRole)

	assert.Error(t, err)
	assert.Nil(t, tokenString)
	assert.Nil(t, expiry)
}

func TestParseAccess_ValidToken(t *testing.T) {
	setup()
	defer teardown()

	tokenService := NewTokenService(NewRedisService(redisClient))

	tokenString, _, err := tokenService.CreateAccess(&Ctx, &userId, &userRole)

	assert.NoError(t, err)
	assert.NotNil(t, tokenString)

	claims, err := tokenService.ParseAccess(*tokenString)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userId, claims.UserId)
	assert.Equal(t, userRole, claims.RoleCode)
}

func TestParseAccess_InvalidToken(t *testing.T) {
	invalidToken := "30608837-e92f-4aad-b45a-ea2d1e4bcd98"

	tokenService := NewTokenService(NewRedisService(redisClient))
	claims, err := tokenService.ParseAccess(invalidToken)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse token")
	assert.Nil(t, claims)
}
