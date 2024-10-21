package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateAccessToken(accessToken string) (bool, int) {
	token, err := jwt.ParseWithClaims(accessToken, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("My$Super@Secret"), nil
	})

	if err != nil {
		return false, 0
	}
	claims, ok := token.Claims.(*JwtCustomClaims)

	if !ok {
		return false, 0
	}

	if claims.TokenType != "access" {
		return false, 0
	}
	return true, claims.UserID
}

func ValidateRefreshToken(refreshToken string) (bool, int) {
	token, err := jwt.ParseWithClaims(refreshToken, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("My$Super@Secret"), nil
	})

	if err != nil {
		return false, 0
	}
	claims, ok := token.Claims.(*JwtCustomClaims)

	if !ok {
		return false, 0
	}

	if claims.TokenType != "refresh" {
		return false, 0
	}
	return true, claims.UserID
}
