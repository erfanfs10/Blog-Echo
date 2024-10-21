package utils

import "github.com/golang-jwt/jwt/v5"

type JwtCustomClaims struct {
	UserID    int    `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}
