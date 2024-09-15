package entity

import "github.com/golang-jwt/jwt/v5"

type AccessTokenClaims struct {
	UserId   string `json:"user_id"`
	RoleCode string `json:"role_code"`
	UUID     string `json:"id"`
	jwt.RegisteredClaims
}

type AccessTokenCached struct {
	AccessUID string `json:"access"`
}

type RefreshTokenClaims struct {
	UserId   string `json:"user_id"`
	RoleCode string `json:"role_code"`
	UUID     string `json:"id"`
	jwt.RegisteredClaims
}

type RefreshTokenCached struct {
	RefreshUID string `json:"refresh"`
}

type ResetTokenClaims struct {
	UserId string `json:"user_id"`
	UUID   string `json:"id"`
	jwt.RegisteredClaims
}

type ResetTokenCached struct {
	ResetUID string `json:"reset"`
}
