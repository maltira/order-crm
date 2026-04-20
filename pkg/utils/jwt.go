package utils

import (
	"errors"
	"order-crm/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID int, roleCode string, roleID int) (string, error) {
	access := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":   userID,
			"role_id":   roleID,
			"role_code": roleCode,
			"exp":       time.Now().Add(config.Env.AccessTokenDuration).Unix(),
		},
	)
	return access.SignedString(config.Env.JWTSecret)
}

func ValidateAccessToken(accessToken string) (jwt.MapClaims, error) {
	if accessToken == "" {
		return nil, errors.New("token is empty")
	}
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		return config.Env.JWTSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

func GenerateRefreshToken(userID int) (string, error) {
	access := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(config.Env.RefreshTokenDuration).Unix(),
		},
	)
	return access.SignedString(config.Env.JWTSecret)
}

func ValidateRefreshToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return config.Env.JWTSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := int(claims["user_id"].(float64))
		return userId, nil
	}

	return 0, errors.New("недействительный refresh токен")
}
