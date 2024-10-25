package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/erfanfs10/Blog-Echo/models"
	"github.com/golang-jwt/jwt/v5"
)

type jwtCustomClaims struct {
	UserID    int    `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int) (models.TokenModel, error) {
	// create access token claims
	accessClaims := &jwtCustomClaims{
		userID,
		"access",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	// create refresh token cliams
	refreshClaims := &jwtCustomClaims{
		userID,
		"refresh",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	// create access and refresh tokens
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	// signed access token with a strong secret
	aToken, err := accessToken.SignedString([]byte("My$Super@Secret"))
	if err != nil {
		return models.TokenModel{}, nil
	}
	// signed refresh token with a strong secret
	rToken, err := refreshToken.SignedString([]byte("My$Super@Secret"))
	if err != nil {
		return models.TokenModel{}, nil
	}
	// create TokenModel with both tokens
	tokenModel := models.TokenModel{
		AccessToken:  aToken,
		RefreshToken: rToken,
	}
	// return them
	return tokenModel, nil
}

func ValidateAccessToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("My$Super@Secret"), nil
	})

	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*jwtCustomClaims)

	if !ok {
		return 0, err
	}

	if claims.TokenType != "access" {
		return 0, err
	}
	return claims.UserID, nil
}

func ValidateRefreshToken(refreshToken string) (int, error) {
	// parse and validate sent refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("My$Super@Secret"), nil
	})
	if err != nil {
		return 0, err
	}
	// extract the data to jwtCustomClaims struct
	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok {
		return 0, err
	}
	// check the token type must be refresh
	if claims.TokenType != "refresh" {
		return 0, errors.New("invalid refresh token")
	}
	return claims.UserID, nil
}
