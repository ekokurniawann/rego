package service

import (
	"belajar-redis/internal/entity"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var mySigningKey = []byte("ikanlele")

type TokenService struct {
	redisService *RedisService
}

func NewTokenService(redisService *RedisService) *TokenService {
	return &TokenService{redisService: redisService}
}

func (t *TokenService) CreateAccess(ctx *context.Context, userId, userRole *string) (*string, *jwt.NumericDate, error) {
	tokenUUID := uuid.NewString()
	claims := entity.AccessTokenClaims{
		UserId:   *userId,
		RoleCode: *userRole,
		UUID:     tokenUUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return nil, nil, fmt.Errorf("can't signed the token")
	}

	tokenString := &ss

	if err := t.redisService.SetAccessToken(*ctx, *userId, tokenUUID); err != nil {
		return nil, nil, fmt.Errorf("can't cache access token: %w", err)
	}

	return tokenString, claims.ExpiresAt, nil
}
