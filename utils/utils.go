package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

	userUuid, ok := claims["user_uuid"].(string)
	if !ok {
		return "", fmt.Errorf("fail to get uuid from jwt")
	}

	return userUuid, nil
}

func GetUserUuid(c *gin.Context) (string, error) {
	authUuidValue, exists := c.Get("user_uuid")
	if !exists {
		err := fmt.Errorf("unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return "", err
	}

	authUuid, ok := authUuidValue.(string)
	if !ok {
		err := fmt.Errorf("invalid user uuid type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return "", err
	}
	return authUuid, nil
}

func GetUserNickname(c *gin.Context) (string, error) {
	authNicknameValue, exists := c.Get("nickname")
	if !exists {
		err := fmt.Errorf("unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return "", err
	}

	authNickname, ok := authNicknameValue.(string)
	if !ok {
		err := fmt.Errorf("invalid user nickname type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return "", err
	}
	return authNickname, nil
}
