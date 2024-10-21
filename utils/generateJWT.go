package utils

import (
	"time"

	"github.com/erfanfs10/Blog-Echo/models"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID int) (models.TokenModel, error) {

	accessClaims := &JwtCustomClaims{
		userID,
		"access",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	refreshClaims := &JwtCustomClaims{
		userID,
		"refresh",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	aToken, err := accessToken.SignedString([]byte("My$Super@Secret"))
	if err != nil {
		return models.TokenModel{}, nil
	}

	rToken, err := refreshToken.SignedString([]byte("My$Super@Secret"))
	if err != nil {
		return models.TokenModel{}, nil
	}
	tokenModel := &models.TokenModel{
		AccessToken:  aToken,
		RefreshToken: rToken,
	}
	return *tokenModel, nil
}
