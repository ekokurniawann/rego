package service

import (
	"belajar-redis/internal/entity"
	"context"
	"encoding/json"
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

func (t *TokenService) ParseAccess(tokenString string) (*entity.AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*entity.AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims or token is not valid")
}

func (t *TokenService) ValidateAccess(ctx *context.Context, claims *entity.AccessTokenClaims) error {
	cacheJSON, err := t.redisService.GetAccessToken(*ctx, claims.UserId)
	if err != nil {
		return fmt.Errorf("error retrieving access token from cache: %w", err)
	}

	if cacheJSON == "" {
		return fmt.Errorf("access token not found for user ID: %s", claims.UserId)
	}

	cachedTokens := &entity.AccessTokenCached{}
	err = json.Unmarshal([]byte(cacheJSON), cachedTokens)
	if err != nil {
		return fmt.Errorf("error unmarshalling cached token data: %w", err)
	}

	if cachedTokens.AccessUID != claims.UUID {
		return fmt.Errorf("token UUID mismatch for user ID: %s", claims.UserId)
	}

	return nil
}
