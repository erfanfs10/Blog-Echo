package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/erfanfs10/Blog-Echo/models"
	"github.com/golang-jwt/jwt/v5"
)

type jwtCustomClaims struct {
	UserID    int32  `json:"user_id"`
	IsActive  bool   `json:"is_active"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int32, IsActive bool) (models.TokenModel, error) {
	// create access token claims
	accessClaims := &jwtCustomClaims{
		userID,
		IsActive,
		"access",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	// create refresh token cliams
	refreshClaims := &jwtCustomClaims{
		userID,
		IsActive,
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

func ValidateAccessToken(accessToken string) (int32, bool, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("My$Super@Secret"), nil
	})
	if err != nil {
		return 0, false, errors.New("invalid token parse")
	}
	// extract the data to jwtCustomClaims struct
	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok {
		return 0, false, errors.New("invalid token claims")
	}
	// check the token type must be access
	if claims.TokenType != "access" {
		return 0, false, errors.New("invalid token type")
	}
	return claims.UserID, claims.IsActive, nil
}

func ValidateRefreshToken(refreshToken string) (int32, bool, error) {
	// parse and validate sent refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("My$Super@Secret"), nil
	})
	if err != nil {
		return 0, false, err
	}
	// extract the data to jwtCustomClaims struct
	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok {
		return 0, false, err
	}
	// check the token type must be refresh
	if claims.TokenType != "refresh" {
		return 0, false, errors.New("invalid refresh token")
	}
	return claims.UserID, claims.IsActive, nil
}
