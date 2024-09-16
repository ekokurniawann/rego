package service

import (
	"belajar-redis/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{client: client}
}

func (r *RedisService) SetAccessToken(ctx context.Context, userId, tokenUUID string) error {
	cachedJson, err := json.Marshal(entity.AccessTokenCached{
		AccessUID: tokenUUID,
	})
	if err != nil {
		return fmt.Errorf("can't marshal access token: %w", err)
	}

	if err := r.client.Set(ctx, fmt.Sprintf("access-token-%s", userId), string(cachedJson), time.Minute*15).Err(); err != nil {
		return fmt.Errorf("can't cache access token: %w", err)
	}

	return nil
}

func (r *RedisService) GetAccessToken(ctx context.Context, userId string) (string, error) {
	cachedJson, err := r.client.Get(ctx, fmt.Sprintf("access-token-%s", userId)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("faied to fetch token from cache: %w", err)
	}

	return cachedJson, nil
}
