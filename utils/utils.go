package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/rainbow96bear/planet_user_server/config"
)

func GetUuidByAccessToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("fail to claims")
	}

	userUuid, ok := claims["userUuid"].(string)
	if !ok {
		return "", fmt.Errorf("fail to get uuid from jwt")
	}

	return userUuid, nil
}
